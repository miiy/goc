package auth

import (
	"context"
	"testing"
)

func TestInjectAndExtractAuthenticatedUser(t *testing.T) {
	ctx := InjectAuthenticatedUser(context.Background(), &AuthenticatedUser{
		ID:       1,
		Username: "test",
	})

	user, err := ExtractAuthenticatedUser(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "test" || user.ID != 1 {
		t.Fatalf("unexpected user: %+v", user)
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
