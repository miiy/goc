package sessions

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/boj/redistore"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type Store = sessions.Store
type Session = sessions.Session
type Options = sessions.Options

const DefaultSessionCookieName = "__Host-session"

var ErrUnsupportedJSONSessionStore = errors.New("sessions: JSON serializer is only supported by Redis store")
var ErrUnsupportedRedisSessionStore = errors.New("sessions: operation is only supported by Redis store")

// SessionManager keeps the configured Session cookie name, store, and cookie
// options together so applications do not pass them through every Session
// operation.
type SessionManager struct {
	cookieName string
	store      sessions.Store
	options    Options
}

type sessionManagerOptions struct {
	cookieName string
}

// SessionManagerOption configures a SessionManager.
type SessionManagerOption interface {
	apply(*sessionManagerOptions)
}

type sessionManagerOptionFunc func(*sessionManagerOptions)

func (f sessionManagerOptionFunc) apply(options *sessionManagerOptions) {
	f(options)
}

// WithCookieName configures the Cookie name used by the SessionManager.
func WithCookieName(name string) SessionManagerOption {
	return sessionManagerOptionFunc(func(options *sessionManagerOptions) {
		options.cookieName = sessionCookieName(name)
	})
}

// NewSessionManager creates a manager backed by a goc Session store.
func NewSessionManager(store sessions.Store, options Options, opts ...SessionManagerOption) *SessionManager {
	managerOptions := sessionManagerOptions{cookieName: DefaultSessionCookieName}
	for _, opt := range opts {
		if opt != nil {
			opt.apply(&managerOptions)
		}
	}
	if store != nil {
		store.Options(options)
	}
	return &SessionManager{
		cookieName: managerOptions.cookieName,
		store:      store,
		options:    options,
	}
}

// CookieName returns the Cookie name used by this manager.
func (m *SessionManager) CookieName() string {
	return m.cookieName
}

// Options returns the default cookie options used by this manager.
func (m *SessionManager) Options() Options {
	return m.options
}

// Middleware loads this manager's Session into Gin.
func (m *SessionManager) Middleware() gin.HandlerFunc {
	return Middleware(m.cookieName, m.store)
}

// Renew saves a fresh Session with the provided ID and values.
func (m *SessionManager) Renew(c *gin.Context, id string, values map[any]any) error {
	return Renew(c, m.store, m.cookieName, id, values)
}

// RenewWithOptions saves a fresh Session with explicit cookie/store options.
func (m *SessionManager) RenewWithOptions(c *gin.Context, id string, values map[any]any, options Options) error {
	return RenewWithOptions(c, m.store, m.cookieName, id, values, &options)
}

// Clear deletes the current server-side Session and expires its Cookie.
func (m *SessionManager) Clear(c *gin.Context) error {
	return Clear(c, m.cookieName, m.options)
}

// NewCookieStore creates a cookie-based session store
func NewCookieStore(secret string) sessions.Store {
	return cookie.NewStore([]byte(secret))
}

// NewRedisStore creates a Redis-backed session store.
func NewRedisStore(size int, network, address string, db int, password string, keyPairs ...[]byte) (sessions.Store, error) {
	return redis.NewStoreWithDB(size, network, address, "", password, strconv.Itoa(db), keyPairs...)
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

// Renew saves a fresh Session with the provided ID and values while ignoring
// any existing cookie for the same name on the request.
func Renew(c *gin.Context, store sessions.Store, name string, id string, values map[any]any) error {
	return RenewWithOptions(c, store, name, id, values, nil)
}

// RenewWithOptions saves a fresh Session with explicit cookie/store options.
func RenewWithOptions(c *gin.Context, store sessions.Store, name string, id string, values map[any]any, options *Options) error {
	name = sessionCookieName(name)
	current, err := store.New(requestWithoutCookie(c.Request, name), name)
	if err != nil {
		return err
	}
	current.ID = id
	if options != nil {
		current.Options = options.ToGorillaOptions()
	}
	for key, value := range values {
		current.Values[key] = value
	}
	return store.Save(c.Request, c.Writer, current)
}

// Clear deletes the current Session values and expires its browser cookie.
func Clear(c *gin.Context, name string, options Options) error {
	name = sessionCookieName(name)
	current := Default(c)
	current.Clear()
	options.MaxAge = -1
	current.Options(options)
	if err := current.Save(); err != nil {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     options.Path,
			Domain:   options.Domain,
			MaxAge:   -1,
			Secure:   options.Secure,
			HttpOnly: options.HttpOnly,
			SameSite: options.SameSite,
		})
		return err
	}
	return nil
}

// Middleware returns a session middleware
func Middleware(name string, store sessions.Store) gin.HandlerFunc {
	return sessions.Sessions(sessionCookieName(name), store)
}

func sessionCookieName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return DefaultSessionCookieName
	}
	return name
}

// Default returns the default session for the context
func Default(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}

// SID returns the current Session ID from the context.
func SID(c *gin.Context) string {
	value, ok := c.Get(sessions.DefaultKey)
	if !ok {
		return ""
	}
	current, ok := value.(sessions.Session)
	if !ok {
		return ""
	}
	// The Session ID is loaded from the backing store, not from the values map.
	sid := current.ID()
	return sid
}

func requestWithoutCookie(request *http.Request, name string) *http.Request {
	cloned := request.Clone(request.Context())
	cloned.Header.Del("Cookie")
	for _, cookie := range request.Cookies() {
		if cookie.Name != name {
			cloned.AddCookie(cookie)
		}
	}
	return cloned
}
