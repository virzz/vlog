package vlog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sync"
)

const (
	timeFormat   = "[15:04:05.000] "
	reset        = "\033[0m"
	colorTime    = "\033[90m" // LightGray
	colorError   = "\033[91m" // LightRed
	colorWarn    = "\033[93m" // LightYellow
	colorDebug   = "\033[95m" // LightMagenta
	colorInfo    = "\033[96m" // LightCyan
	colorSource  = "\033[32m" // Green
	colorMessage = "\033[97m" // White
	colorAttrs   = "\033[34m" // Blue
)

type PrettyHandler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := &bytes.Buffer{}
	// [time] Level [file:line]? message
	buf.WriteString(colorTime)
	buf.WriteString(r.Time.Format(timeFormat)) // time
	switch r.Level {
	case slog.LevelDebug:
		buf.WriteString(colorDebug)
	case slog.LevelInfo:
		buf.WriteString(colorInfo)
	case slog.LevelWarn:
		buf.WriteString(colorWarn)
	case slog.LevelError:
		buf.WriteString(colorError)
	}
	buf.WriteString(r.Level.String())
	buf.WriteString(" ")
	if r.Level == slog.LevelDebug || r.Level == slog.LevelError {
		buf.WriteString(colorSource)
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf.WriteString(fmt.Sprintf("[ %s:%d ] ", f.File, f.Line))
	}
	buf.WriteString(colorMessage)
	buf.WriteString(r.Message)
	if r.NumAttrs() > 0 {
		buf.WriteString(colorAttrs)
		r.Attrs(func(a slog.Attr) bool {
			buf.WriteString(fmt.Sprintf(" %s = %v", a.Key, a.Value))
			return true
		})
	}
	buf.WriteString(reset)
	fmt.Fprintln(os.Stderr, buf.String())
	return nil
}

func NewPrettyHandler(opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}
	return &PrettyHandler{
		b: b,
		m: &sync.Mutex{},
		h: slog.NewTextHandler(b, &slog.HandlerOptions{Level: slog.LevelDebug}),
	}
}
