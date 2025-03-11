package vlog

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log = slog.New(NewPrettyHandler(nil))
	bws *zapcore.BufferedWriteSyncer
)

func G() *slog.Logger { return Log }

func BufferedWriteSyncer(filename string, maxSize, maxBackups, maxAge int) *zapcore.BufferedWriteSyncer {
	bws = &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
		}),
		FlushInterval: time.Minute,
	}
	return bws
}

func Sync() {
	if bws != nil {
		bws.Sync()
	}
}

func New(filename string, ws ...io.Writer) error {
	if ws == nil {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		ws = []io.Writer{f}
	}
	Log = slog.New(NewMultiHandler(
		NewPrettyHandler(nil),
		slog.NewJSONHandler(
			io.MultiWriter(ws...),
			&slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}),
	))
	return nil
}

func NewV2(ws []io.Writer, opts *slog.HandlerOptions) {
	if ws == nil {
		panic("io.Writer is nil")
	}
	if opts == nil {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}
	}
	Log = slog.New(NewMultiHandler(
		NewPrettyHandler(nil),
		slog.NewJSONHandler(io.MultiWriter(ws...), opts),
	))
}
