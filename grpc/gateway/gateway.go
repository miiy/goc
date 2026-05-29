package gateway

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// Endpoint describes a gRPC endpoint
type Endpoint struct {
	Addr string
}

// ServiceConfig defines a gRPC service configuration
type ServiceConfig struct {
	Name     string
	Endpoint Endpoint
	Register RegisterHandler
}

// Options is a set of options to be passed to Run
type Options struct {
	// Addr is the address to listen
	Addr string

	// Services is a list of gRPC services to register
	Services []ServiceConfig

	// APIPrefix is the URL prefix for API routes (default: "/api/")
	APIPrefix string

	// Mux is a list of options to be passed to the gRPC-Gateway multiplexer
	Mux []runtime.ServeMuxOption

	// Middlewares applied to the HTTP handler (e.g. CORS)
	Middlewares []func(http.Handler) http.Handler

	// TLSConfig for mTLS
	TLSConfig *tls.Config
}

// RegisterHandler registers the http handlers to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
type RegisterHandler func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

// Run starts a HTTP server and blocks while running if successful.
// The server will be shutdown when "ctx" is canceled.
func Run(ctx context.Context, opts Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Set default API prefix
	if opts.APIPrefix == "" {
		opts.APIPrefix = "/api/"
	}

	// Default mux options
	if opts.Mux == nil {
		opts.Mux = []runtime.ServeMuxOption{
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
					UseProtoNames:   true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}),
		}
	}

	// Create a gateway mux instance
	gwMux := runtime.NewServeMux(opts.Mux...)

	// Create a http mux instance
	httpMux := http.NewServeMux()

	// TLS credentials
	var creds credentials.TransportCredentials
	if opts.TLSConfig != nil {
		creds = credentials.NewTLS(opts.TLSConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	// Register handlers for each service
	for _, svc := range opts.Services {
		conn, err := grpc.NewClient(svc.Endpoint.Addr, grpc.WithTransportCredentials(creds))
		if err != nil {
			return err
		}
		go func() {
			<-ctx.Done()
			if err := conn.Close(); err != nil {
				log.Printf("Failed to close a client connection to the gRPC server: %v", err)
			}
		}()

		if err := svc.Register(ctx, gwMux, conn); err != nil {
			return err
		}

		log.Printf("registered %s service -> %s", svc.Name, svc.Endpoint.Addr)
	}

	// Mount gateway mux to API prefix. The generated grpc-gateway mux matches
	// the full HTTP path from proto annotations, so keep the API prefix intact.
	httpMux.Handle(opts.APIPrefix, gwMux)

	// Health check endpoint
	httpMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	handler := http.Handler(httpMux)
	for _, mw := range opts.Middlewares {
		handler = mw(handler)
	}

	s := &http.Server{
		Addr:    opts.Addr,
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		log.Printf("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown http server: %v", err)
		}
	}()

	log.Printf("gateway listening on %s", opts.Addr)
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
