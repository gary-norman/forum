package main

import (
	"errors"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"net/http"
)

var AuthErr = errors.New("not authenticated")

func (app *app) isAuthenticated(r *http.Request) error {
	login := r.FormValue("username")
	var user *models.User
	user, err := app.users.GetUserFromLogin(login)
	if err != nil {
		return err
	}

	// Get the Session Token from the request cookie
	st, err := r.Cookie("session_token")
	if st == nil {
		return errors.New("no session token")
	}
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		fmt.Printf("st.Value: %v\nerr: %v\nuser SessionToken: %v\n", st.Value, err, user.SessionToken)
		return AuthErr
	}
	csrf, _ := r.Cookie("csrf_token")

	// Get the CSRF Token from the headers
	csrfToken := r.Header.Get("X-CSRF-Token")
	if csrfToken == "" || csrfToken != user.CSRFToken {
		fmt.Printf("csrf: %v\ncsrfToken: %v\nuser CSRFToken: %v\n", csrf.Value, csrfToken, user.CSRFToken)
		return AuthErr
	}

	return nil
}
