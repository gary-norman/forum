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
	stmt := "INSERT INTO Users (Username, Email_address, HashedPassword, SessionToken, CsrfToken, Avatar, Banner, Description, UserType, Created, Is_flagged) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, DateTime('now'), 0)"
	_, err := m.DB.Exec(stmt, username, email, password, sessionToken, csrfToken, avatar, banner, description)
	return err
}

func (m *UserModel) GetUserFromLogin(login string) (*models.User, error) {
	ErrorMsgs := models.CreateErrorMessages()

	if m == nil || m.DB == nil {
		return nil, errors.New(fmt.Sprintf(ErrorMsgs.UserModel, "GetUserFromLogin", login))
	}
	username, email := login, login
	fmt.Printf("username: %v\n", username)
	usernameExists, queryUserErr := m.QueryUserNameExists(username)
	if queryUserErr != nil {
		log.Printf("username: %v not found", username)
	}
	fmt.Printf("usernameExists: %v\n", usernameExists)
	emailExists, queryEmailErr := m.QueryUserEmailExists(email)
	if queryEmailErr != nil {
		log.Printf("email: %v not found", email)
	}
	var user *models.User
	if usernameExists != true && emailExists != true {
		return nil, errors.New("username or email not found")
	}

	if usernameExists == true {
		user, _ = m.GetUserByUsername(username)
	}

	if emailExists == true {
		user, _ = m.GetUserByEmail(email)
	}

	return user, nil
}

func (m *UserModel) QueryUserNameExists(username string) (bool, error) {
	ErrorMsgs := models.CreateErrorMessages()
	if m == nil || m.DB == nil {
		fmt.Printf(ErrorMsgs.UserModel, "QueryUserNameExists", username)
		return false, errors.New("UserModel or DB is nil")
	}
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE Username = ?", username).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return true, nil
	}
	return false, errors.New("username not found")
}
func (m *UserModel) QueryUserEmailExists(email string) (bool, error) {
	ErrorMsgs := models.CreateErrorMessages()
	if m == nil || m.DB == nil {
		fmt.Printf(ErrorMsgs.UserModel, "QueryUserEmailExists", email)
		return false, errors.New("UserModel or DB is nil")
	}
	var count int
	//ErrorMsgs := models.CreateErrorMessages()
	err := m.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE Email_address = ?", email).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return true, nil
	}
	return false, errors.New("email not found")
}

func (m *UserModel) GetUserByUsername(username string) (*models.User, error) {
	ErrorMsgs := models.CreateErrorMessages()
	username = strings.TrimSpace(username)
	if m == nil || m.DB == nil {
		fmt.Printf(ErrorMsgs.UserModel, "GetUserByUsername", username)
		log.Printf("UserModel or DB is nil")
	}
	// Query to fetch user data by username
	stmt, err := m.DB.Prepare("SELECT ID, Username, HashedPassword, SessionToken, CsrfToken FROM Users WHERE Username = ? LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			log.Printf(ErrorMsgs.Close, "stmt", "getUserByUsername")
		}
	}(stmt) // Prepared statements take up server resources and should be closed after use.
	// Create a User instance to store the result
	var user models.User
	err = stmt.QueryRow(username).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.SessionToken,
		&user.CSRFToken)
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
func (m *UserModel) GetUserByEmail(email string) (*models.User, error) {
	ErrorMsgs := models.CreateErrorMessages()
	email = strings.TrimSpace(email)
	if m == nil || m.DB == nil {
		fmt.Printf(ErrorMsgs.UserModel, "GetUserByEmail", email)
		log.Printf("UserModel or DB is nil")
	}
	// Query to fetch user data by username
	stmt, err := m.DB.Prepare("SELECT ID, HashedPassword, Email_address FROM Users WHERE Email_address = ? LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			log.Printf(ErrorMsgs.Close, "stmt", "getUserByUsername")
		}
	}(stmt) // Prepared statements take up server resources and should be closed after use.
	// Create a User instance to store the result
	var user models.User
	err = stmt.QueryRow(email).Scan(
		&user.ID,
		&user.Login.HashedPassword,
		&user.Login.Email)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found
			log.Printf(ErrorMsgs.NoRows, email, "getUserByUsername")
			return nil, nil
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &user, nil
}

func (m *UserModel) All() ([]models.User, error) {
	ErrorMsgs := models.CreateErrorMessages()
	stmt := "SELECT ID, Username, Email_address, HashedPassword, SessionToken, CsrfToken, Avatar, Banner, Description, UserType, Created, Is_flagged FROM Users ORDER BY ID DESC"
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
