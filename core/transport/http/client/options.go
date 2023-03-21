package client

import (
	"github.com/iobrother/zoo/core/registry"
)

type Option interface {
	apply(c *Client)
}

type optionFunc func(c *Client)

func (f optionFunc) apply(c *Client) {
	f(c)
}

func WithServiceName(name string) Option {
	return optionFunc(func(c *Client) {
		c.serviceName = name
	})
}

func WithServiceAddr(addr string) Option {
	return optionFunc(func(c *Client) {
		c.serviceAddr = addr
	})
}

func Registry(r registry.Registry) Option {
	return optionFunc(func(c *Client) {
		c.registry = r
	})
}
