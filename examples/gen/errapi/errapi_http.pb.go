// Code generated by protoc-gen-zoo-http. DO NOT EDIT.
// versions:
// - protoc-gen-zoo-http v0.1.0
// - protoc                (unknown)
// source: errapi/errapi.proto

package errapi

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	server "github.com/iobrother/zoo/core/transport/http/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = context.TODO
var _ = gin.New
var _ = server.NewServer

type ErrAPIHTTPService interface {
	// SayHello ...
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
	// TestError ...
	TestError(context.Context, *ErrorRequest) (*ErrorReply, error)
}

func RegisterErrAPIHTTPService(g *gin.RouterGroup, svc ErrAPIHTTPService) {
	r := g.Group("")
	r.GET("/hello/:name", _ErrAPI_SayHello0_HTTP_Handler(svc))
	r.GET("/error/:name", _ErrAPI_TestError0_HTTP_Handler(svc))
}

func _ErrAPI_SayHello0_HTTP_Handler(svc ErrAPIHTTPService) gin.HandlerFunc {
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
		if c.ContentType() == "application/x-protobuf" {
			c.ProtoBuf(200, rsp)
		} else {
			c.JSON(200, rsp)
		}
	}
}
func _ErrAPI_TestError0_HTTP_Handler(svc ErrAPIHTTPService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := &server.Context{Context: ctx}
		shouldBind := func(req *ErrorRequest) error {
			if err := c.ShouldBindQuery(req); err != nil {
				return err
			}
			if err := c.ShouldBindUri(req); err != nil {
				return err
			}
			return nil
		}

		var err error
		var req ErrorRequest
		var rsp *ErrorReply

		if err = shouldBind(&req); err != nil {
			c.SetError(err)
			return
		}
		rsp, err = svc.TestError(c.Request.Context(), &req)
		if err != nil {
			c.SetError(err)
			return
		}
		if c.ContentType() == "application/x-protobuf" {
			c.ProtoBuf(200, rsp)
		} else {
			c.JSON(200, rsp)
		}
	}
}
