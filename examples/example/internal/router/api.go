package router

import (
	"github.com/iobrother/zoo/core/transport/http/server"
	"github.com/iobrother/zoo/examples/example/internal/app"
	"github.com/iobrother/zoo/examples/gen/example"
)

func RegisterAPI(s *server.Server) {
	appContext := app.Context()
	g := s.Group("")
	example.RegisterExampleHTTPService(g, appContext.Service.Example)
}
