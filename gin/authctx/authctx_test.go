package authctx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
)

func TestSetAndCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)

	SetUser(c, &gocauth.AuthenticatedUser{ID: "42", Username: "alice"})

	user, ok := CurrentUser(c)
	if !ok {
		t.Fatal("expected user")
	}
	if user.ID != "42" || user.Username != "alice" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestCurrentUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	SetUser(c, &gocauth.AuthenticatedUser{ID: "42", Username: "alice"})

	id, ok := CurrentUserID(c)
	if !ok {
		t.Fatal("expected user id")
	}
	if id != "42" {
		t.Fatalf("expected 42, got %q", id)
	}
}

func TestCurrentUserIDRejectsEmptyID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	SetUser(c, &gocauth.AuthenticatedUser{Username: "alice"})

	id, ok := CurrentUserID(c)
	if ok {
		t.Fatal("expected no user id")
	}
	if id != "" {
		t.Fatalf("expected empty id, got %q", id)
	}
}

func TestCurrentUserInt64ID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	SetUser(c, &gocauth.AuthenticatedUser{ID: "42", Username: "alice"})

	id, ok := CurrentUserInt64ID(c)
	if !ok {
		t.Fatal("expected user id")
	}
	if id != 42 {
		t.Fatalf("expected 42, got %d", id)
	}
}

func TestCurrentUserInt64IDRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	SetUser(c, &gocauth.AuthenticatedUser{ID: "abc", Username: "alice"})

	id, ok := CurrentUserInt64ID(c)
	if ok {
		t.Fatal("expected no user id")
	}
	if id != 0 {
		t.Fatalf("expected zero id, got %d", id)
	}
}
