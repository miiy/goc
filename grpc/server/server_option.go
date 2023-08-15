package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	loggerpkg "github.com/miiy/goc/grpc/interceptor/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

func DefaultServerOption(logger *zap.Logger, authFunc auth.AuthFunc, matcher selector.Matcher) []grpc.ServerOption {
	// Define customfunc to handle panic
	grpcPanicRecoveryHandler := func(p any) (err error) {
		grpclog.Error("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}

	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			// Order matters e.g. tracing interceptor have to create span first for the later exemplars to work.
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFunc), matcher),
			logging.UnaryServerInterceptor(loggerpkg.InterceptorLogger(logger), loggerOpts...),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFunc), matcher),
			logging.StreamServerInterceptor(loggerpkg.InterceptorLogger(logger), loggerOpts...),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	}
}
