package psql

// This package handles all interactions with the PostgreSQL database

import (
	"errors"
	"fmt"
	"time"

	"../which"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	DB *sqlx.DB
}

func (store UserStore) Insert(user which.User) error {
	_, err := store.DB.Exec("INSERT INTO users VALUES ($1, $2) ON CONFLICT (user_id) DO NOTHING", user.ID, user.Email)
	return err
}

type SessionStore struct {
	DB *sqlx.DB
}

func (store SessionStore) Insert(session which.Session) error {
	_, err := store.DB.Exec(
		"INSERT INTO sessions VALUES ($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT session_pkey DO UPDATE SET session_id=EXCLUDED.session_id, created=EXCLUDED.created, expires=EXCLUDED.expires",
		session.ID,
		session.UserID,
		time.Now(),
		time.Now().Add(which.CookieDuration))
	return err
}

func (store SessionStore) Fetch(sessionID string) (which.Session, error) {
	var session which.Session
	err := store.DB.Get(&session, "SELECT session_id, user_id, created, expires FROM sessions WHERE session_id=$1", sessionID)
	if err != nil {
		return session, fmt.Errorf("failed to fetch session ID: '%s', error: %v", sessionID, err)
	}
	return session, nil
}

type QuestionStore struct {
	DB          *sqlx.DB
	OptionStore which.OptionStore
}

func (store QuestionStore) Insert(q which.Question) (string, error) {
	questionID := ""
	rows, err := store.DB.NamedQuery("INSERT INTO questions (name, user_id, type) VALUES (:name, :user_id, :type) RETURNING question_id", q)
	if err != nil {
		return "", fmt.Errorf("failed to insert question: %v", err)
	}
	if rows.Next() {
		rows.Scan(&questionID)
		rows.Close()
	} else {
		return "", errors.New("no questionID received from insert query")
	}
	if len(questionID) > 0 {
		for _, option := range q.Options {
			option.QuestionID = questionID
			err = store.OptionStore.Insert(*option)
			if err != nil {
				return "", fmt.Errorf("incomplete options insert: %v", err)
			}
		}
		return questionID, nil
	} else {
		return "", errors.New("insert query returned empty questionID")
	}
}

func (store QuestionStore) Fetch(questionID string) (which.Question, error) {
	q := which.Question{ID: questionID}
	err := store.DB.Get(&q, "SELECT user_id, type, name FROM questions WHERE question_id=$1", questionID)
	if err != nil {
		return q, fmt.Errorf("failed to fetch question; ID: '%s', error: %v", questionID, err)
	}
	opts, err := store.OptionStore.FetchAll(questionID)
	if err != nil {
		return q, fmt.Errorf("failed to fetch question options: %v", err)
	}
	q.Options = opts
	return q, nil
}

func (store QuestionStore) FetchAuthoredBy(userID string) ([]*which.Question, error) {
	qs := []*which.Question{}
	err := store.DB.Select(&qs, "SELECT question_id, name FROM questions WHERE user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch questions authored by '%s', error: %v", userID, err)
	}
	return qs, nil
}

type OptionStore struct {
	DB *sqlx.DB
}

func (store OptionStore) Insert(o which.Option) error {
	if len(o.Text) > 0 { // exclude blank options TODO: more appropriate place to put this logic?
		_, err := store.DB.Exec("INSERT INTO options VALUES ($1, $2, $3)", o.ID, o.QuestionID, o.Text)
		return err
	}
	return nil
}

func (store OptionStore) FetchAll(questionID string) ([]*which.Option, error) {
	options := []*which.Option{}
	err := store.DB.Select(&options, "SELECT option_id, text FROM options WHERE question_id=$1", questionID)
	if err != nil {
		return options, fmt.Errorf("failed to fetch options for question ID: '%s', error: %v", questionID, err)
	}
	return options, nil
}

type BallotStore struct {
	DB *sqlx.DB
	VoteStore which.VoteStore
}

func (store BallotStore) Update(ballot which.Ballot) (int, error) {
	ballotID := 0
	// yikes!
	rows, err := store.DB.NamedQuery(
		`WITH ins AS (
			INSERT INTO ballots (ballot_id, question_id, user_id)
				VALUES(:ballot_id, :question_id, :user_id) 
			ON CONFLICT ON CONSTRAINT ballots_pkey 
				DO NOTHING 
			RETURNING ballot_id
		)
		SELECT ballot_id FROM ins
		UNION ALL
		SELECT ballot_id FROM ballots
			WHERE question_id=:question_id
			AND user_id=:user_id
		LIMIT 1`, 
		ballot)
	if err != nil {
		return -1, fmt.Errorf("failed to insert ballot: %v", err)
	}
	if rows.Next() {
		rows.Scan(&ballotID)
		rows.Close()
	} else {
		return -1, errors.New("no ballotID received from insert query")
	}
	if ballotID > 0 {
		for _, vote := range ballot.Votes {
			vote.BallotID = ballotID
			err = store.VoteStore.Update(*vote)
			if err != nil {
				return -1, fmt.Errorf("incomplete votes insert: %v", err)
			}
		}
		return ballotID, nil
	} else {
		return -1, errors.New("insert query returned empty ballotID")
	}
}

func (store BallotStore) FetchAll(questionID string) ([]*which.Ballot, error) {
	ballots := []*which.Ballot{}
	err := store.DB.Select(&ballots, "SELECT ballot_id, question_id, user_id FROM ballots WHERE question_id=$1", questionID)
	if err != nil {
		return ballots, fmt.Errorf("failed to fetch votes for question ID: %s, error: %v", questionID, err)
	}
	return ballots, nil
}

type VoteStore struct {
	DB *sqlx.DB
}

// TODO: bulk insert
func (store VoteStore) Update(vote which.Vote) error {
	_, err := store.DB.NamedExec("INSERT INTO votes VALUES(:ballot_id, :option_id, :state) ON CONFLICT ON CONSTRAINT votes_pkey DO UPDATE SET state=EXCLUDED.state", vote)
	return err
}
