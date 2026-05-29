package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin/sessions"
)

const SessionKeyAuthUser = "goc.auth"

func SessionAuthenticationMiddleware(redirectPath string) gin.HandlerFunc {
	if redirectPath == "" {
		redirectPath = "/register"
	}

	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user, ok := session.Get(SessionKeyAuthUser).(*gocauth.AuthenticatedUser)
		if !ok || user == nil {
			ctx.Redirect(http.StatusFound, redirectPath)
			ctx.Abort()
			return
		}

		setAuthUser(ctx, user)
		ctx.Next()
	}
}
