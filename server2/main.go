package main

// TODO: cut down on all the log.Fatal

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"./handlers" // TODO: switch to import from github
)

// TODO: read these from a config file
const (
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	cookieDuration    = 24 * time.Hour
)

var (
	db               *sqlx.DB
	googleAuthConfig *oauth2.Config
)

func init() {

}

func main() {
	http.ListenAndServe(":8080", handlers.NewRoot())
}

type User struct {
	ID      string `json:"user_id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type UserStore interface {
	Insert(User) error
}

type Session struct {
	ID      string    `db:"session_id"`
	UserID  string    `db:"user_id"`
	Created time.Time `db:"created"`
	Expires time.Time `db:"expires"`
}

type SessionStore interface {
	Insert(Session) error
	Fetch(sessionID string) (Session, error)
}

type Question struct {
	ID      string   `json:"question_id" db:"question_id"`
	UserID  string   `json:"-" db:"user_id"`
	Name    string   `json:"name" db:"name"`
	Type    string   `json:"type" db:"type"`
	Options []Option `json:"options" db:"-"`
}

type QuestionStore interface {
	Insert(Question) (questionID string, err error)
	Fetch(questionID string) (Question, error)
	FetchAll(userID string) ([]*Question, error)
}

type Option struct {
	ID         int    `json:"option_id" db:"option_id"`
	Text       string `json:"text" db:"text"`
	QuestionID string `json:"-" db:"question_id"`
}

type OptionStore interface {
	Insert(Option) error
}
