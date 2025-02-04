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
	userDataServer := http.FileServer(http.Dir("./db/userdata"))
	mux.Handle("/db/userdata/", http.StripPrefix("/db/userdata", userDataServer))
	mux.HandleFunc("/", app.getHome)
	mux.HandleFunc("GET /posts/create", app.createPost)
	mux.HandleFunc("POST /posts/create", app.storePost)
	mux.HandleFunc("POST /channels/create", app.storeChannel)
	mux.HandleFunc("/posts/create", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("POST /register", app.register)
	mux.HandleFunc("POST /login", app.login)
	mux.HandleFunc("POST /logout", app.logout)
	mux.HandleFunc("POST /protected", app.protected)
	mux.HandleFunc("POST /store-reaction", app.storeReaction)
	mux.HandleFunc("POST /edituser", app.editUserDetails)
	mux.HandleFunc("POST /channels/join", app.storeMembership)
	mux.HandleFunc("POST /channels/add-rules/{channelId}", app.CreateAndInsertRule)

	return mux
}
