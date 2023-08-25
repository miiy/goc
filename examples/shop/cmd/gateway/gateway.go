package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/miiy/goc/component/auth/api/v1"
	postv1 "github.com/miiy/goc/examples/shop/gen/shop/post/v1"
	confpkg "github.com/miiy/goc/examples/shop/internal/config"
	"github.com/miiy/goc/grpc/gateway"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
)

var (
	conf     = flag.String("c", "./configs/default.yaml", "config file")
	endpoint = flag.String("endpoint", "localhost:50051", "endpoint of the gRPC service")
)

func main() {
	flag.Parse()
	// conf
	config, err := confpkg.NewConfig(*conf)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	opts := gateway.Options{
		Addr: config.Server.Http.Addr,
		GRPCServer: gateway.Endpoint{
			Addr: *endpoint,
		},
		OpenAPIDir: "",
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
				},
			}),
		},
		RegisterHandler: []gateway.RegisterHandler{
			authpb.RegisterAuthHandler,
			postv1.RegisterPostServiceHandler,
			gateway.RegisterUploadHandler,
			gateway.RegisterHealthzHandler,
		},
		TlsConfig: gateway.MTLSConfig(
			config.GrpcClient.Tls.ServerName,
			config.GrpcClient.Tls.CertFile,
			config.GrpcClient.Tls.KeyFile,
			config.GrpcClient.Tls.CaFile,
		),
	}

	if err := gateway.Run(ctx, opts); err != nil {
		log.Fatal(err)
	}
}
