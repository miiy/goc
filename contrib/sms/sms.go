// Package sms provides an SMS sender interface, a development log implementation,
// and verification-code generation.
package sms

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
)

// Sender sends an SMS (e.g. a verification code) to a phone number. Implementations:
// LogSender (development), or a real provider such as aliyun/tencent (production).
type Sender interface {
	Send(ctx context.Context, phone, content string) error
}

// GenerateCode returns a cryptographically random n-digit numeric code with uniform
// distribution (no modulo bias). E.g. GenerateCode(6) -> "048291".
func GenerateCode(digits int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := cryptorand.Int(cryptorand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", digits, n.Int64()), nil
}
