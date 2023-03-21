package main

import (
	"fmt"

	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/transport/http/server"
)

// curl http://127.0.0.1:5180/hello/zoo
func main() {
	app := zoo.New(zoo.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(s *server.Server) error {
	s.GetEx("/hello/:name", func(c *server.Context) {
		c.String(200, fmt.Sprintf("hello %s!", c.Param("name")))
	})

	return nil
}
