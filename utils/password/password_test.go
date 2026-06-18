package password

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{name: "valid letters and digits", password: "password123", wantErr: nil},
		{name: "boundary min length valid", password: "pass1234", wantErr: nil},
		{name: "too short", password: "ab1", wantErr: ErrTooShort},
		{name: "digits only", password: "12345678", wantErr: ErrTooWeak},
		{name: "letters only", password: "abcdefgh", wantErr: ErrTooWeak},
		{name: "too long", password: strings.Repeat("a1", 33), wantErr: ErrTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func TestHash(t *testing.T) {
	hashed, err := Hash("password123")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if hashed == "" {
		t.Fatal("Hash() returned empty hash")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte("password123")); err != nil {
		t.Fatalf("hash does not match password: %v", err)
	}
}
