package vlog

import (
	"context"
	"log/slog"
)

type MultiHandler struct {
	handlers []slog.Handler
}

var _ slog.Handler = (*MultiHandler)(nil)

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if !h.Enabled(ctx, level) {
			return false
		}
	}
	return true
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range m.handlers {
		if !handler.Enabled(ctx, r.Level) {
			continue
		}
		err := handler.Handle(ctx, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewMultiHandler(hs ...slog.Handler) slog.Handler {
	return &MultiHandler{handlers: hs}
}
