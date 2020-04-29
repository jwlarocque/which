package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"../which"
	"golang.org/x/oauth2"
)

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

// returns a Session if the request has a valid session cookie, else error
// TODO: somewhat awkward - retrieves session from store using request data
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

func randomString() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("failed to generate random string: %v\n", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
