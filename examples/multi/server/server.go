package main

import (
	"context"
	"fmt"

	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/log"
	httpserver "github.com/iobrother/zoo/core/transport/http/server"
	"github.com/iobrother/zoo/core/transport/http/server/middleware/error_response"
	"github.com/iobrother/zoo/examples/gen/greeter"
	"github.com/smallnest/rpcx/server"
)

func main() {
	app := zoo.New(
		zoo.InitRpcServer(InitRpcServer),
		zoo.InitHttpServer(InitHttpServer),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitRpcServer(s *server.Server) error {
	if err := s.RegisterName("Greeter", &GreeterImpl{}, ""); err != nil {
		return err
	}
	return nil
}

type GreeterImpl struct{}

func (s *GreeterImpl) SayHello(ctx context.Context, req *greeter.HelloRequest, rsp *greeter.HelloReply) error {
	*rsp = greeter.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}

func InitHttpServer(s *httpserver.Server) error {
	s.Use(error_response.ErrorResponse())
	g := s.Group("")
	greeter.RegisterGreeterHTTPService(g, &HttpGreeter{})

	return nil
}

type HttpGreeter struct {
}

func (s *HttpGreeter) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	rsp := &greeter.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return rsp, nil
}
