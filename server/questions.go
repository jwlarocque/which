package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// QuestionsHandler delegates requests dealing with lists of and new questions
// TODO: does this path make sense? (consider "/qs/list" => "/qs" and "/qs/new" => "/new" or similar)
type QuestionsHandler struct {
	ListHandler        *ListHandler
	NewQuestionHandler *NewQuestionHandler
	QuestionHandler    *QuestionHandler
	NewVoteHandler     *NewVoteHandler
}

// TODO: consider moving user auth here to reduce repeated code (problem: passing around userID)
func (handler *QuestionsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	head, tail := shiftPath(req.URL.Path)
	if head == "list" { // TODO: not sure root qs/ should also be the path for the user's list of questions
		handler.ListHandler.ServeHTTP(resp, req)
	} else if head == "new" {
		handler.NewQuestionHandler.ServeHTTP(resp, req)
	} else if head == "vote" {
		handler.NewVoteHandler.ServeHTTP(resp, req)
	} else if head == "q" {
		req.URL.Path = tail
		handler.QuestionHandler.ServeHTTP(resp, req)
	} else {
		http.Error(resp, "Questions endpoint does not exist.", http.StatusNotFound)
	}
}

type QuestionHandler struct{}

func (handler *QuestionHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	questionID := req.URL.Path[1:]
	q, err := fetchQuestionAndOpts(questionID)
	if err != nil {
		log.Fatalf("Unable to fetch question from db: %v\n", err)
		http.Error(resp, "questions endpoint does not exist", http.StatusNotFound)
		return
	}
	if err = json.NewEncoder(resp).Encode(q); err != nil {
		log.Fatalf("Unable to encode question as JSON: %v\n", err)
	}
}

// ListHandler calls fetchQuestions in db.go and writes the response as JSON
type ListHandler struct{}

func (handler *ListHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	userID := userIDFromSession(req)
	if len(userID) > 0 {
		// TODO: fetch questions authored by this user
		questions, err := fetchQuestions(userID)
		if err != nil {
			log.Fatalf("Unable to fetch questions from db: %v\n", err)
			http.Error(resp, "unable to fetch questions from database", http.StatusInternalServerError) // TODO: more precise error
			return
		}
		if err = json.NewEncoder(resp).Encode(questions); err != nil {
			log.Fatalf("Unable to encode questions as JSON: %v\n", err)
			// TODO: this shouldn't fail, but if it does, what does the client see?
			return
		}
	} else {
		http.Error(resp, "You are not authorized to access questions.", http.StatusUnauthorized)
	}
}

type NewQuestionHandler struct{}

// TODO: break some of this out into helper functions
// TODO: streamline error handling
func (handler *NewQuestionHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Verify that user is logged in
	userID := userIDFromSession(req)
	if len(userID) > 0 {
		var q question
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(resp, "failed to read request body", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(data, &q)
		if err != nil {
			log.Println(err.Error())
			http.Error(resp, "failed to unmarshal json from request", http.StatusInternalServerError)
			return
		}

		q.User_ID = userID
		q.ID, err = insertQuestion(q)
		if err != nil {
			log.Println(err.Error())
			http.Error(resp, "{\"message\": \"failed to insert new question into database\", \"ok\": \"false\"}", http.StatusInternalServerError)
			return
		}
		// TODO: IN_PROGRESS: respond with new question data (then have frontend add it to the list)
		// assume the db query worked because we didn't see an error, send q back to client
		if err = json.NewEncoder(resp).Encode(q); err != nil {
			log.Println(err.Error())
		}
	} else {
		http.Error(resp, "You are not authorized to submit new questions.", http.StatusUnauthorized)
		return
	}
}

type NewVoteHandler struct{}

func (handler *NewVoteHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	user_ID := userIDFromSession(req)
	if len(user_ID) > 0 {
		var votes approvalVotes
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(resp, "failed to read request body", http.StatusInternalServerError)
			return
		}
		log.Printf("%v\n", string(data))
		err = json.Unmarshal(data, &votes)
		if err != nil {
			log.Println(err.Error())
			http.Error(resp, "failed to unmarshal json from request", http.StatusInternalServerError)
			return
		}

		// TODO: stutter (votes.Votes)
		for _, vote := range votes.Votes {
			vote.User_ID = user_ID
			vote.Question_ID = votes.Question_ID
		}
		err = updateVotes(votes)
		if err != nil {
			log.Println(err.Error())
			http.Error(resp, "{\"message\": \"failed to insert new question into database\", \"ok\": \"false\"}", http.StatusInternalServerError)
			return
		}
		// TODO: respond
		fmt.Fprintf(resp, "{\"ok\": \"true\"}\n")
	} else {
		http.Error(resp, "You are not authorized to vote.", http.StatusUnauthorized)
		return
	}
}

type approvalVote struct {
	Option_ID   int    `json:"id" db:"option_id"`
	State       bool   `json:"state" db:"state"`
	Question_ID string `json:"-" db:"question_id"`
	User_ID     string `json:"-" db:"user_id"`
}

type approvalVotes struct {
	Votes       []*approvalVote `json:"votes"`
	Question_ID string          `json:"question_id"`
}

type option struct {
	ID          int    `json:"id" db:"option_id"`
	Text        string `json:"text" db:"text"`
	Question_ID string `json:"-" db:"question_id"` // matches parent Question ID
}

type question struct {
	Name    string   `json:"name" db:"name"`
	Type    string   `json:"type" db:"type"`
	Options []option `json:"options" db:"-"`               // not stored in DB (obviously)
	ID      string   `json:"question_id" db:"question_id"` // generated by database
	User_ID string   `json:"-" db:"user_id"`               // retrieved from session cookie
}
