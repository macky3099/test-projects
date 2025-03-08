package service

import (
	"codeassign/models"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.UserDetails) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestAddUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.UserDetails")).Return(nil)

	service := NewUserService(mockRepo)

	details := &UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}

	result, err := service.AddUser(details)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, details.Name, result.Name)
	assert.Equal(t, details.PAN, result.PAN)
	assert.Equal(t, details.Mobile, result.Mobile)
	assert.Equal(t, details.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestAddUser_InvalidPAN(t *testing.T) {

	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	details := &UserDetails{
		Name:   "John Doe",
		PAN:    "invalid",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}

	result, err := service.AddUser(details)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, validator.ValidationErrors{}, err)
	mockRepo.AssertNotCalled(t, "CreateUser")
}

func TestAddUser_InvalidMobile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	details := &UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 1234567890,
		Email:  "john.doe@example.com",
	}

	result, err := service.AddUser(details)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, validator.ValidationErrors{}, err)
	mockRepo.AssertNotCalled(t, "CreateUser")
}

func TestAddUser_InvalidEmail(t *testing.T) {
	// Setup mock repository
	mockRepo := new(MockUserRepository)

	service := NewUserService(mockRepo)
	details := &UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 9876543210,
		Email:  "invalid-email",
	}

	result, err := service.AddUser(details)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, validator.ValidationErrors{}, err)
	mockRepo.AssertNotCalled(t, "CreateUser")
}

func TestAddUser_RepositoryError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.UserDetails")).Return(errors.New("database error"))

	service := NewUserService(mockRepo)

	details := &UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}

	result, err := service.AddUser(details)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestPanValidator(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("pan", panValidator)

	type TestStruct struct {
		PAN string `validate:"pan"`
	}

	// Valid PAN
	valid := TestStruct{PAN: "ABCDE1234F"}
	err := validate.Struct(valid)
	assert.NoError(t, err)

	// Invalid PAN - wrong format
	invalid1 := TestStruct{PAN: "ABC123"}
	err = validate.Struct(invalid1)
	assert.Error(t, err)

	// Invalid PAN - lowercase letters
	invalid2 := TestStruct{PAN: "abcde1234f"}
	err = validate.Struct(invalid2)
	assert.Error(t, err)

	// Invalid PAN - wrong pattern
	invalid3 := TestStruct{PAN: "12345ABCDE"}
	err = validate.Struct(invalid3)
	assert.Error(t, err)
}

func TestMobileValidator(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("mobile", mobileValidator)

	type TestStruct struct {
		Mobile int `validate:"mobile"`
	}

	// Valid mobile - starts with 9
	valid1 := TestStruct{Mobile: 9876543210}
	err := validate.Struct(valid1)
	assert.NoError(t, err)

	// Valid mobile - starts with 8
	valid2 := TestStruct{Mobile: 8765432109}
	err = validate.Struct(valid2)
	assert.NoError(t, err)

	// Valid mobile - starts with 7
	valid3 := TestStruct{Mobile: 7654321098}
	err = validate.Struct(valid3)
	assert.NoError(t, err)

	// Valid mobile - starts with 6
	valid4 := TestStruct{Mobile: 6543210987}
	err = validate.Struct(valid4)
	assert.NoError(t, err)

	// Invalid mobile - starts with 5
	invalid1 := TestStruct{Mobile: 5432109876}
	err = validate.Struct(invalid1)
	assert.Error(t, err)

	// Invalid mobile - too short
	invalid2 := TestStruct{Mobile: 987654321}
	err = validate.Struct(invalid2)
	assert.Error(t, err)

	// Invalid mobile - too long
	invalid3 := TestStruct{Mobile: 98765432109}
	err = validate.Struct(invalid3)
	assert.Error(t, err)
}
