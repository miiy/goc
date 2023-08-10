package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/miiy/goc/auth/jwt"
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
	claims, err := jwtAuth.ParseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token")
	}

	user, err := arepo.FirstByUsername(ctx, claims.Username)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token")
	}
	authUser := jwt.AuthUser{
		Id:       user.ID,
		Username: claims.Username,
	}

	// WARNING: in production define your own type to avoid context collisions
	newCtx := context.WithValue(ctx, "auth.user", authUser)
	// NOTE: You can also pass the token in the context for further interceptors or gRPC service code.
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
