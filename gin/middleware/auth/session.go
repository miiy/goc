package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
	"github.com/miiy/goc/gin/sessions"
)

const SessionKeyAuthUser = "goc.auth"

func SessionUser(value any) (*gocauth.AuthenticatedUser, bool) {
	values, ok := value.(map[string]any)
	if !ok {
		return nil, false
	}
	return sessionUserFromMap(values)
}

func SessionAuthenticationMiddleware(redirectPath string) gin.HandlerFunc {
	if redirectPath == "" {
		redirectPath = "/register"
	}

	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user, ok := SessionUser(session.Get(SessionKeyAuthUser))
		if !ok {
			ctx.Redirect(http.StatusFound, redirectPath)
			ctx.Abort()
			return
		}

		setAuthUser(ctx, user)
		ctx.Next()
	}
}

func sessionUserFromMap(values map[string]any) (*gocauth.AuthenticatedUser, bool) {
	username, _ := values["username"].(string)
	if username == "" {
		return nil, false
	}

	id, _ := values["id"].(float64)
	return &gocauth.AuthenticatedUser{ID: int64(id), Username: username}, true
}
