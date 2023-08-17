package gateway

import (
	"context"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	Mux []runtime.ServeMuxOption

	// RegisterHandler registers the http handlers to "mux".
	RegisterHandler []RegisterHandler
}

// RegisterHandler registers the http handlers to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
type RegisterHandler func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

// Run starts a HTTP server and blocks while running if successful.
// The server will be shutdown when "ctx" is canceled.
func Run(ctx context.Context, opts Options) error {
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

	// Create a mux instance
	mux := runtime.NewServeMux(opts.Mux...)

	// Register Handlers
	for _, f := range opts.RegisterHandler {
		if err = f(ctx, mux, conn); err != nil {
			return err
		}
	}

	s := &http.Server{
		Addr:    opts.Addr,
		Handler: mux,
	}
	go func() {
		<-ctx.Done()
		glog.Infof("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			glog.Errorf("Failed to shutdown http server: %v", err)
		}
	}()

	glog.Infof("Starting listening at %s", opts.Addr)
	return s.ListenAndServe()
}
