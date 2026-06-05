package csrf

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ginsessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type testSession struct {
	values  map[interface{}]interface{}
	saveErr error
	saved   bool
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
	s.saved = true
	return s.saveErr
}

func newTestContext(method, target, body string) (*gin.Context, *httptest.ResponseRecorder, *testSession) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	if target == "" {
		target = "/"
	}
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	c.Request = httptest.NewRequest(method, target, reader)
	if body != "" {
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	session := newTestSession()
	c.Set(ginsessions.DefaultKey, session)
	return c, recorder, session
}

func TestTokenCreatesAndSavesSessionToken(t *testing.T) {
	c, _, session := newTestContext(http.MethodGet, "/", "")

	token := Token(c)
	if token == "" {
		t.Fatal("expected token")
	}
	if got := session.Get(SessionKey); got != token {
		t.Fatalf("expected session token %q, got %#v", token, got)
	}
	if !session.saved {
		t.Fatal("expected session to be saved")
	}
}

func TestTokenUsesExistingSessionToken(t *testing.T) {
	c, _, session := newTestContext(http.MethodGet, "/", "")
	session.Set(SessionKey, "existing")

	token := Token(c)
	if token != "existing" {
		t.Fatalf("expected existing token, got %q", token)
	}
	if session.saved {
		t.Fatal("expected existing token not to save session")
	}
}

func TestTokenCachesTokenInContext(t *testing.T) {
	c, _, session := newTestContext(http.MethodGet, "/", "")

	first := Token(c)
	session.saved = false
	second := Token(c)

	if first == "" || second != first {
		t.Fatalf("expected cached token %q, got %q", first, second)
	}
	if session.saved {
		t.Fatal("expected cached token not to save session")
	}
}

func TestTokenReportsSaveError(t *testing.T) {
	c, _, session := newTestContext(http.MethodGet, "/", "")
	session.saveErr = errors.New("save failed")

	if token := Token(c); token != "" {
		t.Fatalf("expected empty token on save error, got %q", token)
	}
	if len(c.Errors) != 1 {
		t.Fatalf("expected one context error, got %d", len(c.Errors))
	}
}

func TestMiddlewareAllowsSafeMethods(t *testing.T) {
	c, _, _ := newTestContext(http.MethodGet, "/", "")

	Middleware()(c)

	if c.IsAborted() {
		t.Fatal("expected safe method not to abort")
	}
}

func TestMiddlewareRejectsMissingToken(t *testing.T) {
	c, recorder, session := newTestContext(http.MethodPost, "/", "")
	session.Set(SessionKey, "expected")

	Middleware()(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}

func TestMiddlewareAcceptsFormToken(t *testing.T) {
	c, _, session := newTestContext(http.MethodPost, "/", "_csrf=expected")
	session.Set(SessionKey, "expected")

	Middleware()(c)

	if c.IsAborted() {
		t.Fatal("expected matching form token not to abort")
	}
}

func TestMiddlewareAcceptsHeaderToken(t *testing.T) {
	c, _, session := newTestContext(http.MethodPost, "/", "")
	c.Request.Header.Set(HeaderName, "expected")
	session.Set(SessionKey, "expected")

	Middleware()(c)

	if c.IsAborted() {
		t.Fatal("expected matching header token not to abort")
	}
}

func TestMiddlewareSupportsCustomOptions(t *testing.T) {
	c, _, session := newTestContext(http.MethodPost, "/", "csrf_token=expected")
	session.Set("custom_csrf", "expected")

	Middleware(
		WithFieldName("csrf_token"),
		WithSessionKey("custom_csrf"),
	)(c)

	if c.IsAborted() {
		t.Fatal("expected matching custom token not to abort")
	}
}

func TestMiddlewareUsesCustomUnauthorizedHandler(t *testing.T) {
	c, recorder, session := newTestContext(http.MethodPost, "/", "")
	session.Set(SessionKey, "expected")

	Middleware(WithUnauthorized(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "bad csrf")
	}))(c)

	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}
