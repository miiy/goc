package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/miiy/goc/gin"
	"go.uber.org/zap"
)

func TestNewAppliesHTTPTimeoutOptions(t *testing.T) {
	server := New(
		WithReadHeaderTimeout(time.Second),
		WithReadTimeout(2*time.Second),
		WithWriteTimeout(3*time.Second),
		WithIdleTimeout(4*time.Second),
	)

	if server.readHeaderTimeout != time.Second || server.readTimeout != 2*time.Second || server.writeTimeout != 3*time.Second || server.idleTimeout != 4*time.Second {
		t.Fatalf("server timeouts = %+v", server)
	}
}

func TestHandlerServesRegisteredRoutes(t *testing.T) {
	server := New(WithLogger(zap.NewNop()))
	server.RegisterRouter(func(r *gin.Engine) {
		r.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "hello"})
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	server.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRunContextReturnsListenError(t *testing.T) {
	server := New(WithLogger(zap.NewNop()))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.RunContext(ctx, "127.0.0.1:-1"); err == nil {
		t.Fatal("expected listen error")
	}
}

func TestRunContextPanicsOnNilContext(t *testing.T) {
	server := New(WithLogger(zap.NewNop()))

	requirePanic(t, func() {
		server.RunContext(nil, "127.0.0.1:-1")
	})
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
