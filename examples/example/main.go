package main

import (
	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/transport/http/server"
	"github.com/iobrother/zoo/examples/example/internal/app"
	"github.com/iobrother/zoo/examples/example/internal/router"
)

func main() {
	app := zoo.New(
		zoo.BeforeStart(app.Init),
		zoo.InitHttpServer(func(s *server.Server) error {
			router.Setup(s)
			return nil
		}))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
