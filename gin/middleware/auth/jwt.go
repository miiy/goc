package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gocauth "github.com/miiy/goc/auth"
	gocjwt "github.com/miiy/goc/auth/jwt"
	jwtrequest "github.com/golang-jwt/jwt/v5/request"
)

func JWTAuthenticationMiddleware(jwtAuth *gocjwt.JWTAuth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := jwtrequest.AuthorizationHeaderExtractor.ExtractToken(ctx.Request)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, err := jwtAuth.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		setAuthUser(ctx, &gocauth.AuthenticatedUser{
			Username: claims.Username,
		})
		ctx.Next()
	}
}
