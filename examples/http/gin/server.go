package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/iobrother/zoo/core/transport/http/server"
)

func main() {
	srv := server.NewServer()
	srv.Use(func(c *gin.Context) {
		log.Println("Use")
	})

	srv.UseEx(func(c *server.Context) {
		log.Println("UseEx")
	})

	srv.UseEx(func(c *server.Context) {
	})

	srv.GET("/foo", func(c *gin.Context) {
		c.String(200, "foo")
	})

	srv.GetEx("/bar", func(c *server.Context) {
		c.String(200, "bar")
	})

	_ = srv.Run(":5180")
}
