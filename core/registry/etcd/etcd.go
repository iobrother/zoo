package etcd

import (
	"context"
	"encoding/json"
	"path"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/registry"
	"github.com/iobrother/zoo/core/util/backoff"
)

type etcdRegistry struct {
	client   *clientv3.Client
	leaseID  clientv3.LeaseID
	cfg      clientv3.Config
	basePath string
	ttl      int64
	timeout  time.Duration
	maxRetry int
	ctx      context.Context
}

func (e *etcdRegistry) Register(service *registry.Service) error {
	err := e.register(service)
	if err != nil {
		return err
	}

	go e.keepalive(service)

	return nil
}

func (e *etcdRegistry) serviceKey(service *registry.Service) string {
	return e.basePath + "/" + service.Name + "/" + service.Scheme + "@" + service.Address
}

func (e *etcdRegistry) register(service *registry.Service) error {
	err := e.grant()
	if err != nil {
		return err
	}
	b, err := json.Marshal(service)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	_, err = e.client.Put(ctx,
		e.serviceKey(service),
		string(b),
		clientv3.WithLease(e.leaseID),
	)
	return err
}

func (e *etcdRegistry) Deregister(service *registry.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	_, err := e.client.Delete(ctx,
		e.serviceKey(service),
	)
	return err
}

func (e *etcdRegistry) keepalive(service *registry.Service) {
	kch, err := e.client.KeepAlive(context.Background(), e.leaseID)
	if err != nil {
		e.leaseID = 0
	}

	for {
		if e.leaseID == 0 {
			for i := 0; i < e.maxRetry; i++ {
				if err = e.register(service); err != nil {
					time.Sleep(backoff.Do(i + 1))
					continue
				}

				kch, err = e.client.KeepAlive(context.Background(), e.leaseID)
				if err == nil {
					break
				}

				time.Sleep(backoff.Do(i + 1))
			}
		}

		select {
		case rsp := <-kch:
			if rsp == nil {
				e.leaseID = 0
				continue
			}
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *etcdRegistry) grant() error {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	rsp, err := e.client.Grant(ctx, e.ttl)
	if err != nil {
		log.Error(err)
		return err
	}
	e.leaseID = rsp.ID
	return nil
}

func NewRegistry(opts ...Option) registry.Registry {
	e := &etcdRegistry{
		cfg: clientv3.Config{
			Endpoints: []string{"127.0.0.1:2379"},
		},
		basePath: DefaultBasePath,
		ttl:      DefaultTTL,
		timeout:  DefaultTimeout,
		maxRetry: DefaultMaxRetry,
		ctx:      context.Background(),
	}
	for _, opt := range opts {
		opt.apply(e)
	}

	if err := e.init(); err != nil {
		log.Error(err)
		return nil
	}

	return e
}

func (e *etcdRegistry) init() error {
	cli, err := clientv3.New(e.cfg)
	if err != nil {
		return err
	}
	e.client = cli
	return nil
}

func (e *etcdRegistry) GetService(name string) ([]*registry.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	prefix := path.Join(e.basePath, name)
	rsp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}

	var services []*registry.Service
	for _, v := range rsp.Kvs {
		if service := decode(v.Value); service != nil {
			if service.Name == name {
				services = append(services, service)
			}
		}
	}

	return services, nil
}

func (e *etcdRegistry) Watch(service string) registry.Watcher {
	p := path.Join(e.basePath, service)
	return newWatcher(e.ctx, e.client, p, service)
}
