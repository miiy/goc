package authctx

import (
	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
)

// SetUser stores the authenticated user in the request context.
func SetUser(ctx *gin.Context, user *gocauth.AuthenticatedUser) {
	if ctx == nil || ctx.Request == nil {
		return
	}
	ctx.Request = ctx.Request.WithContext(gocauth.InjectAuthenticatedUser(ctx.Request.Context(), user))
}

// CurrentUser retrieves the authenticated user from the request context.
func CurrentUser(ctx *gin.Context) (*gocauth.AuthenticatedUser, bool) {
	if ctx == nil || ctx.Request == nil {
		return nil, false
	}
	user, err := gocauth.ExtractAuthenticatedUser(ctx.Request.Context())
	return user, err == nil
}

// CurrentUserID retrieves the authenticated user ID from the request context.
func CurrentUserID(ctx *gin.Context) (string, bool) {
	user, ok := CurrentUser(ctx)
	if !ok || user.ID == "" {
		return "", false
	}
	return user.ID, true
}

// CurrentUserInt64ID retrieves the authenticated user ID as int64 from the request
// context.
func CurrentUserInt64ID(ctx *gin.Context) (int64, bool) {
	user, ok := CurrentUser(ctx)
	if !ok {
		return 0, false
	}
	id, err := user.Int64ID()
	if err != nil {
		return 0, false
	}
	return id, true
}
