package logger

import (
	"io"
	"log/slog"
	"os"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger(handler slog.Handler) Logger {
	if handler == nil {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	return &SlogLogger{
		logger: slog.New(handler),
	}
}

func NewSlogJSONLogger(writer io.Writer, opts *slog.HandlerOptions) Logger {
	if writer == nil {
		writer = os.Stdout
	}

	return NewSlogLogger(slog.NewJSONHandler(writer, opts))
}

func NewSlogTextLogger(writer io.Writer, opts *slog.HandlerOptions) Logger {
	if writer == nil {
		writer = os.Stdout
	}

	return NewSlogLogger(slog.NewTextHandler(writer, opts))
}

func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *SlogLogger) With(args ...any) Logger {
	return &SlogLogger{
		logger: l.logger.With(args...),
	}
}
