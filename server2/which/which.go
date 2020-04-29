package which

// package which defines data types which are shared between several other
// packages (for example, psql provides a UserStore implementation which is
// passed to handlers in main.go and used by them to insert request data into
// the database)

import "time"

// TODO: read these from a config file
const (
	OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	CookieDuration    = 24 * time.Hour
)

// server-DB structs, mostly for auth

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

// structs which are also used for communication with the client

type Question struct {
	ID      string    `json:"question_id" db:"question_id"`
	UserID  string    `json:"-" db:"user_id"`
	Name    string    `json:"name" db:"name"`
	Type    string    `json:"type" db:"type"`
	Options []*Option `json:"options" db:"-"`
}

type QuestionStore interface {
	Insert(Question) (questionID string, err error)
	Fetch(questionID string) (Question, error)
	FetchAuthoredBy(userID string) ([]*Question, error)
}

type Option struct {
	ID         int    `json:"option_id" db:"option_id"`
	Text       string `json:"text" db:"text"`
	QuestionID string `json:"-" db:"question_id"`
}

type OptionStore interface {
	Insert(Option) error
	FetchAll(questionID string) ([]*Option, error)
}

type Votes struct {
	QuestionID string  `json:"question_id"`
	Votes      []*Vote `json:"votes"`
}

type VotesStore interface {
	Update(Votes) error
}

type Vote struct {
	OptionID   int    `json:"id" db:"option_id"`
	QuestionID string `json:"-" db:"question_id"`
	UserID     string `json:"-" db:"user_id"` // user who submitted this vote
	State      int    `json:"state" db:"state"`
}