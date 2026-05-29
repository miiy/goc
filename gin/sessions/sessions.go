package sessions

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type Store = sessions.Store
type Session = sessions.Session

// NewCookieStore creates a cookie-based session store
func NewCookieStore(secret string) sessions.Store {
	return cookie.NewStore([]byte(secret))
}

// NewRedisStore creates a Redis-backed session store.
func NewRedisStore(size int, network, address, password string, keyPairs ...[]byte) (sessions.Store, error) {
	return redis.NewStore(size, network, address, "", password, keyPairs...)
}

// Middleware returns a session middleware
func Middleware(name string, store sessions.Store) gin.HandlerFunc {
	return sessions.Sessions(name, store)
}

// Default returns the default session for the context
func Default(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}
