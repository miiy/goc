package logger

import (
	"errors"
	"github.com/miiy/goc/logger/zap"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tLogger, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	defer tLogger.Sync()
	logger.Info("info msg.")
	Debug("debug msg.")
	Info("info msg.")
	Warn("warn msg.")
}

func TestNewLoggerWithLevel(t *testing.T) {
	tLogger, err := NewLogger(WithLevel(DebugLevel))
	if err != nil {
		t.Fatal(err)
	}
	defer tLogger.Sync()
	logger.Debug("debug msg.")
	logger.Info("info msg.")
	logger.Error("error msg.", zap.Error(errors.New("error info")))
}
