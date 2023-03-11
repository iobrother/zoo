package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iobrother/zoo"
	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/transport/http"
	"github.com/iobrother/zoo/core/transport/rpc/client"
	"github.com/iobrother/zoo/examples/gen/api/greeter"
)

// curl http://127.0.0.1:5180/hello/zoo
func main() {
	app := zoo.New(zoo.InitHttpServer(InitHttpServer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitHttpServer(s *http.Server) error {
	s.GET("/hello/:name", func(c *gin.Context) {
		cc, err := client.NewClient(
			client.WithServiceName("Greeter"),
			client.WithServiceAddr("127.0.0.1:5188"),
			client.Tracing(true),
		)
		if err != nil {
			log.Error(err)
			return
		}
		cli := greeter.NewGreeterClient(cc.GetXClient())

		args := &greeter.HelloRequest{
			Name: c.Param("name"),
		}

		log.Infof(args.Name)

		reply, err := cli.SayHello(c.Request.Context(), args)
		if err != nil {
			log.Error(err)
			return
		}

		c.String(200, reply.Message)
	})

	return nil
}
