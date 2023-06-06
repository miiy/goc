package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Default logger.
	logger *ZapLogger
)

type Field = zap.Field

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	DPanic(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

type ZapLogger struct {
	Level            Level
	OutputPaths      []string
	ErrorOutputPaths []string
	Logger           *zap.Logger
}

// expose each method of the default instance to the user as a package-level function
var (
	Debug  = logger.Debug
	Info   = logger.Info
	Warn   = logger.Warn
	Error  = logger.Error
	DPanic = logger.DPanic
	Panic  = logger.Panic
	Fatal  = logger.Fatal
)

func NewLogger(opts ...Option) (*ZapLogger, error) {
	logger = &ZapLogger{
		Level:            InfoLevel,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	for _, o := range opts {
		o(logger)
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(logger.Level)
	zapCfg.OutputPaths = logger.OutputPaths
	zapCfg.ErrorOutputPaths = logger.ErrorOutputPaths
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	logger.Logger = l

	return logger, nil
}

func Default() *ZapLogger {
	return logger
}

func (log *ZapLogger) Debug(msg string, fields ...Field) {
	log.Logger.Debug(msg, fields...)
}

func (log *ZapLogger) Info(msg string, fields ...Field) {
	log.Logger.Info(msg, fields...)
}

func (log *ZapLogger) Warn(msg string, fields ...Field) {
	log.Logger.Warn(msg, fields...)
}

func (log *ZapLogger) Error(msg string, fields ...Field) {
	log.Logger.Error(msg, fields...)
}

func (log *ZapLogger) DPanic(msg string, fields ...Field) {
	log.Logger.DPanic(msg, fields...)
}

func (log *ZapLogger) Panic(msg string, fields ...Field) {
	log.Logger.Panic(msg, fields...)
}

func (log *ZapLogger) Fatal(msg string, fields ...Field) {
	log.Logger.Fatal(msg, fields...)
}

func (log *ZapLogger) Sync() error {
	return log.Logger.Sync()
}
