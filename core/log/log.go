package log

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Enabled(lvl zapcore.Level) bool
	V(lvl int) bool
	Logger() *zap.Logger
	SetLevel(lv Level)

	Debug(v any, fields ...Field)
	Debugf(format string, v ...interface{})
	Debugw(v any, keysAndValues ...interface{})
	Info(v any, fields ...Field)
	Infof(format string, v ...interface{})
	Infow(v any, keysAndValues ...interface{})
	Warn(v any, fields ...Field)
	Warnf(format string, v ...interface{})
	Warnw(v any, keysAndValues ...interface{})
	Error(v any, fields ...Field)
	Errorf(format string, v ...interface{})
	Errorw(v any, keysAndValues ...interface{})
	DPanic(v any, fields ...Field)
	DPanicf(format string, v ...interface{})
	DPanicw(v any, keysAndValues ...interface{})
	Panic(v any, fields ...Field)
	Panicf(format string, v ...interface{})
	Panicw(v any, keysAndValues ...interface{})
	Fatal(v any, fields ...Field)
	Fatalf(format string, v ...interface{})
	Fatalw(v any, keysAndValues ...interface{})

	WithField(fields ...Field) Logger
	WithName(name string) Logger
	WithContext(ctx context.Context, keys ...string) Logger
	WithMap(fields map[string]any) Logger
	Flush() error
}

type logger struct {
	l           *zap.Logger
	lv          *zap.AtomicLevel
	development bool
	addCaller   bool
	callSkip    int
}

var _ Logger = &logger{}

func NewTee(writers []io.Writer, level Level, opts ...Option) Logger {
	l := &logger{callSkip: 1}
	lv := zap.NewAtomicLevelAt(level)
	l.lv = &lv
	for _, opt := range opts {
		opt.apply(l)
	}

	cfg := zap.NewProductionConfig()
	if l.development {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}
	var enc zapcore.Encoder
	if l.development {
		enc = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	} else {
		enc = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	}

	cores := make([]zapcore.Core, 0, len(writers))
	for _, w := range writers {
		core := zapcore.NewCore(
			enc,
			zapcore.AddSync(w),
			lv,
		)
		cores = append(cores, core)
	}

	options := []zap.Option{zap.WithCaller(l.addCaller)}
	if l.development {
		options = append(options, zap.Development())
	}
	options = append(options, zap.AddCallerSkip(l.callSkip))

	l.l = zap.New(
		zapcore.NewTee(cores...),
		options...,
	)

	return l
}

func New(writer io.Writer, level Level, opts ...Option) Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	return NewTee([]io.Writer{writer}, level, opts...)
}

func (l *logger) Logger() *zap.Logger {
	return l.l
}

func (l *logger) Flush() error {
	return l.l.Sync()
}

func (l *logger) SetLevel(lv Level) {
	l.lv.SetLevel(lv)
}

// Enabled returns true if the given level is at or above this level.
func (l *logger) Enabled(lvl zapcore.Level) bool {
	return l.lv.Enabled(lvl)
}

// V returns true if the given level is at or above this level.
// same as Enabled
func (l *logger) V(lvl int) bool {
	return l.lv.Enabled(zapcore.Level(lvl))
}

// WithContext return log with inject context.
func (l *logger) WithContext(ctx context.Context, keys ...string) Logger {
	nl := l.clone()
	for _, k := range keys {
		if v := ctx.Value(k); v != nil {
			nl.l = nl.l.With(zap.Any(k, v))
		}
	}

	return nl
}

func (l *logger) WithField(fields ...Field) Logger {
	nl := l.clone()
	nl.l = nl.l.With(fields...)
	return nl
}

func (l *logger) WithName(name string) Logger {
	nl := l.clone()
	nl.l = nl.l.Named(name)
	return nl
}

func (l *logger) WithMap(m map[string]any) Logger {
	fields := make([]Field, 0, len(m))
	for k, v := range m {
		fields = append(fields, zap.Any(k, v))
	}
	nl := l.clone()
	nl.l = nl.l.With(fields...)
	return nl
}

