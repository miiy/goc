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
