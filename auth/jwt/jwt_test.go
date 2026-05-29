package jwt

import (
	"testing"
)

const testSecret = "123"

func newTestJwt() *JWTAuth {
	return NewJWTAuth(&Options{
		Secret:    testSecret,
		Issuer:    "goc",
		ExpiresIn: 20,
	})
}

func TestJWTAuth_CreateToken(t *testing.T) {
	jwtAuth := newTestJwt()
	token, err := jwtAuth.CreateToken("test")
	if err != nil {
		t.Fatal(err)
	}

	c, err := jwtAuth.ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if c.Username != "test" {
		t.Fatalf("unexpected claims: %+v", c)
	}
}
