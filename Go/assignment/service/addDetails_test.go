package service

import (
	"bytes"
	"codeassign/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) AddUser(details *UserDetails) (*models.UserDetails, error) {
	args := m.Called(details)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserDetails), args.Error(1)
}

func TestAddDetails_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	expectedResult := &models.UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}
	mockService.On("AddUser", mock.AnythingOfType("*service.UserDetails")).Return(expectedResult, nil)

	handler := NewHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}
	jsonData, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest("POST", "/api/add", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AddDetails(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var responseBody models.UserDetails
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.Name, responseBody.Name)
	assert.Equal(t, reqBody.PAN, responseBody.PAN)
	assert.Equal(t, reqBody.Mobile, responseBody.Mobile)
	assert.Equal(t, reqBody.Email, responseBody.Email)

	mockService.AssertExpectations(t)
}

func TestAddDetails_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	handler := NewHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	invalidJSON := []byte(`{"name": "John Doe", "pan": "ABCDE1234F", "mobile": 9876543210, "email": }`) // Missing value for email
	c.Request, _ = http.NewRequest("POST", "/api/add", bytes.NewBuffer(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AddDetails(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var responseBody models.CommonError
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), responseBody.Status)
	assert.Equal(t, "Error in decoding json", responseBody.Response)

	mockService.AssertNotCalled(t, "AddUser")
}

func TestAddDetails_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)

	validate := validator.New()
	validationErr := validate.Struct(struct {
		PAN string `validate:"required"`
	}{})
	mockService.On("AddUser", mock.AnythingOfType("*service.UserDetails")).Return(nil, validationErr)
	handler := NewHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := UserDetails{
		Name:   "John Doe",
		PAN:    "invalid",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}
	jsonData, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest("POST", "/api/add", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AddDetails(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var responseBody models.CommonError
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), responseBody.Status)
	assert.Contains(t, responseBody.Response, "Validation failed for field")

	mockService.AssertExpectations(t)
}

func TestAddDetails_OtherError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	mockService.On("AddUser", mock.AnythingOfType("*service.UserDetails")).Return(nil, errors.New("database error"))

	handler := NewHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	reqBody := UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}
	jsonData, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest("POST", "/api/add", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.AddDetails(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
