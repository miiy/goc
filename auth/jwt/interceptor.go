package jwt

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserProvider interface {
	RetrieveByUsername(ctx context.Context, username string) (*AuthUser, error)
}

func (j *JWTAuth) GrpcAuthenticateInterceptor(p UserProvider) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		ctx = context.WithValue(ctx, "jwtAuth", j)
		ctx = context.WithValue(ctx, "authUserProvider", p)
		return GrpcAuthFunc(ctx)
	}
}

func GrpcAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	jwtAuth, ok := ctx.Value("jwtAuth").(*JWTAuth)
	if !ok {
		return nil, status.New(codes.Internal, "jwtAuth from context error").Err()
	}

	claims, err := jwtAuth.ParseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", claims)

	aup, ok := ctx.Value("authUserProvider").(UserProvider)
	user, err := aup.RetrieveByUsername(ctx, claims.Username)

	if err != nil {
		return nil, status.New(codes.Unauthenticated, err.Error()).Err()
	}
	authUser := &AuthUser{
		Id:       user.Id,
		Username: user.Username,
	}

	// WARNING: in production define your own type to avoid context collisions
	newCtx := context.WithValue(ctx, "auth.user", authUser)

	return newCtx, nil
}
