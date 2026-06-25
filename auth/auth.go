package auth

import (
	"context"
	"errors"
	"strconv"
)

// AuthenticatedUser represents the authenticated user stored in context.
type AuthenticatedUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Int64ID returns the authenticated user ID as int64 for applications that use
// numeric user IDs internally.
func (u *AuthenticatedUser) Int64ID() (int64, error) {
	if u == nil || u.ID == "" {
		return 0, ErrAuthenticatedUserNotFound
	}
	id, err := strconv.ParseInt(u.ID, 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrInvalidAuthenticatedUserID
	}
	return id, nil
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

// ErrInvalidAuthenticatedUserID is returned when an authenticated user's ID cannot be used as a positive int64.
var ErrInvalidAuthenticatedUserID = errors.New("auth: invalid authenticated user id")

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
