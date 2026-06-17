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
	"google.golang.org/grpc/metadata"
)

func TestAuthenticationMiddlewareAuthenticates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authenticator := func(ctx context.Context, token string) (*auth.AuthenticatedUser, error) {
		if token != "valid" {
			return nil, errors.New("invalid token")
		}
		return &auth.AuthenticatedUser{ID: 42, Username: "alice"}, nil
	}

	var gotUser *auth.AuthenticatedUser
	var gotMD metadata.MD
	r := gin.New()
	r.Use(AuthenticationMiddleware(authenticator, WithMetadataPropagation()))
	r.GET("/private", func(c *gin.Context) {
		gotUser, _ = auth.ExtractAuthenticatedUser(c.Request.Context())
		gotMD, _ = metadata.FromOutgoingContext(c.Request.Context())
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer valid")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if gotUser == nil || gotUser.ID != 42 || gotUser.Username != "alice" {
		t.Fatalf("expected injected user {42 alice}, got %+v", gotUser)
	}
	if v := gotMD.Get(auth.AuthenticatedUserIDMetadataKey); len(v) == 0 || v[0] != "42" {
		t.Fatalf("expected metadata x-auth-user-id=42, got %v", v)
	}
	if v := gotMD.Get(auth.AuthenticatedUsernameMetadataKey); len(v) == 0 || v[0] != "alice" {
		t.Fatalf("expected metadata x-auth-username=alice, got %v", v)
	}
}

func TestAuthenticationMiddlewareRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authenticator := func(context.Context, string) (*auth.AuthenticatedUser, error) {
		t.Fatal("authenticator should not be called when token is missing")
		return nil, nil
	}

	r := gin.New()
	r.Use(AuthenticationMiddleware(authenticator, WithMetadataPropagation()))
	r.GET("/private", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthenticationMiddlewareRejectsAuthenticatorError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var gotMessage string
	authenticator := func(context.Context, string) (*auth.AuthenticatedUser, error) {
		return nil, errors.New("token revoked")
	}
	customHandler := func(ctx *gin.Context, message string) {
		gotMessage = message
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": message})
	}

	r := gin.New()
	r.Use(AuthenticationMiddleware(authenticator, WithUnauthorized(customHandler)))
	r.GET("/private", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer whatever")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
	if gotMessage != "token revoked" {
		t.Fatalf("expected message token revoked, got %q", gotMessage)
	}
}

func TestAuthenticationMiddlewareRejectsNilUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authenticator := func(context.Context, string) (*auth.AuthenticatedUser, error) {
		return nil, nil // authenticator succeeded but returned no user
	}

	r := gin.New()
	r.Use(AuthenticationMiddleware(authenticator, WithMetadataPropagation()))
	r.GET("/private", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer whatever")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
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

func TestBearerToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns token from Authorization header", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
		c.Request.Header.Set("Authorization", "Bearer abc.def.ghi")

		got, err := BearerToken(c)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != "abc.def.ghi" {
			t.Fatalf("expected token abc.def.ghi, got %q", got)
		}
	})

	t.Run("returns error when header missing", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)

		got, err := BearerToken(c)
		if !errors.Is(err, request.ErrNoTokenInRequest) {
			t.Fatalf("expected ErrNoTokenInRequest, got %v", err)
		}
		if got != "" {
			t.Fatalf("expected empty token, got %q", got)
		}
	})

	t.Run("returns error when request is nil", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = nil

		got, err := BearerToken(c)
		if !errors.Is(err, request.ErrNoTokenInRequest) {
			t.Fatalf("expected ErrNoTokenInRequest, got %v", err)
		}
		if got != "" {
			t.Fatalf("expected empty token, got %q", got)
		}
	})
}
