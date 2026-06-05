package sqids

import (
	"errors"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	enc := MustNew()

	original := int64(42)
	hash, err := enc.Encode(original)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if len(hash) < 8 {
		t.Fatalf("expected min length 8, got %d", len(hash))
	}

	decoded, err := enc.Decode(hash)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded != original {
		t.Fatalf("expected %d, got %d", original, decoded)
	}
}

func TestEncodeDifferentIDs(t *testing.T) {
	enc := MustNew()

	results := make(map[string]int64)
	for _, id := range []int64{1, 100, 999999, 1 << 40} {
		hash, err := enc.Encode(id)
		if err != nil {
			t.Fatalf("encode %d: %v", id, err)
		}
		if existing, ok := results[hash]; ok {
			t.Fatalf("collision: %d and %d both encode to %q", existing, id, hash)
		}
		results[hash] = id

		decoded, err := enc.Decode(hash)
		if err != nil {
			t.Fatalf("decode %q: %v", hash, err)
		}
		if decoded != id {
			t.Fatalf("expected %d, got %d", id, decoded)
		}
	}
}

func TestEncodeRejectsNegativeID(t *testing.T) {
	enc := MustNew()

	hash, err := enc.Encode(-1)
	if !errors.Is(err, ErrNegativeID) {
		t.Fatalf("expected ErrNegativeID, got %v", err)
	}
	if hash != "" {
		t.Fatalf("expected empty hash, got %q", hash)
	}
}

func TestDecodeInvalidHash(t *testing.T) {
	enc := MustNew()

	_, err := enc.Decode("")
	if err == nil {
		t.Fatal("expected error for empty hash")
	}
}

func TestMustEncode(t *testing.T) {
	enc := MustNew()
	hash := enc.MustEncode(1)
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
}

func TestWithMinLength(t *testing.T) {
	enc := MustNew(WithMinLength(16))
	hash, err := enc.Encode(1)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if len(hash) < 16 {
		t.Fatalf("expected min length 16, got %d", len(hash))
	}
}

func TestWithAlphabet(t *testing.T) {
	enc := MustNew(WithAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	hash, err := enc.Encode(42)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	decoded, err := enc.Decode(hash)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded != 42 {
		t.Fatalf("expected 42, got %d", decoded)
	}
}

func TestNewReturnsInvalidAlphabetError(t *testing.T) {
	enc, err := New(WithAlphabet("ab"))
	if err == nil {
		t.Fatal("expected error for invalid alphabet")
	}
	if enc != nil {
		t.Fatalf("expected nil encoder, got %#v", enc)
	}
}

func TestMustNewPanicsOnInvalidAlphabet(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for invalid alphabet")
		}
	}()

	_ = MustNew(WithAlphabet("ab"))
}
