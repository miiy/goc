package auth

import (
	"context"
)

type AuthenticatedUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserProvider interface {
	FirstByIdentifier(ctx context.Context, identifier string) (*AuthenticatedUser, error)
}
