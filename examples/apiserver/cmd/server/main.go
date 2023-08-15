package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	authpb "github.com/miiy/goc/component/auth/api/v1"
	"github.com/miiy/goc/grpc/server"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	conf := flag.String("c", "./configs/default.yaml", "config file")
	flag.Parse()

	ctx := context.Background()
	app, cleanup, err := initApp(*conf)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// set logger
	logger := app.Logger().ZapLogger()
	grpclog.SetLoggerV2(zapgrpc.NewLogger(logger))

	// set server
	serverOpts := server.DefaultServerOption(
		logger,
		authFunc(app.JWTAuth(), app.UserProvider()),
		selector.MatchFunc(authMatchFunc),
	)
	s, err := server.NewServer(serverOpts...)
	if err != nil {
		grpclog.Fatal("Failed to create server", err)
	}

	// register service
	authpb.RegisterAuthServer(s, app.AuthServer())

	// serve
	if err = s.Serve(ctx, app.Config().Server.Grpc.Addr); err != nil {
		grpclog.Fatal("failed to serve: %v", err)
	}
}
