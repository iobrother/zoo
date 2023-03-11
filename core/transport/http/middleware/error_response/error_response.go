package error_response

import (
	"github.com/iobrother/zoo/core/errors"
	"github.com/iobrother/zoo/core/transport/http"
)

func ErrorResponse() http.HandlerFunc {
	return func(c *http.Context) {
		defer func() {
			if c.GetError() != nil {
				e := errors.FromError(c.GetError())
				c.JSON(500, e)
				c.Abort()
			}
		}()

		c.Next()
	}
}
