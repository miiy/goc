package session

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions/filesystem"
	"github.com/gin-gonic/gin"
	"github.com/miiy/goc/auth"
	ginauth "github.com/miiy/goc/gin/auth"
	gocsessions "github.com/miiy/goc/gin/sessions"
	"google.golang.org/grpc/metadata"
)

const testSessionCookieName = "test-session"

func TestAuthenticateRejectsUnknownSession(t *testing.T) {
	store, cookie := protectedSession(t, "unknown-sid")
	resolverCalled := false
	router := gin.New()
	router.Use(gocsessions.Middleware(testSessionCookieName, store))
	router.Use(Authenticate(func(_ context.Context, sid string) (*auth.AuthenticatedSession, error) {
		resolverCalled = true
		if sid != "unknown-sid" {
			t.Fatalf("expected unknown-sid, got %q", sid)
		}
		return nil, errors.New("unknown session")
	}))
	router.GET("/private", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	request := httptest.NewRequest(http.MethodGet, "/private", nil)
	request.AddCookie(cookie)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if !resolverCalled {
		t.Fatal("expected resolver to be called")
	}
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestAuthenticateRejectsMissingSessionID(t *testing.T) {
	router := gin.New()
	router.Use(Authenticate(func(context.Context, string) (*auth.AuthenticatedSession, error) {
		t.Fatal("resolver should not be called without a Session ID")
		return nil, nil
	}))
	router.GET("/private", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/private", nil))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestAuthenticateInjectsTrustedSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store, cookie := protectedSession(t, "session-id")
	if cookie.Value == "session-id" || strings.Contains(cookie.Value, "session-id") {
		t.Fatalf("protected Cookie exposed raw sid: %q", cookie.Value)
	}
	r := gin.New()
	r.Use(gocsessions.Middleware(testSessionCookieName, store))
	r.Use(Authenticate(func(_ context.Context, sid string) (*auth.AuthenticatedSession, error) {
		if sid != "session-id" {
			return nil, errors.New("unexpected sid")
		}
		return &auth.AuthenticatedSession{User: auth.AuthenticatedUser{ID: "42", Username: "alice"}, SIDKey: "sid-key", ClientType: "web"}, nil
	}, WithMetadataPropagation()))
	r.GET("/private", func(c *gin.Context) {
		user, ok := ginauth.User(c)
		if !ok || user.ID != "42" || user.Username != "alice" {
			t.Fatalf("unexpected user: %+v, %v", user, ok)
		}
		session, ok := ginauth.Session(c)
		if !ok || session.SIDKey != "sid-key" || session.ClientType != "web" {
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

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.AddCookie(cookie)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestAuthenticateUsesCustomSessionCookieName(t *testing.T) {
	store, cookie := protectedSessionWithCookieName(t, "session-id", "custom-session")
	resolverCalled := false
	router := gin.New()
	router.Use(gocsessions.Middleware("custom-session", store))
	router.Use(Authenticate(func(_ context.Context, sid string) (*auth.AuthenticatedSession, error) {
		resolverCalled = true
		if sid != "session-id" {
			t.Fatalf("expected session-id, got %q", sid)
		}
		return &auth.AuthenticatedSession{User: auth.AuthenticatedUser{ID: "42", Username: "alice"}}, nil
	}))
	router.GET("/private", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	request := httptest.NewRequest(http.MethodGet, "/private", nil)
	request.AddCookie(cookie)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if !resolverCalled {
		t.Fatal("expected resolver to be called")
	}
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, response.Code)
	}
}

func protectedSession(t *testing.T, sid string) (gocsessions.Store, *http.Cookie) {
	t.Helper()
	return protectedSessionWithCookieName(t, sid, testSessionCookieName)
}

func protectedSessionWithCookieName(t *testing.T, sid string, name string) (gocsessions.Store, *http.Cookie) {
	t.Helper()
	name = strings.TrimSpace(name)
	store := filesystem.NewStore(
		t.TempDir(),
		[]byte("0123456789abcdef0123456789abcdef"),
		[]byte("abcdef0123456789abcdef0123456789"),
	)
	store.Options(gocsessions.Options{Path: "/", MaxAge: 600, Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	current, err := store.New(request, name)
	if err != nil {
		t.Fatalf("new Session: %v", err)
	}
	current.ID = sid
	if err := store.Save(request, response, current); err != nil {
		t.Fatalf("save Session: %v", err)
	}
	for _, cookie := range response.Result().Cookies() {
		if cookie.Name == name && cookie.Value != "" {
			return store, cookie
		}
	}
	t.Fatal("response did not contain a protected Session Cookie")
	return nil, nil
}
