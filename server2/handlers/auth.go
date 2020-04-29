package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"../which" // TODO: import from github
)

// == Auth root handler ================================

type Auth struct {
	LoginHandler    *Login
	CallbackHandler *Callback
	StatusHandler   *Status
	LogoutHandler   *Logout

	config *oauth2.Config
}

func NewAuth(userStore which.UserStore, sessionStore which.SessionStore) *Auth {
	auth := &Auth{}

	auth.config = &oauth2.Config{ // TODO: read from config
		RedirectURL:  "http://localhost:8080/auth/callback",
		ClientID:     os.Getenv("WHICH_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("WHICH_GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	auth.LoginHandler = &Login{
		userStore:    userStore,
		sessionStore: sessionStore,
		config:       auth.config,
	}

	auth.CallbackHandler = &Callback{
		userStore:    userStore,
		sessionStore: sessionStore,
		config:       auth.config,
	}

	auth.StatusHandler = &Status{
		sessionStore: sessionStore,
	}

	return auth
}

func (handler *Auth) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	head, _ := shiftPath(req.URL.Path)
	if head == "login" {
		handler.LoginHandler.ServeHTTP(resp, req)
	} else if head == "callback" {
		handler.CallbackHandler.ServeHTTP(resp, req)
	} else if head == "status" {
		handler.StatusHandler.ServeHTTP(resp, req)
	} else if head == "logout" {
		handler.LogoutHandler.ServeHTTP(resp, req)
	} else {
		http.Error(resp, "auth endpoint does not exist", 404)
	}
}

// == Login handler ================================

// == Callback handler ================================

// == Status handler ================================

type Status struct {
	sessionStore which.SessionStore
}

func (handler *Status) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_, err := sessionFromRequest(req, handler.sessionStore)
	log.Println(err)
	fmt.Fprintf(resp, "{\"authed\": \"%v\"}\n", err == nil)
}

// == Logout handler ================================

type Logout struct{}

func (handler *Logout) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	addCookie(resp, "session", "", time.Now())
	http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
}

// == Shared helpers ================================

// TODO: this method of state tracking seems really overcomplicated
type state struct {
	QueryString string `json:"query"`
	RandomState string `json:"state"`
}

// encodes state to JSON to base64
func encodeState(s state) (string, error) {
	stateString, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString([]byte(stateString)), nil
}

// decodes state from JSON from base64
func decodeState(stateB64 string) (state, error) {
	queryState := state{}
	stateBytes, err := base64.URLEncoding.DecodeString(stateB64)
	if err != nil {
		return queryState, err
	}
	json.Unmarshal(stateBytes, &queryState)
	return queryState, nil
}

// checks that req cookie and formvalue states match,
// then returns that state
func stateFromCallback(req *http.Request) (state, error) {
	stateCookie, err := req.Cookie("state")
	if err != nil {
		return state{}, err
	}
	if req.FormValue("state") != stateCookie.Value {
		return state{}, err
	}
	return decodeState(stateCookie.Value)
}
