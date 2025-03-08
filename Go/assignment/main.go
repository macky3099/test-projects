package main

import (
	"codeassign/controller"
	"codeassign/di"
	"codeassign/security"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize dependency container
	container := di.NewContainer()
	
	router := gin.Default()
	router.Use(security.JsonLoggerMiddleware())
	
	// Register routes with handler from container
	controller.RegisterRoutes(router, container.Handler)
	
	server := &http.Server{
		Addr:           ":8888",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}