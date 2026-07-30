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
	reader := strings.NewReader(body)
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

func TestTokenReusesSessionToken(t *testing.T) {
	c, _, session := newTestContext(http.MethodGet, "/", "")

	first := Token(c)
	session.saved = false
	second := Token(c)

	if first == "" || second != first {
		t.Fatalf("expected session token %q, got %q", first, second)
	}
	if session.saved {
		t.Fatal("expected existing token not to save session")
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
	if token := session.Get(SessionKey); token != nil {
		t.Fatalf("expected failed token to be removed, got %#v", token)
	}
}

func TestTokenReportsGenerationError(t *testing.T) {
	c, _, session := newTestContext(http.MethodGet, "/", "")

	token := Token(c, optionFunc(func(opts *options) {
		opts.tokenGenerator = func() (string, error) {
			return "", errors.New("generate failed")
		}
	}))

	if token != "" {
		t.Fatalf("expected empty token on generation error, got %q", token)
	}
	if len(c.Errors) != 1 {
		t.Fatalf("expected one context error, got %d", len(c.Errors))
	}
	if session.saved {
		t.Fatal("expected session not to be saved")
	}
}

func TestRotateTokenReplacesExistingToken(t *testing.T) {
	c, _, session := newTestContext(http.MethodGet, "/", "")
	session.Set(SessionKey, "existing")

	token := RotateToken(c)

	if token == "" || token == "existing" {
		t.Fatalf("expected a new token, got %q", token)
	}
	if got := session.Get(SessionKey); got != token {
		t.Fatalf("expected session token %q, got %#v", token, got)
	}
	if !session.saved {
		t.Fatal("expected session to be saved")
	}
}

func TestRotateTokenPreservesExistingTokenOnSaveError(t *testing.T) {
	c, _, session := newTestContext(http.MethodGet, "/", "")
	session.Set(SessionKey, "existing")
	session.saveErr = errors.New("save failed")

	if token := RotateToken(c); token != "" {
		t.Fatalf("expected empty token on save error, got %q", token)
	}
	if got := session.Get(SessionKey); got != "existing" {
		t.Fatalf("expected existing token to be preserved, got %#v", got)
	}
}

func TestMiddlewareAllowsSafeMethods(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
		t.Run(method, func(t *testing.T) {
			c, _, _ := newTestContext(method, "/", "")

			Middleware()(c)

			if c.IsAborted() {
				t.Fatal("expected safe method not to abort")
			}
		})
	}
}

func TestCheckerValidatesWithoutContinuingGinChain(t *testing.T) {
	checker, err := NewChecker()
	if err != nil {
		t.Fatalf("new checker: %v", err)
	}
	c, _, session := newTestContext(http.MethodPost, "/", "")
	c.Request.Header.Set(HeaderName, "expected")
	session.Set(SessionKey, "expected")

	if !checker.Check(c) {
		t.Fatal("expected matching token to pass")
	}
	if c.IsAborted() {
		t.Fatal("expected valid request not to abort")
	}
}

func TestCheckerRejectsInvalidToken(t *testing.T) {
	checker, err := NewChecker()
	if err != nil {
		t.Fatalf("new checker: %v", err)
	}
	c, recorder, session := newTestContext(http.MethodPost, "/", "")
	c.Request.Header.Set(HeaderName, "invalid")
	session.Set(SessionKey, "expected")

	if checker.Check(c) {
		t.Fatal("expected invalid token to fail")
	}
	if !c.IsAborted() {
		t.Fatal("expected request to be aborted")
	}
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}

func TestNewCheckerRejectsInvalidTrustedOrigin(t *testing.T) {
	checker, err := NewChecker(WithTrustedOrigins("app.example"))
	if err == nil {
		t.Fatal("expected invalid trusted origin error")
	}
	if checker != nil {
		t.Fatal("expected checker to be nil")
	}
}

func TestCheckerValidatesWithPreparedTrustedOrigins(t *testing.T) {
	checker, err := NewChecker(WithTrustedOrigins("https://app.example"))
	if err != nil {
		t.Fatalf("new checker: %v", err)
	}
	c, _, session := newTestContext(http.MethodPost, "/", "")
	c.Request.Header.Set(HeaderName, "expected")
	c.Request.Header.Set("Origin", "https://app.example")
	c.Request.Header.Set("Sec-Fetch-Site", "cross-site")
	session.Set(SessionKey, "expected")

	if !checker.Check(c) {
		t.Fatal("expected trusted origin with matching token to pass")
	}
	if c.IsAborted() {
		t.Fatal("expected valid request not to abort")
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

func TestMiddlewareRejectsInvalidToken(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			c, recorder, session := newTestContext(method, "/", "")
			c.Request.Header.Set(HeaderName, "invalid")
			session.Set(SessionKey, "expected")

			Middleware()(c)

			if !c.IsAborted() {
				t.Fatal("expected request to be aborted")
			}
			if recorder.Code != http.StatusForbidden {
				t.Fatalf("expected status 403, got %d", recorder.Code)
			}
		})
	}
}

func TestMiddlewareValidatesRequestOrigin(t *testing.T) {
	tests := []struct {
		name         string
		origin       string
		secFetchSite string
		want         int
	}{
		{name: "same origin", origin: "http://example.com", want: http.StatusOK},
		{name: "foreign origin", origin: "https://attacker.example", want: http.StatusForbidden},
		{name: "null origin", origin: "null", want: http.StatusForbidden},
		{name: "same-origin fetch", secFetchSite: "same-origin", want: http.StatusOK},
		{name: "browser navigation", secFetchSite: "none", want: http.StatusOK},
		{name: "same-site fetch", secFetchSite: "same-site", want: http.StatusForbidden},
		{name: "cross-site fetch", secFetchSite: "cross-site", want: http.StatusForbidden},
		{name: "no source headers", want: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, recorder, session := newTestContext(http.MethodPost, "/", "")
			c.Request.Header.Set(HeaderName, "expected")
			if tt.origin != "" {
				c.Request.Header.Set("Origin", tt.origin)
			}
			if tt.secFetchSite != "" {
				c.Request.Header.Set("Sec-Fetch-Site", tt.secFetchSite)
			}
			session.Set(SessionKey, "expected")

			Middleware()(c)

			if recorder.Code != tt.want {
				t.Fatalf("expected status %d, got %d", tt.want, recorder.Code)
			}
		})
	}
}

func TestMiddlewareAllowsTrustedOrigin(t *testing.T) {
	c, _, session := newTestContext(http.MethodPost, "/", "")
	c.Request.Header.Set(HeaderName, "expected")
	c.Request.Header.Set("Origin", "https://app.example")
	c.Request.Header.Set("Sec-Fetch-Site", "cross-site")
	session.Set(SessionKey, "expected")

	Middleware(WithTrustedOrigins("https://app.example"))(c)

	if c.IsAborted() {
		t.Fatal("expected trusted origin not to abort")
	}
}

func TestMiddlewareRejectsInvalidTrustedOrigin(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected invalid trusted origin to panic")
		}
	}()

	Middleware(WithTrustedOrigins("app.example"))
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
