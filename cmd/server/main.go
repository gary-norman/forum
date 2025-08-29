package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gary-norman/forum/internal/app"
	"github.com/gary-norman/forum/internal/http/routes"
	"github.com/gary-norman/forum/internal/models"
	"github.com/gary-norman/forum/internal/view"

	_ "github.com/mattn/go-sqlite3"
)

func ErrorMsgs() *models.Errors {
	return models.CreateErrorMessages()
}

func Colors() *models.Colors {
	return models.CreateColors()
}

func main() {
	// pprof server for profiling at http://localhost:6060/debug/pprof/
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Initialize the app
	appInstance, cleanup, err := app.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup() // Ensure DB closes on normal exit

	// instantiate temphelper with the initialized app and initialise the templates
	th := view.TempHelper{
		App: appInstance,
	}
	th.Init()

	// Create the RouteHandler with a pointer to the app
	// rh := h.NewRouteHandler(appInstance)
	// mux := rh.Routes()

	router := routes.NewRouter(appInstance)

	// TODO figure this out
	// Count users and create new admin account if none exist
	// userCount, err := app.App.Users.CountUsers(db)

	port := 8888
	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	// Log server listening messages
	address := "http://localhost" + addr
	fmt.Printf(ErrorMsgs().KeyValuePair, "Starting server on port", port)
	fmt.Printf(ErrorMsgs().ConnSuccess, address)

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf(ErrorMsgs().ConnInit, srv.Addr, "srv.ListenAndServe")
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println(Colors().Green + "Stopped serving new connections." + Colors().Reset)
	}()

	// Handle shutdown signals (Ctrl+C, system shutdown)
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)
	<-shutdownSignal

	// create cancellation signal and timeout
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	// shut down the server
	log.Println("Shutting down gracefully...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf(ErrorMsgs().Shutdown, err)
	}
	log.Println(Colors().Green + "Graceful shutdown complete." + Colors().Reset)
}
