package accesslog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	goclogger "github.com/miiy/goc/logger"
	uberzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func newTestLogger() (goclogger.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.InfoLevel)
	return &goclogger.ZapLogger{Logger: uberzap.New(core)}, logs
}

func TestNewLogsMetadataByDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, logs := newTestLogger()
	router := gin.New()
	router.Use(New(log))
	router.POST("/items", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			t.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"received": string(body)})
	})

	body := `{"name":"goc"}`
	req := httptest.NewRequest(http.MethodPost, "/items?page=2", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, http.StatusOK)
	}
	fields := logs.All()[0].ContextMap()
	want := map[string]any{
		"method":      http.MethodPost,
		"path":        "/items",
		"query":       "page=2",
		"request_uri": "/items?page=2",
		"full-path":   "/items",
		"status":      int64(http.StatusOK),
	}
	for key, value := range want {
		if fields[key] != value {
			t.Errorf("field %q: got %#v, want %#v", key, fields[key], value)
		}
	}
	if _, ok := fields["request-body"]; ok {
		t.Error("request-body should not be logged by default")
	}
	if _, ok := fields["response-body"]; ok {
		t.Error("response-body should not be logged by default")
	}
	if _, ok := fields["latency-us"]; !ok {
		t.Error("latency-us should be logged")
	}
}

func TestNewLogsRequestAndResponseBodiesWhenEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, logs := newTestLogger()
	router := gin.New()
	router.Use(New(
		log,
		WithRequestBodyLogging(true),
		WithResponseBodyLogging(true),
	))
	router.POST("/items", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			t.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"received": string(body)})
	})

	body := `{"name":"goc"}`
	req := httptest.NewRequest(http.MethodPost, "/items?page=2", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, http.StatusOK)
	}
	fields := logs.All()[0].ContextMap()
	want := map[string]any{
		"request-body":  body,
		"response-body": `{"received":"{\"name\":\"goc\"}"}`,
	}
	for key, value := range want {
		if fields[key] != value {
			t.Errorf("field %q: got %#v, want %#v", key, fields[key], value)
		}
	}
}

func TestNewBodyOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, logs := newTestLogger()
	router := gin.New()
	router.Use(New(
		log,
		WithRequestBodyLogging(false),
		WithResponseBodyLogging(false),
	))
	router.POST("/items", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			t.Fatal(err)
		}
		c.String(http.StatusOK, string(body))
	})

	body := "request-body"
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Body.String() != body {
		t.Fatalf("handler body: got %q, want %q", recorder.Body.String(), body)
	}
	fields := logs.All()[0].ContextMap()
	if _, ok := fields["request-body"]; ok {
		t.Error("request-body should not be logged")
	}
	if _, ok := fields["response-body"]; ok {
		t.Error("response-body should not be logged")
	}
}

func TestNewExcludesDefaultBodyPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, logs := newTestLogger()
	router := gin.New()
	router.Use(New(
		log,
		WithRequestBodyLogging(true),
		WithResponseBodyLogging(true),
		WithExcludedBodyPathPrefixes("/health/"),
	))
	router.POST("/uploads/file", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodPost, "/uploads/file", strings.NewReader("file-content"))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	fields := logs.All()[0].ContextMap()
	if _, ok := fields["request-body"]; ok {
		t.Error("request-body should not be logged for the default excluded path")
	}
	if _, ok := fields["response-body"]; ok {
		t.Error("response-body should not be logged for the default excluded path")
	}
}

func TestNewUsesRequestBodySanitizer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, logs := newTestLogger()
	router := gin.New()
	router.Use(New(
		log,
		WithRequestBodyLogging(true),
		WithRequestBodySanitizer(func(c *gin.Context, body []byte) ([]byte, error) {
			if c.ContentType() != "application/json" {
				t.Errorf("content type: got %q, want %q", c.ContentType(), "application/json")
			}
			return []byte(`{"password":"****"}`), nil
		}),
	))
	router.POST("/items", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			t.Fatal(err)
		}
		c.String(http.StatusOK, string(body))
	})

	body := `{"password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Body.String() != body {
		t.Fatalf("handler body: got %q, want %q", recorder.Body.String(), body)
	}
	if got := logs.All()[0].ContextMap()["request-body"]; got != `{"password":"****"}` {
		t.Errorf("request-body: got %#v, want %q", got, `{"password":"****"}`)
	}
}

func TestNewUsesResponseBodySanitizer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, logs := newTestLogger()
	router := gin.New()
	router.Use(New(
		log,
		WithResponseBodyLogging(true),
		WithResponseBodySanitizer(func(_ *gin.Context, body []byte) ([]byte, error) {
			return []byte(`{"token":"****"}`), nil
		}),
	))
	router.GET("/items", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"token": "secret"})
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	router.ServeHTTP(recorder, req)

	if recorder.Body.String() != `{"token":"secret"}` {
		t.Fatalf("response body: got %q, want %q", recorder.Body.String(), `{"token":"secret"}`)
	}
	if got := logs.All()[0].ContextMap()["response-body"]; got != `{"token":"****"}` {
		t.Errorf("response-body: got %#v, want %q", got, `{"token":"****"}`)
	}
}

func TestNewLogsOneEntryForRequestErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	log, logs := newTestLogger()
	router := gin.New()
	router.Use(New(log))
	router.GET("/items", func(c *gin.Context) {
		_ = c.Error(errors.New("first error"))
		_ = c.Error(errors.New("second error"))
		c.Status(http.StatusInternalServerError)
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	router.ServeHTTP(recorder, req)

	if logs.Len() != 1 {
		t.Fatalf("log entries: got %d, want 1", logs.Len())
	}
	entry := logs.All()[0]
	if entry.Level != zapcore.ErrorLevel {
		t.Fatalf("log level: got %s, want %s", entry.Level, zapcore.ErrorLevel)
	}
	errorsField, ok := entry.ContextMap()["errors"].(string)
	if !ok {
		t.Fatal("errors field should be logged")
	}
	if !strings.Contains(errorsField, "first error") || !strings.Contains(errorsField, "second error") {
		t.Fatalf("errors field: got %q", errorsField)
	}
}
