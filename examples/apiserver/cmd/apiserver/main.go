package main

import (
	"context"
	"flag"
	iapp "github.com/miiy/goc/examples/apiserver/app"
	"github.com/miiy/goc/server/grpc"
	pb "github.com/miiy/goc/service/auth/api/v1"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	conf := flag.String("c", "./configs/default.yaml", "config file")
	flag.Parse()

	ctx := context.Background()
	app, cleanup, err := iapp.InitApp(*conf)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	grpclog.SetLoggerV2(zapgrpc.NewLogger(app.Logger.ZapLogger()))

	s, err := grpc.NewServer(ctx)
	if err != nil {
		grpclog.Fatal("Failed to create server", err)
	}

	pb.RegisterAuthServiceServer(s, app.AuthServer)

	if err = s.Serve("0.0.0.0:50051"); err != nil {
		grpclog.Fatal("failed to serve: %v", err)
	}
}
