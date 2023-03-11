package main

import (
	"context"

	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/transport/rpc/client"
	"github.com/iobrother/zoo/examples/gen/api/greeter"
)

func main() {
	c, err := client.NewClient(
		client.WithServiceName("Greeter"),
		client.BasePath("/zoo"),
		client.EtcdAddr([]string{"127.0.0.1:2379"}),
	)
	if err != nil {
		log.Error(err)
		return
	}
	cli := greeter.NewGreeterClient(c.GetXClient())

	req := &greeter.HelloRequest{
		Name: "zoo",
	}

	rsp, err := cli.SayHello(context.Background(), req)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("reply: %s", rsp.Message)
}
