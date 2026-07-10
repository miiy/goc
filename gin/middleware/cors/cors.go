package cors

import (
	gincors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config = gincors.Config

var DefaultAllowOrigins = []string{
	"http://localhost:3000",
	"http://localhost:8080",
	"http://127.0.0.1:3000",
	"http://127.0.0.1:8080",
}

func DefaultConfig() Config {
	return gincors.DefaultConfig()
}

func New(allowOrigins ...string) gin.HandlerFunc {
	if len(allowOrigins) == 0 {
		allowOrigins = DefaultAllowOrigins
	}

	config := DefaultConfig()
	config.AllowOrigins = append([]string(nil), allowOrigins...)
	config.AddAllowHeaders("Authorization")
	return NewWithConfig(config)
}

func NewWithConfig(config Config) gin.HandlerFunc {
	return gincors.New(config)
}
