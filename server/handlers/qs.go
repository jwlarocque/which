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
		// TODO: FIXME: runoff ballots must have a vote for every option
	case which.QTypePlurality:
		any := false
		for _, vote := range ballot.Votes {
			if vote != nil && vote.State > 0 {
				if any {
					return false
				}
				any = true
			}
		}
	default:
		return false
	}
	return true
}

func (handler *NewVote) countUpdateVotes(ballot which.Ballot, oldBallot which.Ballot, question which.Question) error {
	if question.Type == which.QTypeApproval || question.Type == which.QTypePlurality {
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

func (handler *NewVote) recomputeVotes(question which.Question) error {
	var results []*which.Result
	ballots, err := handler.ballotStore.FetchAll(question.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve all ballots: %v", err)
	}
	if question.Type == which.QTypeApproval || question.Type == which.QTypePlurality {
		results = handler.recomputeSingleRoundVotes(question, ballots)
	} else if question.Type == which.QTypeRunoff {
		results = handler.recomputeRunoffVotes(question, ballots)
	} else {
		// TODO: implement recompute for plurality question type
		return fmt.Errorf("question type %v not yet implemented", question.Type)
	}

	err = handler.resultStore.RemoveAll(question.ID)
	if err != nil {
		return fmt.Errorf("failed to remove previous results: %v", err)
	}
	for _, result := range results {
		err = handler.resultStore.Update(*result)
		if err != nil {
			return fmt.Errorf("failed to store computed results: %v", err)
		}
	}
	return nil
}

func (handler *NewVote) recomputeSingleRoundVotes(question which.Question, ballots []*which.Ballot) []*which.Result {
	var results []*which.Result
	votes := make(map[int]int)

	for _, opt := range question.Options {
		votes[opt.ID] = 0
	}

	for _, ballot := range ballots {
		for _, vote := range ballot.Votes {
			if vote.State > 0 {
				votes[vote.OptionID]++
			}
		}
	}

	for optID, voteCount := range votes {
		results = append(results, &which.Result{
			QuestionID: question.ID,
			RoundNum:   0,
			OptionID:   optID,
			NumVotes:   voteCount,
		})
	}

	return results
}

// TODO: this function is enormous, break it down
// TODO: this function is quite expensive, and runs every time someone votes (the preceding DB
//       reads are expensive too).  Can something be done about this?
func (handler *NewVote) recomputeRunoffVotes(question which.Question, ballots []*which.Ballot) []*which.Result {
	var results []*which.Result

	curVotes := make(map[int]int)
	eliminated := make(map[int]bool)

	var bestCount, worstCount int
	var worstOptID int

	round := 0
	// eliminate options until one has a majority
	for bestCount < len(ballots)/2+1 {
		if round > 10 {
			log.Fatalln("uh oh, infinite loop")
		}

		curVotes = make(map[int]int)
		for _, opt := range question.Options {
			if !eliminated[opt.ID] {
				curVotes[opt.ID] = 0
			}
		}

		// count the top ranked non-eliminated vote on each ballot
		// into curVotes
		for _, ballot := range ballots {
			topOptID := 0
			highest := -1
			for _, vote := range ballot.Votes {
				if !eliminated[vote.OptionID] && vote.State > highest {
					highest = vote.State
					topOptID = vote.OptionID
				}
			}
			curVotes[topOptID]++
		}

		bestCount = -1
		worstCount = -1
		worstOptID = -1
		// find the options with the most and least votes in this round
		for optID, votes := range curVotes {
			// fill best and worst counts with the first votes
			if bestCount == -1 {
				bestCount = votes
				worstCount = votes
				worstOptID = optID
			} else if votes > bestCount {
				bestCount = votes
			} else if votes < worstCount {
				worstCount = votes
				worstOptID = optID
			}

			if votes > 0 {
				// append the option's votes to results for this round
				results = append(results, &which.Result{
					QuestionID: question.ID,
					RoundNum:   round,
					OptionID:   optID,
					NumVotes:   votes,
				})
			}
		}

		if worstOptID < 0 {
			log.Fatalln("uh oh, eliminated nothing and didn't reach a majority") // TODO: actually handle/consider case
		} else {
			eliminated[worstOptID] = true
		}

		round++
	}
	return results
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
