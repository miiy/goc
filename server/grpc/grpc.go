package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/miiy/goc/auth/jwt"
	authpb "github.com/miiy/goc/service/auth/api/v1"
	"github.com/miiy/goc/service/auth/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
)

type Server interface {
	RegisterService(sd *grpc.ServiceDesc, ss interface{})
	Serve(address string) error
}

type gRPCServer struct {
	context context.Context
	grpc.ServiceRegistrar
	server *grpc.Server
}

type ServerOption = grpc.ServerOption

func NewServer(ctx context.Context, matcher selector.Matcher, jwtAuth *jwt.JWTAuth, arepo repository.AuthRepository) (Server, error) {
	// Setup custom auth.
	authFn := func(ctx context.Context) (context.Context, error) {
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
			Id:       user.Id,
			Username: claims.Username,
		}

		// WARNING: in production define your own type to avoid context collisions
		newCtx := context.WithValue(ctx, "auth.user", authUser)
		// NOTE: You can also pass the token in the context for further interceptors or gRPC service code.
		return newCtx, nil
	}

	// Define customfunc to handle panic
	grpcPanicRecoveryHandler := func(p any) (err error) {
		grpclog.Error("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), matcher),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFn), matcher),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)

	return &gRPCServer{
		context: ctx,
		server:  grpcSrv,
	}, nil
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

func (s *gRPCServer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.server.RegisterService(sd, ss)
}

func (s *gRPCServer) Serve(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			grpclog.Infoln("shutting down gRPC server...")

			s.server.GracefulStop()

			<-s.context.Done()
		}
	}()

	grpclog.Infof("server listening at ", lis.Addr().String())
	return s.server.Serve(lis)
}
