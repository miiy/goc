package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserClaims struct {
	Username string
	jwt.RegisteredClaims
}

type Options struct {
	Secret    string `yaml:"secret"`
	Issuer    string `yaml:"issuer"`
	ExpiresIn int64  `yaml:"expiresIn"`
}

type JWTAuth struct {
	options *Options
}

func NewJWTAuth(o *Options) *JWTAuth {
	return &JWTAuth{
		options: o,
	}
}

func (j *JWTAuth) CreateClaims(username string) *UserClaims {
	ep := time.Second * time.Duration(j.options.ExpiresIn)
	now := time.Now()
	// set claims
	return &UserClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.options.Issuer,
			Subject:   username,
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ep)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
}

func (j *JWTAuth) CreateTokenByClaims(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.options.Secret))
}

func (j *JWTAuth) CreateToken(username string) (string, error) {
	c := j.CreateClaims(username)
	return j.CreateTokenByClaims(c)
}

func (j *JWTAuth) ParseToken(tokenString string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.options.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}
