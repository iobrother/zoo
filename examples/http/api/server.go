package main

import (
	"context"
	"fmt"

	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/transport/http"
	"github.com/iobrother/zoo/core/transport/http/middleware/error_response"
	"github.com/iobrother/zoo/examples/gen/greeter"
)

// curl http://127.0.0.1:5180/hello/zoo
func main() {
	app := zoo.New(zoo.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(r *http.Server) error {
	r.UseEx(error_response.ErrorResponse())
	greeter.RegisterGreeterHTTPService(r, &GreeterImpl{})

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
