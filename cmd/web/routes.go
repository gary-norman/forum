package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./assets/css"))
	mux.Handle("/assets/css/", http.StripPrefix("/assets/css", fileServer))
	mux.HandleFunc("/", app.getHome)
	// Use a single route for /posts/create and distinguish based on HTTP method
	mux.HandleFunc("/posts/create", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.createPost(w, r) // Handle GET requests
		case http.MethodPost:
			app.storePost(w, r) // Handle POST requests
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}
