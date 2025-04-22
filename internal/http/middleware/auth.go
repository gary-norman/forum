package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gary-norman/forum/internal/app"
	"github.com/gary-norman/forum/internal/models"
)

// Create a custom key type to avoid conflicts in context
type contextKey string

const userContextKey = contextKey("currentUser")

// Middleware to add the user to the request context
func WithUser(next http.Handler, app *app.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var currentUser *models.User
		cookie, err := r.Cookie("username")
		if err != nil || cookie.Value == "" {
			log.Printf("Error getting cookie: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		user, err := app.Users.GetUserByUsername(cookie.Value, "WithUser")
		if err != nil {
			log.Printf("No user found: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		currentUser = user

		// Store user in context
		ctx := context.WithValue(r.Context(), userContextKey, currentUser)
		// Pass modified request with context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
