package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenType = "at+jwt"
)

var (
	ErrInvalidConfig = errors.New("auth: invalid JWT configuration")
	ErrInvalidClaims = errors.New("auth: invalid JWT claims")
	ErrInvalidToken  = errors.New("auth: invalid JWT")

	errInvalidTokenType = errors.New("invalid token type")
)

// VerifierConfig configures token verification.
type VerifierConfig struct {
	Issuer       string
	Audience     string
	PublicKeyPEM string
	Leeway       time.Duration
}

// IssuerConfig configures token issuance. PrivateKeyPEM must contain a PKCS1
// or PKCS8 RSA private key.
type IssuerConfig struct {
	Issuer              string
	Audience            string
	PrivateKeyPEM       string
	AccessTokenLifetime time.Duration
}

// Claims contains the identity and session data carried by an access token.
// Subject is the stable user identifier; Username is display-only.
type Claims struct {
	Username   string `json:"username"`
	SessionID  string `json:"sid"`
	ClientID   string `json:"client_id"`
	ClientType string `json:"client_type"`
	jwt.RegisteredClaims
}

// Verifier verifies application access tokens with a trusted public key.
type Verifier struct {
	issuer    string
	audience  string
	publicKey *rsa.PublicKey
	leeway    time.Duration
}

// Issuer signs access tokens. Verification is provided independently by
// Verifier.
type Issuer struct {
	issuer     string
	audience   string
	privateKey *rsa.PrivateKey
	lifetime   time.Duration
}

// NewVerifier validates verification config and parses the trusted public key.
func NewVerifier(config VerifierConfig) (*Verifier, error) {
	if config.Issuer == "" || config.Audience == "" || config.PublicKeyPEM == "" || config.Leeway < 0 {
		return nil, ErrInvalidConfig
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(config.PublicKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("%w: public key: %v", ErrInvalidConfig, err)
	}
	return &Verifier{
		issuer:    config.Issuer,
		audience:  config.Audience,
		publicKey: publicKey,
		leeway:    config.Leeway,
	}, nil
}

// NewIssuer validates signing config and parses its active RS256 private key.
func NewIssuer(config IssuerConfig) (*Issuer, error) {
	if config.Issuer == "" || config.Audience == "" || config.PrivateKeyPEM == "" || config.AccessTokenLifetime < jwt.TimePrecision {
		return nil, ErrInvalidConfig
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.PrivateKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("%w: private key: %v", ErrInvalidConfig, err)
	}
	return &Issuer{
		issuer:     config.Issuer,
		audience:   config.Audience,
		privateKey: privateKey,
		lifetime:   config.AccessTokenLifetime,
	}, nil
}

// CreateClaims creates access-token claims for an authenticated application session.
func (i *Issuer) CreateClaims(subject string, username string, sessionID string, clientID string, clientType string) (*Claims, error) {
	if subject == "" {
		return nil, ErrInvalidClaims
	}

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(i.lifetime)
	issuedAtClaim := jwt.NewNumericDate(issuedAt)
	expiresAtClaim := jwt.NewNumericDate(expiresAt)
	jwtID, err := generateJWTID()
	if err != nil {
		return nil, fmt.Errorf("auth: generate JWT ID: %w", err)
	}
	return &Claims{
		Username:   username,
		SessionID:  sessionID,
		ClientID:   clientID,
		ClientType: clientType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.issuer,
			Subject:   subject,
			Audience:  jwt.ClaimStrings{i.audience},
			IssuedAt:  issuedAtClaim,
			ExpiresAt: expiresAtClaim,
			ID:        jwtID,
		},
	}, nil
}

// CreateTokenByClaims signs claims with RS256.
func (i *Issuer) CreateTokenByClaims(claims jwt.Claims) (string, error) {
	if claims == nil {
		return "", ErrInvalidClaims
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["typ"] = accessTokenType
	value, err := token.SignedString(i.privateKey)
	if err != nil {
		return "", fmt.Errorf("auth: sign JWT: %w", err)
	}
	return value, nil
}

// CreateToken creates an access token for an authenticated application session.
func (i *Issuer) CreateToken(subject, username, sessionID, clientID, clientType string) (string, error) {
	claims, err := i.CreateClaims(subject, username, sessionID, clientID, clientType)
	if err != nil {
		return "", err
	}
	return i.CreateTokenByClaims(claims)
}

// ParseToken parses and validates an encoded access token.
func (v *Verifier) ParseToken(value string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(value, claims, func(*jwt.Token) (any, error) {
		return v.publicKey, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(v.leeway),
	)
	if err != nil {
		return nil, invalidTokenError(err)
	}
	if tokenType, ok := token.Header["typ"].(string); !ok || tokenType != accessTokenType {
		return nil, invalidTokenError(errInvalidTokenType)
	}
	return claims, nil
}

func invalidTokenError(err error) error {
	return fmt.Errorf("%w: %v", ErrInvalidToken, err)
}

func generateJWTID() (string, error) {
	var value [32]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value[:]), nil
}
