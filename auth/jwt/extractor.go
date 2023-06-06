package jwt

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoTokenInRequest = status.New(codes.Unauthenticated, "Unauthenticated").Err()
)

func AuthUserFromContext(ctx context.Context) (*AuthUser, error) {
	user, ok := ctx.Value("auth.user").(*AuthUser)
	if !ok {
		return nil, ErrNoTokenInRequest
	}
	return user, nil
}
