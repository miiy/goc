package auth

import (
	"context"
	"testing"
)

func TestInjectAndExtractAuthenticatedUser(t *testing.T) {
	ctx := InjectAuthenticatedUser(context.Background(), &AuthenticatedUser{
		ID:       "1",
		Username: "test",
	})

	user, err := ExtractAuthenticatedUser(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "test" || user.ID != "1" {
		t.Fatalf("unexpected user: %+v", user)
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

func TestExtractAuthenticatedUserMissing(t *testing.T) {
	_, err := ExtractAuthenticatedUser(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInjectNilUser(t *testing.T) {
	ctx := InjectAuthenticatedUser(context.Background(), nil)
	_, err := ExtractAuthenticatedUser(ctx)
	if err == nil {
		t.Fatal("expected error for nil user")
	}
}

func TestInjectNilContext(t *testing.T) {
	ctx := InjectAuthenticatedUser(nil, &AuthenticatedUser{Username: "test"})
	user, err := ExtractAuthenticatedUser(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "test" {
		t.Fatalf("expected username test, got %q", user.Username)
	}
}

func TestExtractNilContext(t *testing.T) {
	_, err := ExtractAuthenticatedUser(nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
