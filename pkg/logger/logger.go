package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/lmittmann/tint"
)

type MultiHandler struct {
	Handlers []slog.Handler
}

func (m *MultiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range m.Handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.Handlers {
		_ = h.Handle(ctx, r)
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.Handlers))
	for i, h := range m.Handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{Handlers: newHandlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.Handlers))
	for i, h := range m.Handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{Handlers: newHandlers}
}

func Setup(homeDir, dataHome, logFileName string) error {
	logFilePath := filepath.Join(dataHome, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	textHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.DateTime,
	})

	jsonHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	multi := &MultiHandler{Handlers: []slog.Handler{textHandler, jsonHandler}}
	logger := slog.New(multi)
	slog.SetDefault(logger)

	return nil
}
