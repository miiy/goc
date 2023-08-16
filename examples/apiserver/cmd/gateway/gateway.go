package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	authpb "github.com/miiy/goc/component/auth/api/v1"
	"github.com/miiy/goc/grpc/gateway"
)

var (
	endpoint   = flag.String("endpoint", "localhost:50051", "endpoint of the gRPC service")
	network    = flag.String("network", "tcp", `one of "tcp" or "unix". Must be consistent to -endpoint`)
	openAPIDir = flag.String("openapi_dir", "examples/internal/proto/examplepb", "path to the directory which contains OpenAPI definitions")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	opts := gateway.Options{
		Addr: ":8080",
		GRPCServer: gateway.Endpoint{
			Network: *network,
			Addr:    *endpoint,
		},
		OpenAPIDir: *openAPIDir,
	}

	handlers := []gateway.GatewayHandler{
		authpb.RegisterAuthHandler,
	}
	if err := gateway.Run(ctx, opts, handlers...); err != nil {
		glog.Fatal(err)
	}
}
