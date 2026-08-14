package zap

import (
	"strings"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Ginzap(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: func(c *gin.Context) (fields []zapcore.Field) {
			// log request ID
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request-id", requestID))
			}

			fields = append(fields,
				zap.String("host", c.Request.Host),
				zap.String("remote-addr", c.Request.RemoteAddr),
				zap.String("authorization", maskAuthorization(c.Request.Header.Get("Authorization"))),
				zap.String("full-path", c.FullPath()),
				zap.String("sign", c.Request.Header.Get("sign")),
				zap.String("ts", c.Request.Header.Get("ts")),
				zap.Int("size", c.Writer.Size()),
			)

			return
		},
	})
}

func RecoveryWithZap(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return ginzap.RecoveryWithZap(logger, stack)
}

// maskAuthorization returns "Bearer ****" for Bearer tokens, or "****" for any non-empty value.
func maskAuthorization(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 2 {
		return parts[0] + " ****"
	}
	return "****"
}
