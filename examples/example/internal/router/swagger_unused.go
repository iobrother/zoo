//go:build !swag
// +build !swag

package router

import "github.com/gin-gonic/gin"

func Swagger(gin.IRouter) {}
