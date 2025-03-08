package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestErrHandler(t *testing.T) {
	// Test with a simple error message
	errorMsg := "Something went wrong"
	result := ErrHandler(errorMsg)

	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Status)
	assert.Equal(t, errorMsg, result.Response)
	assert.WithinDuration(t, time.Now(), result.Created, 2*time.Second)

	// Test with an empty error message
	emptyError := ErrHandler("")
	assert.Equal(t, "", emptyError.Response)
	assert.Equal(t, int64(1), emptyError.Status)
}

func TestSuccessResponse(t *testing.T) {
	// Test with a simple success message
	successMsg := "Operation successful"
	result := SuccessResponse(successMsg)

	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Status)
	assert.Equal(t, successMsg, result.Response)
	assert.WithinDuration(t, time.Now(), result.Created, 2*time.Second)

	// Test with an empty success message
	emptySuccess := SuccessResponse("")
	assert.Equal(t, "", emptySuccess.Response)
	assert.Equal(t, int64(0), emptySuccess.Status)
}

func TestUserDetailsStruct(t *testing.T) {
	user := UserDetails{
		Name:   "John Doe",
		PAN:    "ABCDE1234F",
		Mobile: 9876543210,
		Email:  "john.doe@example.com",
	}

	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "ABCDE1234F", user.PAN)
	assert.Equal(t, 9876543210, user.Mobile)
	assert.Equal(t, "john.doe@example.com", user.Email)
}
