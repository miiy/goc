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
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

	config := app.Config()

	// set logger
	logger := app.Logger().ZapLogger()
	grpclog.SetLoggerV2(zapgrpc.NewLogger(logger))

	// grpc server options
	var serverOpts []grpc.ServerOption
	// mTLS
	serverOpts = append(serverOpts,
		server.WithMTLS(
			config.Server.Grpc.Tls.CertFile,
			config.Server.Grpc.Tls.KeyFile,
			config.Server.Grpc.Tls.CaFile,
		),
	)
	// interceptor
	serverOpts = append(serverOpts, server.DefaultInterceptor(
		logger,
		authFunc(app.JWTAuth(), app.UserProvider()),
		selector.MatchFunc(authMatchFunc),
	)...)

	// run server
	err = server.Run(ctx, server.Options{
		Network:      "tcp",
		Addr:         app.Config().Server.Grpc.Addr,
		ServerOption: serverOpts,
		RegisterService: func(s server.ServiceRegistrar) {
			healthpb.RegisterHealthServer(s, health.NewServer())
			authpb.RegisterAuthServer(s, app.AuthServer())
			echopb.RegisterEchoServer(s, echoSrv.NewEchoServiceServer())
		},
	})
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
