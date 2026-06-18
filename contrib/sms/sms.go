// Package sms provides an SMS sender interface, a development log implementation,
// and verification-code generation.
package sms

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// Sender sends an SMS to a phone number.
type Sender interface {
	Send(ctx context.Context, phone, content string) error
}

// GenerateCode returns a cryptographically random n-digit numeric code with uniform
// distribution (no modulo bias). E.g. GenerateCode(6) -> "048291".
func GenerateCode(digits int) (string, error) {
	if digits <= 0 {
		return "", fmt.Errorf("digits must be positive")
	}
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := cryptorand.Int(cryptorand.Reader, max)
	if err != nil {
		return "", err
	}
	code := n.String()
	if len(code) < digits {
		code = strings.Repeat("0", digits-len(code)) + code
	}
	return code, nil
}
