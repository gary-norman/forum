package main

import (
	"database/sql"
	"github.com/gary-norman/forum/internal/sqlite"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type app struct {
	posts *sqlite.PostModel
}

func main() {
	db, err := sql.Open("sqlite3", "./forum_database.db")
	if err != nil {
		log.Fatal(err)
	}

	app := app{
		posts: &sqlite.PostModel{
			DB: db,
		},
	}

	srv := http.Server{
		Addr:    ":8989",
		Handler: app.routes(),
	}

	log.Printf("Listening on %v", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		return
	}
}
