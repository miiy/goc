package jwtauth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
	gocauth "github.com/miiy/goc/auth"
	"google.golang.org/grpc/metadata"
)

// Token extracts the bearer token from the Authorization header.
func Token(ctx *gin.Context) (string, error) {
	if ctx == nil || ctx.Request == nil {
		return "", request.ErrNoTokenInRequest
	}
	return request.BearerExtractor{}.ExtractToken(ctx.Request)
}

// UserResolver resolves a JWT bearer token to an authenticated user.
type UserResolver func(ctx context.Context, token string) (*gocauth.AuthenticatedUser, error)

// Authenticate authenticates a JWT bearer token with an application resolver and
// stores the resolved user in the request context.
func Authenticate(resolveUser UserResolver, opts ...Option) gin.HandlerFunc {
	o := newOptions(opts...)

	return func(ctx *gin.Context) {
		token, err := Token(ctx)
		if err != nil {
			o.handleUnauthorized(ctx, err.Error())
			return
		}

		if resolveUser == nil {
			o.handleUnauthorized(ctx, "jwt user resolver not configured")
			return
		}

		user, err := resolveUser(ctx.Request.Context(), token)
		if err != nil {
			o.handleUnauthorized(ctx, err.Error())
			return
		}
		if user == nil {
			o.handleUnauthorized(ctx, "authenticated user not found")
			return
		}

		c := gocauth.InjectAuthenticatedUser(ctx.Request.Context(), user)
		if o.metadataPropagation {
			c = metadata.AppendToOutgoingContext(c,
				gocauth.AuthenticatedUserIDMetadataKey, user.ID,
				gocauth.AuthenticatedUsernameMetadataKey, user.Username,
			)
		}
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}

// Option configures Authenticate.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(opts *options) {
	f(opts)
}

type options struct {
	unauthorizedHandler func(ctx *gin.Context, message string)
	metadataPropagation bool
}

func newOptions(opts ...Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

func (o *options) handleUnauthorized(ctx *gin.Context, message string) {
	if o.unauthorizedHandler != nil {
		o.unauthorizedHandler(ctx, message)
		return
	}
	ctx.AbortWithStatus(http.StatusUnauthorized)
}

// WithUnauthorized sets a custom 401 handler.
func WithUnauthorized(fn func(ctx *gin.Context, message string)) Option {
	return optionFunc(func(opts *options) {
		opts.unauthorizedHandler = fn
	})
}

// WithMetadataPropagation forwards the authenticated identity to downstream gRPC
// calls via metadata.
func WithMetadataPropagation() Option {
	return optionFunc(func(opts *options) {
		opts.metadataPropagation = true
	})
}
