package server

import (
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	goccredentials "github.com/miiy/goc/grpc/credentials"
	loggerpkg "github.com/miiy/goc/grpc/interceptor/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

func WithMTLS(certFilePath, keyFilePath, caFilePath string) (grpc.ServerOption, error) {
	creds, err := goccredentials.NewServerMTLS(certFilePath, keyFilePath, caFilePath)
	if err != nil {
		return nil, err
	}
	return grpc.Creds(creds), nil
}

func DefaultInterceptor(logger *zap.Logger, authFunc auth.AuthFunc, matcher selector.Matcher) []grpc.ServerOption {
	grpcPanicRecoveryHandler := func(p any) (err error) {
		grpclog.Error("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}

	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(loggerpkg.InterceptorLogger(logger), loggerOpts...),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(loggerpkg.InterceptorLogger(logger), loggerOpts...),
	}
	if authFunc != nil && matcher != nil {
		unaryInterceptors = append(unaryInterceptors,
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFunc), matcher),
		)
		streamInterceptors = append(streamInterceptors,
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFunc), matcher),
		)
	}
	unaryInterceptors = append(unaryInterceptors,
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
	)
	streamInterceptors = append(streamInterceptors,
		recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
	)

	return []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	}
}
