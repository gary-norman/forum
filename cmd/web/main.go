package main

import (
	"database/sql"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/sqlite"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type app struct {
	posts *sqlite.PostModel
}

func main() {
	ErrorMsgs := models.CreateErrorMessages()

	db, err := sql.Open("sqlite3", "./forum_database.db")
	if err != nil {
		fmt.Printf(ErrorMsgs.Open, "./forum_database.db", "sql.Open")
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

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Printf(ErrorMsgs.ConnInit, srv.Addr, "srv.ListenAndServe")
		return
	} else {
		log.Printf("Listening on %v", srv.Addr)
	}
}
