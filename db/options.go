package db

import "time"

// An Option configures a Server.
type Option interface {
	apply(*options)
}

type options struct {
	connMaxLifetime time.Duration
	maxIdleConns    int
	maxOpenConns    int
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(o *options)

var defaultOption = options{
	connMaxLifetime: time.Minute * 3,
	maxIdleConns:    10,
	maxOpenConns:    100,
}

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithConnMaxLifetime(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.connMaxLifetime = t
	})
}

func WithMaxIdleConns(n int) Option {
	return optionFunc(func(c *options) {
		c.maxIdleConns = n
	})
}

func WithMaxOpenConns(n int) Option {
	return optionFunc(func(c *options) {
		c.maxOpenConns = n
	})
}
