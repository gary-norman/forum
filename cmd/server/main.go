package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gary-norman/forum/internal/app"
	h "github.com/gary-norman/forum/internal/http/routes"
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
	// Initialize the app
	a, cleanup, err := app.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup() // Ensure DB closes on normal exit

	// instantiate temphelper with the initialized app and initialise the templates
	th := view.TempHelper{
		App: a,
	}
	th.Init()

	// Create the RouteHandler with a pointer to the app
	rh := h.NewRouteHandler(a)

	// Handle shutdown signals (Ctrl+C, system shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	port := 8989
	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: rh.Routes(),
	}

	go func() {
		// Log server listening messages
		log.Printf(ErrorMsgs().KeyValuePair, "Starting server on port", port)
		address := "http://localhost" + addr
		log.Printf(ErrorMsgs().ConnSuccess, address)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf(ErrorMsgs().ConnInit, srv.Addr, "srv.ListenAndServe")
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Printf(Colors().Green + "Stopped serving new connections." + Colors().Reset)
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
	log.Printf(Colors().Green + "Graceful shutdown complete." + Colors().Reset)
}
