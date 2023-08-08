package gin

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

func WithLogger() Option {
	return optionFunc(func(opts *Server) {
	})
}
