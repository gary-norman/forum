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
	users     *sqlite.UserModel
	posts     *sqlite.PostModel
	reactions *sqlite.ReactionModel
	reactionStatus *sqlite.ReactionModel
	saved     *sqlite.SavedModel
	mods      *sqlite.ModModel
	comments  *sqlite.CommentModel
	images    *sqlite.ImageModel
	channels  *sqlite.ChannelModel
	flags     *sqlite.FlagModel
	loyalty   *sqlite.LoyaltyModel
	members   *sqlite.MembershipModel
	muted     *sqlite.MutedChannelModel
	cookies   *CookieModel
}

func ErrorMsgs() *models.Errors {
	return models.CreateErrorMessages()
}

func main() {
	db, err := sql.Open("sqlite3", "./forum_database.db")
	if err != nil {
		log.Fatal(ErrorMsgs().Open, "./forum_database.db", "sql.Open", err)
	}

	app := app{
		posts: &sqlite.PostModel{
			DB: db,
		},
		users: &sqlite.UserModel{
			DB: db,
		},
		cookies: &CookieModel{
			DB: db,
		},
		reactions: &sqlite.ReactionModel{
			DB: db,
		},
		reactionStatus: &sqlite.ReactionModel{
			DB: db,
		},
	}
	// Initialise templates if (app *app) is a receiver for
	//the init() function that sets up custom go template functions
	app.init()

	srv := http.Server{
		Addr:    ":8989",
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Printf(ErrorMsgs().ConnInit, srv.Addr, "srv.ListenAndServe")
		return
	}
	log.Printf("Listening on %v", srv.Addr)
}
