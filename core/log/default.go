package log

import (
	"os"
)

var defaultLogger = New(os.Stderr, InfoLevel, WithCaller(true))

var (
	GetLogger   = defaultLogger.Logger
	SetLevel    = defaultLogger.SetLevel
	Enabled     = defaultLogger.Enabled
	V           = defaultLogger.V
	WithContext = defaultLogger.WithContext
	WithField   = defaultLogger.WithField
	WithMap     = defaultLogger.WithMap
	WithName    = defaultLogger.WithName

	Debug   = defaultLogger.Debug
	Debugf  = defaultLogger.Debugf
	Debugw  = defaultLogger.Debugw
	Info    = defaultLogger.Info
	Infof   = defaultLogger.Infof
	Infow   = defaultLogger.Infow
	Warn    = defaultLogger.Warn
	Warnf   = defaultLogger.Warnf
	Warnw   = defaultLogger.Warnw
	Error   = defaultLogger.Error
	Errorf  = defaultLogger.Errorf
	Errorw  = defaultLogger.Errorw
	DPanic  = defaultLogger.DPanic
	DPanicf = defaultLogger.DPanicf
	DPanicw = defaultLogger.DPanicw
	Panic   = defaultLogger.Panic
	Panicf  = defaultLogger.Panicf
	Panicw  = defaultLogger.Panicw
	Fatal   = defaultLogger.Fatal
	Fatalf  = defaultLogger.Fatalf
	Fatalw  = defaultLogger.Fatalw

	Flush = defaultLogger.Flush
)

func ResetDefault(l Logger) {
	defaultLogger = l

	GetLogger = defaultLogger.Logger
	SetLevel = defaultLogger.SetLevel
	Enabled = defaultLogger.Enabled
	V = defaultLogger.V
	WithContext = defaultLogger.WithContext
	WithField = defaultLogger.WithField
	WithMap = defaultLogger.WithMap
	WithName = defaultLogger.WithName

	Debug = defaultLogger.Debug
	Debugf = defaultLogger.Debugf
	Info = defaultLogger.Info
	Infof = defaultLogger.Infof
	Warn = defaultLogger.Warn
	Warnf = defaultLogger.Warnf
	Error = defaultLogger.Error
	Errorf = defaultLogger.Errorf
	DPanic = defaultLogger.DPanic
	DPanicf = defaultLogger.DPanicf
	Panic = defaultLogger.Panic
	Panicf = defaultLogger.Panicf
	Fatal = defaultLogger.Fatal
	Fatalf = defaultLogger.Fatalf

	Flush = defaultLogger.Flush
}

func Default() Logger {
	return defaultLogger
}
