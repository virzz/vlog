package vlog_test

import (
	"testing"

	"github.com/virzz/vlog"
)

func TestNew(t *testing.T) {
	vlog.New("test.log")
	vlog.G().Debug("Debug")
	vlog.G().Error("Error")
	vlog.G().Info("Info")
	vlog.G().Warn("Warn")
	vlog.Log.Warn("Warn")
}
