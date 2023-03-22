package client

import (
	"context"
	"sync"

	"github.com/iobrother/zoo/core/registry"
)

type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}

type Resolver struct {
	target   *Target
	ctx      context.Context
	cancel   context.CancelFunc
	r        registry.Registry
	wg       sync.WaitGroup
	balancer Balancer
}

type Address struct {
	Addr     string
	Weight   int
	Metadata map[string]string
}

func newResolver(r registry.Registry, s Balancer, target *Target) *Resolver {
	ctx, cancel := context.WithCancel(context.Background())
	rr := &Resolver{
		balancer: s,
		target:   target,
		ctx:      ctx,
		cancel:   cancel,
		r:        r,
	}

	rr.wg.Add(1)
	rr.update()
	go rr.watcher()

	return rr
}

func (r *Resolver) watcher() {
	defer r.wg.Done()
	w := r.r.Watch(r.target.Endpoint)

	go func() {
		defer w.Stop()
		select {
		case <-r.ctx.Done():
			return
		}
	}()

	for {
		res, err := w.Next()
		if err != nil {
			// TODO: maybe return here
			break
		}

		if res.Service == nil {
			continue
		}
		r.update()
	}
}

func (r *Resolver) Close() {
	r.cancel()
	r.wg.Wait()
}

func (r *Resolver) update() {
	var addrs []*Address
	services, _ := r.r.GetService(r.target.Endpoint)
	for _, service := range services {
		if service.Scheme == "http" {
			addr := &Address{
				Addr:     service.Address,
				Weight:   service.Weight,
				Metadata: service.Metadata,
			}
			addrs = append(addrs, addr)
		}
	}

	r.balancer.Update(addrs)
}
