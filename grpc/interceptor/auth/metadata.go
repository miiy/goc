package auth

import (
	"context"
	"strconv"
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
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authenticated user id: %v", err)
	}
	if id <= 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid authenticated user id")
	}

	username, ok := metadataValue(md, gauth.AuthenticatedUsernameMetadataKey)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing authenticated username")
	}

	return gauth.InjectAuthenticatedUser(ctx, &gauth.AuthenticatedUser{
		ID:       id,
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
