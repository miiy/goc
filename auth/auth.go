package auth

import (
	"context"
	"errors"
	"strconv"
)

// AuthenticatedUser represents the user principal within an authenticated session.
type AuthenticatedUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// AuthenticatedSession is the trusted session information obtained from an
// authenticated Cookie Session or access token. SIDKey is the non-secret lookup
// key derived from the raw session ID; the raw session ID must not be stored here.
type AuthenticatedSession struct {
	User       AuthenticatedUser `json:"user"`
	SIDKey     string            `json:"sid_key"`
	ClientType string            `json:"client_type"`
	ClientID   string            `json:"client_id"`
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
	// AuthenticatedUsernameMetadataKey is the gRPC metadata key used to forward a display name.
	// Authorization must use AuthenticatedUserIDMetadataKey only.
	AuthenticatedUsernameMetadataKey = "x-auth-username"
)

// ErrAuthenticatedUserNotFound is returned when an authenticated user is nil or
// has no identifier.
var ErrAuthenticatedUserNotFound = errors.New("auth: authenticated user not found")

// ErrAuthenticatedSessionNotFound is returned when no authenticated session is found in context.
var ErrAuthenticatedSessionNotFound = errors.New("auth: authenticated session not found")

// ErrInvalidAuthenticatedUserID is returned when an authenticated user's ID cannot be used as a positive int64.
var ErrInvalidAuthenticatedUserID = errors.New("auth: invalid authenticated user id")
