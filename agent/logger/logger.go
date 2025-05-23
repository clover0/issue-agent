package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)

	AddPrefix(prefix string) Logger
	SetColor(color Color) Logger
	SetOutput(out io.Writer) Logger
}

func NewDefaultLogger(level string) Logger {
	opt := &slog.HandlerOptions{
		Level: slogLevel(level),
	}

	l := slog.New(slog.NewTextHandler(os.Stdout, opt))
	return DefaultLogger{
		logger: *l,
		opt:    opt,
	}
}

type DefaultLogger struct {
	logger slog.Logger
	opt    *slog.HandlerOptions

	// prefix is used to prefix the log message
	prefix string
}

func (l DefaultLogger) Info(msg string, args ...any) {
	l.logger.Info(l.prefix+msg, args...)
}

func (l DefaultLogger) Error(msg string, args ...any) {
	l.logger.Error(l.prefix+msg, args...)
}

func (l DefaultLogger) Debug(msg string, args ...any) {
	l.logger.Debug(l.prefix+msg, args...)
}

func (l DefaultLogger) AddPrefix(prefix string) Logger {
	l.prefix += prefix
	return l
}

func (l DefaultLogger) SetColor(color Color) Logger {
	panic("SetColor is not implemented in DefaultLogger")
}

func (l DefaultLogger) SetOutput(out io.Writer) Logger {
	ll := slog.New(slog.NewTextHandler(out, l.opt))
	return DefaultLogger{
		logger: *ll,
		opt:    l.opt,
		prefix: l.prefix,
	}
}

func slogLevel(l string) slog.Level {
	switch l {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "error":
		return slog.LevelError
	case "":
		return slog.LevelInfo
	default:
		panic(fmt.Sprintf("unknown log level: %s", l))
	}
}
