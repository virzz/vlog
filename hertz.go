package vlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) { fn(cfg) }

type config struct {
	level              *slog.LevelVar
	withLevel          bool
	handlerOptions     *slog.HandlerOptions
	withHandlerOptions bool
	output             io.Writer
}

func hLevelToSLevel(level hlog.Level) (lvl slog.Level) {
	switch level {
	case hlog.LevelTrace:
		lvl = LevelTrace
	case hlog.LevelDebug:
		lvl = slog.LevelDebug
	case hlog.LevelInfo:
		lvl = slog.LevelInfo
	case hlog.LevelWarn:
		lvl = slog.LevelWarn
	case hlog.LevelNotice:
		lvl = LevelNotice
	case hlog.LevelError:
		lvl = slog.LevelError
	case hlog.LevelFatal:
		lvl = LevelFatal
	default:
		lvl = slog.LevelWarn
	}
	return
}

func defaultConfig() *config {
	lvl := &slog.LevelVar{}
	lvl.Set(hLevelToSLevel(hlog.LevelInfo))
	handlerOptions := &slog.HandlerOptions{Level: lvl}
	return &config{
		level:              lvl,
		withLevel:          false,
		handlerOptions:     handlerOptions,
		withHandlerOptions: false,
		output:             os.Stdout,
	}
}

func WithLevel(lvl *slog.LevelVar) Option {
	return option(func(cfg *config) {
		cfg.level = lvl
		cfg.withLevel = true
	})
}

func WithHandlerOptions(opts *slog.HandlerOptions) Option {
	return option(func(cfg *config) {
		cfg.handlerOptions = opts
		cfg.withHandlerOptions = true
	})
}

func WithOutput(writer io.Writer) Option {
	return option(func(cfg *config) {
		cfg.output = writer
	})
}

const (
	LevelTrace  = slog.Level(-8)
	LevelNotice = slog.Level(2)
	LevelFatal  = slog.Level(12)
)

var _ hlog.FullLogger = (*HLog)(nil)

func NewHLog(opts ...Option) *HLog {
	config := defaultConfig()
	for _, opt := range opts {
		opt.apply(config)
	}
	if !config.withLevel && config.withHandlerOptions && config.handlerOptions.Level != nil {
		lvl := &slog.LevelVar{}
		lvl.Set(config.handlerOptions.Level.Level())
		config.level = lvl
	}
	config.handlerOptions.Level = config.level

	var replaceAttrDefined bool
	if config.handlerOptions.ReplaceAttr == nil {
		replaceAttrDefined = false
	} else {
		replaceAttrDefined = true
	}
	replaceFun := config.handlerOptions.ReplaceAttr
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.LevelKey {
			level := a.Value.Any().(slog.Level)
			switch level {
			case LevelTrace:
				a.Value = slog.StringValue("Trace")
			case slog.LevelDebug:
				a.Value = slog.StringValue("Debug")
			case slog.LevelInfo:
				a.Value = slog.StringValue("Info")
			case LevelNotice:
				a.Value = slog.StringValue("Notice")
			case slog.LevelWarn:
				a.Value = slog.StringValue("Warn")
			case slog.LevelError:
				a.Value = slog.StringValue("Error")
			case LevelFatal:
				a.Value = slog.StringValue("Fatal")
			default:
				a.Value = slog.StringValue("Warn")
			}
		}
		if replaceAttrDefined {
			return replaceFun(groups, a)
		} else {
			return a
		}
	}
	config.handlerOptions.ReplaceAttr = replaceAttr
	return &HLog{l: Log.WithGroup("hlog"), cfg: config}
}

// Logger slog impl
type HLog struct {
	l   *slog.Logger
	cfg *config
}

func (l *HLog) Logger() *slog.Logger { return l.l }

func (l *HLog) SetLevel(level hlog.Level) {
	lvl := hLevelToSLevel(level)
	l.cfg.level.Set(lvl)
}
func (l *HLog) SetOutput(writer io.Writer) {
	l.cfg.output = writer
	l.l = slog.New(slog.NewJSONHandler(writer, l.cfg.handlerOptions))
}
func (l *HLog) log(level hlog.Level, v ...any) {
	l.l.Log(context.TODO(), hLevelToSLevel(level), fmt.Sprint(v...))
}
func (l *HLog) logf(level hlog.Level, format string, kvs ...any) {
	l.l.Log(context.TODO(), hLevelToSLevel(level), fmt.Sprintf(format, kvs...))
}
func (l *HLog) ctxLogf(level hlog.Level, ctx context.Context, format string, v ...any) {
	l.l.Log(ctx, hLevelToSLevel(level), fmt.Sprintf(format, v...))
}
func (l *HLog) Trace(v ...any)             { l.log(hlog.LevelTrace, v...) }
func (l *HLog) Debug(v ...any)             { l.log(hlog.LevelDebug, v...) }
func (l *HLog) Info(v ...any)              { l.log(hlog.LevelInfo, v...) }
func (l *HLog) Notice(v ...any)            { l.log(hlog.LevelNotice, v...) }
func (l *HLog) Warn(v ...any)              { l.log(hlog.LevelWarn, v...) }
func (l *HLog) Error(v ...any)             { l.log(hlog.LevelError, v...) }
func (l *HLog) Fatal(v ...any)             { l.log(hlog.LevelFatal, v...) }
func (l *HLog) Tracef(f string, v ...any)  { l.logf(hlog.LevelTrace, f, v...) }
func (l *HLog) Debugf(f string, v ...any)  { l.logf(hlog.LevelDebug, f, v...) }
func (l *HLog) Infof(f string, v ...any)   { l.logf(hlog.LevelInfo, f, v...) }
func (l *HLog) Noticef(f string, v ...any) { l.logf(hlog.LevelNotice, f, v...) }
func (l *HLog) Warnf(f string, v ...any)   { l.logf(hlog.LevelWarn, f, v...) }
func (l *HLog) Errorf(f string, v ...any)  { l.logf(hlog.LevelError, f, v...) }
func (l *HLog) Fatalf(f string, v ...any)  { l.logf(hlog.LevelFatal, f, v...) }

func (l *HLog) CtxTracef(ctx context.Context, f string, v ...any) {
	l.ctxLogf(hlog.LevelDebug, ctx, f, v...)
}
func (l *HLog) CtxDebugf(ctx context.Context, f string, v ...any) {
	l.ctxLogf(hlog.LevelDebug, ctx, f, v...)
}
func (l *HLog) CtxInfof(ctx context.Context, f string, v ...any) {
	l.ctxLogf(hlog.LevelInfo, ctx, f, v...)
}
func (l *HLog) CtxNoticef(ctx context.Context, f string, v ...any) {
	l.ctxLogf(hlog.LevelNotice, ctx, f, v...)
}
func (l *HLog) CtxWarnf(ctx context.Context, f string, v ...any) {
	l.ctxLogf(hlog.LevelWarn, ctx, f, v...)
}
func (l *HLog) CtxErrorf(ctx context.Context, f string, v ...any) {
	l.ctxLogf(hlog.LevelError, ctx, f, v...)
}
func (l *HLog) CtxFatalf(ctx context.Context, f string, v ...any) {
	l.ctxLogf(hlog.LevelFatal, ctx, f, v...)
}
