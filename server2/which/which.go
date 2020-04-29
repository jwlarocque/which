package which

import "time"

// TODO: read these from a config file
const (
	OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	CookieDuration    = 24 * time.Hour
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
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
