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

// QType is a question type/voting method "enum"
type QType int

const (
	// QTypeApproval indicates an approval question
	QTypeApproval = iota
	// QTypeRunoff indicates a ranked choice/instant runoff question
	QTypeRunoff
	// QTypePlurality indicates a plurality voting/single selection question
	QTypePlurality
)

func (t QType) String() string {
	return [...]string{"Approval", "Runoff", "Plurality"}[t]
}

// == server-DB structs, mostly for auth ================================

// User has information about a user
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// UserStore is implemented by structs which provide a way to insert Users into a database
type UserStore interface {
	Insert(User) error
}

// Session authenticates a user UserID and has creation and expiration times
type Session struct {
	ID      string    `db:"session_id"`
	UserID  string    `db:"user_id"`
	Created time.Time `db:"created"`
	Expires time.Time `db:"expires"`
}

// SessionStore is implemented by structs which provide a way to retrieve and update Sessions in a database
type SessionStore interface {
	Insert(Session) error
	Fetch(sessionID string) (Session, error)
}

// == structs which are also used for communication with the client ================================

// Question has an author User UserID, question text Name, a voting method Type, a list of Option, and a list of result Round
type Question struct {
	ID      string    `json:"question_id" db:"question_id"`
	UserID  string    `json:"-" db:"user_id"`
	Name    string    `json:"name" db:"name"`
	Type    QType     `json:"type" db:"type"`
	Options []*Option `json:"options" db:"-"`
	Results []*Result `json:"results" db:"-"`
}

// QuestionStore is implemented by structs which provide a way to retrieve and update Questions in a database
type QuestionStore interface {
	Insert(Question) (questionID string, err error)
	Fetch(questionID string) (Question, error)
	FetchAuthoredBy(userID string) ([]*Question, error)
	Remove(questionID string) error
}

// Option has the text of an option for a Question QuestionID
type Option struct {
	ID         int    `json:"option_id" db:"option_id"` // TODO: rename to option_id in json
	Text       string `json:"text" db:"text"`
	QuestionID string `json:"-" db:"question_id"`
}

// OptionStore is implemented by structs which provide a way to retrieve and update Options in a database
type OptionStore interface {
	Insert(Option) error
	FetchAll(questionID string) ([]*Option, error)
	RemoveAll(questionID string) error
}

// Ballot has a list of Votes submitted by a user UserID on a question QuestionID
type Ballot struct {
	ID         int     `json:"ballot_id" db:"ballot_id"`
	QuestionID string  `json:"question_id" db:"question_id"`
	UserID     string  `json:"-" db:"user_id"`
	Votes      []*Vote `json:"votes"`
}

// BallotStore is implemented by structs which provide a way to retrieve and update Ballots in a database
type BallotStore interface {
	Update(Ballot) (int, error)
	Fetch(questionID string, userID string) (Ballot, error)
	FetchAll(questionID string) ([]*Ballot, error)
	RemoveAll(questionID string) error
}

// Vote has the ranking or selectedness given to Option OptionID on the ballot BallotID
type Vote struct {
	BallotID int `json:"-" db:"ballot_id"`
	OptionID int `json:"option_id" db:"option_id"`
	State    int `json:"state" db:"state"`
}

// VoteStore is implemented by structs which provide a way to retrieve and update Votes in a database
type VoteStore interface {
	Update(Vote) error
	FetchAll(ballotID int) ([]*Vote, error)
}

// Result has the number of votes an option OptionID has during round RoundNum
type Result struct {
	QuestionID string `json:"-" db:"question_id"`
	RoundNum   int    `json:"round_num" db:"round_num"`
	// TODO: can I make this a foreign key even though options has primary key (option_id, question_id)?
	OptionID int `json:"option_id" db:"option_id"`
	NumVotes int `json:"num_votes" db:"num_votes"`
}

// ResultStore is implemented by structs which provide a way to retrieve and update Results in a database
type ResultStore interface {
	Update(Result) error
	FetchAll(questionID string) ([]*Result, error)
	RemoveAll(questionID string) error
}
