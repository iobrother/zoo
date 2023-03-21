package etcd

import (
	"context"
	"time"
)

const (
	DefaultBasePath = "/zoo"
	DefaultTTL      = 15
	DefaultTimeout  = time.Second * 5
	DefaultMaxRetry = 10
)

type Option interface {
	apply(e *etcdRegistry)
}

type optionFunc func(e *etcdRegistry)

func (f optionFunc) apply(e *etcdRegistry) {
	f(e)
}

func Addrs(addrs ...string) Option {
	return optionFunc(func(e *etcdRegistry) {
		e.cfg.Endpoints = addrs
	})
}

func Timeout(t time.Duration) Option {
	return optionFunc(func(e *etcdRegistry) {
		e.timeout = t
	})
}

func RegisterTTL(ttl int64) Option {
	return optionFunc(func(e *etcdRegistry) {
		e.ttl = ttl
	})
}

func MaxRetry(maxRetry int) Option {
	return optionFunc(func(e *etcdRegistry) {
		e.maxRetry = maxRetry
	})
}

func Auth(username, password string) Option {
	return optionFunc(func(e *etcdRegistry) {
		e.cfg.Username = username
		e.cfg.Password = password
	})
}

func Context(ctx context.Context) Option {
	return optionFunc(func(e *etcdRegistry) {
		e.ctx = ctx
	})
}

func BasePath(basePath string) Option {
	return optionFunc(func(e *etcdRegistry) {
		e.basePath = basePath
	})
}
