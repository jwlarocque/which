package psql

import (
	"time"

	"../which"
	"github.com/jmoiron/sqlx"
)

// This package handles all interactions with the PostgreSQL database

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
		time.Now().Add(which.CookieDuration)) // TODO: read expiration from config
	return err
}

func (store SessionStore) Fetch(sessionID string) (which.Session, error) {
	var session which.Session
	err := store.DB.Get(&session, "SELECT session_id, user_id, created, expires FROM sessions WHERE session_id=$1", sessionID)
	if err != nil {
		return session, err
	}
	return session, nil
}
