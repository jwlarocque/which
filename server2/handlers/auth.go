package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type Login struct {
	userStore    which.UserStore
	sessionStore which.SessionStore
	config       *oauth2.Config
}

func (handler *Login) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_, err := sessionFromRequest(req, handler.sessionStore)
	if err != nil {
		// no valid session, redirect to oauth
		state, err := encodeState(state{QueryString: req.URL.Query().Encode(), RandomState: randomString()})
		if err != nil {
			log.Printf("Unable to generate random state string: %v\n", err)
			http.Redirect(resp, req, "/", http.StatusInternalServerError)
		}
		addCookie(resp, "state", state, time.Now().Add(time.Minute)) // TODO: abstract out this duration?
		http.Redirect(resp, req, handler.config.AuthCodeURL(state), http.StatusTemporaryRedirect)
	} else {
		// already logged in
		http.Redirect(resp, req, "/", http.StatusOK)
	}
}

// == Callback handler ================================

type Callback struct {
	userStore    which.UserStore
	sessionStore which.SessionStore
	config       *oauth2.Config
}

// TODO: this function is still rather long
func (handler *Callback) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// check that callback and cookie states match
	s, err := stateFromCallback(req)
	if err != nil {
		log.Printf("State cookie/formvalue mismatch, possibly due to CSRF attack: %v\n", err)
		http.Redirect(resp, req, "/", http.StatusForbidden)
		return
	}

	// cookie matches state, proceed
	user, err := createUserFromCallback(req, handler.config)
	if err != nil {
		log.Printf("Unable to create User from request: %v\n", err)
		http.Redirect(resp, req, "/", http.StatusInternalServerError)
		return
	}

	err = handler.userStore.Insert(user)
	if err != nil {
		log.Printf("Unable to add User to store: %v\n", err)
		http.Redirect(resp, req, "/", http.StatusInternalServerError)
		return
	}

	session := which.Session{
		ID:      randomString(),
		UserID:  user.ID,
		Created: time.Now(),
		Expires: time.Now().Add(which.CookieDuration),
	}
	err = handler.sessionStore.Insert(session)
	if err != nil {
		log.Printf("Unable to add Session to store: %v\n", err)
		http.Redirect(resp, req, "/", http.StatusInternalServerError)
		return
	}

	addCookie(resp, "session", session.ID, session.Expires)
	http.Redirect(resp, req, "/?"+s.QueryString, http.StatusTemporaryRedirect)
}

// TODO: passing config down through both of these function is awkward
func createUserFromCallback(req *http.Request, config *oauth2.Config) (which.User, error) {
	userData, err := getUserData(req.FormValue("code"), config)
	if err != nil {
		return which.User{}, err
	}
	var user which.User
	err = json.Unmarshal(userData, &user)
	if err != nil {
		return which.User{}, err
	}
	return user, nil
}

// getUserData exchanges an OAuth code for user data from Google
func getUserData(authCode string, config *oauth2.Config) ([]byte, error) {
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %s", err.Error())
	}
	response, err := http.Get(which.OauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo from Google Oauth: %s", err.Error())
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body) // TODO: improve this?  why using ioutil.ReadAll?
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err.Error())
	}
	return data, nil
}

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

func randomString() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("failed to generate random string: %v\n", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// returns a Session if the request has a valid session cookie, else error
// TODO: somewhat awkward - retrieves session from store using id from request data
func sessionFromRequest(req *http.Request, store which.SessionStore) (which.Session, error) {
	sessionCookie, err := req.Cookie("session")
	if err != nil {
		return which.Session{}, fmt.Errorf("failed to get session cookie from request: %v", err)
	}
	session, err := store.Fetch(sessionCookie.Value)
	if err != nil {
		return which.Session{}, fmt.Errorf("no session matching cookie: %v", err)
	}
	if session.ID != sessionCookie.Value {
		return which.Session{}, fmt.Errorf("session ID from db didn't match cookie (this should be impossible)")
	}
	if time.Now().After(session.Expires) {
		return which.Session{}, fmt.Errorf("session expired")
	}
	return session, nil
}