func (l *logger) Debug(v any, fields ...Field) {
	if !l.lv.Enabled(DebugLevel) {
		return
	}
	l.l.Debug(fmt.Sprint(v), fields...)
}

func (l *logger) Debugf(format string, args ...any) {
	if !l.lv.Enabled(DebugLevel) {
		return
	}
	l.l.Sugar().Debugf(format, args...)
}

func (l *logger) Debugw(v any, keysAndValues ...interface{}) {
	if !l.lv.Enabled(DebugLevel) {
		return
	}
	l.l.Sugar().Debugw(fmt.Sprint(v), keysAndValues...)
}

func (l *logger) Info(v any, fields ...Field) {
	if !l.lv.Enabled(InfoLevel) {
		return
	}
	l.l.Info(fmt.Sprint(v), fields...)
}

func (l *logger) Infof(format string, args ...any) {
	if !l.lv.Enabled(InfoLevel) {
		return
	}
	l.l.Sugar().Infof(format, args...)
}

func (l *logger) Infow(v any, keysAndValues ...interface{}) {
	if !l.lv.Enabled(InfoLevel) {
		return
	}
	l.l.Sugar().Infow(fmt.Sprint(v), keysAndValues...)
}

func (l *logger) Warn(v any, fields ...Field) {
	if !l.lv.Enabled(WarnLevel) {
		return
	}
	l.l.Warn(fmt.Sprint(v), fields...)
}

func (l *logger) Warnf(format string, args ...any) {
	if !l.lv.Enabled(WarnLevel) {
		return
	}
	l.l.Sugar().Warnf(format, args...)
}

func (l *logger) Warnw(v any, keysAndValues ...interface{}) {
	l.l.Sugar().Warnw(fmt.Sprint(v), keysAndValues...)
}

func (l *logger) Error(v any, fields ...Field) {
	if !l.lv.Enabled(ErrorLevel) {
		return
	}
	l.l.Error(fmt.Sprint(v), fields...)
}

func (l *logger) Errorf(format string, args ...any) {
	if !l.lv.Enabled(ErrorLevel) {
		return
	}
	l.l.Sugar().Errorf(format, args...)
}

func (l *logger) Errorw(v any, keysAndValues ...interface{}) {
	if !l.lv.Enabled(ErrorLevel) {
		return
	}
	l.l.Sugar().Errorw(fmt.Sprint(v), keysAndValues...)
}

func (l *logger) DPanic(v any, fields ...Field) {
	if !l.lv.Enabled(DPanicLevel) {
		return
	}
	l.l.DPanic(fmt.Sprint(v), fields...)
}

func (l *logger) DPanicf(format string, args ...any) {
	if !l.lv.Enabled(DPanicLevel) {
		return
	}
	l.l.Sugar().DPanicf(format, args...)
}

func (l *logger) DPanicw(v any, keysAndValues ...interface{}) {
	if !l.lv.Enabled(DPanicLevel) {
		return
	}
	l.l.Sugar().DPanicw(fmt.Sprint(v), keysAndValues...)
}

func (l *logger) Panic(v any, fields ...Field) {
	if !l.lv.Enabled(PanicLevel) {
		return
	}
	l.l.Panic(fmt.Sprint(v), fields...)
}

func (l *logger) Panicf(format string, args ...any) {
	if !l.lv.Enabled(PanicLevel) {
		return
	}
	l.l.Sugar().Panicf(format, args...)
}

func (l *logger) Panicw(v any, keysAndValues ...interface{}) {
	l.l.Sugar().Panicw(fmt.Sprint(v), keysAndValues...)
}

func (l *logger) Fatal(v any, fields ...Field) {
	if !l.lv.Enabled(FatalLevel) {
		return
	}
	l.l.Fatal(fmt.Sprint(v), fields...)
}

func (l *logger) Fatalf(format string, args ...any) {
	if !l.lv.Enabled(FatalLevel) {
		return
	}
	l.l.Sugar().Fatalf(format, args...)
}

func (l *logger) Fatalw(v any, keysAndValues ...interface{}) {
	l.l.Sugar().Fatalw(fmt.Sprint(v), keysAndValues...)
}

func (l *logger) clone() *logger {
	copied := *l
	return &copied
}
