package auth

/*
func JWTAuthenticationMiddleware(jwtAuth *jwt.JWTAuth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := jwt.HeaderToken(ctx.Request)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, err := jwtAuth.ParseToken(token)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		user := &AuthUser{
			Username: claims.Username,
		}
		ctx.Set("auth", user)
		ctx.Next()
	}
}
*/
