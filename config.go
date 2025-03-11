package vlog

import (
	"io"
	"log/slog"
)

type Option interface {
	Apply(cfg *Config)
}

type OptionFunc func(cfg *Config)

func (fn OptionFunc) Apply(cfg *Config) { fn(cfg) }

type Config struct {
	Level     *slog.LevelVar
	WithLevel bool
	Opts      *slog.HandlerOptions
	WithOpts  bool
	Output    io.Writer
}

func WithLevel(lvl *slog.LevelVar) Option {
	return OptionFunc(func(cfg *Config) {
		cfg.Level, cfg.WithLevel = lvl, true
	})
}

func WithHandlerOptions(opts *slog.HandlerOptions) Option {
	return OptionFunc(func(cfg *Config) {
		cfg.Opts, cfg.WithOpts = opts, true
	})
}

func WithOutput(writer io.Writer) Option {
	return OptionFunc(func(cfg *Config) { cfg.Output = writer })
}
