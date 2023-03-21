package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/iobrother/zoo/core/errors"
	"github.com/iobrother/zoo/core/registry"
)

type Client struct {
	cc          *resty.Client
	serviceName string
	serviceAddr string
	registry    registry.Registry
	mws         []Middleware

	balancer Balancer
	resolver *Resolver
}

func NewClient(opts ...Option) *Client {
	cli := &Client{
		cc:       resty.New(),
		mws:      []Middleware{},
		balancer: &Wrr{},
	}

	for _, opt := range opts {
		opt.apply(cli)
	}

	if cli.registry != nil {
		cli.resolver = newResolver(cli.registry, cli.balancer, &Target{
			Scheme:    "http",
			Authority: "",
			Endpoint:  cli.serviceName,
		})
	}

	return cli
}

func (c *Client) Use(mws ...Middleware) *Client {
	return c
}

func (c *Client) Invoke(ctx context.Context, method, path string, req, rsp any) error {
	url := "http://"
	if c.serviceAddr != "" {
		url += c.serviceAddr
	} else {
		node, err := c.balancer.Select()
		if err != nil {
			return err
		}
		url += node.Addr
	}

	url += path

	h := func(ctx context.Context, in any) (any, error) {
		r := c.cc.R().SetContext(ctx)
		// TODO: encode
		r.SetHeader("Content-Type", "application/json")
		b, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		r = r.SetBody(string(b))

		result, err := r.Execute(method, url)
		if err != nil {
			return nil, err
		}

		if result.IsError() {
			return nil, errors.Parse(string(result.Body()))
		}

		if err := json.Unmarshal(result.Body(), rsp); err != nil {
			return nil, err
		}

		return rsp, nil
	}

	if len(c.mws) > 0 {
		Chain(c.mws...)(h)
	}

	_, err := h(ctx, req)

	return err
}

func (c *Client) Post(ctx context.Context, path string, req, rsp any) error {
	return c.Invoke(ctx, http.MethodPost, path, req, rsp)
}

func (c *Client) Get(ctx context.Context, path string, req, rsp any) error {
	return c.Invoke(ctx, http.MethodGet, path, req, rsp)
}

func (c *Client) Delete(ctx context.Context, path string, req, rsp any) error {
	return c.Invoke(ctx, http.MethodDelete, path, req, rsp)
}

func (c *Client) Patch(ctx context.Context, path string, req, rsp any) error {
	return c.Invoke(ctx, http.MethodPatch, path, req, rsp)
}

func (c *Client) Put(ctx context.Context, path string, req, rsp any) error {
	return c.Invoke(ctx, http.MethodPut, path, req, rsp)
}

func (c *Client) Options(ctx context.Context, path string, req, rsp any) error {
	return c.Invoke(ctx, http.MethodOptions, path, req, rsp)
}

func (c *Client) Head(ctx context.Context, path string, req, rsp any) error {
	return c.Invoke(ctx, http.MethodHead, path, req, rsp)
}
