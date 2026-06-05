package auth

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	ginsessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
)

type testSession struct {
	values map[interface{}]interface{}
}

func newTestSession() *testSession {
	return &testSession{
		values: make(map[interface{}]interface{}),
	}
}

func (s *testSession) ID() string { return "" }

func (s *testSession) Get(key interface{}) interface{} {
	return s.values[key]
}

func (s *testSession) Set(key interface{}, val interface{}) {
	s.values[key] = val
}

func (s *testSession) Delete(key interface{}) {
	delete(s.values, key)
}

func (s *testSession) Clear() {
	for key := range s.values {
		delete(s.values, key)
	}
}

func (s *testSession) AddFlash(interface{}, ...string) {}

func (s *testSession) Flashes(...string) []interface{} {
	return nil
}

func (s *testSession) Options(ginsessions.Options) {}

func (s *testSession) Save() error {
	return nil
}

func TestSessionAuthenticationMiddlewareRedirectsWhenMissingAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Set(ginsessions.DefaultKey, newTestSession())

	SessionAuthenticationMiddleware("/login")(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if c.Writer.Status() != 302 {
		t.Fatalf("expected redirect status 302, got %d", c.Writer.Status())
	}
	if location := c.Writer.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}

func TestSessionAuthenticationMiddlewareSetsMapAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	session := newTestSession()
	session.Set(SessionKeyAuthUser, map[string]any{"id": float64(1), "username": "user"})
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Set(ginsessions.DefaultKey, session)

	SessionAuthenticationMiddleware("/login")(c)

	user, ok := GetAuthUser(c)
	if !ok {
		t.Fatal("expected auth user in request context")
	}
	if user.ID != 1 || user.Username != "user" {
		t.Fatalf("unexpected user: %#v", user)
	}

	ctxUser, err := gocauth.ExtractAuthenticatedUser(c.Request.Context())
	if err != nil {
		t.Fatalf("expected auth user in request context: %v", err)
	}
	if ctxUser.ID != 1 || ctxUser.Username != "user" {
		t.Fatalf("unexpected request context user: %#v", ctxUser)
	}
}

func TestSessionAuthenticationMiddlewareRedirectsWhenAuthTypeIsInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	session := newTestSession()
	session.Set(SessionKeyAuthUser, map[string]any{"id": float64(1)})
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Set(ginsessions.DefaultKey, session)

	SessionAuthenticationMiddleware("/login")(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if c.Writer.Status() != http.StatusFound {
		t.Fatalf("expected redirect status 302, got %d", c.Writer.Status())
	}
	if location := c.Writer.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}

func TestSessionAuthenticationMiddlewareRedirectsWhenAuthIDIsInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	session := newTestSession()
	session.Set(SessionKeyAuthUser, map[string]any{"id": float64(0), "username": "user"})
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Set(ginsessions.DefaultKey, session)

	SessionAuthenticationMiddleware("/login")(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if c.Writer.Status() != http.StatusFound {
		t.Fatalf("expected redirect status 302, got %d", c.Writer.Status())
	}
	if location := c.Writer.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}

func TestSessionUserAcceptsCommonIDTypes(t *testing.T) {
	tests := []struct {
		name string
		id   any
	}{
		{name: "int", id: int(1)},
		{name: "int64", id: int64(1)},
		{name: "float32", id: float32(1)},
		{name: "float64", id: float64(1)},
		{name: "string", id: "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, ok := SessionUser(map[string]any{"id": tt.id, "username": "user"})
			if !ok {
				t.Fatal("expected session user")
			}
			if user.ID != 1 || user.Username != "user" {
				t.Fatalf("unexpected user: %#v", user)
			}
		})
	}
}

func TestSessionUserRejectsInvalidFloatID(t *testing.T) {
	tests := []struct {
		name string
		id   any
	}{
		{name: "float32 fraction", id: float32(1.5)},
		{name: "float64 fraction", id: float64(1.9)},
		{name: "nan", id: math.NaN()},
		{name: "positive infinity", id: math.Inf(1)},
		{name: "negative infinity", id: math.Inf(-1)},
		{name: "int64 overflow", id: float64(1 << 63)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, ok := SessionUser(map[string]any{"id": tt.id, "username": "user"})
			if ok {
				t.Fatalf("expected invalid session user, got %#v", user)
			}
		})
	}
}

func TestSessionAuthenticationMiddlewareUsesDefaultRedirectPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Set(ginsessions.DefaultKey, newTestSession())

	SessionAuthenticationMiddleware("")(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if c.Writer.Status() != http.StatusFound {
		t.Fatalf("expected redirect status 302, got %d", c.Writer.Status())
	}
	if location := c.Writer.Header().Get("Location"); location != "/register" {
		t.Fatalf("expected redirect to /register, got %q", location)
	}
}
