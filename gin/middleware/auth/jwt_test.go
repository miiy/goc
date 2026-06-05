package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/miiy/goc/auth"
)

func newTestJWTAuth() *auth.JWTAuth {
	return auth.NewJWTAuth(&auth.Options{
		Secret:    "secret",
		Issuer:    "goc",
		ExpiresIn: 60,
	})
}

func TestJWTAuthenticationMiddlewareSetsUser(t *testing.T) {
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

	u, ok := GetAuthUser(c)
	if !ok {
		t.Fatal("expected auth user in request context")
	}
	if u.Username != "user" {
		t.Fatalf("expected username user, got %q", u.Username)
	}

	ctxUser, err := auth.ExtractAuthenticatedUser(c.Request.Context())
	if err != nil {
		t.Fatalf("expected auth user in request context: %v", err)
	}
	if ctxUser.Username != "user" {
		t.Fatalf("expected username user, got %q", ctxUser.Username)
	}

}

func TestJWTAuthenticationMiddlewareWithUserResolverSetsResolvedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtAuth := newTestJWTAuth()
	token, err := jwtAuth.CreateToken("alice")
	if err != nil {
		t.Fatalf("expected token, got error: %v", err)
	}

	var gotUsername string
	var gotToken string
	var calledCount int
	resolver := func(_ context.Context, claims *auth.UserClaims, tokenString string) (*auth.AuthenticatedUser, error) {
		gotUsername = claims.Username
		gotToken = tokenString
		calledCount++
		return &auth.AuthenticatedUser{
			ID:       42,
			Username: "alice",
		}, nil
	}

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	JWTAuthenticationMiddleware(jwtAuth, WithUserResolver(resolver))(c)

	if calledCount != 1 {
		t.Fatalf("expected resolver to be called once, got %d", calledCount)
	}
	if gotUsername != "alice" {
		t.Fatalf("expected resolver username alice, got %q", gotUsername)
	}
	if gotToken != token {
		t.Fatalf("expected resolver token %q, got %q", token, gotToken)
	}

	u, ok := GetAuthUser(c)
	if !ok {
		t.Fatal("expected auth user in request context")
	}
	if u.ID != 42 || u.Username != "alice" {
		t.Fatalf("unexpected auth user: %+v", u)
	}
}

func TestJWTAuthenticationMiddlewareWithUserResolverRejectsMissingUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtAuth := newTestJWTAuth()
	token, err := jwtAuth.CreateToken("disabled")
	if err != nil {
		t.Fatalf("expected token, got error: %v", err)
	}

	var gotMessage string
	customHandler := func(ctx *gin.Context, message string) {
		gotMessage = message
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": message})
	}
	resolver := func(context.Context, *auth.UserClaims, string) (*auth.AuthenticatedUser, error) {
		return nil, errors.New("user disabled")
	}

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	JWTAuthenticationMiddleware(jwtAuth, WithUserResolver(resolver), WithUnauthorized(customHandler))(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if c.Writer.Status() != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", c.Writer.Status())
	}
	if gotMessage != "user disabled" {
		t.Fatalf("expected error message user disabled, got %q", gotMessage)
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

func TestJWTAuthenticationMiddlewareWithUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var gotMessage string
	customHandler := func(ctx *gin.Context, message string) {
		gotMessage = message
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": message})
	}

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)

	JWTAuthenticationMiddleware(newTestJWTAuth(), WithUnauthorized(customHandler))(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if gotMessage != request.ErrNoTokenInRequest.Error() {
		t.Fatalf("expected error message %q, got %q", request.ErrNoTokenInRequest.Error(), gotMessage)
	}
}

func TestJWTAuthenticationMiddlewareWithAfterAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtAuth := newTestJWTAuth()
	token, err := jwtAuth.CreateToken("alice")
	if err != nil {
		t.Fatalf("expected token, got error: %v", err)
	}

	var gotClaims *auth.UserClaims
	var gotToken string
	hook := func(ctx *gin.Context, claims *auth.UserClaims, tokenString string) {
		gotClaims = claims
		gotToken = tokenString
	}

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Request.Header.Set("Authorization", "Bearer "+token)

	JWTAuthenticationMiddleware(jwtAuth, WithAfterAuth(hook))(c)

	if gotClaims == nil || gotClaims.Username != "alice" {
		t.Fatalf("expected claims with username alice, got %+v", gotClaims)
	}
	if gotToken != token {
		t.Fatalf("expected token %q, got %q", token, gotToken)
	}
}

func TestGetAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns user when set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		user := &auth.AuthenticatedUser{Username: "bob"}
		c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
		c.Request = c.Request.WithContext(auth.InjectAuthenticatedUser(c.Request.Context(), user))

		got, ok := GetAuthUser(c)
		if !ok {
			t.Fatal("expected user to be found")
		}
		if got.Username != "bob" {
			t.Fatalf("expected username bob, got %q", got.Username)
		}
	})

	t.Run("returns false when not set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)

		got, ok := GetAuthUser(c)
		if ok {
			t.Fatal("expected user to not be found")
		}
		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})
}

func TestGetAuthUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns ID when set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		user := &auth.AuthenticatedUser{ID: 42, Username: "bob"}
		c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
		c.Request = c.Request.WithContext(auth.InjectAuthenticatedUser(c.Request.Context(), user))

		got, ok := GetAuthUserID(c)
		if !ok {
			t.Fatal("expected user ID to be found")
		}
		if got != 42 {
			t.Fatalf("expected user ID 42, got %d", got)
		}
	})

	t.Run("returns false when ID is empty", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		user := &auth.AuthenticatedUser{Username: "bob"}
		c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
		c.Request = c.Request.WithContext(auth.InjectAuthenticatedUser(c.Request.Context(), user))

		got, ok := GetAuthUserID(c)
		if ok {
			t.Fatal("expected user ID to not be found")
		}
		if got != 0 {
			t.Fatalf("expected zero user ID, got %d", got)
		}
	})
}
