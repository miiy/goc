package sessions

import (
	"errors"
	"testing"
)

func TestUseJSONSerializerRejectsUnsupportedStore(t *testing.T) {
	err := UseJSONSerializer(NewCookieStore("secret"))
	if !errors.Is(err, ErrUnsupportedJSONSessionStore) {
		t.Fatalf("expected ErrUnsupportedJSONSessionStore, got %v", err)
	}
}

func TestSetMaxAgeRejectsUnsupportedStore(t *testing.T) {
	err := SetMaxAge(NewCookieStore("secret"), 300)
	if !errors.Is(err, ErrUnsupportedRedisSessionStore) {
		t.Fatalf("expected ErrUnsupportedRedisSessionStore, got %v", err)
	}
}
