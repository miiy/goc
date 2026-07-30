package jwt

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miiy/goc/auth"
	gocgin "github.com/miiy/goc/gin"
	"google.golang.org/grpc/metadata"
)

var ErrResolverUnconfigured = errors.New("jwt: resolver not configured")

// Resolver validates a raw JWT and returns its trusted session information.
type Resolver func(context.Context, string) (*auth.AuthenticatedSession, error)

type options struct {
	unauthorizedHandler func(*gin.Context, error)
	metadataPropagation bool
}

// Option configures JWT authentication.
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

// WithUnauthorized changes the failure response for invalid JWT credentials.
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

// Authenticate validates a Bearer JWT and stores its session in Gin.
// WithMetadataPropagation optionally forwards trusted identity downstream.
func Authenticate(resolve Resolver, opts ...Option) gin.HandlerFunc {
	o := newOptions(opts...)

	return func(c *gin.Context) {
		token, err := gocgin.ExtractToken(c)
		if err != nil {
			o.handleUnauthorized(c, err)
			return
		}
		if resolve == nil {
			o.handleUnauthorized(c, ErrResolverUnconfigured)
			return
		}
		session, err := resolve(c.Request.Context(), token)
		if err != nil {
			o.handleUnauthorized(c, err)
			return
		}

		c.Request = c.Request.WithContext(auth.InjectAuthenticatedSession(c.Request.Context(), session))
		if o.metadataPropagation {
			requestContext := metadata.AppendToOutgoingContext(
				c.Request.Context(),
				auth.AuthenticatedUserIDMetadataKey, session.User.ID,
				auth.AuthenticatedUsernameMetadataKey, session.User.Username,
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
