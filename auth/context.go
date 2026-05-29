package auth

import (
	"context"
	"errors"
)

var ErrAuthenticatedUserNotFound = errors.New("auth: authenticated user not found")

type authenticatedUserContextKey string

const authenticatedUserContextKeyValue authenticatedUserContextKey = "goc.auth.AuthenticatedUser"

func InjectAuthenticatedUser(ctx context.Context, u *AuthenticatedUser) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if u == nil {
		return ctx
	}
	return context.WithValue(ctx, authenticatedUserContextKeyValue, u)
}

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
