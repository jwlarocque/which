package main

// authentication stuff
// TODO: some kind of cron job to remove expired sessions from the database
// TODO: reorganize this file (maybe it can be a package)

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

type AuthHandler struct {
	LoginHandler    *LoginHandler
	CallbackHandler *CallbackHandler
}

func (handler *AuthHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	head, _ := shiftPath(req.URL.Path)
	if head == "login" {
		log.Printf("AuthHandler: %v\n", handler)
		handler.LoginHandler.ServeHTTP(resp, req) // TODO: req URL correct?
	} else if head == "callback" {
		handler.CallbackHandler.ServeHTTP(resp, req)
	} else {
		http.Error(resp, "auth endpoint does not exist", 404)
	}
}

type LoginHandler struct{}

func (handler *LoginHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if len(userIDFromSession(req)) > 0 {
		http.Redirect(resp, req, "/loggedin", http.StatusTemporaryRedirect)
	} else {
		stateString, err := randomString()
		if err != nil {
			log.Fatalf("Unable to generate random state string: %v\n", err)
		}
		addCookie(resp, "state", stateString, time.Minute) // TODO: abstract out this duration?
		url := googleAuthConfig.AuthCodeURL(stateString)
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
		log.Println("no session matching cookie")
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

type CallbackHandler struct{}

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
			fmt.Fprintf(resp, "UserInfo: %s\n", data)
		}
	}
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
	insertUser(info)

	sessionID, err := randomString()
	if err != nil {
		log.Fatalf("Unable to generate random session id: %v\n", err)
	}
	err = insertSession(sessionID, info.ID)
	if err != nil {
		log.Fatalf("Failed to create new session: %v\n", err)
	}

	return sessionID
}

func getUserData(authCode string) ([]byte, error) {
	// TODO: context.Background() ???
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

// randomString returns a b64 string with 32 bytes of randomosity
// TODO: consider removing error return and crashing on failure (shouldn't happen)
func randomString() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
