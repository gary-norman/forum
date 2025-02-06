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
