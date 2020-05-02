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
	GetBallotsHandler *GetBallots
}

func NewQs(sessionStore which.SessionStore, questionStore which.QuestionStore, ballotStore which.BallotStore) *Qs {
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
		ballotStore:   ballotStore,
	}

	qs.GetQuestionHandler = &GetQuestion{
		questionStore: questionStore,
	}

	qs.GetBallotsHandler = &GetBallots{
		ballotStore: ballotStore,
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
	} else if head == "vs" {
		// handle send votes
		handler.GetBallotsHandler.ServeHTTP(resp, req)
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
	ballotStore  which.BallotStore
}

// TODO: alert user that vote failed
func (handler *NewVote) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session, err := sessionFromRequest(req, handler.sessionStore)
	if err != nil {
		http.Error(resp, "not authorized to vote (oof)", http.StatusUnauthorized)
		return
	}
	ballot, err := createBallotFromRequest(req)
	ballot.UserID = session.UserID
	if err != nil {
		log.Printf("failed to create new votes: %v\n", err)
		http.Error(resp, "failed to create new votes", http.StatusInternalServerError)
		return
	}
	_, err = handler.ballotStore.Update(ballot)
	if err != nil {
		log.Printf("failed to insert/update ballot: %v\n", err)
		http.Error(resp, "failed to insert/update ballot", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(resp, "{\"ok\": \"true\"}\n")
}

// TODO: code is identical to createQuestionFromRequest, reduce repetition
func createBallotFromRequest(req *http.Request) (which.Ballot, error) {
	var vs which.Ballot
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

// == GetVotes handler ================================

type GetBallots struct {
	ballotStore which.BallotStore
}

func (handler *GetBallots) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_, tail := shiftPath(req.URL.Path)
	questionID := tail[1:]
	ballots, err := handler.ballotStore.FetchAll(questionID)
	if err != nil {
		log.Printf("failed to fetch ballots: %v\n", err)
		http.Error(resp, "failed to fetch ballots", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(resp).Encode(ballots); err != nil {
		log.Printf("failed to encode ballots as JSON: %v\n", err)
		http.Error(resp, "failed to encode ballots as JSON", http.StatusInternalServerError)
	}
}
