package main

import (
	"context"
	"fmt"

	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/errors"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/transport/http"
	"github.com/iobrother/zoo/core/transport/http/middleware/error_response"
	"github.com/iobrother/zoo/core/transport/rpc/client"
	"github.com/iobrother/zoo/examples/gen/errapi"
)

// curl -w " status=%{http_code}" http://localhost:5180/error/unknown
// curl -w " status=%{http_code}" http://localhost:5180/error/db
// curl -w " status=%{http_code}" http://localhost:5180/error/biz
// curl -w " status=%{http_code}" http://localhost:5180/error/zoo

var cc *client.Client

func main() {
	app := zoo.New(zoo.InitHttpServer(InitHttpServer))

	cc, _ = client.NewClient(client.WithServiceName("ErrAPI"), client.WithServiceAddr("127.0.0.1:5188"))

	if cc == nil {
		log.Fatal("err")
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(s *http.Server) error {
	s.Use(error_response.ErrorResponse())
	g := s.Group("")
	errapi.RegisterErrAPIHTTPService(g, &ErrImpl{})

	return nil
}

type ErrImpl struct {
}

var _ errapi.ErrAPIHTTPService = &ErrImpl{}

func (s *ErrImpl) SayHello(ctx context.Context, req *errapi.HelloRequest) (*errapi.HelloReply, error) {
	cli := errapi.NewErrAPIClient(cc.GetXClient())

	request := errapi.HelloRequest{Name: req.Name}

	reply, err := cli.SayHello(ctx, &request)
	if err != nil {
		err = errors.WrapRpcError(err)
		log.Errorf("%+v", err)
		return nil, err
	}

	rsp := &errapi.HelloReply{
		Message: fmt.Sprintf("hello %s!", reply.Message),
	}
	return rsp, nil
}

func (s *ErrImpl) TestError(ctx context.Context, req *errapi.ErrorRequest) (*errapi.ErrorReply, error) {
	cli := errapi.NewErrAPIClient(cc.GetXClient())
	request := errapi.ErrorRequest{Name: req.Name}

	reply, err := cli.TestError(ctx, &request)
	if err != nil {
		err = errors.WrapRpcError(err)
		log.Errorf("%+v", err)
		return nil, err
	}

	rsp := &errapi.ErrorReply{
		Message: reply.Message,
	}
	return rsp, nil

}
