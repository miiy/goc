package sms

import (
	"testing"
)

func TestGenerateCode(t *testing.T) {
	code, err := GenerateCode(6)
	if err != nil {
		t.Fatalf("GenerateCode: %v", err)
	}
	if len(code) != 6 {
		t.Fatalf("expected 6 digits, got %d (%q)", len(code), code)
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			t.Fatalf("expected digits only, got %q", code)
		}
	}
}

func TestGenerateCodeDigitCounts(t *testing.T) {
	for _, n := range []int{4, 6, 8} {
		code, err := GenerateCode(n)
		if err != nil {
			t.Fatalf("GenerateCode(%d): %v", n, err)
		}
		if len(code) != n {
			t.Fatalf("expected %d digits, got %d (%q)", n, len(code), code)
		}
	}
}
