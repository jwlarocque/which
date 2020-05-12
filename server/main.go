package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"

	// TODO: import from github
	"./handlers"
	"./psql"
)

var (
	db               *sqlx.DB
	googleAuthConfig *oauth2.Config
)

func main() {
	// connect to database
	db, err := sqlx.Connect("pgx", os.Getenv("WHICH_DB_STRING"))
	if err != nil {
		log.Fatalf("Unable to establish connection to the database: %v\n", err)
	}

	userStore := &psql.UserStore{DB: db}
	sessionStore := &psql.SessionStore{DB: db}
	optionStore := &psql.OptionStore{DB: db}
	voteStore := &psql.VoteStore{DB: db}
	ballotStore := &psql.BallotStore{DB: db, VoteStore: voteStore}
	// TODO: IMPORTANT: reconsider this store nesting?
	//       - it's not obvious from the interface definition that QuestionStore can insert options
	questionStore := &psql.QuestionStore{DB: db, OptionStore: optionStore, BallotStore: ballotStore}
	resultStore := &psql.ResultStore{DB: db}

	rh := handlers.NewRoot(userStore, sessionStore, questionStore, ballotStore, resultStore)

	log.Println("serving!")
	err = handlers.ListenAndServe(":80", rh)
	if err != nil {
		log.Fatalf("Error returned by ListenAndServe: %v\n", err)
	}
}
