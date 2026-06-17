package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
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

	token, err := jwtAuth.CreateToken(42, "user")
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	claims, err := jwtAuth.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("expected user id 42, got %d", claims.UserID)
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

	token, _ := auth1.CreateToken(1, "user")
	_, err := auth2.ParseToken(token)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestParseTokenWithWrongIssuer(t *testing.T) {
	auth1 := NewJWTAuth(&Options{Secret: "secret", Issuer: "other"})
	auth2 := NewJWTAuth(&Options{Secret: "secret", Issuer: "goc"})

	token, _ := auth1.CreateToken(1, "user")
	_, err := auth2.ParseToken(token)
	if err == nil {
		t.Fatal("expected error for wrong issuer")
	}
}

func TestParseTokenRejectsNonHS256(t *testing.T) {
	jwtAuth := newTestJWTAuth()

	// Forge a token signed with HS384 using the same secret.
	forge := jwt.NewWithClaims(jwt.SigningMethodHS384, jwtAuth.CreateClaims(1, "user"))
	signed, err := forge.SignedString([]byte(jwtAuth.options.Secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	if _, err := jwtAuth.ParseToken(signed); err == nil {
		t.Fatal("expected error for HS384 token; ParseToken must only accept HS256")
	}
}
