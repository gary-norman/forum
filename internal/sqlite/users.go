package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gary-norman/forum/internal/models"
)

type UserModel struct {
	DB *sql.DB
}

func ErrorMsgs() *models.Errors {
	return models.CreateErrorMessages()
}

func getNextSortID(db *sql.DB) (int, error) {
	var next int
	err := db.QueryRow(`SELECT COALESCE(MAX(SortID), 0) + 1 FROM Users`).Scan(&next)
	return next, err
}

func (m *UserModel) Insert(id models.UUIDField, username, email, password, avatar, banner, userType string) error {
	nextID, execErr := getNextSortID(m.DB)
	if execErr != nil {
		log.Fatal("Error getting next SortID:", execErr)
	}

	// FIXME this prepare statement is unnecessary as it is not used in a loop
	stmt, insertErr := m.DB.Prepare("INSERT INTO Users (ID, SortID, Username, EmailAddress, HashedPassword, Avatar, Banner, UserType, Created, IsFlagged) VALUES (?, ?, ?, ?, ?, ?, ?, ?, DateTime('now'), 0)")
	if insertErr != nil {
		log.Printf(ErrorMsgs().Query, username, insertErr)
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			log.Printf(ErrorMsgs().Close, "stmt", "insert", closErr)
		}
	}(stmt) // Prepared statements take up server resources and should be closed after use.
	_, err := stmt.Exec(id, nextID, username, email, password, avatar, banner, userType)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Edit(user *models.User) error {
	stmt, prepErr := m.DB.Prepare("UPDATE Users SET Username = ?, EmailAddress = ?, HashedPassword = ?, SessionToken = ?, CsrfToken = ?, Avatar = ?, Banner = ?, Description = ? WHERE ID = ?")
	if prepErr != nil {
		log.Printf(ErrorMsgs().Query, "Users", prepErr)
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			log.Printf(ErrorMsgs().Close, "stmt", "edit", closErr)
		}
	}(stmt) // Prepared statements take up server resources and should be closed after use.
	_, err := stmt.Exec(user.Username, user.Email, user.HashedPassword, user.SessionToken, user.CSRFToken, user.Avatar, user.Banner, user.Description, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Delete(user *models.User) error {
	stmt, prepErr := m.DB.Prepare("DELETE FROM Users WHERE ID = ?")
	if prepErr != nil {
		log.Printf(ErrorMsgs().Query, "Users", prepErr)
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			log.Printf(ErrorMsgs().Close, "stmt", "delete", closErr)
		}
	}(stmt) // Prepared statements take up server resources and should be closed after use.
	_, err := stmt.Exec(user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) GetUserFromLogin(login, calledBy string) (*models.User, error) {
	if m == nil || m.DB == nil {
		return nil, fmt.Errorf(ErrorMsgs().UserModel, "GetUserFromLogin", login)
	}
	username, email := login, login
	var loginType string
	usernameQuery, ok, _ := m.QueryUserNameExists(username)
	if ok {
		loginType = usernameQuery
	}
	emailQuery, ok, _ := m.QueryUserEmailExists(email)
	if ok {
		loginType = emailQuery
	}
	switch loginType {
	case "username":
		user, err := m.GetUserByUsername(username, "GetUserFromLogin")
		if err != nil {
			return nil, err
		} else {
			log.Printf(ErrorMsgs().KeyValuePair, "Successfully found user by username", user.Username)
			return user, nil
		}
	case "email":
		user, err := m.GetUserByEmail(email, "GetUserFromLogin")
		if err != nil {
			return nil, err
		} else {
			log.Printf(ErrorMsgs().KeyValuePair, "Successfully found user by email", user.Username)
			return user, nil
		}
	default:
		return nil, fmt.Errorf("user: %v not found", login)
	}
}

func (m *UserModel) QueryUserNameExists(username string) (string, bool, error) {
	if m == nil || m.DB == nil {
		err := fmt.Errorf("error connecting to database: %s", "QueryUserNameExists")
		return "", false, err

	}
	var count int
	queryErr := m.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE Username = ?", username).Scan(&count)
	if queryErr != nil {
		log.Printf(ErrorMsgs().Query, username, queryErr)
		return "", false, queryErr
	}
	if count > 0 {
		return "username", true, queryErr
	}
	return "", false, queryErr
}

func (m *UserModel) QueryUserEmailExists(email string) (string, bool, error) {
	if m == nil || m.DB == nil {
		err := fmt.Errorf("error connecting to database: %s", "QueryUserEmailExists")
		return "", false, err
	}
	var count int
	queryErr := m.DB.QueryRow("SELECT COUNT(*) FROM Users WHERE EmailAddress = ?", email).Scan(&count)
	if queryErr != nil {
		log.Printf(ErrorMsgs().Query, email, queryErr)
		return "", false, queryErr
	}
	if count > 0 {
		return "username", true, queryErr
	}
	return "", false, queryErr
}

// TODO unify these functions to accept parameters

func (m *UserModel) GetUserByUsername(username, calledBy string) (*models.User, error) {
	username = strings.TrimSpace(username)
	if m == nil || m.DB == nil {
		log.Printf(ErrorMsgs().UserModel, "GetUserByUsername", username)
	}
	// Query to fetch user data by username
	stmt, prepErr := m.DB.Prepare("SELECT * FROM Users WHERE Username = ? LIMIT 1")
	if prepErr != nil {
		log.Printf(ErrorMsgs().Query, username, prepErr)
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			log.Printf(ErrorMsgs().Close, "stmt", "getUserByUsername", closErr)
		}
	}(stmt) // Prepared statements take up server resources and should be closed after use.
	// Create a User instance to store the result
	var user models.User
	queryErr := stmt.QueryRow(username).Scan(
		&user.ID,
		&user.SortID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.Banner,
		&user.Description,
		&user.Usertype,
		&user.Created,
		&user.IsFlagged,
		&user.SessionToken,
		&user.CSRFToken,
		&user.HashedPassword)
	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			return nil, fmt.Errorf(ErrorMsgs().NoRows, username, calledBy, queryErr)
		}
		return nil, fmt.Errorf(ErrorMsgs().Query, username, queryErr)
	}

	return &user, nil
}

