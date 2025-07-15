package jwt

import (
	"testing"
)

const testSecret = "123"

func newTestJwt() *JWTAuth {
	return NewJWTAuth(&Options{
		Secret:    testSecret,
		ExpiresIn: 20,
	})
}

func TestJWTAuth_CreateToken(t *testing.T) {
	jwtAuth := newTestJwt()
	token, err := jwtAuth.CreateToken("test")
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
	c, err := jwtAuth.ParseToken(token)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", c)
}

func TestJWTAuth_ParseToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6InRlc3QiLCJleHAiOjE2ODU5MzU0MzB9.6g7zsZ0P3D85R6YjlizRDHkEyX-wU-eYV81hSeHzPAg"
	jwtAuth := newTestJwt()
	c, err := jwtAuth.ParseToken(token)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", c)
}
