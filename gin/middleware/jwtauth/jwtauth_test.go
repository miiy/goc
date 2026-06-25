package jwtauth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
	gocauth "github.com/miiy/goc/auth"
	"google.golang.org/grpc/metadata"
)

func TestAuthenticateAuthenticates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resolveUser := func(ctx context.Context, token string) (*gocauth.AuthenticatedUser, error) {
		if token != "valid" {
			return nil, errors.New("invalid token")
		}
		return &gocauth.AuthenticatedUser{ID: "42", Username: "alice"}, nil
	}

	var gotUser *gocauth.AuthenticatedUser
	var gotMD metadata.MD
	r := gin.New()
	r.Use(Authenticate(resolveUser, WithMetadataPropagation()))
	r.GET("/private", func(c *gin.Context) {
		gotUser, _ = gocauth.ExtractAuthenticatedUser(c.Request.Context())
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
	if gotUser == nil || gotUser.ID != "42" || gotUser.Username != "alice" {
		t.Fatalf("expected injected user {42 alice}, got %+v", gotUser)
	}
	if v := gotMD.Get(gocauth.AuthenticatedUserIDMetadataKey); len(v) == 0 || v[0] != "42" {
		t.Fatalf("expected metadata x-auth-user-id=42, got %v", v)
	}
	if v := gotMD.Get(gocauth.AuthenticatedUsernameMetadataKey); len(v) == 0 || v[0] != "alice" {
		t.Fatalf("expected metadata x-auth-username=alice, got %v", v)
	}
}

func TestAuthenticateRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resolveUser := func(context.Context, string) (*gocauth.AuthenticatedUser, error) {
		t.Fatal("resolver should not be called when token is missing")
		return nil, nil
	}

	r := gin.New()
	r.Use(Authenticate(resolveUser, WithMetadataPropagation()))
	r.GET("/private", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthenticateRejectsResolverError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var gotMessage string
	resolveUser := func(context.Context, string) (*gocauth.AuthenticatedUser, error) {
		return nil, errors.New("token revoked")
	}
	customHandler := func(ctx *gin.Context, message string) {
		gotMessage = message
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": message})
	}

	r := gin.New()
	r.Use(Authenticate(resolveUser, WithUnauthorized(customHandler)))
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

func TestAuthenticateRejectsNilUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resolveUser := func(context.Context, string) (*gocauth.AuthenticatedUser, error) {
		return nil, nil
	}

	r := gin.New()
	r.Use(Authenticate(resolveUser, WithMetadataPropagation()))
	r.GET("/private", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer whatever")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns token from Authorization header", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
		c.Request.Header.Set("Authorization", "Bearer abc.def.ghi")

		got, err := Token(c)
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

		got, err := Token(c)
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

		got, err := Token(c)
		if !errors.Is(err, request.ErrNoTokenInRequest) {
			t.Fatalf("expected ErrNoTokenInRequest, got %v", err)
		}
		if got != "" {
			t.Fatalf("expected empty token, got %q", got)
		}
	})
}
