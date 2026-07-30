package auth

import (
	"context"
	"strings"

	gauth "github.com/miiy/goc/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// MetadataAuthFunc extracts authenticated user info from incoming gRPC metadata
// and creates a user-only authenticated session. Session metadata stays at the
// HTTP edge until downstream services need it.
func MetadataAuthFunc(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	return injectAuthenticatedSessionFromUserMetadata(ctx, md)
}

func injectAuthenticatedSessionFromUserMetadata(ctx context.Context, md metadata.MD) (context.Context, error) {
	userID, ok := singleMetadataValue(md, gauth.AuthenticatedUserIDMetadataKey)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing authenticated user id")
	}

	username, ok := singleMetadataValue(md, gauth.AuthenticatedUsernameMetadataKey)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing authenticated username")
	}

	return gauth.InjectAuthenticatedSession(ctx, &gauth.AuthenticatedSession{
		User: gauth.AuthenticatedUser{ID: userID, Username: username},
	}), nil
}

func singleMetadataValue(md metadata.MD, key string) (string, bool) {
	values := md.Get(key)
	if len(values) != 1 {
		return "", false
	}
	value := strings.TrimSpace(values[0])
	return value, value != ""
}
