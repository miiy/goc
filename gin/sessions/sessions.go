package sessions

import (
	"errors"

	"github.com/boj/redistore"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type Store = sessions.Store
type Session = sessions.Session

var ErrUnsupportedJSONSessionStore = errors.New("sessions: JSON serializer is only supported by Redis store")
var ErrUnsupportedRedisSessionStore = errors.New("sessions: operation is only supported by Redis store")

// NewCookieStore creates a cookie-based session store
func NewCookieStore(secret string) sessions.Store {
	return cookie.NewStore([]byte(secret))
}

// NewRedisStore creates a Redis-backed session store.
func NewRedisStore(size int, network, address, password string, keyPairs ...[]byte) (sessions.Store, error) {
	return redis.NewStore(size, network, address, "", password, keyPairs...)
}

// UseJSONSerializer configures a Redis-backed session store to serialize
// the whole session values map as JSON instead of gob.
func UseJSONSerializer(store sessions.Store) error {
	rediStore, err := redis.GetRedisStore(store)
	if err != nil {
		return ErrUnsupportedJSONSessionStore
	}
	rediStore.SetSerializer(redistore.JSONSerializer{})
	return nil
}

// SetMaxAge configures the Redis-backed session max age in seconds. The value
// controls both the browser cookie max age and the Redis key TTL.
func SetMaxAge(store sessions.Store, maxAge int) error {
	rediStore, err := redis.GetRedisStore(store)
	if err != nil {
		return ErrUnsupportedRedisSessionStore
	}
	rediStore.SetMaxAge(maxAge)
	return nil
}

// Middleware returns a session middleware
func Middleware(name string, store sessions.Store) gin.HandlerFunc {
	return sessions.Sessions(name, store)
}

// Default returns the default session for the context
func Default(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}
