package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	testIssuer   = "https://auth.example.com"
	testAudience = "nova-api"
)

func testKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}
	return key
}

func testPublicKey(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()
	value, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: value}))
}

func testPrivateKey(key *rsa.PrivateKey) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}))
}

func testVerifierConfig(t *testing.T, key *rsa.PrivateKey) VerifierConfig {
	t.Helper()
	return VerifierConfig{
		Issuer:       testIssuer,
		Audience:     testAudience,
		PublicKeyPEM: testPublicKey(t, key),
	}
}

func testIssuerConfig(key *rsa.PrivateKey) IssuerConfig {
	return IssuerConfig{
		Issuer:              testIssuer,
		Audience:            testAudience,
		PrivateKeyPEM:       testPrivateKey(key),
		AccessTokenLifetime: time.Minute,
	}
}

func testVerifier(t *testing.T, key *rsa.PrivateKey) *Verifier {
	t.Helper()
	verifier, err := NewVerifier(testVerifierConfig(t, key))
	if err != nil {
		t.Fatalf("new JWT verifier: %v", err)
	}
	return verifier
}

func testIssuerForKey(t *testing.T, key *rsa.PrivateKey) *Issuer {
	t.Helper()
	issuer, err := NewIssuer(testIssuerConfig(key))
	if err != nil {
		t.Fatalf("new JWT issuer: %v", err)
	}
	return issuer
}

func TestIssuerCreatesTokenAcceptedByVerifier(t *testing.T) {
	key := testKey(t)
	issuer := testIssuerForKey(t, key)
	verifier := testVerifier(t, key)

	claims, err := issuer.CreateClaims(
		"42",
		"alice",
		"app-session",
		"nova-ios",
		"app",
	)
	if err != nil {
		t.Fatalf("create claims: %v", err)
	}
	if lifetime := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time); lifetime != time.Minute {
		t.Fatalf("token lifetime = %v, want %v", lifetime, time.Minute)
	}
	value, err := issuer.CreateTokenByClaims(claims)
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	verifiedClaims, err := verifier.ParseToken(value)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if verifiedClaims.Subject != "42" || verifiedClaims.Username != "alice" || verifiedClaims.SessionID != "app-session" ||
		verifiedClaims.ClientID != "nova-ios" || verifiedClaims.ClientType != "app" || verifiedClaims.ID == "" {
		t.Fatalf("unexpected claims: %+v", verifiedClaims)
	}
}

func TestIssuerCreateToken(t *testing.T) {
	key := testKey(t)
	issuer := testIssuerForKey(t, key)
	verifier := testVerifier(t, key)
	value, err := issuer.CreateToken("42", "", "", "", "")
	if err != nil {
		t.Fatalf("create token: %v", err)
	}
	claims, err := verifier.ParseToken(value)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.Subject != "42" || claims.Username != "" || claims.SessionID != "" || claims.ClientID != "" || claims.ClientType != "" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestNewVerifierRejectsInvalidConfig(t *testing.T) {
	key := testKey(t)
	publicKey := testPublicKey(t, key)
	tests := []struct {
		name   string
		config VerifierConfig
	}{
		{name: "missing profile", config: VerifierConfig{PublicKeyPEM: publicKey}},
		{name: "missing public key", config: VerifierConfig{Issuer: testIssuer, Audience: testAudience}},
		{name: "negative leeway", config: VerifierConfig{Issuer: testIssuer, Audience: testAudience, PublicKeyPEM: publicKey, Leeway: -time.Second}},
		{name: "invalid public key", config: VerifierConfig{Issuer: testIssuer, Audience: testAudience, PublicKeyPEM: "not PEM"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewVerifier(test.config); !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("error = %v, want invalid config", err)
			}
		})
	}
}

func TestNewIssuerRejectsInvalidConfig(t *testing.T) {
	key := testKey(t)
	tests := []struct {
		name   string
		config IssuerConfig
	}{
		{name: "missing issuer", config: IssuerConfig{Audience: testAudience, PrivateKeyPEM: testPrivateKey(key), AccessTokenLifetime: time.Minute}},
		{name: "missing audience", config: IssuerConfig{Issuer: testIssuer, PrivateKeyPEM: testPrivateKey(key), AccessTokenLifetime: time.Minute}},
		{name: "missing private key", config: IssuerConfig{Issuer: testIssuer, Audience: testAudience, AccessTokenLifetime: time.Minute}},
		{name: "short lifetime", config: IssuerConfig{Issuer: testIssuer, Audience: testAudience, PrivateKeyPEM: testPrivateKey(key), AccessTokenLifetime: jwt.TimePrecision / 2}},
		{name: "invalid private key", config: IssuerConfig{Issuer: testIssuer, Audience: testAudience, PrivateKeyPEM: "not PEM", AccessTokenLifetime: time.Minute}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewIssuer(test.config); !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("error = %v, want invalid config", err)
			}
		})
	}
}

func TestCreateClaimsRejectsInvalidInput(t *testing.T) {
	issuer := testIssuerForKey(t, testKey(t))
	if _, err := issuer.CreateClaims("", "alice", "app-session", "nova-ios", "app"); !errors.Is(err, ErrInvalidClaims) {
		t.Fatalf("create claims error = %v, want invalid claims", err)
	}
}

func TestVerifierRejectsInvalidToken(t *testing.T) {
	key := testKey(t)
	verifier := testVerifier(t, key)
	claims := validClaims(time.Now())
	tests := []struct {
		name      string
		method    jwt.SigningMethod
		tokenType string
		claims    Claims
	}{
		{name: "wrong type", method: jwt.SigningMethodRS256, tokenType: "JWT", claims: claims},
		{name: "wrong method", method: jwt.SigningMethodHS256, tokenType: accessTokenType, claims: claims},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := signJWT(t, key, test.method, test.tokenType, test.claims)
			if _, err := verifier.ParseToken(value); !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("error = %v, want invalid JWT", err)
			}
		})
	}
}

func TestVerifierRejectsWrongPublicKey(t *testing.T) {
	verifier := testVerifier(t, testKey(t))
	value := signJWT(t, testKey(t), jwt.SigningMethodRS256, accessTokenType, validClaims(time.Now()))
	if _, err := verifier.ParseToken(value); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("error = %v, want invalid JWT", err)
	}
}

func validClaims(now time.Time) Claims {
	return Claims{
		Username: "alice", SessionID: "app-session", ClientID: "nova-ios", ClientType: "app",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: testIssuer, Audience: jwt.ClaimStrings{testAudience}, Subject: "42", ID: "token-id",
			IssuedAt: jwt.NewNumericDate(now.Add(-time.Minute)), ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
		},
	}
}

func signJWT(t *testing.T, key *rsa.PrivateKey, method jwt.SigningMethod, tokenType string, claims Claims) string {
	t.Helper()
	token := jwt.NewWithClaims(method, claims)
	token.Header["typ"] = tokenType
	var signingKey any = key
	if method == jwt.SigningMethodHS256 {
		signingKey = []byte("not-an-rsa-key")
	}
	value, err := token.SignedString(signingKey)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return value
}
