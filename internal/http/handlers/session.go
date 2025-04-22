package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gary-norman/forum/internal/app"
	"github.com/gary-norman/forum/internal/models"
)

var ErrAuth = errors.New("not authenticated")

type SessionHandler struct {
	App *app.App
}

func (s *SessionHandler) IsAuthenticated(r *http.Request, username string) error {
	Colors := models.CreateColors()
	var user *models.User
	user, getUserErr := s.App.Users.GetUserByUsername(username, "isAuthenticated")
	if getUserErr != nil {
		return fmt.Errorf(ErrorMsgs().NotFound, username, "isAuthenticated", getUserErr)
	}
	// Get the Session Token from the request cookie
	st, err := r.Cookie("session_token")
	if st == nil {
		return errors.New("no session token")
	}
	if err != nil || st.Value == "" || st.Value != user.SessionToken {
		fmt.Printf(ErrorMsgs().KeyValuePair, "Cookie SessionToken", st.Value)
		fmt.Printf(ErrorMsgs().KeyValuePair, "Error", err)
		fmt.Printf(ErrorMsgs().KeyValuePair, "User SessionToken", user.SessionToken)
		return ErrAuth
	}
	csrf, _ := r.Cookie("csrf_token")

	// Get the CSRF Token from the headers
	csrfToken := r.Header.Get("x-csrf-token")
	fmt.Printf(ErrorMsgs().KeyValuePair, "Header", r.Header)
	if csrfToken == "" || csrfToken != user.CSRFToken {
		fmt.Printf(ErrorMsgs().KeyValuePair, "Cookie csrfToken", csrf.Value)
		fmt.Printf(ErrorMsgs().KeyValuePair, "Header csrfToken", csrfToken)
		fmt.Printf(ErrorMsgs().KeyValuePair, "User csrfToken", user.CSRFToken)
		fmt.Printf(ErrorMsgs().Divider)
		fmt.Printf(Colors.Blue + "Authorise user: " + Colors.Red + "Failed!\n" + Colors.Reset)
		return ErrAuth
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "Cookie SessionToken", st.Value)
	fmt.Printf(ErrorMsgs().KeyValuePair, "Error", err)
	fmt.Printf(ErrorMsgs().KeyValuePair, "User SessionToken", user.SessionToken)
	fmt.Printf(ErrorMsgs().Divider)
	fmt.Printf(ErrorMsgs().KeyValuePair, "Cookie csrfToken", csrf.Value)
	fmt.Printf(ErrorMsgs().KeyValuePair, "Header csrfToken", csrfToken)
	fmt.Printf(ErrorMsgs().KeyValuePair, "User csrfToken", user.CSRFToken)
	fmt.Printf(ErrorMsgs().Divider)
	fmt.Printf(Colors.Blue + "Authorise user: " + Colors.Green + "Success!\n" + Colors.Reset)
	return nil
}
