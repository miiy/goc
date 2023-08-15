package auth

import (
	"context"
	"errors"
)

const AuthenticatedUserContextKey = "goc.auth.AuthenticatedUser"

func InjectAuthenticatedUser(ctx context.Context, u *AuthenticatedUser) context.Context {
	return context.WithValue(ctx, AuthenticatedUserContextKey, u)
}

func ExtractAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	u, ok := ctx.Value(AuthenticatedUserContextKey).(*AuthenticatedUser)
	if !ok {
		return nil, errors.New("extract goc.auth.AuthenticatedUser error")
	}
	return u, nil
}
