package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../which"
)

// == Qs root handler ================================

type Qs struct {
	ListHandler        *List
	NewQuestionHandler *NewQuestion
	GetQuestionHandler *GetQuestion
	NewVoteHandler     *NewVote
}

func NewQs(sessionStore which.SessionStore, questionStore which.QuestionStore, votesStore which.VotesStore) *Qs {
	qs := &Qs{}

	qs.ListHandler = &List{
		sessionStore:  sessionStore,
		questionStore: questionStore,
	}

	qs.NewQuestionHandler = &NewQuestion{
		sessionStore:  sessionStore,
		questionStore: questionStore,
	}

	qs.NewVoteHandler = &NewVote{
		sessionStore: sessionStore,
		votesStore:   votesStore,
	}

	qs.GetQuestionHandler = &GetQuestion{
		questionStore: questionStore,
	}

	return qs
}

func (handler *Qs) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// TODO: consider handling auth here, because all subhandlers require auth
	//       (would reduce repeated auth code, but some handlers need userID)
	head, tail := shiftPath(req.URL.Path)
	if head == "list" {
		handler.ListHandler.ServeHTTP(resp, req)
	} else if head == "new" {
		handler.NewQuestionHandler.ServeHTTP(resp, req)
	} else if head == "vote" {
		handler.NewVoteHandler.ServeHTTP(resp, req)
	} else if head == "q" {
		req.URL.Path = tail
		handler.GetQuestionHandler.ServeHTTP(resp, req)
	} else {
		http.Error(resp, "qs endpoint does not exist", http.StatusNotFound)
	}
}

// == List handler ================================

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

// == NewQuestion handler ================================

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

// == GetQuestion handler ================================

type GetQuestion struct {
	questionStore which.QuestionStore
}

// note: authorization is not needed to view a question
func (handler *GetQuestion) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	questionID := req.URL.Path[1:]
	q, err := handler.questionStore.Fetch(questionID)
	if err != nil {
		log.Printf("failed to fetch question: %v\n", err)
		http.Error(resp, "failed to fetch question", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(resp).Encode(q); err != nil {
		log.Printf("failed to encode question as JSON: %v\n", err)
		http.Error(resp, "failed to encode question as JSON", http.StatusInternalServerError)
	}
}

// == NewVote handler ================================

type NewVote struct {
	sessionStore which.SessionStore
	votesStore   which.VotesStore
}

func (handler *NewVote) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session, err := sessionFromRequest(req, handler.sessionStore)
	if err != nil {
		http.Error(resp, "not authorized to vote (oof)", http.StatusUnauthorized)
		return
	}
	vs, err := createVotesFromRequest(req)
	if err != nil {
		log.Printf("failed to create new votes: %v\n", err)
		http.Error(resp, "failed to create new votes", http.StatusInternalServerError)
		return
	}
	for _, vote := range vs.Votes {
		vote.UserID = session.UserID
		vote.QuestionID = vs.QuestionID
	}
	err = handler.votesStore.Update(vs)
	if err != nil {
		log.Printf("failed to insert/update new votes: %v\n", err)
		http.Error(resp, "failed to insert/update new votes", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(resp, "{\"ok\": \"true\"}\n")
}

// TODO: code is identical to createQuestionFromRequest, reduce repetition
func createVotesFromRequest(req *http.Request) (which.Votes, error) {
	var vs which.Votes
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return vs, fmt.Errorf("failed to read request body: %v", err)
	}
	err = json.Unmarshal(data, &vs)
	if err != nil {
		return vs, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return vs, nil
}
