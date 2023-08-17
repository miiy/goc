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

type Options struct {
	Network, Addr    string
	ServerOption     []grpc.ServerOption
	ServiceRegistrar func(s grpc.ServiceRegistrar)
}

type ServiceRegistrar = grpc.ServiceRegistrar

func Run(ctx context.Context, opts Options) error {
	server := grpc.NewServer(opts.ServerOption...)
	lis, err := net.Listen(opts.Network, opts.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}

	// register service
	opts.ServiceRegistrar(server)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			grpclog.Infoln("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	grpclog.Infof("server listening at ", lis.Addr().String())
	return server.Serve(lis)
}
