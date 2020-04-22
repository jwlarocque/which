package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// QuestionsHandler calls fetchQuestions in db.go and writes the response as JSON
type QuestionsHandler struct{}

func (handler *QuestionsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	user_id := userIDFromSession(req)
	if len(user_id) > 0 {
		// TODO: fetch questions authored by this user
		questions, err := fetchQuestions(user_id)
		if err != nil {
			log.Fatalf("Unable to fetch questions from db: %v\n", err)
			http.Error(resp, "unable to fetch questions from database", http.StatusInternalServerError) // TODO: more precise error
		}
		if err = json.NewEncoder(resp).Encode(questions); err != nil {
			log.Fatalf("Unable to encode questions as JSON: %v\n", err)
			// TODO: this shouldn't fail, but if it does, what does the client see?
		}
	} else {
		// TODO: not authorized
		http.Error(resp, "You are not authorized to access questions.", http.StatusUnauthorized)
	}
}
