// Code generated by protoc-gen-zoo-http. DO NOT EDIT.
// versions:
// - protoc-gen-zoo-http v0.1.0
// - protoc                (unknown)
// source: greeter/greeter.proto

package greeter

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	server "github.com/iobrother/zoo/core/transport/http/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = context.TODO
var _ = gin.New
var _ = server.NewServer

type GreeterHTTPService interface {
	// SayHello ...
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterHTTPService(g *gin.RouterGroup, svc GreeterHTTPService) {
	r := g.Group("")
	r.GET("/hello/:name", _Greeter_SayHello0_HTTP_Handler(svc))
}

func _Greeter_SayHello0_HTTP_Handler(svc GreeterHTTPService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := &server.Context{Context: ctx}
		shouldBind := func(req *HelloRequest) error {
			if err := c.ShouldBindQuery(req); err != nil {
				return err
			}
			if err := c.ShouldBindUri(req); err != nil {
				return err
			}
			return nil
		}

		var err error
		var req HelloRequest
		var rsp *HelloReply

		if err = shouldBind(&req); err != nil {
			c.SetError(err)
			return
		}
		rsp, err = svc.SayHello(c.Request.Context(), &req)
		if err != nil {
			c.SetError(err)
			return
		}
		c.JSON(200, rsp)
	}
}
