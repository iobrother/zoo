package main

import (
	"context"
	"fmt"

	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/transport/http/server"
	"github.com/iobrother/zoo/core/transport/http/server/middleware/error_response"
	"github.com/iobrother/zoo/examples/gen/greeter"
)

// curl http://127.0.0.1:5180/hello/zoo
func main() {
	app := zoo.New(zoo.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(s *server.Server) error {
	s.Use(error_response.ErrorResponse())
	g := s.Group("")
	greeter.RegisterGreeterHTTPService(g, &GreeterImpl{})

	return nil
}

type GreeterImpl struct {
	greeter.GreeterHTTPService
}

func (s *GreeterImpl) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	rsp := &greeter.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	log.Info(req)
	return rsp, nil
}
