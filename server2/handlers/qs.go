package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../which"
)

type Qs struct {
	ListHandler        *List
	NewQuestionHandler *NewQuestion
	//QuestionHandler    *QuestionHandler
	//NewVoteHandler     *NewVoteHandler
}

func NewQs(sessionStore which.SessionStore, questionStore which.QuestionStore) *Qs {
	qs := &Qs{}

	qs.ListHandler = &List{
		sessionStore:  sessionStore,
		questionStore: questionStore,
	}

	qs.NewQuestionHandler = &NewQuestion{
		sessionStore:  sessionStore,
		questionStore: questionStore,
	}

	return qs
}

func (handler *Qs) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	head, tail := shiftPath(req.URL.Path)
	if head == "list" {
		handler.ListHandler.ServeHTTP(resp, req)
	} else if head == "new" {
		handler.NewQuestionHandler.ServeHTTP(resp, req)
	} else if head == "vote" {
		//handler.NewVoteHandler.ServeHTTP(resp, req)
	} else if head == "q" {
		req.URL.Path = tail
		//handler.QuestionHandler.ServeHTTP(resp, req)
	} else {
		http.Error(resp, "qs endpoint does not exist", http.StatusNotFound)
	}
}

type List struct {
	sessionStore  which.SessionStore
	questionStore which.QuestionStore
}

func (handler *List) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session, err := sessionFromRequest(req, handler.sessionStore)
	if err != nil {
		http.Error(resp, "not authorized to access questions list", http.StatusUnauthorized)
		return
	}
	userID := session.UserID
	questions, err := handler.questionStore.FetchAuthoredBy(userID)
	if err != nil {
		log.Printf("failed to serve questions list: %v\n", err)
		http.Error(resp, "failed to fetch questions", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(resp).Encode(questions); err != nil {
		log.Printf("Unable to encode questions as JSON: %v\n", err)
		// TODO: this shouldn't fail, but if it does, what does the client see?
		return
	}
}

type NewQuestion struct {
	sessionStore  which.SessionStore
	questionStore which.QuestionStore
}

func (handler *NewQuestion) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session, err := sessionFromRequest(req, handler.sessionStore)
	if err != nil {
		http.Error(resp, "not authorized to submit questions", http.StatusUnauthorized)
		return
	}
	q, err := createQuestionFromRequest(req)
	if err != nil {
		log.Printf("failed to create new question: %v\n", err)
		http.Error(resp, "failed to create new question", http.StatusInternalServerError)
		return
	}
	q.UserID = session.UserID
	q.ID, err = handler.questionStore.Insert(q)
	if err != nil {
		log.Printf("failed to insert new question: %v\n", err)
		http.Error(resp, "failed to insert new question", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(resp).Encode(q); err != nil {
		log.Printf("failed to encode new question JSON response: %v\n", err)
	}
}

func createQuestionFromRequest(req *http.Request) (which.Question, error) {
	var q which.Question
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return q, fmt.Errorf("failed to read request body: %v", err)
	}
	err = json.Unmarshal(data, &q)
	if err != nil {
		return q, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return q, nil
}
