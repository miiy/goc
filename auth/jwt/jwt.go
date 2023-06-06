package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/wire"
	"time"
)

type Options struct {
	Secret    string `yaml:"secret"`
	ExpiresIn int64  `yaml:"expires-in"`
}

type AuthUser struct {
	Id       int64
	Username string
}

type JWTAuth struct {
	Options *Options
}

type Claims struct {
	Username string
	jwt.RegisteredClaims
}

func NewJWTAuth(o *Options) *JWTAuth {
	return &JWTAuth{
		Options: o,
	}
}

func (j *JWTAuth) CreateToken(username string) (string, error) {
	tokenExpireDuration := time.Second * time.Duration(j.Options.ExpiresIn)
	// set our claims
	claims := Claims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpireDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Options.Secret))
}

func (j *JWTAuth) ParseToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Options.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}

func (j *JWTAuth) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err == nil {
		return j.CreateToken(claims.Username)
	}
	return "", err
}

var ProviderSet = wire.NewSet(NewJWTAuth)
