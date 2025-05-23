package logger_test

import (
	"bytes"
	"testing"

	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/test/assert"
)

func captureOutput(t *testing.T, printer logger.Logger, f func(logger.Logger)) string {
	t.Helper()

	var buf bytes.Buffer

	f(printer.SetOutput(&buf))

	return buf.String()
}

func TestPrinterDebug(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level    string
		msg      string
		args     []any
		prefix   string
		expected string
	}{
		"debug level prints debug message": {
			level:    "debug",
			msg:      "test message",
			expected: "test message",
		},
		"info level doesn't print debug message": {
			level:    "info",
			msg:      "test message",
			expected: "",
		},
		"error level doesn't print debug message": {
			level:    "error",
			msg:      "test message",
			expected: "",
		},
		"debug with format args": {
			level:    "debug",
			msg:      "test %s %d",
			args:     []any{"message", 123},
			expected: "test message 123",
		},
		"debug with prefix": {
			level:    "debug",
			msg:      "test message",
			prefix:   "PREFIX: ",
			expected: "PREFIX: test message",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			printer := logger.NewPrinter(tt.level)
			printer = printer.AddPrefix(tt.prefix)

			output := captureOutput(t, printer, func(p logger.Logger) {
				p.Debug(tt.msg, tt.args...)
			})

			assert.Equal(t, output, tt.expected)
		})
	}
}

func TestPrinterError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level    string
		msg      string
		args     []any
		prefix   string
		expected string
	}{
		"debug level prints error message": {
			level:    "debug",
			msg:      "test message",
			expected: "test message",
		},
		"info level prints error message": {
			level:    "info",
			msg:      "test message",
			expected: "test message",
		},
		"error level prints error message": {
			level:    "error",
			msg:      "test message",
			expected: "test message",
		},
		"error with format args": {
			level:    "error",
			msg:      "test %s %d",
			args:     []any{"message", 123},
			expected: "test message 123",
		},
		"error with prefix": {
			level:    "error",
			msg:      "test message",
			prefix:   "PREFIX: ",
			expected: "PREFIX: test message",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			printer := logger.NewPrinter(tt.level)
			printer = printer.AddPrefix(tt.prefix)

			output := captureOutput(t, printer, func(p logger.Logger) {
				p.Error(tt.msg, tt.args...)
			})

			assert.Equal(t, output, tt.expected)
		})
	}
}

func TestPrinterAddPrefix(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		initialPrefix string
		addPrefix     string
		msg           string
		expected      string
	}{
		"add prefix to empty": {
			initialPrefix: "",
			addPrefix:     "PREFIX: ",
			msg:           "test",
			expected:      "PREFIX: test",
		},
		"add prefix to existing": {
			initialPrefix: "INITIAL: ",
			addPrefix:     "PREFIX: ",
			msg:           "test",
			expected:      "INITIAL: PREFIX: test",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			printer := logger.NewPrinter("debug")
			printer = printer.AddPrefix(tt.initialPrefix)
			printer = printer.AddPrefix(tt.addPrefix)

			output := captureOutput(t, printer, func(p logger.Logger) {
				p.Debug("test")
			})

			assert.Equal(t, output, tt.expected)
		})
	}
}

func TestPrinterSetColor(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		color    logger.Color
		expected string
	}{
		"set color to green": {
			color:    logger.Green,
			expected: logger.Green.String() + "test message" + logger.Reset.String(),
		},
		"set color to red": {
			color:    logger.Red,
			expected: logger.Red.String() + "test message" + logger.Reset.String(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			printer := logger.NewPrinter("debug")
			printer = printer.SetColor(tt.color)

			output := captureOutput(t, printer, func(p logger.Logger) {
				p.Debug("test message")
			})

			assert.Equal(t, output, tt.expected)
		})
	}
}
