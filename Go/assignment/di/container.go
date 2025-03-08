package di

import (
	"codeassign/config"
	"codeassign/repository"
	"codeassign/service"
)

// Container holds all application dependencies
type Container struct {
	UserRepository repository.UserRepository
	UserService    service.UserService
	Handler        *service.Handler
}

// NewContainer initializes all dependencies
func NewContainer() *Container {
	// Get database connection
	db := config.GetConfig()
	
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	
	// Initialize services
	userService := service.NewUserService(userRepo)
	
	// Initialize handlers
	handler := service.NewHandler(userService)
	
	return &Container{
		UserRepository: userRepo,
		UserService:    userService,
		Handler:        handler,
	}
}