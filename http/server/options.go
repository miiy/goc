package server

import (
	"time"

	"go.uber.org/zap"
)

// An Option configures a Server.
type Option interface {
	apply(*Server)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(server *Server)

func (f optionFunc) apply(server *Server) {
	f(server)
}

// WithDebug configures the Server debug
func WithDebug() Option {
	return optionFunc(func(opts *Server) {
		opts.debug = true
	})
}

func WithAddr(addr string) Option {
	return optionFunc(func(opts *Server) {
		opts.addr = addr
	})
}

func WithLogger(logger *zap.Logger) Option {
	return optionFunc(func(opts *Server) {
		opts.logger = logger
	})
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *Server) {
		if timeout > 0 {
			opts.shutdownTimeout = timeout
		}
	})
}
