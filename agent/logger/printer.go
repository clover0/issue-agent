package logger

import (
	"fmt"
	"log/slog"
)

type Printer struct {
	level slog.Level
}

func NewPrinter(levelStr string) Logger {
	level := slogLevel(levelStr)

	return Printer{level}
}

func (l Printer) Debug(msg string, args ...any) {
	if l.level <= slog.LevelDebug {
		fmt.Printf(msg, args...)
	}
}

func (l Printer) Info(msg string, args ...any) {
	if l.level <= slog.LevelInfo {
		fmt.Printf(msg, args...)
	}
}

func (l Printer) Error(msg string, args ...any) {
	if l.level <= slog.LevelError {
		fmt.Printf(msg, args...)
	}
}
