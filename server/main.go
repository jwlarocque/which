package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Question struct {
	ID   int    `db:"question_id"`
	Name string `db:"name"`
}

func main() {
	db, err := sqlx.Connect("pgx", "dbname=whichdb")
	if err != nil {
		log.Fatalf("Unable to establish connection to the database: %v\n", err)
	}
	questions, err := fetchQuestions(db)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	for _, question := range questions {
		log.Println(question.Name)
	}

	http.HandleFunc("/", Server)
	http.ListenAndServe(":8080", nil)
}

func Server(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "public/index.html")
	} else if r.URL.Path == "/rand" {
		log.Println("TODO: write the thing that goes here")
	} else {
		http.ServeFile(w, r, "public/"+r.URL.Path[1:])
	}
}
