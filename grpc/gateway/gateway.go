package gateway

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
)

// Endpoint describes a gRPC endpoint
type Endpoint struct {
	Addr string
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

	TlsConfig *tls.Config
}

// RegisterHandler registers the http handlers to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
type RegisterHandler func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

// Run starts a HTTP server and blocks while running if successful.
// The server will be shutdown when "ctx" is canceled.
func Run(ctx context.Context, opts Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// tls
	var creds credentials.TransportCredentials
	if opts.TlsConfig != nil {
		creds = credentials.NewTLS(opts.TlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.DialContext(ctx, opts.GRPCServer.Addr, grpc.WithTransportCredentials(creds))
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
	mux := http.NewServeMux()

	// Create a gateway mux instance
	gwMux := runtime.NewServeMux(opts.Mux...)

	//mux.Handle("/", gwMux)

	// Register Handlers
	for _, f := range opts.RegisterHandler {
		if err = f(ctx, gwMux, conn); err != nil {
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

func MTLSConfig(serverName, certFilePath, keyFilePath, caFilePath string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		log.Fatalf("failed to load client cert: %v", err)
	}

	ca := x509.NewCertPool()
	caBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Fatalf("failed to read ca cert %q: %v", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		log.Fatalf("failed to parse %q", caFilePath)
	}

	return &tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}
}
