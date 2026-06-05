package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var DefaultAllowOrigins = []string{
	"http://localhost:3000",
	"http://localhost:8080",
	"http://127.0.0.1:3000",
	"http://127.0.0.1:8080",
}

func New(allowOrigins ...string) gin.HandlerFunc {
	if len(allowOrigins) == 0 {
		allowOrigins = DefaultAllowOrigins
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = append([]string(nil), allowOrigins...)
	config.AddAllowHeaders("Authorization")
	return NewWithConfig(config)
}

func NewWithConfig(config cors.Config) gin.HandlerFunc {
	return cors.New(config)
}
