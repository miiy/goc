// Package sms provides an SMS sender interface, a development log implementation,
// and verification-code generation.
package sms

import (
	"context"
	"fmt"
)

// Logger is the minimal logger contract used by LogSender.
type Logger interface {
	Info(msg string)
}

// LogSender is a development Sender that logs the message instead of sending it.
type LogSender struct {
	logger Logger
}

// NewLogSender creates a LogSender that logs messages via the given logger.
func NewLogSender(logger Logger) *LogSender {
	return &LogSender{logger: logger}
}

func (s *LogSender) Send(_ context.Context, phone, content string) error {
	if s.logger == nil {
		return nil
	}
	s.logger.Info(fmt.Sprintf("sms (log sender, not actually sent): phone=%s content=%s", phone, content))
	return nil
}
