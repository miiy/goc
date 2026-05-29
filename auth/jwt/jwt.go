package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidSigningMethod = errors.New("jwt auth: invalid signing method")

type UserClaims struct {
	Username string `json:"username"`
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
	return &JWTAuth{options: o}
}

func (j *JWTAuth) CreateClaims(username string) *UserClaims {
	now := time.Now()
	claims := &UserClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.options.Issuer,
			Subject:   username,
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	if j.options.ExpiresIn > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Second * time.Duration(j.options.ExpiresIn)))
	}
	return claims
}

func (j *JWTAuth) CreateTokenByClaims(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.options.Secret))
}

func (j *JWTAuth) CreateToken(username string) (string, error) {
	return j.CreateTokenByClaims(j.CreateClaims(username))
}

func (j *JWTAuth) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrInvalidSigningMethod
	}
	return []byte(j.options.Secret), nil
}

func (j *JWTAuth) ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, j.keyFunc)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}

