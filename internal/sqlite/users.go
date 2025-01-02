package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gary-norman/forum/internal/models"
	"log"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, email, password, sessionToken, csrfToken, avatar, banner, description string) error {
	stmt := "INSERT INTO Users (Username, Email_address, HashedPassword, SessionToken, Csrf_token, Avatar, Banner, Description, UserType, Created, Is_flagged) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, DateTime('now'), 0)"
	_, err := m.DB.Exec(stmt, username, email, password, sessionToken, csrfToken, avatar, banner, description)
	return err
}

func (m *UserModel) QueryUserEmailExists(email string) bool {
	if m == nil || m.DB == nil {
		log.Printf("UserModel or DB is nil")
		return false
	}
	ErrorMsgs := models.CreateErrorMessages()
	query := "SELECT 1 FROM Users WHERE Email_address = ? LIMIT 1"
	var exists int
	err := m.DB.QueryRow(query, email).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No rows returned means the user doesn't exist
			log.Printf(ErrorMsgs.NoRows, email, "QueryUserEmailExists")
			return false
		}
		// Return other errors
		log.Printf(ErrorMsgs.Query, email, err)
		return false
	}

	// User exists
	return true
}
func (m *UserModel) QueryUserNameExists(username string) bool {
	if m == nil || m.DB == nil {
		log.Printf("UserModel or DB is nil")
		return false
	}
	query := "SELECT 1 FROM Users WHERE Username = ? LIMIT 1"
	var exists int
	err := m.DB.QueryRow(query, username).Scan(&exists)

	ErrorMsgs := models.CreateErrorMessages()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No rows returned means the user doesn't exist
			log.Printf(ErrorMsgs.NoRows, username, "QueryUserNameExists")
			return false
		}
		// Return other errors
		log.Printf(ErrorMsgs.Query, username, err)
		return false
	}

	// User exists
	return true
}

func (m *UserModel) GetUserByEmail(email string) (*models.User, error) {
	if m == nil || m.DB == nil {
		log.Printf("UserModel or DB is nil")
	}
	// Query to fetch user data by email
	query := "SELECT * FROM Users WHERE Email_address = ? LIMIT 1"
	row := m.DB.QueryRow(query, email)

	// Create a User instance to store the result
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Login.Email, &user.Login.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found
			ErrorMsgs := models.CreateErrorMessages()
			log.Printf(ErrorMsgs.NoRows, email, "getUserByEmail")
			return nil, nil
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &user, nil
}
func (m *UserModel) GetUserByUsername(username string) (*models.User, error) {
	if m == nil || m.DB == nil {
		log.Printf("UserModel or DB is nil")
	}
	// Query to fetch user data by username
	query := "SELECT * FROM Users WHERE Username = ? LIMIT 1"
	username = strings.TrimSpace(strings.ToLower(username))
	row := m.DB.QueryRow(query, username)
	fmt.Printf("query: %s\nrow: %v\n", username, row)
	// Create a User instance to store the result
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Login.Email, &user.Login.HashedPassword)
	ErrorMsgs := models.CreateErrorMessages()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found
			log.Printf(ErrorMsgs.NoRows, username, "getUserByUsername")
			return nil, nil
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &user, nil
}

func (m *UserModel) All() ([]models.User, error) {
	ErrorMsgs := models.CreateErrorMessages()
	stmt := "SELECT ID, Username, Email_address, HashedPassword, SessionToken, Csrf_token, Avatar, Banner, Description, UserType, Created, Is_flagged FROM Users ORDER BY ID DESC"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs.Close, "rows", "All")
		}
	}()

	var Users []models.User
	for rows.Next() {
		p := models.User{}
		err = rows.Scan(&p.ID, &p.Username, &p.Email, &p.HashedPassword, &p.SessionToken, &p.CSRFToken, &p.Avatar, &p.Banner, &p.Description, &p.Usertype, &p.Created, &p.IsFlagged)
		if err != nil {
			return nil, err
		}
		Users = append(Users, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Users, nil
}
