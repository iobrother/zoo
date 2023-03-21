package registry

import "context"

type Registry interface {
	Register(service *Service) error
	Deregister(service *Service) error
}

type Service struct {
	Scheme   string
	Name     string
	Address  string
	Weight   int
	Metadata map[string]string
}

type Discovery interface {
	GetService(ctx context.Context, name string) ([]*Service, error)
	Watch(ctx context.Context, name string) (Watcher, error)
}

type Watcher interface {
	Next() ([]*Service, error)
	Stop() error
}
