package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todolist/controller"
	"todolist/repository"
	"todolist/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	todoRepo := repository.NewInMemoryTodoRepository()
	todoService := service.NewTodoService(todoRepo)
	todoController := controller.NewTodoController(todoService)

	e.POST("/todos", todoController.CreateTodo)
	e.GET("/todos/:id", todoController.GetTodoByID)
	e.GET("/todos", todoController.ListTodos)
	e.PUT("/todos/:id", todoController.UpdateTodo)
	e.DELETE("/todos/:id", todoController.DeleteTodo)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s\n", port)
		if err := e.StartServer(srv); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
