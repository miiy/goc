package sessionauth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ginsessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/miiy/goc/gin/authctx"
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

func testContextWithSession(session *testSession) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
	c.Set(ginsessions.DefaultKey, session)
	return c
}

func TestAuthenticateAcceptsStringID(t *testing.T) {
	session := newTestSession()
	session.Set(SessionKeyUser, map[string]any{"id": "1", "username": "user"})

	c := testContextWithSession(session)
	Authenticate()(c)

	user, ok := authctx.CurrentUser(c)
	if !ok {
		t.Fatal("expected session user")
	}
	if user.ID != "1" || user.Username != "user" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestAuthenticateRejectsEmptyOrNonStringID(t *testing.T) {
	tests := []struct {
		name string
		id   any
	}{
		{name: "empty", id: ""},
		{name: "number type", id: float64(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := newTestSession()
			session.Set(SessionKeyUser, map[string]any{"id": tt.id, "username": "user"})

			c := testContextWithSession(session)
			Authenticate()(c)

			user, ok := authctx.CurrentUser(c)
			if ok {
				t.Fatalf("expected invalid session user, got %#v", user)
			}
			if !c.IsAborted() {
				t.Fatal("expected request to be aborted")
			}
		})
	}
}

func TestAuthenticateSetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	session := newTestSession()
	session.Set(SessionKeyUser, map[string]any{"id": "1", "username": "user"})

	c := testContextWithSession(session)
	Authenticate()(c)

	user, ok := authctx.CurrentUser(c)
	if !ok {
		t.Fatal("expected user in request context")
	}
	if user.ID != "1" || user.Username != "user" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestLoadSessionUserSetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	session := newTestSession()
	session.Set(SessionKeyUser, map[string]any{"id": "1", "username": "user"})

	c := testContextWithSession(session)
	LoadSessionUser()(c)

	user, ok := authctx.CurrentUser(c)
	if !ok {
		t.Fatal("expected user in request context")
	}
	if user.ID != "1" || user.Username != "user" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestLoadSessionUserAllowsAnonymousRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c := testContextWithSession(newTestSession())
	LoadSessionUser()(c)

	if c.IsAborted() {
		t.Fatal("expected request to continue")
	}
	if _, ok := authctx.CurrentUser(c); ok {
		t.Fatal("expected no context user")
	}
}

func TestAuthenticateRedirectsWhenMissingAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c := testContextWithSession(newTestSession())
	Authenticate(WithRedirect("/login"))(c)

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

func TestAuthenticateReturnsUnauthorizedByDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c := testContextWithSession(newTestSession())
	Authenticate()(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if c.Writer.Status() != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", c.Writer.Status())
	}
}

func TestAuthenticateUsesCustomSessionKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	session := newTestSession()
	session.Set("custom-auth", map[string]any{"id": "9", "username": "custom"})

	c := testContextWithSession(session)
	Authenticate(WithSessionKey("custom-auth"))(c)

	user, ok := authctx.CurrentUser(c)
	if !ok {
		t.Fatal("expected user in request context")
	}
	if user.ID != "9" || user.Username != "custom" {
		t.Fatalf("unexpected user: %#v", user)
	}
}
