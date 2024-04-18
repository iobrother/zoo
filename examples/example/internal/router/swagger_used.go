//go:build swag
// +build swag

package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"

	"zcash-go/gen/docs"
)

func Swagger(r gin.IRouter) {
	swag.Register(swag.Name, new(docs.Admin))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
