package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Printer struct {
	level slog.Level

	// prefix is a string that will be added to the beginning of each log message
	prefix string

	// out is log output
	out io.Writer

	// colorize is a function that colorizes the log message
	colorize ColorFunc
}

func NewPrinter(levelStr string) Logger {
	level := slogLevel(levelStr)

	return Printer{
		level:    level,
		colorize: func(s string) string { return s },
		out:      os.Stdout,
	}
}

func (p Printer) Debug(msg string, args ...any) {
	if p.level <= slog.LevelDebug {
		_, _ = fmt.Fprintf(p.out, p.colorize(p.prefix+msg), args...)
	}
}

func (p Printer) Info(msg string, args ...any) {
	if p.level <= slog.LevelInfo {
		_, _ = fmt.Fprintf(p.out, p.colorize(p.prefix+msg), args...)
	}
}

func (p Printer) Error(msg string, args ...any) {
	if p.level <= slog.LevelError {
		_, _ = fmt.Fprintf(p.out, p.colorize(p.prefix+msg), args...)
	}
}

func (p Printer) AddPrefix(prefix string) Logger {
	p.prefix += prefix
	return p
}

func (p Printer) SetColor(color Color) Logger {
	p.colorize = GetColorize(color)
	return p
}

func (p Printer) SetOutput(out io.Writer) Logger {
	p.out = out
	return p
}
