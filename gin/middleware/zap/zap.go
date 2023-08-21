package zap

import (
	"bytes"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"time"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	if count, err := w.body.Write(b); err != nil {
		return count, err
	}
	return w.ResponseWriter.Write(b)
}

func ResponseBodyBuffer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// response body
		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w
		// next
		c.Next()
	}
}

func Ginzap(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: func(c *gin.Context) (fields []zapcore.Field) {
			// log request ID
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request-id", requestID))
			}

			fields = append(fields, zap.String("host", c.Request.Host))
			fields = append(fields, zap.String("remote-addr", c.Request.RemoteAddr))
			fields = append(fields, zap.String("authorization", c.Request.Header.Get("Authorization")))
			fields = append(fields, zap.String("full-path", c.FullPath()))
			fields = append(fields, zap.Any("request-header", c.Request.Header))

			// request body
			if strings.HasPrefix(c.Request.URL.Path, "/uploads/") {
				return
			}
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
			fields = append(fields, zap.String("request-body", string(body)))

			w, ok := c.Writer.(*responseBodyWriter)
			if ok {
				fields = append(fields, zap.String("response-body", w.body.String()))
			}

			return
		},
	})
}

func RecoveryWithZap(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return ginzap.RecoveryWithZap(logger, stack)
}
