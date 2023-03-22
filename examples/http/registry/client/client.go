package main

import (
	"context"

	"github.com/iobrother/zoo/core/log"
	"github.com/iobrother/zoo/core/registry/etcd"
	"github.com/iobrother/zoo/core/transport/http/client"
)

type HelloRequest struct {
	Name string `json:"name,omitempty"`
}

type HelloResponse struct {
	Message string `json:"message,omitempty"`
}

func main() {
	r := etcd.NewRegistry()

	cli := client.NewClient(
		client.WithServiceName("example"),
		client.Registry(r),
	)

	for i := 0; i < 5; i++ {
		req := HelloRequest{Name: "zoo"}
		rsp := HelloResponse{}

		if err := cli.Invoke(context.Background(), "POST", "/hello", &req, &rsp); err != nil {
			log.Error(err)
			return
		}

		log.Info(rsp)
	}
}
