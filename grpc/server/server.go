package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
)

type Server struct {
	server *grpc.Server
	grpc.ServiceRegistrar
}

func DefaultServerOption() []grpc.ServerOption {
	// Define customfunc to handle panic
	grpcPanicRecoveryHandler := func(p any) (err error) {
		grpclog.Error("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFn), matcher),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFn), matcher),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	}
}

func NewServer(opt ...grpc.ServerOption) (*Server, error) {
	server := grpc.NewServer(opt...)
	return &Server{
		server: server,
	}, nil
}

func (s *Server) Serve(ctx context.Context, address string) error {
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

			<-ctx.Done()
		}
	}()

	grpclog.Infof("server listening at ", lis.Addr().String())
	return s.server.Serve(lis)
}

func (s *Server) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.server.RegisterService(sd, ss)
}
