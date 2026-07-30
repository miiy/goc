package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
)

func TestUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)

	injectSession(c, &gocauth.AuthenticatedSession{
		User: gocauth.AuthenticatedUser{ID: "42", Username: "alice"},
	})

	user, ok := User(c)
	if !ok {
		t.Fatal("expected user")
	}
	if user.ID != "42" || user.Username != "alice" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestSessionPanicsOnNilContext(t *testing.T) {
	requirePanic(t, func() {
		Session(nil)
	})
}

func TestSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	injectSession(c, &gocauth.AuthenticatedSession{
		User:       gocauth.AuthenticatedUser{ID: "42", Username: "alice"},
		SIDKey:     "session-key",
		ClientType: "app",
		ClientID:   "ios",
	})

	session, ok := Session(c)
	if !ok || session.SIDKey != "session-key" || session.User.ID != "42" || session.User.Username != "alice" ||
		session.ClientType != "app" || session.ClientID != "ios" {
		t.Fatalf("unexpected session: %#v", session)
	}
	user, ok := User(c)
	if !ok || user.ID != "42" || user.Username != "alice" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestSessionAllowsUserOnlySession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	injectSession(c, &gocauth.AuthenticatedSession{
		User: gocauth.AuthenticatedUser{ID: "43", Username: "bob"},
	})

	session, ok := Session(c)
	if !ok || session.User.ID != "43" || session.User.Username != "bob" || session.SIDKey != "" {
		t.Fatalf("unexpected session: %#v", session)
	}
	user, ok := User(c)
	if !ok || user.ID != "43" || user.Username != "bob" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	injectSession(c, &gocauth.AuthenticatedSession{
		User: gocauth.AuthenticatedUser{ID: "42", Username: "alice"},
	})

	id, ok := UserID(c)
	if !ok {
		t.Fatal("expected user id")
	}
	if id != "42" {
		t.Fatalf("expected 42, got %q", id)
	}
}

func TestUserIDRejectsEmptyID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	injectSession(c, &gocauth.AuthenticatedSession{
		User: gocauth.AuthenticatedUser{Username: "alice"},
	})

	id, ok := UserID(c)
	if ok {
		t.Fatal("expected no user id")
	}
	if id != "" {
		t.Fatalf("expected empty id, got %q", id)
	}
}

func TestUserInt64ID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	injectSession(c, &gocauth.AuthenticatedSession{
		User: gocauth.AuthenticatedUser{ID: "42", Username: "alice"},
	})

	id, ok := UserInt64ID(c)
	if !ok {
		t.Fatal("expected user id")
	}
	if id != 42 {
		t.Fatalf("expected 42, got %d", id)
	}
}

func TestUserInt64IDRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	injectSession(c, &gocauth.AuthenticatedSession{
		User: gocauth.AuthenticatedUser{ID: "abc", Username: "alice"},
	})

	id, ok := UserInt64ID(c)
	if ok {
		t.Fatal("expected no user id")
	}
	if id != 0 {
		t.Fatalf("expected zero id, got %d", id)
	}
}

func injectSession(ctx *gin.Context, session *gocauth.AuthenticatedSession) {
	ctx.Request = ctx.Request.WithContext(gocauth.InjectAuthenticatedSession(ctx.Request.Context(), session))
}

func requirePanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}
