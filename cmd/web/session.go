package main

import (
	"errors"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"net/http"
)

var AuthErr = errors.New("not authenticated")

func (app *app) isAuthenticated(r *http.Request, username string) error {
	ErrorMsgs := models.CreateErrorMessages()
	Colors := models.CreateColors()
	var user *models.User
	user, getUserErr := app.users.GetUserByUsername(username, "isAuthenticated")
	if getUserErr != nil {
		return errors.New(fmt.Sprintf(ErrorMsgs.NotFound, "user", username, "isAuthenticated", getUserErr))
	}
	// Get the Session Token from the request cookie
	st, err := r.Cookie("session_token")
	if st == nil {
		return errors.New("no session token")
	}
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		fmt.Printf(ErrorMsgs.KeyValuePair, "Cookie SessionToken", st.Value)
		fmt.Printf(ErrorMsgs.KeyValuePair, "Error", err)
		fmt.Printf(ErrorMsgs.KeyValuePair, "User SessionToken", user.SessionToken)
		return AuthErr
	}
	csrf, _ := r.Cookie("csrf_token")

	// Get the CSRF Token from the headers
	csrfToken := r.Header.Get("x-csrf-token")
	fmt.Printf(ErrorMsgs.KeyValuePair, "Header", r.Header)
	if csrfToken == "" || csrfToken != user.CSRFToken {
		fmt.Printf(ErrorMsgs.KeyValuePair, "Cookie csrfToken", csrf.Value)
		fmt.Printf(ErrorMsgs.KeyValuePair, "Header csrfToken", csrfToken)
		fmt.Printf(ErrorMsgs.KeyValuePair, "User csrfToken", user.CSRFToken)
		fmt.Printf(ErrorMsgs.Divider)
		fmt.Printf(Colors.Blue + "Authorise user: " + Colors.Red + "Failed!\n" + Colors.Reset)
		return AuthErr
	}
	fmt.Printf(ErrorMsgs.KeyValuePair, "Cookie SessionToken", st.Value)
	fmt.Printf(ErrorMsgs.KeyValuePair, "Error", err)
	fmt.Printf(ErrorMsgs.KeyValuePair, "User SessionToken", user.SessionToken)
	fmt.Printf(ErrorMsgs.Divider)
	fmt.Printf(ErrorMsgs.KeyValuePair, "Cookie csrfToken", csrf.Value)
	fmt.Printf(ErrorMsgs.KeyValuePair, "Header csrfToken", csrfToken)
	fmt.Printf(ErrorMsgs.KeyValuePair, "User csrfToken", user.CSRFToken)
	fmt.Printf(ErrorMsgs.Divider)
	fmt.Printf(Colors.Blue + "Authorise user: " + Colors.Green + "Success!\n" + Colors.Reset)
	return nil
}
