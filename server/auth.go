package main

// authentication stuff
// TODO: some kind of cron job to remove expired sessions from the database
// TODO: reorganize this file, it's getting a bit out of hand (maybe it can be a package)
// TODO: review auth security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// -- Auth root handler --

type AuthHandler struct {
	LoginHandler    *LoginHandler
	CallbackHandler *CallbackHandler
	StatusHandler   *StatusHandler
	LogoutHandler   *LogoutHandler
}

func (handler *AuthHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	head, _ := shiftPath(req.URL.Path)
	if head == "login" {
		handler.LoginHandler.ServeHTTP(resp, req) // TODO: req URL correct?
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

// -- Login handler (redirect to OAuth) --

type LoginHandler struct{}

func (handler *LoginHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if len(userIDFromSession(req)) > 0 {
		http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
	} else {
		state, err := encodeState(state{QueryString: req.URL.Query().Encode(), RandomState: randomString()})
		if err != nil {
			log.Fatalf("Unable to generate random state string: %v\n", err)
		}
		addCookie(resp, "state", state, time.Minute) // TODO: abstract out this duration?
		url := googleAuthConfig.AuthCodeURL(state)
		http.Redirect(resp, req, url, http.StatusTemporaryRedirect)
	}
}

// TODO: consider error handling other than returning false
func userIDFromSession(req *http.Request) string {
	log.Println("checking if logged in...")
	sessionCookie, err := req.Cookie("session")
	if err != nil {
		log.Println("failed to get session cookie from request")
		return ""
	}
	session, err := fetchSession(sessionCookie.Value)
	if err != nil {
		log.Printf("no session matching cookie: %s\n", sessionCookie.Value)
		return ""
	}
	if session.ID != sessionCookie.Value {
		log.Println("session ID from db didn't match cookie (this should be impossible)")
		return ""
	}
	if time.Now().After(session.Expires) {
		log.Println("session expired")
		return ""
	}
	return session.User_ID
}

// -- Callback (from OAuth) handler --

type CallbackHandler struct{}

// TODO: this function is way out of hand
func (handler *CallbackHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	stateCookie, err := req.Cookie("state")
	if err != nil {
		log.Printf("Auth callback missing cookie: %v\n", err)
		http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
	}
	if req.FormValue("state") != stateCookie.Value {
		// TODO: alert user to possible CSRF attack?
		http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
	} else { // cookie matches state, proceed
		data, err := getUserData(req.FormValue("code"))
		if err != nil {
			log.Println(err.Error())
			http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
		} else {
			sessionID := createSession(data)
			addCookie(resp, "session", sessionID, cookieDuration)
			s, err := decodeState(req.FormValue("state"))
			if err != nil {
				log.Printf("unable to decode state from callback: %v\n", err)
				http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
			} else {
				if len(s.QueryString) > 0 {
					http.Redirect(resp, req, "/?"+s.QueryString, http.StatusTemporaryRedirect)
				} else {
					http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
				}
			}

		}
	}
}

// getUserData exchanges an OAuth code for user data from Google
func getUserData(authCode string) ([]byte, error) {
	token, err := googleAuthConfig.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
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

// -- Status handler (provides auth status without taking any other action) --
// Use exclusively for frontend rendering decisions.

type StatusHandler struct{}

func (handler *StatusHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Printf("status: %v\n", len(userIDFromSession(req)) > 0)
	fmt.Fprintf(resp, "{\"authed\": \"%v\"}\n", len(userIDFromSession(req)) > 0)
}

type userInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func createSession(data []byte) string {
	var info userInfo
	err := json.Unmarshal(data, &info)
	if err != nil {
		log.Fatalf("Failed to unmarshal user data: %v\n", err)
	}
	err = insertUser(info)
	if err != nil {
		log.Fatalf("Unable to insert user info into database: %v\n", err)
	}

	sessionID := randomString()
	err = insertSession(sessionID, info.ID)
	if err != nil {
		log.Fatalf("Failed to create new session: %v\n", err)
	}

	return sessionID
}

type LogoutHandler struct{}

func (handler *LogoutHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// TODO: this is dumb, change addCookie to using an expiration time
	addCookie(resp, "session", "", -24*time.Hour)
	http.Redirect(resp, req, "/", http.StatusTemporaryRedirect)
}

// randomString returns a b64 string with 32 bytes of randomosity
// TODO: consider removing error return and crashing on failure (shouldn't happen)
func randomString() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("failed to generate random string: %v\n", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// TODO: this seems really overcomplicated
type state struct {
	QueryString string `json:"query"`
	RandomState string `json:"state"`
}

func encodeState(s state) (string, error) {
	stateString, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString([]byte(stateString)), nil
}

func decodeState(stateB64 string) (state, error) {
	queryState := state{}
	stateBytes, err := base64.URLEncoding.DecodeString(stateB64)
	if err != nil {
		return queryState, err
	}
	json.Unmarshal(stateBytes, &queryState)
	return queryState, nil
}
