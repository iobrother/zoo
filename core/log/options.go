package log

type Option interface {
	apply(logger *logger)
}

type optionFunc func(logger *logger)

func (f optionFunc) apply(log *logger) {
	f(log)
}

func WithCaller(enabled bool) Option {
	return optionFunc(func(l *logger) {
		l.addCaller = enabled
	})
}

func WithCallerSkip(skip int) Option {
	return optionFunc(func(l *logger) {
		l.callSkip = skip
	})
}

func Development() Option {
	return optionFunc(func(l *logger) {
		l.development = true
	})
}
