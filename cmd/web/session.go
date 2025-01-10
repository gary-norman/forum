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
		fmt.Printf(Colors.Blue+"st.Value: "+Colors.White+"%v\n"+Colors.Blue+"err: "+Colors.White+"%v\n"+Colors.Blue+
			"user SessionToken: "+Colors.White+"%v\n"+Colors.Reset, st.Value, err, user.SessionToken)
		return AuthErr
	}
	csrf, _ := r.Cookie("csrf_token")

	// Get the CSRF Token from the headers
	csrfToken := r.Header.Get("x-csrf-token")
	if csrfToken == "" || csrfToken != user.CSRFToken {
		fmt.Printf(Colors.Blue+"cookie csrfToken: "+Colors.White+"%v\n"+Colors.Blue+"header csrfToken: "+Colors.White+
			"%v\n"+Colors.Blue+"user csrfToken: "+Colors.White+"%v\n"+Colors.Blue+"Authorise user: "+Colors.Red+"Failed!\n"+
			Colors.Reset, csrf.Value, csrfToken, user.CSRFToken)
		return AuthErr
	}
	fmt.Printf(
		Colors.Blue+"cookie SessionToken: "+Colors.White+"%v\n"+Colors.Blue+"user SessionToken: "+Colors.White+"%v\n"+
			Colors.Blue+"cookie csrfToken: "+Colors.White+"%v\n "+Colors.Blue+"header csrfToken: "+Colors.White+"%v\n"+
			Colors.Blue+"user CSRFToken: "+Colors.White+"%v\n"+Colors.Blue+"Authorise user: "+Colors.Green+"Success!\n"+Colors.Reset,
		st.Value, user.SessionToken, csrf.Value, csrfToken, user.CSRFToken)

	return nil
}
