// Package password provides password strength validation.
//
// The rules enforce a minimum length, a maximum length that stays within bcrypt's
// 72-byte hashing limit, and basic complexity (must contain both letters and digits).
package password

import (
	"errors"
	"unicode"
)

const (
	// MinLen is the minimum acceptable password length in bytes.
	MinLen = 8
	// MaxLen is the maximum acceptable password length in bytes, kept within
	// bcrypt's 72-byte limit to avoid silent truncation.
	MaxLen = 64
)

var (
	// ErrTooShort is returned when the password is shorter than MinLen.
	ErrTooShort = errors.New("password must be at least 8 characters")
	// ErrTooLong is returned when the password is longer than MaxLen.
	ErrTooLong = errors.New("password must be at most 64 characters")
	// ErrTooWeak is returned when the password lacks letters or digits.
	ErrTooWeak = errors.New("password must contain both letters and numbers")
)

// Validate checks password length (8-64 bytes) and basic complexity (must contain
// both letters and digits). It returns nil for a valid password.
func Validate(password string) error {
	n := len(password)
	if n < MinLen {
		return ErrTooShort
	}
	if n > MaxLen {
		return ErrTooLong
	}
	hasLetter, hasDigit := false, false
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return ErrTooWeak
	}
	return nil
}
