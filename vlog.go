package vlog

import (
	"fmt"
	"log/slog"
	"os"
)

var Log = slog.New(NewPrettyHandler(nil))

func G() *slog.Logger { return Log }

func New(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	Log = slog.New(NewMultiHandler(
		NewPrettyHandler(nil),
		slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}),
	))
	return nil
}
