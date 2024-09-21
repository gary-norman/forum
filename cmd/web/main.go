package main

import (
	"database/sql"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "./forum_database.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)

	srv := http.Server{
		Addr:    ":8989",
		Handler: routes(),
	}

	log.Printf("Listening on %v", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		return
	}
}
