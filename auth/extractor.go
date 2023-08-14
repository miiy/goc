package auth

import (
	"context"
	"errors"
	"github.com/miiy/goc/auth/jwt"
)

const JWTAuthContextKey = "jwt.JWTAuth"
const AuthenticatedUserContextKey = "auth.authenticatedUser"
const UserProviderContextKey = "auth.userProvider"

func InjectJWTAuth(ctx context.Context, jwt *jwt.JWTAuth) context.Context {
	return context.WithValue(ctx, JWTAuthContextKey, jwt)
}

func ExtractJWTAuth(ctx context.Context) (*jwt.JWTAuth, error) {
	j, ok := ctx.Value(JWTAuthContextKey).(*jwt.JWTAuth)
	if !ok {
		return nil, errors.New("extract jwt.JWTAuth error")
	}
	return j, nil
}

func InjectAuthenticatedUser(ctx context.Context, u *AuthenticatedUser) context.Context {
	return context.WithValue(ctx, AuthenticatedUserContextKey, u)
}

func ExtractAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	u, ok := ctx.Value(AuthenticatedUserContextKey).(*AuthenticatedUser)
	if !ok {
		return nil, errors.New("extract auth.AuthenticatedUser error")
	}
	return u, nil
}

func InjectUserProvider(ctx context.Context, up UserProvider) context.Context {
	return context.WithValue(ctx, UserProviderContextKey, up)
}

func ExtractUserProvider(ctx context.Context) (UserProvider, error) {
	u, ok := ctx.Value(UserProviderContextKey).(UserProvider)
	if !ok {
		return nil, errors.New("extract auth.UserProvider error")
	}
	return u, nil
}
