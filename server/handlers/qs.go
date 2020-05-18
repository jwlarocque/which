package handlers

// TODO: this file is getting out of hand, especially the Qs handler.  Do something about it.

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
	ListHandler           *List
	NewQuestionHandler    *NewQuestion
	GetQuestionHandler    *GetQuestion
	DeleteQuestionHandler *DeleteQuestion
	NewVoteHandler        *NewVote
	GetResultsHandler     *GetResults
	GetBallotHandler      *GetBallot
}

func NewQs(sessionStore which.SessionStore, questionStore which.QuestionStore, ballotStore which.BallotStore, resultStore which.ResultStore) *Qs {
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
		sessionStore:  sessionStore,
		questionStore: questionStore,
		ballotStore:   ballotStore,
		resultStore:   resultStore,
	}

	qs.GetQuestionHandler = &GetQuestion{
		questionStore: questionStore,
	}

	qs.DeleteQuestionHandler = &DeleteQuestion{
		sessionStore:  sessionStore,
		questionStore: questionStore,
	}

	qs.GetResultsHandler = &GetResults{
		resultStore: resultStore,
	}

	qs.GetBallotHandler = &GetBallot{
		sessionStore: sessionStore,
		ballotStore:  ballotStore,
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
	} else if head == "del" {
		handler.DeleteQuestionHandler.ServeHTTP(resp, req)
	} else if head == "rs" {
		// handle send votes
		handler.GetResultsHandler.ServeHTTP(resp, req)
	} else if head == "b" {
		handler.GetBallotHandler.ServeHTTP(resp, req)
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
	var q which.Question
	unmarshaled, err := unmarshalStructFromRequest(req, &q)
	q = *unmarshaled.(*which.Question)
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

// == DeleteQuestion handler ================================

type DeleteQuestion struct {
	sessionStore  which.SessionStore
	questionStore which.QuestionStore
}

func (handler *DeleteQuestion) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session, err := sessionFromRequest(req, handler.sessionStore)
	if err != nil {
		http.Error(resp, "not authorized to delete questions", http.StatusUnauthorized)
		return
	}
	questionID, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("failed to read question ID from request body: %v\n", err)
		http.Error(resp, "failed to read question ID from request body", http.StatusInternalServerError)
		return
	}
	q, err := handler.questionStore.Fetch(string(questionID))
	if err != nil {
		log.Printf("failed to verify deletion question ownership: %v\n", err)
		http.Error(resp, "failed to verify ownership", http.StatusInternalServerError)
		return
	}
	if session.UserID != q.UserID {
		log.Printf("question ownership mismatch; session UserID: '%s', question UserID: '%s'", session.UserID, q.UserID)
		http.Error(resp, "question ownership mismatch", http.StatusForbidden)
		return
	}
	err = handler.questionStore.Remove(q.ID)
	if err != nil {
		log.Printf("failed to remove question: %v", err)
		http.Error(resp, "failed to remove question", http.StatusInternalServerError)
		return
	}
	// TODO: what is the idiomatic way to send OK?  this seems weird.
	fmt.Fprintf(resp, "{\"ok\": \"true\"}\n")
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
	sessionStore  which.SessionStore
	questionStore which.QuestionStore
	ballotStore   which.BallotStore
	resultStore   which.ResultStore
}

// TODO: alert user that vote failed (this is a frontend issue...)
func (handler *NewVote) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session, err := sessionFromRequest(req, handler.sessionStore)
	if err != nil {
		http.Error(resp, "not authorized to vote (oof)", http.StatusUnauthorized)
		return
	}
	var ballot which.Ballot
	unmarshaled, err := unmarshalStructFromRequest(req, &ballot)
	ballot = *unmarshaled.(*which.Ballot)
	ballot.UserID = session.UserID
	if err != nil {
		log.Printf("failed to create new votes: %v\n", err)
		http.Error(resp, "failed to create new votes", http.StatusInternalServerError)
		return
	}
	q, err := handler.questionStore.Fetch(ballot.QuestionID)
	if err != nil {
		log.Printf("failed to get question for ballot: %v\n", err)
		http.Error(resp, "failed to get question for ballot", http.StatusInternalServerError)
	}
	if !ballotValid(ballot, q.Type) {
		http.Error(resp, "invalid ballot", http.StatusBadRequest)
		return
	}
	// fetch previous ballot, used to compute delta in countUpdateVotes
	oldBallot, err := handler.ballotStore.Fetch(q.ID, ballot.UserID)
	if err != nil {
		// TODO: better flow than using -1 as invalid ballot
		oldBallot = which.Ballot{ID: -1}
	}
	_, err = handler.ballotStore.Update(ballot)
	if err != nil {
		log.Printf("failed to insert/update ballot: %v\n", err)
		http.Error(resp, "failed to insert/update ballot", http.StatusInternalServerError)
		return
	}
	handler.countUpdateVotes(ballot, oldBallot, q)
	//fmt.Fprintf(resp, "{\"ok\": \"true\"}\n")
}

