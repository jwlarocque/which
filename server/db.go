package main

import (
	"fmt"
	"log"
	"time"
)

// functions in this file directly handle interactions with the database

// note: Question struct is in questions.go
// TODO: consider organizing code into files by data type rather than
//       function, to avoid this (i.e., eliminate this db.go file)
// TODO: errors in this function are really bad, because they probably mean
//       a partial db write.  So, I dunno, fix yo bugs.
func insertQuestion(q question) (string, error) {
	question_id := ""
	rows, err := db.NamedQuery("INSERT INTO questions (name, user_id, type) VALUES (:name, :user_id, :type) RETURNING question_id", q)
	if err != nil {
		return "", err
	}
	if rows.Next() {
		rows.Scan(&question_id)
	} else {
		// TODO: better error handling
		return "", fmt.Errorf("no question_id received from db query")
	}
	if len(question_id) > 0 {
		// TODO: check for duplicated option ids (might happen if the client is messing with us)
		for _, option := range q.Options {
			if len(option.Text) > 0 { // exclude blank options
				option.Question_ID = question_id // TODO: maybe just pass this to insertOption
				err = insertOption(option)
				if err != nil {
					return "", err
				}
			}
		}
		return question_id, nil
	} else {
		// TODO: return better error (no/bad question id returned from INSERT INTO questions)
		return "", fmt.Errorf("uh oh")
	}
}

func insertOption(o option) error {
	_, err := db.Exec("INSERT INTO options VALUES ($1, $2, $3)", o.ID, o.Question_ID, o.Text)
	return err
}

func fetchQuestionAndOpts(question_id string) (question, error) {
	var q question
	q.ID = question_id
	err := db.Get(&q, "SELECT user_id, type, name FROM questions WHERE question_id=$1", question_id)
	if err != nil {
		return q, err
	}
	err = db.Select(&(q.Options), "SELECT option_id, text FROM options WHERE question_id=$1", question_id)
	if err != nil {
		return q, err
	}
	return q, nil
}

func fetchQuestions(user_id string) ([]*question, error) {
	questions := []*question{}
	err := db.Select(&questions, "SELECT question_id, name FROM questions WHERE user_id=$1", user_id)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

// TODO: this delete then insert strategy for updates is not ideal.  switch to:
//       store all options for each user, with a boolean for whether they voted for it (then just UPDATE)
//       (yes, far more storage required, but I'm pretty sure it's worth it.)
func updateVotes(user_ID string, question_ID string, votes approvalVotes) error {
	_, err := db.Exec("DELETE FROM approval_votes WHERE question_id=$1 AND user_id=$2", question_ID, user_ID)
	if err != nil {
		return err
	}
	if len(votes.Votes) > 0 {
		_, err = db.NamedExec("INSERT INTO approval_votes VALUES(:option_id, :question_id, :user_id) ON CONFLICT DO NOTHING", votes.Votes)
	}
	return err
}

func insertUser(info userInfo) error {
	_, err := db.Exec("INSERT INTO users VALUES ($1, $2) ON CONFLICT (user_id) DO NOTHING", info.ID, info.Email)
	return err
}

// TODO: why did I add this stub?
func fetchUser() {

}

func insertSession(sessionID string, userID string) error {
	// clear existing sessions for this user
	// TODO: consider changin this behavior or moving the logic to auth.go
	_, err := db.Exec("DELETE FROM sessions WHERE user_id=$1", userID)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"INSERT INTO sessions VALUES ($1, $2, $3, $4)",
		sessionID,
		userID,
		time.Now(),
		time.Now().Add(24*time.Hour))
	return err
}

type Session struct {
	ID      string    `db:"session_id"`
	User_ID string    `db:"user_id"`
	Created time.Time `db:"created"`
	Expires time.Time `db:"expires"`
}

func fetchSession(sessionID string) (Session, error) {
	var session Session
	err := db.Get(&session, "SELECT session_id, user_id, created, expires FROM sessions WHERE session_id=$1", sessionID)
	if err != nil {
		log.Printf("error!: %v, maybe session: %v\n", err, session)
		return session, err
	}
	return session, nil
}
