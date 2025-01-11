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

func ErrorMsgs() *models.Errors {
	return models.CreateErrorMessages()
}

func (m *UserModel) Insert(username, email, password, sessionToken, csrfToken, avatar, banner, description string) error {
	stmt := "INSERT INTO Users (Username, Email_address, HashedPassword, SessionToken, CsrfToken, Avatar, Banner, Description, UserType, Created, Is_flagged) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, DateTime('now'), 0)"
	_, err := m.DB.Exec(stmt, username, email, password, sessionToken, csrfToken, avatar, banner, description)
	return err
}

func (m *UserModel) GetUserFromLogin(login, calledBy string) (*models.User, error) {
	Colors := models.CreateColors()
	if m == nil || m.DB == nil {
		return nil, errors.New(fmt.Sprintf(ErrorMsgs().UserModel, "GetUserFromLogin", login))
	}
	username, email := login, login
	fmt.Printf(Colors.Blue+"username called by %v:"+Colors.White+" %v\n"+Colors.Reset, calledBy, username)
	usernameExists, queryUserErr := m.QueryUserNameExists(username)
	if queryUserErr != nil {
		log.Printf(ErrorMsgs().NotFound, "username", username, "GetUserFromLogin", queryUserErr)
	}
	fmt.Printf(Colors.Blue+"usernameExists called by %v:"+Colors.White+" %v\n"+Colors.Reset, calledBy, usernameExists)
	emailExists, queryEmailErr := m.QueryUserEmailExists(email)
	if queryEmailErr != nil {
		log.Printf(ErrorMsgs().NotFound, "email", email, "GetUserFromLogin", queryEmailErr)
	}
	var user *models.User
	if usernameExists != true && emailExists != true {
		return nil, errors.New(fmt.Sprintf("Username & Email: %v do not exist", login))
	}
	var failure string
	if usernameExists == true {
		user, _ = m.GetUserByUsername(username, "GetUserFromLogin")
		failure = "email"
	}

	if emailExists == true {
		user, _ = m.GetUserByEmail(email, "GetUserFromLogin")
		failure = "username"
	}
	log.Printf(Colors.Blue+"*Either username or email should fail as only 1 is entered at login. In this case, "+Colors.White+"%v"+Colors.Red+" failed"+Colors.Blue+" as expected.", failure)
	return user, nil
}

func (m *UserModel) QueryUserNameExists(username string) (bool, error) {
	Colors := models.CreateColors()
	if m == nil || m.DB == nil {
		return false, errors.New(fmt.Sprintf(ErrorMsgs().UserModel, "QueryUserNameExists", username))
	}
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE Username = ?", username).Scan(&count)
	if err != nil {
		log.Fatal(ErrorMsgs().Query, username, err)
	}
	if count > 0 {
		return true, nil
	}
	return false, errors.New(fmt.Sprintf(Colors.Red+"Username does not exist: "+Colors.White+"%v"+Colors.Reset, username))
}
func (m *UserModel) QueryUserEmailExists(email string) (bool, error) {
	Colors := models.CreateColors()
	if m == nil || m.DB == nil {
		return false, errors.New(fmt.Sprintf(ErrorMsgs().UserModel, "QueryUserEmailExists", email))
	}
	var count int
	err := m.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE Email_address = ?", email).Scan(&count)
	if err != nil {
		log.Fatal(ErrorMsgs().Query, email, err)
	}
	if count > 0 {
		return true, nil
	}
	return false, errors.New(fmt.Sprintf(Colors.Red+"Email does not exist: "+Colors.White+"%v"+Colors.Reset, email))
}

func (m *UserModel) GetUserByUsername(username, calledBy string) (*models.User, error) {
	username = strings.TrimSpace(username)
	if m == nil || m.DB == nil {
		log.Printf(fmt.Sprintf(ErrorMsgs().UserModel, "GetUserByUsername", username))
	}
	// Query to fetch user data by username
	stmt, prepErr := m.DB.Prepare("SELECT ID, Username, HashedPassword, SessionToken, CsrfToken FROM Users WHERE Username = ? LIMIT 1")
	if prepErr != nil {
		log.Fatal(ErrorMsgs().Query, username, prepErr)
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			log.Printf(ErrorMsgs().Close, "stmt", "getUserByUsername")
		}
	}(stmt) // Prepared statements take up server resources and should be closed after use.
	// Create a User instance to store the result
	var user models.User
	queryErr := stmt.QueryRow(username).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.SessionToken,
		&user.CSRFToken)
	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			// No user found
			log.Printf(ErrorMsgs().NoRows, username, "getUserByUsername, called by: "+calledBy)
			return nil, nil
		}
		return nil, fmt.Errorf(ErrorMsgs().Query, username, queryErr)
	}

	return &user, nil
}
func (m *UserModel) GetUserByEmail(email, calledBy string) (*models.User, error) {
	email = strings.TrimSpace(email)
	if m == nil || m.DB == nil {
		log.Printf(fmt.Sprintf(ErrorMsgs().UserModel, "GetUserByEmail", email))
	}
	// Query to fetch user data by username
	stmt, prepErr := m.DB.Prepare("SELECT ID, HashedPassword, Email_address FROM Users WHERE Email_address = ? LIMIT 1")
	if prepErr != nil {
		log.Fatal(ErrorMsgs().Query, email, prepErr)
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			log.Printf(ErrorMsgs().Close, "stmt", "getUserByUsername")
		}
	}(stmt) // Prepared statements take up server resources and should be closed after use.
	// Create a User instance to store the result
	var user models.User
	queryErr := stmt.QueryRow(email).Scan(
		&user.ID,
		&user.HashedPassword,
		&user.Email)
	fmt.Printf(ErrorMsgs().Query, email, queryErr)
	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			// No user found
			log.Printf(ErrorMsgs().NoRows, email, "getUserByUsername")
			return nil, nil
		}
		return nil, fmt.Errorf(ErrorMsgs().Query, email, queryErr)
	}

	return &user, nil
}

func (m *UserModel) All() ([]models.User, error) {
	stmt := "SELECT ID, Username, Email_address, HashedPassword, SessionToken, CsrfToken, Avatar, Banner, Description, UserType, Created, Is_flagged FROM Users ORDER BY ID DESC"
	rows, queryErr := m.DB.Query(stmt)
	if queryErr != nil {
		return nil, errors.New(fmt.Sprintf(ErrorMsgs().Query, "Users", queryErr))
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, "rows", "All")
		}
	}()

	var Users []models.User
	for rows.Next() {
		p := models.User{}
		scanErr := rows.Scan(&p.ID, &p.Username, &p.Email, &p.HashedPassword, &p.SessionToken, &p.CSRFToken, &p.Avatar, &p.Banner, &p.Description, &p.Usertype, &p.Created, &p.IsFlagged)
		if scanErr != nil {
			return nil, scanErr
		}
		Users = append(Users, p)
	}

	if queryErr = rows.Err(); queryErr != nil {
		return nil, errors.New(fmt.Sprintf(ErrorMsgs().Query, "Users", queryErr))
	}

	return Users, nil
}
