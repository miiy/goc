package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
	gocjwt "github.com/miiy/goc/auth/jwt"
)

func newTestJWTAuth() *gocjwt.JWTAuth {
	return gocjwt.NewJWTAuth(&gocjwt.Options{
		Secret:    "secret",
		Issuer:    "goc",
		ExpiresIn: 60,
	})
}

func TestJWTAuthenticationMiddlewareSetsAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtAuth := newTestJWTAuth()
	token, err := jwtAuth.CreateToken("user")
	if err != nil {
		t.Fatalf("expected token, got error: %v", err)
	}

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	JWTAuthenticationMiddleware(jwtAuth)(c)

	auth, exists := c.Get(AuthUserKey)
	if !exists {
		t.Fatal("expected auth in context")
	}

	user, ok := auth.(*gocauth.AuthenticatedUser)
	if !ok {
		t.Fatalf("expected *AuthenticatedUser, got %T", auth)
	}
	if user.Username != "user" {
		t.Fatalf("expected username user, got %q", user.Username)
	}

	ctxUser, err := gocauth.ExtractAuthenticatedUser(c.Request.Context())
	if err != nil {
		t.Fatalf("expected auth user in request context: %v", err)
	}
	if ctxUser != user {
		t.Fatalf("expected request context user, got %v", ctxUser)
	}
}

func TestJWTAuthenticationMiddlewareRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)

	JWTAuthenticationMiddleware(newTestJWTAuth())(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if c.Writer.Status() != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", c.Writer.Status())
	}
}

func TestJWTAuthenticationMiddlewareRejectsInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid")

	JWTAuthenticationMiddleware(newTestJWTAuth())(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if c.Writer.Status() != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", c.Writer.Status())
	}
}
