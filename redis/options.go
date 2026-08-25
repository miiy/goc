package redis

import (
	"time"

	upstream "github.com/redis/go-redis/v9"
)

// Option configures the Redis client connection pool and timeouts.
type Option interface {
	apply(*upstream.UniversalOptions)
}

type optionFunc func(*upstream.UniversalOptions)

func (f optionFunc) apply(options *upstream.UniversalOptions) {
	f(options)
}

func WithPoolSize(size int) Option {
	return optionFunc(func(options *upstream.UniversalOptions) {
		options.PoolSize = size
	})
}

func WithMinIdleConns(count int) Option {
	return optionFunc(func(options *upstream.UniversalOptions) {
		options.MinIdleConns = count
	})
}

func WithMaxIdleConns(count int) Option {
	return optionFunc(func(options *upstream.UniversalOptions) {
		options.MaxIdleConns = count
	})
}

func WithDialTimeout(timeout time.Duration) Option {
	return optionFunc(func(options *upstream.UniversalOptions) {
		options.DialTimeout = timeout
	})
}

func WithReadTimeout(timeout time.Duration) Option {
	return optionFunc(func(options *upstream.UniversalOptions) {
		options.ReadTimeout = timeout
	})
}

func WithWriteTimeout(timeout time.Duration) Option {
	return optionFunc(func(options *upstream.UniversalOptions) {
		options.WriteTimeout = timeout
	})
}
