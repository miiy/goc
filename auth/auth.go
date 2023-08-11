package auth

import (
	"context"
)

type AuthenticatedUser struct {
	ID       int64
	Username string
}

type UserProvider interface {
	FirstByIdentifier(ctx context.Context, identifier string) (*AuthenticatedUser, error)
}
