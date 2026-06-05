package sqids

import (
	sq "github.com/sqids/sqids-go"
)

// DefaultAlphabet is the default alphabet used for encoding.
var DefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Encoder provides hashid encoding/decoding for int64 IDs.
type Encoder struct {
	sqids *sq.Sqids
}

// New creates a new Encoder with the given options.
func New(opts ...Option) (*Encoder, error) {
	o := &options{
		alphabet:  DefaultAlphabet,
		minLength: 8,
	}
	for _, opt := range opts {
		opt.apply(o)
	}

	s, err := sq.New(sq.Options{
		Alphabet:  o.alphabet,
		MinLength: o.minLength,
	})
	if err != nil {
		return nil, err
	}
	return &Encoder{sqids: s}, nil
}

// MustNew creates a new Encoder, panicking if the options are invalid.
func MustNew(opts ...Option) *Encoder {
	enc, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return enc
}

// Encode converts an int64 ID to a hashid string.
func (e *Encoder) Encode(id int64) (string, error) {
	if id < 0 {
		return "", ErrNegativeID
	}
	return e.sqids.Encode([]uint64{uint64(id)})
}

// MustEncode converts an int64 ID to a hashid string, panicking on error.
func (e *Encoder) MustEncode(id int64) string {
	s, err := e.Encode(id)
	if err != nil {
		panic(err)
	}
	return s
}

// Decode converts a hashid string back to the original int64 ID.
// Returns an error if the hashid is invalid.
func (e *Encoder) Decode(hash string) (int64, error) {
	nums := e.sqids.Decode(hash)
	if len(nums) == 0 {
		return 0, ErrInvalidHash
	}
	return int64(nums[0]), nil
}
