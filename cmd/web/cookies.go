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

func (m *CookieModel) CreateCookies(w http.ResponseWriter, user *models.User) {
	ErrorMsgs := models.CreateErrorMessages()
	sessionToken := models.GenerateToken(32)
	csrfToken := models.GenerateToken(32)

	// Set Session Token and CSRF Token cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: false,
	})
	// Store tokens in the database
	err := m.UpdateCookies(user, sessionToken, csrfToken)
	if err != nil {
		log.Printf(ErrorMsgs.Cookies, err)
	}
}

func (m *CookieModel) UpdateCookies(user *models.User, sessionToken, csrfToken string) error {
	ErrorMsgs := models.CreateErrorMessages()
	if m == nil || m.DB == nil {
		fmt.Printf(ErrorMsgs.UserModel, "UpdateCookies", user.Username)
		return errors.New("UserModel or DB is nil")
	}
	var stmt string
	fmt.Printf("Updating DB Cookies for: %v\n", user.Username)
	stmt = "UPDATE Users SET SessionToken = ?, CsrfToken = ? WHERE Username = ?"
	result, err := m.DB.Exec(stmt, sessionToken, csrfToken, user.Username)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	fmt.Printf("result.RowsAffected: %v\n", rows)
	return err
}

func (m *CookieModel) DeleteCookies(user *models.User) error {
	stmt := "UPDATE Users SET SessionToken = '', CsrfToken = '' WHERE Username = ?"
	result, err := m.DB.Exec(stmt, user.Username)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	fmt.Printf("result.RowsAffected: %v\n", rows)
	return err
}
