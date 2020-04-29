package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"../which"
	"golang.org/x/oauth2"
)

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
