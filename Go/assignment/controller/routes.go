package controller

import (
	"codeassign/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, handler *service.Handler) {
	routes := router.Group("/api")
	routes.POST("/add", handler.AddDetails)
}