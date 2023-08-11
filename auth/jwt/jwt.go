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

type JWTAuth struct {
	options *Options
}

func NewJWTAuth(o *Options) *JWTAuth {
	return &JWTAuth{
		options: o,
	}
}

func (j *JWTAuth) CreateClaims(subject string) jwt.Claims {
	ep := time.Second * time.Duration(j.options.ExpiresIn)
	now := time.Now()
	// set claims
	return &jwt.RegisteredClaims{
		Issuer:    j.options.Issuer,
		Subject:   subject,
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ep)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
}

func (j *JWTAuth) CreateTokenByClaims(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.options.Secret)
}

func (j *JWTAuth) CreateToken(subject string) (string, error) {
	c := j.CreateClaims(subject)
	return j.CreateTokenByClaims(c)
}

func (j *JWTAuth) Parse(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.options.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return token.Claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}
