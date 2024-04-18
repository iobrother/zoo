package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/iobrother/zoo/core/transport/http/server"
	"github.com/iobrother/zoo/core/transport/http/server/middleware/error_response"
	"github.com/iobrother/zoo/examples/example/pkg/validate"
)

func Setup(s *server.Server) {
	Swagger(s)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation(v)
	}

	s.NoMethod(func(ctx *gin.Context) {
		ctx.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	})

	s.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})

	s.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.Use(error_response.ErrorResponse())

	RegisterAPI(s)
}
