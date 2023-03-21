package main

import (
	"fmt"
	
	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/errors"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/transport/http/server"
)

func main() {
	app := zoo.New(zoo.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

type HelloRequest struct {
	Name string `json:"name,omitempty"`
}

type HelloResponse struct {
	Message string `json:"message,omitempty"`
}

func InitHttpServer(s *server.Server) error {
	s.PostEx("/hello", func(c *server.Context) {
		req := HelloRequest{}
		if err := c.ShouldBind(&req); err != nil {
			e := errors.FromError(err)
			c.JSON(500, e)
			c.Abort()
			return
		}

		rsp := HelloResponse{
			Message: fmt.Sprintf("hello %s!", req.Name),
		}
		c.JSON(200, rsp)
	})
	return nil
}
