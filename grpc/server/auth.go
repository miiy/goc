package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	gauth "github.com/miiy/goc/auth"
	authpb "github.com/miiy/goc/service/auth/api/v1"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Setup custom auth.
func authFn(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	jwtAuth, err := gauth.ExtractJWTAuth(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "extract jwt.JWTAuth error: %v", err)
	}
	claims, err := jwtAuth.Parse(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	subject, err := claims.GetSubject()
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	userProvider, err := gauth.ExtractUserProvider(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "extract auth.userProvider error: %v", err)
	}
	user, err := userProvider.FirstByIdentifier(ctx, subject)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token")
	}

	newCtx := gauth.InjectAuthenticatedUser(ctx, user)
	return newCtx, nil
}

func authMatcher(ctx context.Context, c interceptors.CallMeta) bool {
	// health check
	if healthpb.Health_ServiceDesc.ServiceName == c.Service {
		return false
	}
	// auth service
	var fullMethodNames []string
	for _, v := range []string{"Login", "Register", "UsernameCheck", "EmailCheck", "PhoneCheck"} {
		fullMethodNames = append(fullMethodNames, fmt.Sprintf("/%s/%s", authpb.AuthService_ServiceDesc.ServiceName, v))
	}
	// avatar service

	for _, v := range fullMethodNames {
		if c.FullMethod() == v {
			return false
		}
	}

	return true
}
