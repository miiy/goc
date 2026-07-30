package gin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwtrequest "github.com/golang-jwt/jwt/v5/request"
)

func TestUserAgentUsesFirstHeaderValue(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header["User-Agent"] = []string{"first", "second"}

	if got := UserAgent(c); got != "first" {
		t.Fatalf("UserAgent = %q, want first", got)
	}
}

func TestUserAgentLimitsLength(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("User-Agent", strings.Repeat("测", MaxUserAgentLength+1))

	got := UserAgent(c)
	if len([]rune(got)) != MaxUserAgentLength {
		t.Fatalf("UserAgent length = %d, want %d", len([]rune(got)), MaxUserAgentLength)
	}
}

func TestUserAgentPanicsOnNilContext(t *testing.T) {
	requirePanic(t, func() {
		UserAgent(nil)
	})
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name    string
		values  []string
		want    string
		wantErr error
	}{
		{name: "valid", values: []string{"Bearer token"}, want: "token"},
		{name: "case insensitive scheme", values: []string{"bearer token"}, want: "token"},
		{name: "missing", wantErr: jwtrequest.ErrNoTokenInRequest},
		{name: "duplicate uses first", values: []string{"Bearer first", "Bearer second"}, want: "first"},
		{name: "wrong scheme", values: []string{"Basic token"}, wantErr: jwtrequest.ErrNoTokenInRequest},
		{name: "missing token", values: []string{"Bearer"}, wantErr: jwtrequest.ErrNoTokenInRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, _ := CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodGet, "/private", nil)
			if test.values != nil {
				c.Request.Header["Authorization"] = test.values
			}
			got, err := ExtractToken(c)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("error = %v, want %v", err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("token = %q, want %q", got, test.want)
			}
		})
	}
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
