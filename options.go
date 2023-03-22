package zoo

import (
	httpserver "github.com/iobrother/zoo/core/transport/http/server"
	"github.com/iobrother/zoo/core/transport/rpc/server"
)

type Options struct {
	InitRpcServer  server.InitRpcServerFunc
	InitHttpServer httpserver.InitHttpServerFunc
	BeforeStart    []func() error
	AfterStart     []func() error
	BeforeStop     []func() error
	AfterStop      []func() error
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func InitRpcServer(f server.InitRpcServerFunc) Option {
	return func(o *Options) {
		o.InitRpcServer = f
	}
}

func InitHttpServer(f httpserver.InitHttpServerFunc) Option {
	return func(o *Options) {
		o.InitHttpServer = f
	}
}

func BeforeStart(f func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, f)
	}
}

func AfterStart(f func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, f)
	}
}

func BeforeStop(f func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, f)
	}
}

func AfterStop(f func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, f)
	}
}
