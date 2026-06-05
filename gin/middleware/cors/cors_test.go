package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gincors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func TestNewUsesDefaultAllowOrigins(t *testing.T) {
	router := testRouter(New())

	resp := performCORSRequest(router, "http://localhost:3000")
	if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:3000")
	}
}

func TestNewRejectsOriginOutsideDefaultAllowOrigins(t *testing.T) {
	router := testRouter(New())

	resp := performCORSRequest(router, "http://localhost:8081")
	if resp.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusForbidden)
	}
}

func TestNewRejectsUnconfiguredOrigin(t *testing.T) {
	router := testRouter(New("https://app.example.com"))

	resp := performCORSRequest(router, "https://evil.example.com")
	if resp.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusForbidden)
	}
	if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", got)
	}
}

func TestNewAllowsConfiguredOrigin(t *testing.T) {
	router := testRouter(New("https://app.example.com"))

	resp := performCORSRequest(router, "https://app.example.com")
	if got := resp.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "https://app.example.com")
	}
}

func TestNewWithConfigSupportsCredentials(t *testing.T) {
	config := gincors.DefaultConfig()
	config.AllowOrigins = []string{"https://app.example.com"}
	config.AllowCredentials = true
	router := testRouter(NewWithConfig(config))

	resp := performCORSRequest(router, "https://app.example.com")
	if got := resp.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q, want true", got)
	}
}

func testRouter(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware)
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return router
}

func performCORSRequest(handler http.Handler, origin string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", origin)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	return resp
}
