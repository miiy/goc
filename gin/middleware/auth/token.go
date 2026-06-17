package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/miiy/goc/auth"
	"google.golang.org/grpc/metadata"
)

// BearerToken extracts the bearer token from the Authorization header of the request.
// It returns the token and the extraction error (request.ErrNoTokenInRequest when the
// header is missing or not a Bearer token).
func BearerToken(c *gin.Context) (string, error) {
	if c == nil || c.Request == nil {
		return "", request.ErrNoTokenInRequest
	}
	return request.BearerExtractor{}.ExtractToken(c.Request)
}

// Authenticator resolves a bearer token to an authenticated user. It is the
// customization point of AuthenticationMiddleware; callers wire it to whatever
// verifies the token (e.g. an auth-service RPC).
type Authenticator func(ctx context.Context, token string) (*auth.AuthenticatedUser, error)

// AuthenticationMiddleware returns a gin middleware that authenticates each request
// by resolving its bearer token via authenticator. On success the user is injected
// into the request context; use WithMetadataPropagation to also forward the identity
// to downstream gRPC calls via metadata (x-auth-user-id, x-auth-username). On failure
// the unauthorized handler runs.
func AuthenticationMiddleware(authenticator Authenticator, opts ...Option) gin.HandlerFunc {
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}

	unauthorized := o.unauthorizedHandler
	if unauthorized == nil {
		unauthorized = func(ctx *gin.Context, _ string) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}

	return func(ctx *gin.Context) {
		token, err := BearerToken(ctx)
		if err != nil {
			unauthorized(ctx, err.Error())
			return
		}

		user, err := authenticator(ctx.Request.Context(), token)
		if err != nil {
			unauthorized(ctx, err.Error())
			return
		}
		if user == nil {
			unauthorized(ctx, "authenticated user not found")
			return
		}

		c := auth.InjectAuthenticatedUser(ctx.Request.Context(), user)
		if o.metadataPropagation {
			c = metadata.AppendToOutgoingContext(c,
				auth.AuthenticatedUserIDMetadataKey, strconv.FormatInt(user.ID, 10),
				auth.AuthenticatedUsernameMetadataKey, user.Username,
			)
		}
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}

// Option configures AuthenticationMiddleware behavior.
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

// WithUnauthorized sets a custom handler for 401 responses.
func WithUnauthorized(fn func(ctx *gin.Context, message string)) Option {
	return optionFunc(func(opts *options) {
		opts.unauthorizedHandler = fn
	})
}

// WithMetadataPropagation forwards the authenticated identity to downstream gRPC
// calls via metadata (x-auth-user-id, x-auth-username). Enable in services that fan
// out to gRPC backends (e.g. an API gateway); leave off for plain HTTP services.
func WithMetadataPropagation() Option {
	return optionFunc(func(opts *options) {
		opts.metadataPropagation = true
	})
}

func setAuthUser(ctx *gin.Context, user *auth.AuthenticatedUser) {
	if ctx == nil || ctx.Request == nil {
		return
	}
	ctx.Request = ctx.Request.WithContext(auth.InjectAuthenticatedUser(ctx.Request.Context(), user))
}

// GetAuthUser retrieves the authenticated user from the request context.
func GetAuthUser(ctx *gin.Context) (*auth.AuthenticatedUser, bool) {
	if ctx == nil || ctx.Request == nil {
		return nil, false
	}
	user, err := auth.ExtractAuthenticatedUser(ctx.Request.Context())
	return user, err == nil
}

// GetAuthUserID retrieves the authenticated user ID from the request context.
func GetAuthUserID(ctx *gin.Context) (int64, bool) {
	user, ok := GetAuthUser(ctx)
	if !ok || user.ID <= 0 {
		return 0, false
	}
	return user.ID, true
}
