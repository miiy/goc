package auth

import (
	"context"
	"errors"
	"testing"

	gauth "github.com/miiy/goc/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestMetadataAuthFunc(t *testing.T) {
	tests := []struct {
		name     string
		metadata metadata.MD
		wantUser *gauth.AuthenticatedUser
		wantCode codes.Code
	}{
		{
			name: "user id and username",
			metadata: metadata.Pairs(
				gauth.AuthenticatedUserIDMetadataKey, "42",
				gauth.AuthenticatedUsernameMetadataKey, "alice",
			),
			wantUser: &gauth.AuthenticatedUser{ID: "42", Username: "alice"},
		},
		{
			name: "missing user id",
			metadata: metadata.Pairs(
				gauth.AuthenticatedUsernameMetadataKey, "alice",
			),
			wantCode: codes.Unauthenticated,
		},
		{
			name: "missing username",
			metadata: metadata.Pairs(
				gauth.AuthenticatedUserIDMetadataKey, "42",
			),
			wantCode: codes.Unauthenticated,
		},
		{
			name: "repeated user id",
			metadata: metadata.Pairs(
				gauth.AuthenticatedUserIDMetadataKey, "42",
				gauth.AuthenticatedUserIDMetadataKey, "43",
			),
			wantCode: codes.Unauthenticated,
		},
		{
			name: "repeated username",
			metadata: metadata.Pairs(
				gauth.AuthenticatedUserIDMetadataKey, "42",
				gauth.AuthenticatedUsernameMetadataKey, "alice",
				gauth.AuthenticatedUsernameMetadataKey, "mallory",
			),
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), tt.metadata)

			ctx, err := MetadataAuthFunc(ctx)
			if tt.wantCode != codes.OK {
				if status.Code(err) != tt.wantCode {
					t.Fatalf("expected code %v, got %v", tt.wantCode, status.Code(err))
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			session, err := gauth.ExtractAuthenticatedSession(ctx)
			if err != nil {
				t.Fatalf("expected authenticated session, got %v", err)
			}
			user := &session.User
			if user.ID != tt.wantUser.ID || user.Username != tt.wantUser.Username {
				t.Fatalf("unexpected authenticated user: %+v", user)
			}
			if session.SIDKey != "" || session.ClientType != "" || session.ClientID != "" {
				t.Fatalf("gRPC unexpectedly propagated session metadata: %+v", session)
			}
		})
	}
}

func TestMetadataAuthFuncAcceptsOpaqueUserID(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		gauth.AuthenticatedUserIDMetadataKey, "user-abc",
		gauth.AuthenticatedUsernameMetadataKey, "alice",
	))

	ctx, err := MetadataAuthFunc(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	session, err := gauth.ExtractAuthenticatedSession(ctx)
	if err != nil {
		t.Fatalf("expected authenticated session, got %v", err)
	}
	user := &session.User
	if user.ID != "user-abc" {
		t.Fatalf("expected opaque user id, got %q", user.ID)
	}
}

func TestMetadataAuthFuncMissingMetadata(t *testing.T) {
	ctx, err := MetadataAuthFunc(context.Background())
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected unauthenticated, got %v", status.Code(err))
	}
	if ctx != nil {
		t.Fatalf("expected nil context, got %v", ctx)
	}
	if errors.Is(err, gauth.ErrAuthenticatedSessionNotFound) {
		t.Fatal("expected gRPC status error")
	}
}
