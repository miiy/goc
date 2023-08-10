package grpc_interceptor

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/miiy/goc/auth/jwt"
	"github.com/miiy/goc/auth/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcAuthenticateInterceptor(j *jwt.JWTAuth, p repository.AuthenticateRepository) auth.AuthFunc {
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

	aup, ok := ctx.Value("authUserProvider").(repository.AuthenticateRepository)
	user, err := aup.RetrieveByIdentifier(ctx, "username", claims.Username)

	if err != nil {
		return nil, status.New(codes.Unauthenticated, err.Error()).Err()
	}
	authUser := &jwt.AuthUser{
		Id:       user.ID,
		Username: user.Username,
	}

	// WARNING: in production define your own type to avoid context collisions
	newCtx := context.WithValue(ctx, "auth.user", authUser)

	return newCtx, nil
}
