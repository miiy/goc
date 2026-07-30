package sessions

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	ginsessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gsessions "github.com/gorilla/sessions"
)

func TestUseJSONSerializerRejectsUnsupportedStore(t *testing.T) {
	err := UseJSONSerializer(NewCookieStore("secret"))
	if !errors.Is(err, ErrUnsupportedJSONSessionStore) {
		t.Fatalf("expected ErrUnsupportedJSONSessionStore, got %v", err)
	}
}

func TestSetMaxAgeRejectsUnsupportedStore(t *testing.T) {
	err := SetMaxAge(NewCookieStore("secret"), 300)
	if !errors.Is(err, ErrUnsupportedRedisSessionStore) {
		t.Fatalf("expected ErrUnsupportedRedisSessionStore, got %v", err)
	}
}

func TestSessionManagerDefaultsBlankCookieName(t *testing.T) {
	for _, name := range []string{"", "  "} {
		manager := NewSessionManager(nil, Options{}, WithCookieName(name))
		if manager.CookieName() != DefaultSessionCookieName {
			t.Fatalf("blank manager cookie name = %q, want %q", manager.CookieName(), DefaultSessionCookieName)
		}
	}
}

func TestSessionManagerKeepsConfigAndDelegates(t *testing.T) {
	store := &recordingStore{}
	manager := NewSessionManager(store, Options{Path: "/", MaxAge: 600, HttpOnly: true}, WithCookieName("custom-session"))

	if manager.CookieName() != "custom-session" {
		t.Fatalf("manager cookie name = %q, want custom-session", manager.CookieName())
	}
	if manager.Options().MaxAge != 600 {
		t.Fatalf("manager options = %+v, want MaxAge 600", manager.Options())
	}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: "custom-session", Value: "old"})
	request.AddCookie(&http.Cookie{Name: "other", Value: "keep"})
	c.Request = request

	if err := manager.RenewWithOptions(c, "new-id", map[any]any{"key": "value"}, Options{Path: "/", MaxAge: 60}); err != nil {
		t.Fatalf("renew managed Session: %v", err)
	}
	if store.saved == nil || store.saved.ID != "new-id" || store.saved.Values["key"] != "value" {
		t.Fatalf("saved Session = %+v", store.saved)
	}
	if cookies := store.newRequest.Cookies(); len(cookies) != 1 || cookies[0].Name != "other" {
		t.Fatalf("renew request cookies = %+v, want only unrelated cookie", cookies)
	}
}

func TestSIDReturnsCurrentSessionID(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ginsessions.DefaultKey, &testSession{id: "session-id"})

	if got := SID(c); got != "session-id" {
		t.Fatalf("SID = %q, want session-id", got)
	}
	requirePanic(t, func() {
		SID(nil)
	})
	c, _ = gin.CreateTestContext(httptest.NewRecorder())
	if got := SID(c); got != "" {
		t.Fatalf("context without Session SID = %q, want empty", got)
	}
}

func TestRenewIgnoresExistingCookieAndSavesFreshSession(t *testing.T) {
	store := &recordingStore{}
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: "session", Value: "old"})
	request.AddCookie(&http.Cookie{Name: "other", Value: "keep"})
	c.Request = request

	if err := Renew(c, store, "session", "new-id", map[any]any{"key": "value"}); err != nil {
		t.Fatalf("renew Session: %v", err)
	}
	if store.saved == nil {
		t.Fatal("session was not saved")
	}
	if store.saved.ID != "new-id" || store.saved.Values["key"] != "value" {
		t.Fatalf("saved Session = %+v", store.saved)
	}
	cookies := store.newRequest.Cookies()
	if len(cookies) != 1 || cookies[0].Name != "other" || cookies[0].Value != "keep" {
		t.Fatalf("renew request cookies = %+v, want only unrelated cookie", cookies)
	}
}

func TestRenewWithOptionsSavesFreshSessionOptions(t *testing.T) {
	store := &recordingStore{}
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	options := Options{Path: "/", MaxAge: 60, Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode}

	if err := RenewWithOptions(c, store, "session", "new-id", map[any]any{"key": "value"}, &options); err != nil {
		t.Fatalf("renew Session: %v", err)
	}
	if store.saved == nil || store.saved.Options == nil {
		t.Fatal("session options were not saved")
	}
	if store.saved.Options.MaxAge != 60 || store.saved.Options.Path != "/" || !store.saved.Options.Secure || !store.saved.Options.HttpOnly {
		t.Fatalf("saved Session options = %+v, want explicit options", store.saved.Options)
	}
}

func TestClearExpiresSessionCookie(t *testing.T) {
	session := newTestSession()
	session.Set("key", "value")
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ginsessions.DefaultKey, session)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	if err := Clear(c, "session", Options{Path: "/", MaxAge: 600, HttpOnly: true, SameSite: http.SameSiteLaxMode}); err != nil {
		t.Fatalf("clear Session: %v", err)
	}
	if session.Get("key") != nil {
		t.Fatalf("Session value was not cleared: %#v", session.Get("key"))
	}
	if !session.saved {
		t.Fatal("Session was not saved")
	}
	if session.options.MaxAge != -1 || session.options.Path != "/" || !session.options.HttpOnly {
		t.Fatalf("Session options = %+v, want expired cookie options", session.options)
	}
}

func TestClearWritesExpiredCookieWhenSaveFails(t *testing.T) {
	session := newTestSession()
	session.saveErr = errors.New("save failed")
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Set(ginsessions.DefaultKey, session)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	err := Clear(c, "session", Options{Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	if err == nil {
		t.Fatal("clear Session error = nil, want save error")
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "session" || cookies[0].MaxAge != -1 || cookies[0].Value != "" {
		t.Fatalf("fallback cookies = %+v, want expired session cookie", cookies)
	}
}

type recordingStore struct {
	newRequest *http.Request
	saved      *gsessions.Session
}

func (s *recordingStore) Get(request *http.Request, name string) (*gsessions.Session, error) {
	return s.New(request, name)
}

func (s *recordingStore) New(request *http.Request, name string) (*gsessions.Session, error) {
	s.newRequest = request
	return &gsessions.Session{Values: map[interface{}]interface{}{}, Options: &gsessions.Options{}}, nil
}

func (s *recordingStore) Save(_ *http.Request, _ http.ResponseWriter, session *gsessions.Session) error {
	s.saved = session
	return nil
}

func (s *recordingStore) Options(Options) {}

func requirePanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}
