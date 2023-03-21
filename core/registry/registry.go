package registry

type Registry interface {
	Register(service *Service) error
	Deregister(service *Service) error
	GetService(name string) ([]*Service, error)
	Watch(name string) Watcher
}

type Service struct {
	Scheme   string
	Name     string
	Address  string
	Weight   int
	Metadata map[string]string
}

type Watcher interface {
	Next() (*Result, error)
	Stop()
}

type Result struct {
	Action  string
	Service *Service
}
