package auth

import (
	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
)

// Session retrieves the authenticated session from the request context.
func Session(ctx *gin.Context) (*gocauth.AuthenticatedSession, bool) {
	return gocauth.Session(ctx.Request.Context())
}

// User retrieves the user from the authenticated session.
func User(ctx *gin.Context) (*gocauth.AuthenticatedUser, bool) {
	return gocauth.User(ctx.Request.Context())
}

// UserID retrieves the authenticated user ID from the request context.
func UserID(ctx *gin.Context) (string, bool) {
	return gocauth.UserID(ctx.Request.Context())
}

// UserInt64ID retrieves the authenticated user ID as int64 from the request
// context.
func UserInt64ID(ctx *gin.Context) (int64, bool) {
	return gocauth.UserInt64ID(ctx.Request.Context())
}
