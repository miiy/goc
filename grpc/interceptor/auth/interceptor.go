package auth2

import (
	"context"
	"github.com/miiy/goc/auth"
	"github.com/miiy/goc/auth/jwt"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(jwt *jwt.JWTAuth, userProvider auth.UserProvider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		newCtx := auth.InjectJWTAuth(ctx, jwt)
		newCtx = auth.InjectUserProvider(ctx, userProvider)
		return handler(newCtx, req)
	}
}

func StreamServerInterceptor(jwt *jwt.JWTAuth, userProvider auth.UserProvider) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		newCtx := auth.InjectJWTAuth(ctx, jwt)
		newCtx = auth.InjectUserProvider(ctx, userProvider)

		ws := &wrappedStream{
			ServerStream: ss,
			ctx:          newCtx,
		}
		// Continue execution of handler after ensuring a valid token.
		return handler(srv, ws)
	}
}

// https://github.com/grpc/grpc-go/blob/master/orca/call_metrics.go

// wrappedStream wraps the grpc.ServerStream received by the streaming
// interceptor. Overrides only the Context() method to return a context which
// contains a reference to the CallMetricsRecorder corresponding to this
// stream.
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
