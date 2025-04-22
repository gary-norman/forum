package handlers

import (
	"net/http"

	mw "github.com/gary-norman/forum/internal/http/middleware"
)

// type RouteHandler struct {
// 	App  *app.App
// 	Main *MainHandler
// }

func (r *RouteHandler) Routes() http.Handler {
	s := r
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
	mux.HandleFunc("POST /register", s.Auth.Register)
	mux.HandleFunc("POST /login", s.Auth.Login)
	mux.HandleFunc("POST /logout", s.Auth.Logout)
	mux.HandleFunc("POST /protected", s.Auth.Protected)
	mux.Handle("/", mw.WithUser(http.HandlerFunc(s.Home.GetHome), s.App))
	mux.Handle("/channels/{channelId}", mw.WithUser(http.HandlerFunc(s.Channel.GetThisChannel), s.App))
	mux.HandleFunc("GET /posts/create", s.Post.CreatePost)
	mux.HandleFunc("/search", s.Search.Search)
	mux.Handle("/posts/{postId}", mw.WithUser(http.HandlerFunc(s.Post.GetThisPost), s.App))
	mux.Handle("/users/{userId}", mw.WithUser(http.HandlerFunc(s.User.GetThisUser), s.App))
	mux.Handle("POST /posts/create", mw.WithUser(http.HandlerFunc(s.Post.StorePost), s.App))
	mux.Handle("POST /channels/create", mw.WithUser(http.HandlerFunc(s.Channel.StoreChannel), s.App))
	mux.HandleFunc("/posts/create", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.Handle("POST /store-reaction", mw.WithUser(http.HandlerFunc(s.Reaction.StoreReaction), s.App))
	mux.Handle("POST /edituser", mw.WithUser(http.HandlerFunc(s.User.EditUserDetails), s.App))
	mux.Handle("POST /channels/join", mw.WithUser(http.HandlerFunc(s.Channel.StoreMembership), s.App))
	mux.Handle("POST /channels/add-rules/{channelId}", mw.WithUser(http.HandlerFunc(s.Channel.CreateAndInsertRule), s.App))
	mux.Handle("POST /store-comment", mw.WithUser(http.HandlerFunc(s.Comment.StoreComment), s.App))

	return mux
}
