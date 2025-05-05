package service

import (
	"context"

	"github.com/gary-norman/forum/internal/models"
)

// following tutorial at https://medium.com/@moeid_72300/elevate-your-golang-tests-with-database-mocking-a-step-by-step-guide-ee961da7600

type UserRepository interface {
	EditUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error
	// Additional methods related to users...
}

var defaultUserType = "user"

func NewUser(username, email, password string) (*models.User, error) {
	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		return nil, err
	}
	id := models.NewUUIDField()

	return &models.User{
		ID:          id,
		Username:    username,
		Avatar:      "noimage_" + username,
		Banner:      "default.png",
		Description: "",
		Usertype:    defaultUserType,
		Login: models.Login{
			Email:          email,
			SessionToken:   "",
			CSRFToken:      "",
			HashedPassword: hashedPassword,
		},
	}, nil
}
