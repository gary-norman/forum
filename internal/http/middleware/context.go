package middleware

import (
	"context"

	"github.com/gary-norman/forum/internal/models"
)

// getUserFromContext retrieves the user from the context
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(userContextKey).(*models.User)
	if !ok || user == nil {
		return nil, false
	}
	return user, true
}
