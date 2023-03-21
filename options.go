package zoo

import (
	server2 "github.com/iobrother/zoo/core/transport/http/server"
	"github.com/iobrother/zoo/core/transport/rpc/server"
)

type BeforeFunc func() error

type Options struct {
	InitRpcServer  server.InitRpcServerFunc
	InitHttpServer server2.InitHttpServerFunc
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

func InitHttpServer(f server2.InitHttpServerFunc) Option {
	return func(o *Options) {
		o.InitHttpServer = f
	}
}

func Before(f BeforeFunc) Option {
	return func(o *Options) {
		o.Before = f
	}
}
