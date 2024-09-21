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

	log.Printf("Listening on %v", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}
