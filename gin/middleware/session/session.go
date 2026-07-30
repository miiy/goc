package session

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin/sessions"
	"google.golang.org/grpc/metadata"
)

var ErrSessionResolverUnconfigured = errors.New("session: session resolver not configured")
var ErrSessionIDMissing = errors.New("session: session id missing")

// Resolver validates a session ID and returns its trusted session information.
type Resolver func(context.Context, string) (*auth.AuthenticatedSession, error)

type options struct {
	unauthorizedHandler func(*gin.Context, error)
	metadataPropagation bool
}

// Option configures session authentication.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(options *options) {
	f(options)
}

func newOptions(opts ...Option) *options {
	config := &options{}
	for _, option := range opts {
		option.apply(config)
	}
	return config
}

// WithUnauthorized changes the failure response for invalid session credentials.
func WithUnauthorized(handler func(*gin.Context, error)) Option {
	return optionFunc(func(options *options) {
		if handler != nil {
			options.unauthorizedHandler = handler
		}
	})
}

// WithMetadataPropagation forwards the authenticated identity to downstream
// gRPC calls. Authentication context injection does not imply propagation.
func WithMetadataPropagation() Option {
	return optionFunc(func(options *options) {
		options.metadataPropagation = true
	})
}

// Authenticate resolves a Cookie Session and stores it in Gin.
// WithMetadataPropagation optionally forwards trusted identity downstream.
func Authenticate(resolve Resolver, opts ...Option) gin.HandlerFunc {
	o := newOptions(opts...)

	return func(c *gin.Context) {
		sid := sessions.SID(c)
		if sid == "" {
			o.handleUnauthorized(c, ErrSessionIDMissing)
			return
		}
		if resolve == nil {
			o.handleUnauthorized(c, ErrSessionResolverUnconfigured)
			return
		}
		authenticatedSession, err := resolve(c.Request.Context(), sid)
		if err != nil {
			o.handleUnauthorized(c, err)
			return
		}

		c.Request = c.Request.WithContext(auth.InjectAuthenticatedSession(c.Request.Context(), authenticatedSession))
		if o.metadataPropagation {
			requestContext := metadata.AppendToOutgoingContext(
				c.Request.Context(),
				auth.AuthenticatedUserIDMetadataKey, authenticatedSession.User.ID,
				auth.AuthenticatedUsernameMetadataKey, authenticatedSession.User.Username,
			)
			c.Request = c.Request.WithContext(requestContext)
		}
		c.Next()
	}
}

func (o *options) handleUnauthorized(c *gin.Context, err error) {
	if o.unauthorizedHandler != nil {
		c.Abort()
		o.unauthorizedHandler(c, err)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
