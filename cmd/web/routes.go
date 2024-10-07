package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()
	cssServer := http.FileServer(http.Dir("./assets/css"))
	mux.Handle("/assets/css/", http.StripPrefix("/assets/css", cssServer))
	jsServer := http.FileServer(http.Dir("./assets/js"))
	mux.Handle("/assets/js/", http.StripPrefix("/assets/js", jsServer))
	iconServer := http.FileServer(http.Dir("./assets/icons"))
	mux.Handle("/assets/icons/", http.StripPrefix("/assets/icons", iconServer))
	imageServer := http.FileServer(http.Dir("./assets/images"))
	mux.Handle("/assets/images/", http.StripPrefix("/assets/images", imageServer))
	fontServer := http.FileServer(http.Dir("./assets/fonts"))
	mux.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts", fontServer))
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
