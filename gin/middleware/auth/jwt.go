package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/miiy/goc/auth"
)

// JWTAuthenticationMiddleware returns a gin middleware that validates JWT Bearer tokens.
func JWTAuthenticationMiddleware(jwtAuth *auth.JWTAuth, opts ...Option) gin.HandlerFunc {
	o := &jwtOptions{}
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
		tokenString, err := request.BearerExtractor{}.ExtractToken(ctx.Request)
		if err != nil {
			unauthorized(ctx, err.Error())
			return
		}

		claims, err := jwtAuth.ParseToken(tokenString)
		if err != nil {
			unauthorized(ctx, "invalid auth token")
			return
		}

		var user *auth.AuthenticatedUser
		if o.userResolver != nil {
			user, err = o.userResolver(ctx.Request.Context(), claims, tokenString)
			if err != nil {
				unauthorized(ctx, err.Error())
				return
			}
			if user == nil {
				unauthorized(ctx, "invalid auth user")
				return
			}
		} else {
			user = &auth.AuthenticatedUser{Username: claims.Username}
		}
		setAuthUser(ctx, user)

		if o.afterAuth != nil {
			o.afterAuth(ctx, claims, tokenString)
		}

		ctx.Next()
	}
}

// --- Options ---

// Option configures JWT middleware behavior.
type Option interface {
	apply(*jwtOptions)
}

// UserResolver resolves the current authenticated user after JWT validation.
type UserResolver func(ctx context.Context, claims *auth.UserClaims, tokenString string) (*auth.AuthenticatedUser, error)

type optionFunc func(*jwtOptions)

func (f optionFunc) apply(opts *jwtOptions) {
	f(opts)
}

type jwtOptions struct {
	unauthorizedHandler func(ctx *gin.Context, message string)
	afterAuth           func(ctx *gin.Context, claims *auth.UserClaims, tokenString string)
	userResolver        UserResolver
}

// WithUnauthorized sets a custom handler for 401 responses.
func WithUnauthorized(fn func(ctx *gin.Context, message string)) Option {
	return optionFunc(func(opts *jwtOptions) {
		opts.unauthorizedHandler = fn
	})
}

// WithAfterAuth registers a callback invoked after successful authentication.
func WithAfterAuth(fn func(ctx *gin.Context, claims *auth.UserClaims, tokenString string)) Option {
	return optionFunc(func(opts *jwtOptions) {
		opts.afterAuth = fn
	})
}

// WithUserResolver sets a resolver used to re-check the current user after JWT validation.
func WithUserResolver(fn UserResolver) Option {
	return optionFunc(func(opts *jwtOptions) {
		opts.userResolver = fn
	})
}

// --- Context helpers ---

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
