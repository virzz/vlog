package vlog_test

import (
	"log/slog"
	"testing"

	"github.com/virzz/vlog"
)

func TestPretty(t *testing.T) {
	log := slog.New(vlog.NewPrettyHandler(nil))
	log.Debug("Debug message", "a", "1")
	log.Error("Error message", "err", "orzzzz")
	log = log.WithGroup("testtttttt")
	log.Warn("Warn message", "a", "1", "b", "2")
	log = log.With("with", "----")
	log.Info("Info message")
}
