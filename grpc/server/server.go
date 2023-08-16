package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
	"os"
	"os/signal"
)

type Server struct {
	server *grpc.Server
	grpc.ServiceRegistrar
}

func NewServer(opt ...grpc.ServerOption) (*Server, error) {
	server := grpc.NewServer(opt...)
	return &Server{
		server: server,
	}, nil
}

func (s *Server) Serve(ctx context.Context, network, address string) error {
	lis, err := net.Listen(network, address)
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