func ballotValid(ballot which.Ballot, qType which.QType) bool {
	switch qType {
	case which.QTypeApproval:
		for _, vote := range ballot.Votes {
			if vote != nil && vote.State > 1 {
				return false
			}
		}
	case which.QTypeRunoff:
		seen := make(map[int]bool)
		for _, vote := range ballot.Votes {
			if vote.State < 1 || vote.State > len(ballot.Votes) {
				return false
			}
			if seen[vote.State] {
				return false
			}
			seen[vote.State] = true
		}
	default:
		return false
	}
	return true
}

func (handler *NewVote) countUpdateVotes(ballot which.Ballot, oldBallot which.Ballot, question which.Question) error {
	if question.Type == which.QTypeApproval {
		if oldBallot.ID < 0 {
			// we failed to retrieve the old ballot earlier, assume one might
			// have existed so recompute total from scratch
			return handler.recomputeVotes(question)
		}
		// compute vote delta between oldBallot and ballot and apply it to
		// the Results in Question.Rounds[0]
		return handler.recomputeVotes(question) // TODO: implement delta vote counting
	}
	return handler.recomputeVotes(question)
}

// TODO: this function is quite complex, consider breaking it up
func (handler *NewVote) recomputeVotes(question which.Question) error {
	ballots, err := handler.ballotStore.FetchAll(question.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve all ballots: %v", err)
	}
	if question.Type == which.QTypeApproval {
		results := make(map[int]*which.Result)
		for _, ballot := range ballots {
			for _, vote := range ballot.Votes {
				_, exists := results[vote.OptionID]
				if !exists {
					results[vote.OptionID] = &which.Result{
						QuestionID: question.ID,
						RoundNum:   0,
						OptionID:   vote.OptionID,
						NumVotes:   0,
					}
				}
				if vote.State > 0 {
					results[vote.OptionID].NumVotes++
				}
			}
		}
		for _, result := range results {
			err = handler.resultStore.Update(*result)
			if err != nil {
				return fmt.Errorf("failed to store computed results: %v", err)
			}
		}
	} else {
		// TODO: implement recompute for other question types
		return fmt.Errorf("question type %v not yet implemented", question.Type)
	}
	return nil
}

// == GetResults handler ================================

type GetResults struct {
	resultStore which.ResultStore
}

func (handler *GetResults) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_, tail := shiftPath(req.URL.Path)
	questionID := tail[1:]
	results, err := handler.resultStore.FetchAll(questionID)
	if err != nil {
		log.Printf("failed to fetch results: %v\n", err)
		http.Error(resp, "failed to fetch results", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(resp).Encode(results); err != nil {
		log.Printf("failed to encode results as JSON: %v\n", err)
		http.Error(resp, "failed to encode results as JSON", http.StatusInternalServerError)
	}
}

// == GetBallot handler ================================

type GetBallot struct {
	sessionStore which.SessionStore
	ballotStore  which.BallotStore
}

func (handler *GetBallot) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_, tail := shiftPath(req.URL.Path)
	questionID := tail[1:]
	session, err := sessionFromRequest(req, handler.sessionStore)
	if err != nil {
		http.Error(resp, "not authorized to retrieve a ballot", http.StatusUnauthorized)
		return
	}
	ballot, err := handler.ballotStore.Fetch(questionID, session.UserID)
	if err != nil {
		//log.Printf("failed to fetch ballot: %v\n", err)
		//http.Error(resp, "failed to fetch ballot", http.StatusNotFound)
		return
	}
	if err = json.NewEncoder(resp).Encode(ballot); err != nil {
		log.Printf("failed to encode ballot as JSON: %v\n", err)
		http.Error(resp, "failed to encode ballot as JSON", http.StatusInternalServerError)
	}
}
