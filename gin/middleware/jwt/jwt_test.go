package jwt

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/miiy/goc/auth"
	ginauth "github.com/miiy/goc/gin/auth"
	"google.golang.org/grpc/metadata"
)

func TestAuthenticateInjectsTrustedSession(t *testing.T) {
	var gotToken string
	router := gin.New()
	router.Use(Authenticate(func(_ context.Context, token string) (*auth.AuthenticatedSession, error) {
		gotToken = token
		return &auth.AuthenticatedSession{User: auth.AuthenticatedUser{ID: "42", Username: "alice"}, SIDKey: "sid-key", ClientType: "app", ClientID: "ios"}, nil
	}, WithMetadataPropagation()))
	router.GET("/private", func(c *gin.Context) {
		user, ok := ginauth.User(c)
		if !ok || user.ID != "42" || user.Username != "alice" {
			t.Fatalf("unexpected user: %+v, %v", user, ok)
		}
		session, ok := ginauth.Session(c)
		if !ok || session.SIDKey != "sid-key" || session.ClientType != "app" || session.ClientID != "ios" {
			t.Fatalf("unexpected session: %+v, %v", session, ok)
		}
		outgoing, _ := metadata.FromOutgoingContext(c.Request.Context())
		if got := outgoing.Get(auth.AuthenticatedUserIDMetadataKey); len(got) != 1 || got[0] != "42" {
			t.Fatalf("unexpected ID metadata: %v", got)
		}
		if got := outgoing.Get(auth.AuthenticatedUsernameMetadataKey); len(got) != 1 || got[0] != "alice" {
			t.Fatalf("unexpected username metadata: %v", got)
		}
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/private", nil)
	request.Header.Set("Authorization", "Bearer token")
	request.Header.Set("Proxy-Authorization", "Basic proxy")
	request.AddCookie(&http.Cookie{Name: "session", Value: "cookie"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if gotToken != "token" {
		t.Fatalf("resolver token = %q, want token", gotToken)
	}
}

func TestAuthenticateRejectsResolverError(t *testing.T) {
	wantErr := errors.New("resolve JWT")
	var gotErr error
	router := gin.New()
	router.Use(Authenticate(func(context.Context, string) (*auth.AuthenticatedSession, error) {
		return nil, wantErr
	}, WithUnauthorized(func(c *gin.Context, err error) {
		gotErr = err
		c.Status(http.StatusUnauthorized)
	})))
	router.GET("/private", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	request := httptest.NewRequest(http.MethodGet, "/private", nil)
	request.Header.Set("Authorization", "Bearer token")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized || !errors.Is(gotErr, wantErr) {
		t.Fatalf("status = %d, error = %v", response.Code, gotErr)
	}
}
