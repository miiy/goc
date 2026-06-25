package auth

import (
	"context"
	"strings"

	gauth "github.com/miiy/goc/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// MetadataAuthFunc extracts authenticated user info from incoming gRPC metadata.
func MetadataAuthFunc(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	userID, ok := metadataValue(md, gauth.AuthenticatedUserIDMetadataKey)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing authenticated user id")
	}

	username, ok := metadataValue(md, gauth.AuthenticatedUsernameMetadataKey)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing authenticated username")
	}

	return gauth.InjectAuthenticatedUser(ctx, &gauth.AuthenticatedUser{
		ID:       userID,
		Username: username,
	}), nil
}

func metadataValue(md metadata.MD, key string) (string, bool) {
	values := md.Get(key)
	if len(values) == 0 {
		return "", false
	}
	value := strings.TrimSpace(values[0])
	return value, value != ""
}
