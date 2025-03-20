package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
	"net/http"
	"time"
)

type CookieModel struct {
	DB *sql.DB
}

func (m *CookieModel) CreateCookies(w http.ResponseWriter, user *models.User) error {
	sessionToken := models.GenerateToken(32)
	csrfToken := models.GenerateToken(32)
	expires := time.Now().Add(time.Hour * 24)

	// Set Session, Username, and CSRF Token cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expires,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    user.Username,
		Expires:  expires,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  expires,
		HttpOnly: false,
	})
	// Store tokens in the database
	err := m.UpdateCookies(user, sessionToken, csrfToken)
	if err != nil {
		log.Printf(ErrorMsgs().Cookies, "update", err)
		return err
	}
	return nil
}

// GetLoggedInUser gets the currently logged-in user from the session token and returns the user's struct
func (app *app) GetLoggedInUser(r *http.Request) (*models.User, error) {
	// Get the username from the request cookie
	userCookie, getCookieErr := r.Cookie("username")
	if getCookieErr != nil {
		log.Printf(ErrorMsgs().Cookies, "get", getCookieErr)
		return nil, getCookieErr
	}
	var username string
	if userCookie != nil {
		username = userCookie.Value
	}
	fmt.Printf(ErrorMsgs().KeyValuePair, "Username", username)
	if username == "" {
		return nil, errors.New("no user is logged in")
	}
	user, getUserErr := app.users.GetUserByUsername(username, "GetLoggedInUser")
	if getUserErr != nil {
		return nil, getUserErr
	}
	return user, nil
}

func (m *CookieModel) QueryCookies(w http.ResponseWriter, r *http.Request, user *models.User) bool {
	Colors := models.CreateColors()
	var success bool

	// Get the Session Token from the request cookie
	st, err := r.Cookie("session_token")
	if err != nil {
		log.Printf(ErrorMsgs().Cookies, "query", err)
		return false
	}
	csrf, _ := r.Cookie("csrf_token")

	// Get the CSRF Token from the headers
	csrfToken := r.Header.Get("x-csrf-token")

	stColor, csrfColor := Colors.Red, Colors.Red
	stMatchString, csrfMatchString := "Failed!", "Failed!"
	if st.Value == user.SessionToken {
		stColor = Colors.Green
		stMatchString = "Success!"
		success = true
	} else {
		err := m.DeleteCookies(w, user)
		if err != nil {
			log.Printf("error deleting cookies")
		}
		success = false
	}
	if csrf.Value == csrfToken && csrfToken == user.CSRFToken {
		csrfColor = Colors.Green
		csrfMatchString = "Success!"
	}
	log.Printf(ErrorMsgs().KeyValuePair, "Cookie SessionToken", st.Value)
	log.Printf(ErrorMsgs().KeyValuePair, "User SessionToken", user.SessionToken)
	log.Printf(Colors.Blue+"Session token verficiation: "+stColor+"%v\n"+Colors.Reset, stMatchString)
	fmt.Printf(ErrorMsgs().Divider)
	log.Printf(ErrorMsgs().KeyValuePair, "Cookie csrfToken", csrf.Value)
	log.Printf(ErrorMsgs().KeyValuePair, "Header csrfToken", csrfToken)
	log.Printf(ErrorMsgs().KeyValuePair, "User csrfToken", user.CSRFToken)
	log.Printf(Colors.Blue+"CSRF token verficiation: "+csrfColor+"%v\n"+Colors.Reset, csrfMatchString)

	return success
}

func (m *CookieModel) UpdateCookies(user *models.User, sessionToken, csrfToken string) error {
	Colors := models.CreateColors()
	if m == nil || m.DB == nil {
		fmt.Printf(ErrorMsgs().UserModel, "UpdateCookies", user.Username)
		return errors.New("UserModel or DB is nil")
	}
	var stmt string
	fmt.Printf(Colors.Blue+"Updating DB Cookies for: "+Colors.White+"%v\n"+Colors.Reset, user.Username)
	stmt = "UPDATE Users SET SessionToken = ?, CsrfToken = ? WHERE Username = ?"
	result, err := m.DB.Exec(stmt, sessionToken, csrfToken, user.Username)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	dbUpdated := "Failed!"
	dbUpdatedColor := Colors.Red
	if rows > 0 {
		dbUpdated = "Success!"
		dbUpdatedColor = Colors.Green
	}
	fmt.Printf(Colors.Blue+"Database update: "+dbUpdatedColor+"%v\n", dbUpdated)

	return err
}

func (m *CookieModel) DeleteCookies(w http.ResponseWriter, user *models.User) error {
	Colors := models.CreateColors()
	expires := time.Now().Add(time.Hour - 1000)
	stmt := "UPDATE Users SET SessionToken = '', CsrfToken = '' WHERE Username = ?"
	result, err := m.DB.Exec(stmt, user.Username)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	dbUpdated := "Failed!"
	dbUpdatedColor := Colors.Red
	if rows > 0 {
		dbUpdated = "Success!"
		dbUpdatedColor = Colors.Green
	}
	fmt.Printf(Colors.Blue+"Database update: "+dbUpdatedColor+"%v\n", dbUpdated)
	// Set Session, Username, and CSRF Token cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  expires,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    "",
		Expires:  expires,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  expires,
		HttpOnly: false,
	})
	return err
}
