package sqids

type options struct {
	alphabet  string
	minLength uint8
}

// Option configures an Encoder.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithAlphabet sets a custom alphabet for encoding.
func WithAlphabet(alphabet string) Option {
	return optionFunc(func(o *options) {
		if alphabet != "" {
			o.alphabet = alphabet
		}
	})
}

// WithMinLength sets the minimum length of the encoded string.
func WithMinLength(n uint8) Option {
	return optionFunc(func(o *options) {
		if n > 0 {
			o.minLength = n
		}
	})
}
