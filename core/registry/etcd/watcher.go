package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"path"

	"github.com/iobrother/zoo/core/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type watcher struct {
	w      clientv3.WatchChan
	client *clientv3.Client
	ctx    context.Context
	cancel func()
}

func newWatcher(ctx context.Context, c *clientv3.Client, basePath, service string) *watcher {
	w := &watcher{
		client: c,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	watchPath := path.Join(basePath, service) + "/"
	w.w = c.Watch(ctx, watchPath, clientv3.WithPrefix(), clientv3.WithPrevKV())

	return w
}

func decode(b []byte) *registry.Service {
	var s *registry.Service
	json.Unmarshal(b, &s)
	return s
}

func (w *watcher) Next() (*registry.Result, error) {
	for rsp := range w.w {
		if rsp.Err() != nil {
			return nil, rsp.Err()
		}
		if rsp.Canceled {
			return nil, errors.New("could not get next")
		}
		for _, e := range rsp.Events {
			service := decode(e.Kv.Value)
			var action string

			switch e.Type {
			case clientv3.EventTypePut:
				if e.IsCreate() {
					action = "create"
				} else if e.IsModify() {
					action = "update"
				}
			case clientv3.EventTypeDelete:
				action = "delete"
				service = decode(e.PrevKv.Value)
			}

			if service == nil {
				continue
			}

			return &registry.Result{Action: action, Service: service}, nil
		}
	}

	return nil, errors.New("could not get next")
}

func (w *watcher) Stop() {
	w.cancel()
	_ = w.client.Close()
}
