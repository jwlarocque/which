package main

import (
	"log"
)

// functions in this file handle interactions with the database

type Question struct {
	ID   int    `db:"question_id"`
	Name string `db:"name"`
}

func fetchQuestions() ([]*Question, error) {
	questions := []*Question{}
	err := db.Select(&questions, "select * from questions")
	if err != nil {
		log.Println("Unable to fetch questions")
		return nil, err
	}

	return questions, nil
}