func (m *UserModel) GetUserByEmail(email, calledBy string) (*models.User, error) {
	email = strings.TrimSpace(email)
	if m == nil || m.DB == nil {
		log.Printf(ErrorMsgs().UserModel, "GetUserByEmail", email)
	}
	// Query to fetch user data by username
	stmt, prepErr := m.DB.Prepare("SELECT ID, HashedPassword, EmailAddress FROM Users WHERE EmailAddress = ? LIMIT 1")
	if prepErr != nil {
		log.Fatal(ErrorMsgs().Query, email, prepErr)
	}
	defer func(stmt *sql.Stmt) {
		closErr := stmt.Close()
		if closErr != nil {
			// FIXME this error
			log.Printf(ErrorMsgs().Close, "stmt", "getUserByEmail")
		}
	}(stmt)
	// Create a User instance to store the result
	var user models.User
	queryErr := stmt.QueryRow(email).Scan(
		&user.ID,
		&user.HashedPassword,
		&user.Email)
	fmt.Printf(ErrorMsgs().Query, email, queryErr)
	if queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			return nil, fmt.Errorf(ErrorMsgs().NoRows, email, calledBy, queryErr)
		}
		return nil, fmt.Errorf(ErrorMsgs().Query, email, queryErr)
	}

	return &user, nil
}

func (m *UserModel) GetUserByID(ID int64) (models.User, error) {
	stmt := "SELECT * FROM Users WHERE ID = ?"
	row := m.DB.QueryRow(stmt, ID)
	u := models.User{}
	err := row.Scan(
		&u.ID,
		&u.SortID,
		&u.Username,
		&u.Email,
		&u.Avatar,
		&u.Banner,
		&u.Description,
		&u.Usertype,
		&u.Created,
		&u.IsFlagged,
		&u.SessionToken,
		&u.CSRFToken,
		&u.HashedPassword)
	if err != nil {
		log.Printf(ErrorMsgs().Query, "GetUserFromId", err)
		return u, err
	}
	models.UpdateTimeSince(&u)
	return u, nil
}

