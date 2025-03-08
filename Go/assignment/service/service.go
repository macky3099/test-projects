package service

import (
	"codeassign/models"
	"codeassign/repository"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	AddUser(details *UserDetails) (*models.UserDetails, error)
	// Add other methods as needed
}

// UserServiceImpl implements UserService
type UserServiceImpl struct {
	userRepo repository.UserRepository
	validate *validator.Validate
}

// NewUserService creates a new user service
func NewUserService(repo repository.UserRepository) UserService {
	validate := registerValidator()
	return &UserServiceImpl{
		userRepo: repo,
		validate: validate,
	}
}

// AddUser implements the user addition logic
func (s *UserServiceImpl) AddUser(details *UserDetails) (*models.UserDetails, error) {
	userDetails := &models.UserDetails{
		Name:   details.Name,
		PAN:    details.PAN,
		Mobile: details.Mobile,
		Email:  details.Email,
	}
	
	err := s.validate.Struct(userDetails)
	if err != nil {
		return nil, err
	}
	
	err = s.userRepo.CreateUser(userDetails)
	if err != nil {
		return nil, err
	}
	
	return userDetails, nil
}

func registerValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("pan", panValidator)
	validate.RegisterValidation("mobile", mobileValidator)
	return validate
}

func panValidator(fl validator.FieldLevel) bool {
	panRegex := `^[A-Z]{5}[0-9]{4}[A-Z]{1}$`
	match, _ := regexp.MatchString(panRegex, fl.Field().String())
	return match
}

func mobileValidator(fl validator.FieldLevel) bool {
	mobile := fl.Field().Int()                             // Get the integer value
	mobileStr := fmt.Sprintf("%d", mobile)                 // Convert to string
	mobileRegex := `^[6-9]\d{9}$`                          // Regex for 10-digit number starting with 6-9
	match, _ := regexp.MatchString(mobileRegex, mobileStr) // Match against the regex
	return match
}