package repository

import (
	"codeassign/models"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(user *models.UserDetails) error
	// Add other methods as needed: GetUser, UpdateUser, etc.
}

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	DB *gorm.DB
}

// CreateUser implements the user creation functionality
func (r *GormUserRepository) CreateUser(user *models.UserDetails) error {
	result := r.DB.Create(user)
	return result.Error
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{DB: db}
}