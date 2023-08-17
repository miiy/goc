package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/miiy/goc/component/auth/api/v1"
	echopb "github.com/miiy/goc/examples/apiserver/api/echo/v1"
	"github.com/miiy/goc/grpc/gateway"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
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
		Mux: []runtime.ServeMuxOption{
			//gwruntime.WithMarshalerOption(gwruntime.MIMEWildcard, &gwruntime.JSONPb{
			//	MarshalOptions: protojson.MarshalOptions{
			//		EmitUnpopulated: true,
			//		UseProtoNames:   true,
			//	},
			//	UnmarshalOptions: protojson.UnmarshalOptions{
			//		DiscardUnknown: true,
			//	},
			//}),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &gateway.CustomMarshaler{
				Marshaler: &runtime.JSONPb{
					MarshalOptions: protojson.MarshalOptions{
						EmitUnpopulated: true,
						UseProtoNames:   true,
					},
					UnmarshalOptions: protojson.UnmarshalOptions{
						DiscardUnknown: true,
					},
				}}),
		},
		RegisterHandler: []gateway.RegisterHandler{
			authpb.RegisterAuthHandler,
			echopb.RegisterEchoHandler,
			gateway.RegisterUploadHandler,
		},
	}

	if err := gateway.Run(ctx, opts); err != nil {
		log.Fatal(err)
	}
}
