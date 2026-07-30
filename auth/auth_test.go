package auth

import (
	"context"
	"testing"
)

func TestInjectAndExtractAuthenticatedSession(t *testing.T) {
	ctx := InjectAuthenticatedSession(context.Background(), &AuthenticatedSession{
		User:       AuthenticatedUser{ID: "1", Username: "test"},
		SIDKey:     "session-key",
		ClientType: "app",
		ClientID:   "ios",
	})

	session, err := ExtractAuthenticatedSession(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if session.User.Username != "test" || session.User.ID != "1" || session.SIDKey != "session-key" ||
		session.ClientType != "app" || session.ClientID != "ios" {
		t.Fatalf("unexpected session: %+v", session)
	}
}

func TestUserID(t *testing.T) {
	ctx := InjectAuthenticatedSession(context.Background(), &AuthenticatedSession{
		User: AuthenticatedUser{ID: "42", Username: "alice"},
	})

	id, ok := UserID(ctx)
	if !ok {
		t.Fatal("expected user id")
	}
	if id != "42" {
		t.Fatalf("expected 42, got %q", id)
	}
}

func TestUserIDRejectsMissingUser(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{name: "missing session", ctx: context.Background()},
		{name: "empty id", ctx: InjectAuthenticatedSession(context.Background(), &AuthenticatedSession{User: AuthenticatedUser{Username: "alice"}})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, ok := UserID(tt.ctx)
			if ok {
				t.Fatal("expected no user id")
			}
			if id != "" {
				t.Fatalf("expected empty id, got %q", id)
			}
		})
	}
}

func TestUserInt64ID(t *testing.T) {
	ctx := InjectAuthenticatedSession(context.Background(), &AuthenticatedSession{
		User: AuthenticatedUser{ID: "42", Username: "alice"},
	})

	id, ok := UserInt64ID(ctx)
	if !ok {
		t.Fatal("expected user id")
	}
	if id != 42 {
		t.Fatalf("expected 42, got %d", id)
	}
}

func TestUserInt64IDRejectsInvalidID(t *testing.T) {
	ctx := InjectAuthenticatedSession(context.Background(), &AuthenticatedSession{User: AuthenticatedUser{ID: "abc"}})
	if id, ok := UserInt64ID(ctx); ok || id != 0 {
		t.Fatalf("expected no user id, got %d", id)
	}
}

func TestAuthenticatedUserInt64ID(t *testing.T) {
	user := &AuthenticatedUser{ID: "42"}
	id, err := user.Int64ID()
	if err != nil {
		t.Fatal(err)
	}
	if id != 42 {
		t.Fatalf("expected 42, got %d", id)
	}
}

func TestAuthenticatedUserInt64IDRejectsInvalidID(t *testing.T) {
	tests := []struct {
		name string
		user *AuthenticatedUser
	}{
		{name: "nil user", user: nil},
		{name: "empty id", user: &AuthenticatedUser{}},
		{name: "zero", user: &AuthenticatedUser{ID: "0"}},
		{name: "negative", user: &AuthenticatedUser{ID: "-1"}},
		{name: "not number", user: &AuthenticatedUser{ID: "abc"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if id, err := tt.user.Int64ID(); err == nil {
				t.Fatalf("expected error, got id %d", id)
			}
		})
	}
}

func TestExtractAuthenticatedSessionMissing(t *testing.T) {
	_, err := ExtractAuthenticatedSession(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInjectNilSession(t *testing.T) {
	ctx := InjectAuthenticatedSession(context.Background(), nil)
	_, err := ExtractAuthenticatedSession(ctx)
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestInjectNilContextPanics(t *testing.T) {
	requirePanic(t, func() {
		InjectAuthenticatedSession(nil, &AuthenticatedSession{User: AuthenticatedUser{Username: "test"}})
	})
}

func TestExtractNilContextPanics(t *testing.T) {
	requirePanic(t, func() {
		ExtractAuthenticatedSession(nil)
	})
}

func requirePanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}
