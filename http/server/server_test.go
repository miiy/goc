package server

import (
	"github.com/miiy/goc/gin"
	ginzap "github.com/miiy/goc/gin/middleware/zap"
	"go.uber.org/zap"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	server := New(WithLogger(zap.NewExample()))
	server.Run()
}

func TestWithDebug(t *testing.T) {
	server := New(WithDebug(), WithLogger(zap.NewExample()))
	server.Run()
}

func TestZap(t *testing.T) {
	logger := zap.NewExample()
	server := New(WithDebug(), WithLogger(logger))
	server.Use(ginzap.ResponseBodyBuffer())
	server.Use(ginzap.Ginzap(logger))
	server.Use(ginzap.RecoveryWithZap(logger, true))
	server.RegisterRouter(func(r *gin.Engine) {
		r.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello",
			})
		})
	})
	go func() {
		time.Sleep(1 * time.Second)
		c := http.Client{}
		c.Get("http://127.0.0.1:8080/hello?p1=1&p2=2")
	}()
	server.Run()
}
