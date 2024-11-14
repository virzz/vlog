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
	Level              *slog.LevelVar
	WithLevel          bool
	HandlerOptions     *slog.HandlerOptions
	WithHandlerOptions bool
	Output             io.Writer
}

func WithLevel(lvl *slog.LevelVar) Option {
	return OptionFunc(func(cfg *Config) {
		cfg.Level = lvl
		cfg.WithLevel = true
	})
}

func WithHandlerOptions(opts *slog.HandlerOptions) Option {
	return OptionFunc(func(cfg *Config) {
		cfg.HandlerOptions = opts
		cfg.WithHandlerOptions = true
	})
}

func WithOutput(writer io.Writer) Option {
	return OptionFunc(func(cfg *Config) {
		cfg.Output = writer
	})
}
