package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/miiy/goc/logger/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type Options struct {
	Network, Addr   string
	Logger          *zap.Logger
	ServerOption    []grpc.ServerOption
	RegisterService func(s GRPCServer)
}

type GRPCServer interface {
	grpc.ServiceRegistrar
	GetServiceInfo() map[string]grpc.ServiceInfo
}

func Run(ctx context.Context, opts Options) error {
	if opts.Logger != nil {
		grpclog.SetLoggerV2(zapgrpc.NewLogger(opts.Logger))
	}

	server := grpc.NewServer(opts.ServerOption...)
	lis, err := net.Listen(opts.Network, opts.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// register service
	opts.RegisterService(server)

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

	grpclog.Infof("server listening at %s", lis.Addr().String())
	return server.Serve(lis)
}
