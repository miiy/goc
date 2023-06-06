package logger

type Option func(*ZapLogger)

func WithLevel(l Level) Option {
	return func(log *ZapLogger) {
		log.Level = l
	}
}

func WithOutputPath(path []string) Option {
	return func(log *ZapLogger) {
		log.OutputPaths = path
	}
}

func WithErrOutputPath(path []string) Option {
	return func(log *ZapLogger) {
		log.ErrorOutputPaths = path
	}
}
