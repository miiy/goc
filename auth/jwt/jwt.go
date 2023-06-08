package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Options struct {
	Secret    string `yaml:"secret"`
	Issuer    string `yaml:"issuer"`
	ExpiresIn int64  `yaml:"expiresIn"`
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

func (j *JWTAuth) CreateClaims(username string) *Claims {
	ep := time.Second * time.Duration(j.Options.ExpiresIn)
	// set our claims
	return &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			Issuer:    j.Options.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),
		},
	}
}

func (j *JWTAuth) CreateTokenByClaims(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Options.Secret))
}

func (j *JWTAuth) CreateToken(username string) (string, error) {
	c := j.CreateClaims(username)
	return j.CreateTokenByClaims(c)
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
	if err != nil {
		return "", err
	}
	return j.CreateToken(claims.Username)
}