// TODO accept an interface for any given value
func isValidUserColumn(column string) bool {
	validColumns := map[string]bool{
		"ID":             true,
		"SortID":         true,
		"Username":       true,
		"EmailAddress":   true,
		"HashedPassword": true,
		"SessionToken":   true,
		"CsrfToken":      true,
		"Avatar":         true,
		"Banner":         true,
		"Description":    true,
		"UserType":       true,
		"Created":        true,
		"IsFlagged":      true,
	}
	return validColumns[column]
}

// GetSingleUserValue returns the string of the column specified in output, which should be entered in all lower case
func (m *UserModel) GetSingleUserValue(ID models.UUIDField, searchColumn, outputColumn string) (string, error) {
	if !isValidUserColumn(searchColumn) {
		return "", fmt.Errorf("invalid searchColumn name: %s", searchColumn)
	}
	stmt := fmt.Sprintf(
		"SELECT * FROM Users WHERE %s = ?",
		searchColumn,
	)
	rows, queryErr := m.DB.Query(stmt, ID)
	if queryErr != nil {
		return "", fmt.Errorf(ErrorMsgs().Query, "GetSingleUserValue", queryErr)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, "rows", "All")
		}
	}()
	var user models.User
	if rows.Next() {
		if scanErr := rows.Scan(
			&user.ID, &user.SortID, &user.Username, &user.Email, &user.Avatar, &user.Banner, &user.Description, &user.Usertype,
			&user.Created, &user.IsFlagged, &user.SessionToken, &user.CSRFToken, &user.HashedPassword); scanErr != nil {
			return "", scanErr
		}
	} else {
		return "", fmt.Errorf("no user found")
	}

	// Map searchColumn names to their values
	fields := map[string]any{
		"id":             user.ID,
		"sortId":         user.SortID,
		"username":       user.Username,
		"email":          user.Email,
		"hashedPassword": user.HashedPassword,
		"sessionToken":   user.SessionToken,
		"csrfToken":      user.CSRFToken,
		"avatar":         user.Avatar,
		"banner":         user.Banner,
		"description":    user.Description,
		"usertype":       user.Usertype,
		"created":        user.Created,
		"isFlagged":      user.IsFlagged,
	}

	// Check if outputColumn exists in the map
	value, exists := fields[outputColumn]
	if !exists {
		return "", fmt.Errorf("invalid search Column name: %s", outputColumn)
	}

	// Convert the value to a string (handling different types)
	outputValue := fmt.Sprintf("%v", value)
	return outputValue, nil
}

func (m *UserModel) All() ([]models.User, error) {
	stmt := "SELECT * FROM Users ORDER BY ID DESC"
	rows, queryErr := m.DB.Query(stmt)
	if queryErr != nil {
		return nil, fmt.Errorf(ErrorMsgs().Query, "Users", queryErr)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf(ErrorMsgs().Close, "rows", "All")
		}
	}()

	users := make([]models.User, 0)
	for rows.Next() {
		p, err := parseUserRows(rows)
		if err != nil {
			return nil, fmt.Errorf("error parsing row: %w", err)
		}
		users = append(users, *p)
	}
	return users, nil
}

func parseUserRow(row *sql.Row) (*models.User, error) {
	var user models.User

	if err := row.Scan(
		&user.ID,
		&user.SortID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.Banner,
		&user.Description,
		&user.Usertype,
		&user.Created,
		&user.IsFlagged,
		&user.SessionToken,
		&user.CSRFToken,
		&user.HashedPassword,
	); err != nil {
		log.Printf(ErrorMsgs().Query, "GetUserFromId", err)
		return nil, err
	}
	models.UpdateTimeSince(&user)
	return &user, nil
}

func parseUserRows(rows *sql.Rows) (*models.User, error) {
	var user models.User

	if err := rows.Scan(
		&user.ID,
		&user.SortID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.Banner,
		&user.Description,
		&user.Usertype,
		&user.Created,
		&user.IsFlagged,
		&user.SessionToken,
		&user.CSRFToken,
		&user.HashedPassword,
	); err != nil {
		log.Printf(ErrorMsgs().Query, "GetUserFromId", err)
		return nil, err
	}
	models.UpdateTimeSince(&user)
	return &user, nil
}
