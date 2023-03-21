package error_response

import (
	"github.com/gin-gonic/gin"
	"github.com/iobrother/zoo/core/errors"
	"github.com/iobrother/zoo/core/transport/http/server"
)

//func ErrorResponse() http.HandlerFunc {
//	return func(c *http.Context) {
//		defer func() {
//			err := c.GetError()
//			if err != nil {
//				e := errors.FromError(err)
//				delete(e.Metadata, "_zoo_error_stack")
//				c.JSON(500, e)
//				c.Abort()
//			}
//		}()
//
//		c.Next()
//	}
//}

func ErrorResponse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			c := server.Context{Context: ctx}
			err := c.GetError()
			if err != nil {
				e := errors.FromError(err)
				delete(e.Metadata, "_zoo_error_stack")
				c.JSON(500, e)
				c.Abort()
			}
		}()

		ctx.Next()
	}
}
