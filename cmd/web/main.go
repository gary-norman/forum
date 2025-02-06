package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

type app struct {
	users          *sqlite.UserModel
	posts          *sqlite.PostModel
	reactions      *sqlite.ReactionModel
	reactionStatus *sqlite.ReactionModel
	saved          *sqlite.SavedModel
	mods           *sqlite.ModModel
	comments       *sqlite.CommentModel
	images         *sqlite.ImageModel
	channels       *sqlite.ChannelModel
	flags          *sqlite.FlagModel
	loyalty        *sqlite.LoyaltyModel
	memberships    *sqlite.MembershipModel
	muted          *sqlite.MutedChannelModel
	cookies        *CookieModel
}

func ErrorMsgs() *models.Errors {
	return models.CreateErrorMessages()
}

func main() {
	db, err := sql.Open("sqlite3", "./db/forum_database.db")
	if err != nil {
		log.Fatal(ErrorMsgs().Open, "./db/forum_database.db", "sql.Open", err)
	}

	app := app{
		posts: &sqlite.PostModel{
			DB: db,
		},
		users: &sqlite.UserModel{
			DB: db,
		},
		channels: &sqlite.ChannelModel{
			DB: db,
		},
		memberships: &sqlite.MembershipModel{
			DB: db,
		},
		cookies: &CookieModel{
			DB: db,
		},
		comments: &sqlite.CommentModel{
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
	// the init() function that sets up custom go template functions
	app.init()

	port := 8989
	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}

	go func() {
		// Log server listening messages
		log.Printf(ErrorMsgs().KeyValuePair, "Starting server on port", port)
		log.Printf(ErrorMsgs().ConnSuccess, "http://localhost"+addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf(ErrorMsgs().ConnInit, srv.Addr, "srv.ListenAndServe")
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Printf("Stopped serving new connections.")
	}()

	// set up a channel to listen for kill or interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// create cancellation signal and timeout
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	// shut down the server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf(ErrorMsgs().Shutdown, err)
	}
	log.Printf("Graceful shutdown complete.")
}
