package grpc_interceptor

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/miiy/goc/auth/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserProvider interface {
	RetrieveByUsername(ctx context.Context, username string) (*jwt.AuthUser, error)
}

func GrpcAuthenticateInterceptor(j *jwt.JWTAuth, p UserProvider) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		ctx = context.WithValue(ctx, "jwtAuth", j)
		ctx = context.WithValue(ctx, "authUserProvider", p)
		return GrpcAuthFunc(ctx)
	}
}

func GrpcAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	jwtAuth, ok := ctx.Value("jwtAuth").(*jwt.JWTAuth)
	if !ok {
		return nil, status.New(codes.Internal, "jwtAuth from context error").Err()
	}

	claims, err := jwtAuth.ParseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	//grpc_ctxtags.Extract(ctx).Set("auth.sub", claims)

	aup, ok := ctx.Value("authUserProvider").(UserProvider)
	user, err := aup.RetrieveByUsername(ctx, claims.Username)

	if err != nil {
		return nil, status.New(codes.Unauthenticated, err.Error()).Err()
	}
	authUser := &jwt.AuthUser{
		Id:       user.Id,
		Username: user.Username,
	}

	// WARNING: in production define your own type to avoid context collisions
	newCtx := context.WithValue(ctx, "auth.user", authUser)

	return newCtx, nil
}
