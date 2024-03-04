package vlog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

func log(level slog.Level, msg string, args ...any) {
	ctx := context.Background()
	if !Log.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, log]
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	attrs := make([]slog.Attr, 0, len(args))
	for i := 0; i < len(args)/2; i++ {
		attrs = append(attrs, slog.String(fmt.Sprint(args[i]), fmt.Sprint(args[i+1])))
	}
	r.AddAttrs(attrs...)
	_ = Log.Handler().Handle(ctx, r)
}

func logf(level slog.Level, format string, args ...any) {
	ctx := context.Background()
	if !Log.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, logf, ...]
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, args...), pcs[0])
	_ = Log.Handler().Handle(ctx, r)

}

func Debugf(format string, args ...any) { logf(slog.LevelDebug, format, args...) }
func Infof(format string, args ...any)  { logf(slog.LevelInfo, format, args...) }
func Warnf(format string, args ...any)  { logf(slog.LevelWarn, format, args...) }
func Errorf(format string, args ...any) { logf(slog.LevelError, format, args...) }
func Debug(msg string, args ...any)     { log(slog.LevelDebug, msg, args...) }
func Info(msg string, args ...any)      { log(slog.LevelInfo, msg, args...) }
func Warn(msg string, args ...any)      { log(slog.LevelWarn, msg, args...) }
func Error(msg string, args ...any)     { log(slog.LevelError, msg, args...) }
