package sqids

import "errors"

var (
	ErrInvalidHash = errors.New("sqids: invalid hash")
	ErrNegativeID  = errors.New("sqids: negative id")
)
