// Package hash provides stable digest helpers for non-password values.
package hash

import (
	"crypto/sha256"
	"encoding/base64"
)

// SHA256Base64URL returns the SHA-256 digest of value encoded with unpadded
// base64 URL encoding.
func SHA256Base64URL(value string) string {
	digest := sha256.Sum256([]byte(value))
	return base64.RawURLEncoding.EncodeToString(digest[:])
}
