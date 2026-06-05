package auth

import (
	"context"
	"errors"
)

// AuthenticatedUser represents the authenticated user stored in context.
type AuthenticatedUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// UserProvider looks up users by identifier.
type UserProvider interface {
	FirstByIdentifier(ctx context.Context, identifier string) (*AuthenticatedUser, error)
}

const (
	// AuthenticatedUserIDMetadataKey is the gRPC metadata key used to forward authenticated user IDs.
	AuthenticatedUserIDMetadataKey = "x-auth-user-id"
	// AuthenticatedUsernameMetadataKey is the gRPC metadata key used to forward authenticated usernames.
	AuthenticatedUsernameMetadataKey = "x-auth-username"
)

// ErrAuthenticatedUserNotFound is returned when no authenticated user is found in context.
var ErrAuthenticatedUserNotFound = errors.New("auth: authenticated user not found")

type authenticatedUserContextKey string

const authenticatedUserContextKeyValue authenticatedUserContextKey = "goc.auth.AuthenticatedUser"

// InjectAuthenticatedUser stores the authenticated user in the context.
func InjectAuthenticatedUser(ctx context.Context, u *AuthenticatedUser) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if u == nil {
		return ctx
	}
	return context.WithValue(ctx, authenticatedUserContextKeyValue, u)
}

// ExtractAuthenticatedUser retrieves the authenticated user from the context.
func ExtractAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	if ctx == nil {
		return nil, ErrAuthenticatedUserNotFound
	}

	u, ok := ctx.Value(authenticatedUserContextKeyValue).(*AuthenticatedUser)
	if !ok || u == nil {
		return nil, ErrAuthenticatedUserNotFound
	}
	return u, nil
}
