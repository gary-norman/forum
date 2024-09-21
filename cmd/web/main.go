package main

import (
	"log"
	"net/http"
)

func main() {
	srv := http.Server{
		Addr:    ":8989",
		Handler: routes(),
	}

	log.Println("Listening on :8989")
	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}
