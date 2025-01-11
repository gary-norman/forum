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
	mux.HandleFunc("GET /posts/create", app.createPost)
	mux.HandleFunc("POST /posts/create", app.storePost)
	mux.HandleFunc("/posts.create", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/register", app.register)
	mux.HandleFunc("/login", app.login)
	mux.HandleFunc("/logout", app.logout)
	mux.HandleFunc("/protected", app.protected)

	return mux
}
