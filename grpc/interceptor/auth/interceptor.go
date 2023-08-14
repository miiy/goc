package auth

import (
	"context"
	"errors"
	"github.com/miiy/goc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// https://github.com/grpc/grpc-go/blob/master/examples/features/authentication/server/main.go
// https://github.com/grpc/grpc-go/blob/master/authz/grpc_authz_server_interceptors.go
var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

// valid validates the authorization.
func valid(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	authorization := md["authorization"]
	if len(authorization) < 1 {
		return nil, errors.New("authorization token is not supplied")
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.

	jwtAuth, err := auth.ExtractJWTAuth(ctx)
	if err != nil {
		return nil, err
	}
	claims, err := jwtAuth.Parse(token)
	if err != nil {
		return nil, err
	}
	subject, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}

	userProvider, err := auth.ExtractUserProvider(ctx)
	user, err := userProvider.FirstByIdentifier(ctx, subject)
	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, "auth.user", user)
	return newCtx, nil
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		// The keys within metadata.MD are normalized to lowercase.
		// See: https://godoc.org/google.golang.org/grpc/metadata#New
		newCtx, err := valid(ctx)
		if err != nil {
			return nil, errInvalidToken
		}

		// Continue execution of handler after ensuring a valid token.
		return handler(newCtx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// The keys within metadata.MD are normalized to lowercase.
		// See: https://godoc.org/google.golang.org/grpc/metadata#New
		newCtx, err := valid(ss.Context())
		if err != nil {
			return errInvalidToken
		}
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
