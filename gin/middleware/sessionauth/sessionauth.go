package sessionauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin/authctx"
	"github.com/miiy/goc/gin/sessions"
)

const SessionKeyUser = "goc.auth.user"

func sessionUser(ctx *gin.Context, sessionKey string) (*gocauth.AuthenticatedUser, bool) {
	values, ok := sessions.Default(ctx).Get(sessionKey).(map[string]any)
	if !ok {
		return nil, false
	}

	username, _ := values["username"].(string)
	if username == "" {
		return nil, false
	}

	idText, _ := values["id"].(string)
	if idText == "" {
		return nil, false
	}

	return &gocauth.AuthenticatedUser{ID: idText, Username: username}, true
}

// LoadSessionUser bridges the session user into the request context. It never
// rejects anonymous requests.
func LoadSessionUser(opts ...Option) gin.HandlerFunc {
	o := newOptions(opts...)

	return func(ctx *gin.Context) {
		if user, ok := sessionUser(ctx, o.sessionKey); ok {
			authctx.SetUser(ctx, user)
		}
		ctx.Next()
	}
}

// Authenticate ensures the request has an authenticated session user.
func Authenticate(opts ...Option) gin.HandlerFunc {
	o := newOptions(opts...)

	return func(ctx *gin.Context) {
		user, ok := sessionUser(ctx, o.sessionKey)
		if !ok {
			o.handleUnauthorized(ctx)
			return
		}
		authctx.SetUser(ctx, user)
		ctx.Next()
	}
}

// Option configures session auth middleware.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(opts *options) {
	f(opts)
}

type options struct {
	redirectPath        string
	unauthorizedHandler func(ctx *gin.Context)
	sessionKey          string
}

func newOptions(opts ...Option) *options {
	o := &options{
		sessionKey: SessionKeyUser,
	}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func (o *options) handleUnauthorized(ctx *gin.Context) {
	if o.unauthorizedHandler != nil {
		o.unauthorizedHandler(ctx)
		return
	}
	if o.redirectPath != "" {
		ctx.Redirect(http.StatusFound, o.redirectPath)
		ctx.Abort()
		return
	}
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

// WithRedirect makes Authenticate redirect unauthenticated requests.
func WithRedirect(path string) Option {
	return optionFunc(func(opts *options) {
		opts.redirectPath = path
	})
}

// WithUnauthorized sets a custom unauthenticated request handler.
func WithUnauthorized(fn func(ctx *gin.Context)) Option {
	return optionFunc(func(opts *options) {
		opts.unauthorizedHandler = fn
	})
}

// WithSessionKey changes the session value key used to store the authenticated
// user.
func WithSessionKey(key string) Option {
	return optionFunc(func(opts *options) {
		if key != "" {
			opts.sessionKey = key
		}
	})
}
