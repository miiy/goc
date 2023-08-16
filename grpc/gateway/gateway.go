package gateway

import (
	"context"
	"github.com/golang/glog"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net/http"
)

// Endpoint describes a gRPC endpoint
type Endpoint struct {
	Network, Addr string
}

// Options is a set of options to be passed to Run
type Options struct {
	// Addr is the address to listen
	Addr string

	// GRPCServer defines an endpoint of a gRPC service
	GRPCServer Endpoint

	// OpenAPIDir is a path to a directory from which the server
	// serves OpenAPI specs.
	OpenAPIDir string

	// Mux is a list of options to be passed to the gRPC-Gateway multiplexer
	Mux []gwruntime.ServeMuxOption
}

type Gateway struct {
	options Options
}

type GatewayHandler func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error

// Run starts a HTTP server and blocks while running if successful.
// The server will be shutdown when "ctx" is canceled.
func Run(ctx context.Context, opts Options, handlers ...GatewayHandler) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := dial(ctx, opts.GRPCServer.Network, opts.GRPCServer.Addr)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			glog.Errorf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/openapiv2/", openAPIServer(opts.OpenAPIDir))
	mux.HandleFunc("/healthz", healthzServer(conn))

	//gw, err := newGateway(ctx, conn, opts.Mux)
	//if err != nil {
	//	return err
	//}

	gwmux := gwruntime.NewServeMux(opts.Mux...)

	// 上传文件
	gwmux.HandlePath("POST", "/v1/file/upload", handleBinaryFileUpload)

	// Register Greeter
	for _, f := range handlers {
		if err = f(ctx, gwmux, conn); err != nil {
			return err
		}
	}

	mux.Handle("/", gwmux)

	s := &http.Server{
		Addr:    opts.Addr,
		Handler: allowCORS(mux),
	}
	go func() {
		<-ctx.Done()
		glog.Infof("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			glog.Errorf("Failed to shutdown http server: %v", err)
		}
	}()

	glog.Infof("Starting listening at %s", opts.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
		return err
	}
	return nil
}
