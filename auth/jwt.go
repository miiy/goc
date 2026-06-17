package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth handles JWT token creation and parsing.
type JWTAuth struct {
	options *Options
}

// Options configures JWTAuth.
type Options struct {
	Secret    string `yaml:"secret"`
	Issuer    string `yaml:"issuer"`
	ExpiresIn int64  `yaml:"expiresIn"`
}

// UserClaims represents JWT claims with a user id and username.
type UserClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var ErrInvalidSigningMethod = errors.New("jwt auth: invalid signing method")

// NewJWTAuth creates a new JWTAuth instance.
func NewJWTAuth(o *Options) *JWTAuth {
	return &JWTAuth{options: o}
}

// CreateClaims builds UserClaims for the given user id and username.
func (j *JWTAuth) CreateClaims(userID int64, username string) *UserClaims {
	now := time.Now()
	claims := &UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.options.Issuer,
			Subject:   strconv.FormatInt(userID, 10),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	if j.options.ExpiresIn > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Second * time.Duration(j.options.ExpiresIn)))
	}
	return claims
}

// CreateTokenByClaims creates a signed JWT token from custom claims.
func (j *JWTAuth) CreateTokenByClaims(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.options.Secret))
}

// CreateToken creates a signed JWT token for the given user id and username.
func (j *JWTAuth) CreateToken(userID int64, username string) (string, error) {
	return j.CreateTokenByClaims(j.CreateClaims(userID, username))
}

func (j *JWTAuth) keyFunc(token *jwt.Token) (interface{}, error) {
	// Only HS256 is accepted; tokens signed with any other alg (HS384/HS512,
	// RS256, "none", etc.) are rejected to prevent algorithm confusion.
	if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, ErrInvalidSigningMethod
	}
	return []byte(j.options.Secret), nil
}

// ParseToken parses and validates a JWT token string.
func (j *JWTAuth) ParseToken(tokenString string) (*UserClaims, error) {
	var opts []jwt.ParserOption
	opts = append(opts, jwt.WithExpirationRequired())
	if j.options.Issuer != "" {
		opts = append(opts, jwt.WithIssuer(j.options.Issuer))
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, j.keyFunc, opts...)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}
