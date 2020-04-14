package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// functions in this file handle interactions with the database

func fetchQuestions(db *sqlx.DB) ([]*Question, error) {
	questions := []*Question{}
	err := db.Select(&questions, "select * from questions")
	if err != nil {
		log.Println("Unable to fetch questions")
		return nil, err
	}

	return questions, nil
}
