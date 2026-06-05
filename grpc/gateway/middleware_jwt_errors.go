package gateway

import "errors"

var (
	ErrMissingToken      = errors.New("missing authentication token")
	ErrInvalidTokenFormat = errors.New("invalid authentication token format")
)
