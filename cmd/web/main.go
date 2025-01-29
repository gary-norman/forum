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
	db             *sql.DB // Store DB reference for cleanup
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

//// Global template variable
//var tpl *template.Template
//
//func loadTemplates() error {
//	var err error
//	tpl, err = template.ParseFiles("assets/templates/index.html")
//	return err
//}

func initializeApp() (*app, func(), error) {
	// Open database connection
	db, err := sql.Open("sqlite3", "db/forum_database.db") // Adjust DB path
	if err != nil {
		return nil, nil, err
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		db.Close() // Close DB if ping fails
		return nil, nil, err
	}

	// App instance with DB reference
	appInstance := &app{
		db: db, // Store reference for cleanup
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
		reactions: &sqlite.ReactionModel{
			DB: db,
		},
		reactionStatus: &sqlite.ReactionModel{
			DB: db,
		},
	}
	// Cleanup function to close DB connection
	cleanup := func() {
		log.Println("Closing database connection...")
		if err := db.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
		} else {
			log.Println("Database closed successfully.")
		}
	}

	return appInstance, cleanup, nil
}

func main() {
	// Load templates at startup

	//tpl, err := GetTemplate()
	//if err != nil {
	//	log.Printf(ErrorMsgs().Parse, "./assets/templates/index.html", "main", err)
	//	return
	//}
	//
	//t := tpl.Lookup("index.html")
	//if t == nil {
	//	log.Printf("Template not found: index.html")
	//	return
	//}

	// Initialize the app
	app, cleanup, err := initializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup() // Ensure DB closes on normal exit
	app.init()
	// Handle shutdown signals (Ctrl+C, system shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
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
