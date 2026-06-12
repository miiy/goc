package auth

import (
	"testing"
)

func newTestJWTAuth() *JWTAuth {
	return NewJWTAuth(&Options{
		Secret:    "secret",
		Issuer:    "goc",
		ExpiresIn: 60,
	})
}

func TestCreateAndParseToken(t *testing.T) {
	jwtAuth := newTestJWTAuth()

	token, err := jwtAuth.CreateToken("user")
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	claims, err := jwtAuth.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.Username != "user" {
		t.Fatalf("expected username user, got %q", claims.Username)
	}
}

func TestParseInvalidToken(t *testing.T) {
	_, err := newTestJWTAuth().ParseToken("invalid")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestParseTokenWithWrongSecret(t *testing.T) {
	auth1 := NewJWTAuth(&Options{Secret: "secret1", Issuer: "goc"})
	auth2 := NewJWTAuth(&Options{Secret: "secret2", Issuer: "goc"})

	token, _ := auth1.CreateToken("user")
	_, err := auth2.ParseToken(token)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestParseTokenWithWrongIssuer(t *testing.T) {
	auth1 := NewJWTAuth(&Options{Secret: "secret", Issuer: "other"})
	auth2 := NewJWTAuth(&Options{Secret: "secret", Issuer: "goc"})

	token, _ := auth1.CreateToken("user")
	_, err := auth2.ParseToken(token)
	if err == nil {
		t.Fatal("expected error for wrong issuer")
	}
}
