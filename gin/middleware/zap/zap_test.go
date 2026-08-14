package zap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	uberzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestGinzapAddsRequestFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zapcore.InfoLevel)
	logger := uberzap.New(core)
	router := gin.New()
	router.Use(Ginzap(logger))
	router.GET("/items", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/items?page=2", nil)
	req.Header.Set("sign", "test-signature")
	req.Header.Set("ts", "1723000000")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	entries := logs.All()
	if len(entries) != 1 {
		t.Fatalf("expected one access log entry, got %d", len(entries))
	}

	fields := entries[0].ContextMap()
	want := map[string]any{
		"request_uri": "/items?page=2",
		"sign":        "test-signature",
		"ts":          "1723000000",
		"size":        int64(2),
	}
	for key, value := range want {
		if fields[key] != value {
			t.Errorf("field %q: got %#v, want %#v", key, fields[key], value)
		}
	}
}

func TestGinzapDoesNotLogBodies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, logs := observer.New(zapcore.InfoLevel)
	logger := uberzap.New(core)
	router := gin.New()
	router.Use(Ginzap(logger))
	router.POST("/items", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		c.String(http.StatusOK, string(body))
	})

	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"name":"goc"}`))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Body.String() != `{"name":"goc"}` {
		t.Fatalf("response body: got %q", recorder.Body.String())
	}

	entries := logs.All()
	if len(entries) != 1 {
		t.Fatalf("expected one access log entry, got %d", len(entries))
	}

	fields := entries[0].ContextMap()
	if _, ok := fields["request-body"]; ok {
		t.Error("request-body should not be logged")
	}
	if _, ok := fields["response-body"]; ok {
		t.Error("response-body should not be logged")
	}
}
