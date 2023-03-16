package main

import (
	"context"
	"fmt"

	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/examples/gen/greeter"
	"github.com/smallnest/rpcx/server"
)

func main() {
	app := zoo.New(zoo.InitRpcServer(InitRpcServer))

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
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
