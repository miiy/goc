package auth

import (
	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
)

const AuthUserKey = "auth"

func setAuthUser(ctx *gin.Context, user *gocauth.AuthenticatedUser) {
	ctx.Set(AuthUserKey, user)
	ctx.Request = ctx.Request.WithContext(gocauth.InjectAuthenticatedUser(ctx.Request.Context(), user))
}
