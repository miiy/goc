package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	authpb "github.com/miiy/goc/component/auth/api/v1"
	echopb "github.com/miiy/goc/examples/apiserver/api/echo/v1"
	echoSrv "github.com/miiy/goc/examples/apiserver/echo"
	"github.com/miiy/goc/grpc/server"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
	"log"
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

	// server options
	serverOpts := server.DefaultServerOption(
		logger,
		authFunc(app.JWTAuth(), app.UserProvider()),
		selector.MatchFunc(authMatchFunc),
	)

	// run server
	err = server.Run(ctx, server.Options{
		Network:      "tcp",
		Addr:         app.Config().Server.Grpc.Addr,
		ServerOption: serverOpts,
		ServiceRegistrar: func(s server.ServiceRegistrar) {
			authpb.RegisterAuthServer(s, app.AuthServer())
			echopb.RegisterEchoServer(s, echoSrv.NewEchoServiceServer())
		},
	})
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
