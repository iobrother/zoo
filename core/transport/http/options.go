package http

import "github.com/iobrother/zoo/core/registry"

type Options struct {
	Name           string
	Addr           string
	InitHttpServer InitHttpServerFunc
	Mode           string
	Tracing        bool
	Registry       registry.Registry
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func Name(s string) Option {
	return func(o *Options) {
		o.Name = s
	}
}

func Addr(s string) Option {
	return func(o *Options) {
		o.Addr = s
	}
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

type InitHttpServerFunc func(r *Server) error

func InitHttpServer(f InitHttpServerFunc) Option {
	return func(o *Options) {
		o.InitHttpServer = f
	}
}

func Mode(s string) Option {
	return func(o *Options) {
		o.Mode = s
	}
}

func Tracing(b bool) Option {
	return func(o *Options) {
		o.Tracing = b
	}
}
