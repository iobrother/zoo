package zoo

import (
	"github.com/iobrother/zoo/core/transport/http"
	"github.com/iobrother/zoo/core/transport/rpc/server"
)

type BeforeFunc func() error

type Options struct {
	InitRpcServer  server.InitRpcServerFunc
	InitHttpServer http.InitHttpServerFunc
	Before         BeforeFunc
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

func InitHttpServer(f http.InitHttpServerFunc) Option {
	return func(o *Options) {
		o.InitHttpServer = f
	}
}

func Before(f BeforeFunc) Option {
	return func(o *Options) {
		o.Before = f
	}
}
