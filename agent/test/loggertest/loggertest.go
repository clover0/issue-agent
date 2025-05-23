package loggertest

import (
	"io"

	"github.com/clover0/issue-agent/logger"
)

type testLogger struct{}

func NewTestLogger() logger.Logger {
	return &testLogger{}
}

func (l *testLogger) Info(msg string, args ...any)  {}
func (l *testLogger) Error(msg string, args ...any) {}
func (l *testLogger) Debug(msg string, args ...any) {}
func (l *testLogger) AddPrefix(prefix string) logger.Logger {
	return l
}
func (l *testLogger) SetColor(color logger.Color) logger.Logger {
	return l
}
func (l *testLogger) SetOutput(out io.Writer) logger.Logger { return l }
