package logger_test

import (
	"bytes"
	"testing"

	"github.com/clover0/issue-agent/logger"
	"github.com/clover0/issue-agent/test/assert"
)

func captureLoggerOutput(t *testing.T, l logger.Logger, f func(logger.Logger)) string {
	t.Helper()

	var buf bytes.Buffer

	f(l.SetOutput(&buf))

	return buf.String()
}

func TestDefaultLoggerDebug(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level  string
		msg    string
		args   []any
		prefix string
		want   string
	}{
		"debug level prints debug message": {
			level: "debug",
			msg:   "test message",
			want:  "test message",
		},
		"info level doesn't print debug message": {
			level: "info",
			msg:   "test message",
			want:  "",
		},
		"error level doesn't print debug message": {
			level: "error",
			msg:   "test message",
			want:  "",
		},
		"debug with format args": {
			level: "debug",
			msg:   "test",
			args:  []any{"param1", 123},
			want:  "test",
		},
		"debug with prefix": {
			level:  "debug",
			msg:    "test message",
			prefix: "PREFIX: ",
			want:   "PREFIX: test message",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l := logger.NewDefaultLogger(tt.level)
			l = l.AddPrefix(tt.prefix)

			output := captureLoggerOutput(t, l, func(log logger.Logger) {
				log.Debug(tt.msg, tt.args...)
			})

			assert.Contains(t, output, tt.want)
			if len(tt.args) > 0 {
				assert.Contains(t, output, "param1=123")
			}
		})
	}
}

func TestDefaultLoggerInfo(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level  string
		msg    string
		args   []any
		prefix string
		want   string
	}{
		"debug level prints info message": {
			level: "debug",
			msg:   "test message",
			want:  "test message",
		},
		"info level prints info message": {
			level: "info",
			msg:   "test message",
			want:  "test message",
		},
		"error level doesn't print info message": {
			level: "error",
			msg:   "test message",
			want:  "",
		},
		"info with format args": {
			level: "info",
			msg:   "test",
			args:  []any{"param1", 123},
			want:  "test",
		},
		"info with prefix": {
			level:  "info",
			msg:    "test message",
			prefix: "PREFIX: ",
			want:   "PREFIX: test message",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l := logger.NewDefaultLogger(tt.level)
			l = l.AddPrefix(tt.prefix)

			output := captureLoggerOutput(t, l, func(log logger.Logger) {
				log.Info(tt.msg, tt.args...)
			})

			assert.Contains(t, output, tt.want)
			if len(tt.args) > 0 {
				assert.Contains(t, output, "param1=123")
			}
		})
	}
}

func TestDefaultLoggerError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level  string
		msg    string
		args   []any
		prefix string
		want   string
	}{
		"debug level prints error message": {
			level: "debug",
			msg:   "test message",
			want:  "test message",
		},
		"info level prints error message": {
			level: "info",
			msg:   "test message",
			want:  "test message",
		},
		"error level prints error message": {
			level: "error",
			msg:   "test message",
			want:  "test message",
		},
		"error with format args": {
			level: "error",
			msg:   "test %s %d",
			args:  []any{"param1", 123},
			want:  "test",
		},
		"error with prefix": {
			level:  "error",
			msg:    "test message",
			prefix: "PREFIX: ",
			want:   "PREFIX: test message",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l := logger.NewDefaultLogger(tt.level)
			l = l.AddPrefix(tt.prefix)

			output := captureLoggerOutput(t, l, func(log logger.Logger) {
				log.Error(tt.msg, tt.args...)
			})

			assert.Contains(t, output, tt.want)
			if len(tt.args) > 0 {
				assert.Contains(t, output, "param1=123")
			}
		})
	}
}

func TestDefaultLoggerAddPrefix(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		initialPrefix string
		addPrefix     string
		msg           string
		want          string
	}{
		"add prefix to empty": {
			initialPrefix: "",
			addPrefix:     "PREFIX: ",
			msg:           "test",
			want:          "PREFIX: test",
		},
		"add prefix to existing": {
			initialPrefix: "INITIAL: ",
			addPrefix:     "PREFIX: ",
			msg:           "test",
			want:          "INITIAL: PREFIX: test",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l := logger.NewDefaultLogger("debug")
			l = l.AddPrefix(tt.initialPrefix)
			l = l.AddPrefix(tt.addPrefix)

			output := captureLoggerOutput(t, l, func(log logger.Logger) {
				log.Debug(tt.msg)
			})

			assert.Contains(t, output, tt.want)
		})
	}
}

func TestDefaultLoggerSetColorPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SetColor did not panic")
		}
	}()

	l := logger.NewDefaultLogger("debug")
	l.SetColor(logger.Green)
}
