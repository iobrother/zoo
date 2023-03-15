package main

import (
	"context"
	"fmt"

	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/errors"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/examples/gen/api/errapi"
	"github.com/smallnest/rpcx/server"
)

func main() {
	app := zoo.New(zoo.InitRpcServer(InitRpcServer))

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func InitRpcServer(s *server.Server) error {
	if err := s.RegisterName("ErrAPI", &ErrImpl{}, ""); err != nil {
		return err
	}
	return nil
}

type ErrImpl struct {
}

var _ errapi.ErrAPIAble = &ErrImpl{}

func (s *ErrImpl) SayHello(ctx context.Context, req *errapi.HelloRequest, rsp *errapi.HelloReply) error {
	*rsp = errapi.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}

func (s *ErrImpl) TestError(ctx context.Context, req *errapi.ErrorRequest, rsp *errapi.ErrorReply) error {
	if req.Name == "unknown" {
		return fmt.Errorf("模拟的一个服务器未知错误")
	} else if req.Name == "db" {
		return errors.New(100101, "数据库错误", "这是模拟的一个数据库错误")
	} else if req.Name == "biz" {
		x := errors.Errorf(1000, "kkk", "kdadadf")
		log.Errorf("%+v", x)
		return errors.New(100201, "订单不存在", "订单不存在")
	}

	*rsp = errapi.ErrorReply{
		Message: fmt.Sprintf("[%s] 不是错误", req.Name),
	}

	return nil
}
