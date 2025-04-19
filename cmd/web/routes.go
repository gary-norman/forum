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
	dbImageServer := http.FileServer(http.Dir("./db/userdata/images"))
	mux.Handle("/db/userdata/images/", http.StripPrefix("/db/userdata/images", dbImageServer))
	fontServer := http.FileServer(http.Dir("./assets/fonts"))
	mux.Handle("/assets/fonts/", http.StripPrefix("/assets/fonts", fontServer))
	cursorServer := http.FileServer(http.Dir("./assets/cursors"))
	mux.Handle("/assets/cursors/", http.StripPrefix("/assets/cursors", cursorServer))
	userDataServer := http.FileServer(http.Dir("./db/userdata"))
	mux.Handle("/db/userdata/", http.StripPrefix("/db/userdata", userDataServer))
	mux.HandleFunc("POST /register", app.register)
	mux.HandleFunc("POST /login", app.login)
	mux.HandleFunc("POST /logout", app.logout)
	mux.HandleFunc("POST /protected", app.protected)
	mux.Handle("/", withUser(http.HandlerFunc(app.getHome), app))
	mux.Handle("/channels/{channelId}", withUser(http.HandlerFunc(app.getThisChannel), app))
	mux.HandleFunc("GET /posts/create", app.createPost)
	mux.HandleFunc("/search", app.search)
	mux.Handle("/posts/{postId}", withUser(http.HandlerFunc(app.getThisPost), app))
	mux.Handle("/users/{userId}", withUser(http.HandlerFunc(app.getThisUser), app))
	mux.Handle("POST /posts/create", withUser(http.HandlerFunc(app.storePost), app))
	mux.Handle("POST /channels/create", withUser(http.HandlerFunc(app.storeChannel), app))
	mux.HandleFunc("/posts/create", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.Handle("POST /store-reaction", withUser(http.HandlerFunc(app.storeReaction), app))
	mux.Handle("POST /edituser", withUser(http.HandlerFunc(app.editUserDetails), app))
	mux.Handle("POST /channels/join", withUser(http.HandlerFunc(app.storeMembership), app))
	mux.Handle("POST /channels/add-rules/{channelId}", withUser(http.HandlerFunc(app.CreateAndInsertRule), app))
	mux.Handle("POST /store-comment", withUser(http.HandlerFunc(app.storeComment), app))

	return mux
}
