package main

import (
	"log"
	"time"
)

// functions in this file directly handle interactions with the database

type Question struct {
	ID   int    `db:"question_id"`
	Name string `db:"name"`
}

// assumes questions table exists
func fetchQuestions(user_id string) ([]*Question, error) {
	questions := []*Question{}
	err := db.Select(&questions, "SELECT question_id, name FROM questions WHERE user_id=$1", user_id)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func insertUser(info userInfo) error {
	_, err := db.Exec("INSERT INTO users VALUES ($1, $2) ON CONFLICT (user_id) DO NOTHING", info.ID, info.Email)
	return err
}

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
