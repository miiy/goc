package server

import (
	"context"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	middlewareauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	goccredentials "github.com/miiy/goc/grpc/credentials"
	loggerpkg "github.com/miiy/goc/grpc/interceptor/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	reflectionv1 "google.golang.org/grpc/reflection/grpc_reflection_v1"
	reflectionv1alpha "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

func WithMTLS(certFilePath, keyFilePath, caFilePath string) (grpc.ServerOption, error) {
	creds, err := goccredentials.NewServerMTLS(certFilePath, keyFilePath, caFilePath)
	if err != nil {
		return nil, err
	}
	return grpc.Creds(creds), nil
}

// AuthConfig configures metadata auth for non-anonymous RPCs.
type AuthConfig struct {
	authFunc    middlewareauth.AuthFunc
	requireAuth selector.Matcher
}

// WithAuth installs authFunc for RPCs matched by requireAuth.
func WithAuth(authFunc middlewareauth.AuthFunc, requireAuth selector.Matcher) *AuthConfig {
	return &AuthConfig{
		authFunc:    authFunc,
		requireAuth: requireAuth,
	}
}

// RequireAuthExcept matches every application RPC except anonymousMethods.
// Health and reflection RPCs are always anonymous.
func RequireAuthExcept(anonymousMethods ...string) selector.Matcher {
	anonymousMethodSet := make(map[string]struct{}, len(operationalAnonymousRPCs())+len(anonymousMethods))
	for _, method := range operationalAnonymousRPCs() {
		anonymousMethodSet[method] = struct{}{}
	}
	for _, method := range anonymousMethods {
		if method == "" {
			continue
		}
		anonymousMethodSet[method] = struct{}{}
	}
	return selector.MatchFunc(func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		_, ok := anonymousMethodSet[callMeta.FullMethod()]
		return !ok
	})
}

func operationalAnonymousRPCs() []string {
	return []string{
		healthpb.Health_Check_FullMethodName,
		healthpb.Health_List_FullMethodName,
		healthpb.Health_Watch_FullMethodName,
		reflectionv1.ServerReflection_ServerReflectionInfo_FullMethodName,
		reflectionv1alpha.ServerReflection_ServerReflectionInfo_FullMethodName,
	}
}

// DefaultInterceptor installs logging and recovery for every RPC. When authConfig
// is provided, auth runs only for non-anonymous RPCs selected by WithAuth.
func DefaultInterceptor(logger *zap.Logger, authConfig *AuthConfig) []grpc.ServerOption {
	var authFunc middlewareauth.AuthFunc
	var requireAuth selector.Matcher
	if authConfig != nil {
		authFunc = authConfig.authFunc
		requireAuth = authConfig.requireAuth
	}

	grpcPanicRecoveryHandler := func(p any) (err error) {
		grpclog.Error("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}

	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	// logger
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(loggerpkg.InterceptorLogger(logger), loggerOpts...),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(loggerpkg.InterceptorLogger(logger), loggerOpts...),
	}

	// auth
	if authFunc != nil && requireAuth != nil {
		unaryInterceptors = append(unaryInterceptors,
			selector.UnaryServerInterceptor(middlewareauth.UnaryServerInterceptor(authFunc), requireAuth),
		)
		streamInterceptors = append(streamInterceptors,
			selector.StreamServerInterceptor(middlewareauth.StreamServerInterceptor(authFunc), requireAuth),
		)
	}

	// recovery
	unaryInterceptors = append(unaryInterceptors,
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
	)
	streamInterceptors = append(streamInterceptors,
		recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
	)

	// server options
	return []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}
}
