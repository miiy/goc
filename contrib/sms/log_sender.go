// Package sms provides an SMS sender interface, a development log implementation,
// and verification-code generation.
package sms

import (
	"context"
	"fmt"

	"github.com/miiy/goc/logger/zap"
)

// LogSender is a development Sender that logs the message instead of sending it.
type LogSender struct {
	logger *zap.Logger
}

// NewLogSender creates a LogSender that logs messages via the given logger.
func NewLogSender(logger *zap.Logger) *LogSender {
	return &LogSender{logger: logger}
}

func (s *LogSender) Send(_ context.Context, phone, content string) error {
	s.logger.Info(fmt.Sprintf("sms (log sender, not actually sent): phone=%s content=%s", phone, content))
	return nil
}
