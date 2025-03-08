package models

import (
	"time"

	"gorm.io/gorm"
)

type UserDetails struct {
	gorm.Model
	Name   string `json:"name"`
	PAN    string `json:"pan" validate:"required,pan"`
	Mobile int    `json:"mobile" validate:"required,mobile"`
	Email  string `json:"email" validate:"required,email"`
}

type CommonError struct {
	Status   int64     `json:"status"`
	Response string    `json:"response"`
	Created  time.Time `json:"created"`
}

// ErrHandler returns error message response
func ErrHandler(errmessage string) *CommonError {
	errresponse := CommonError{}
	errresponse.Status = 1
	errresponse.Response = errmessage
	errresponse.Created = time.Now()
	return &errresponse
}

// ErrHandler returns success message response
func SuccessResponse(message string) *CommonError {
	errresponse := CommonError{}
	errresponse.Status = 0
	errresponse.Response = message
	errresponse.Created = time.Now()
	return &errresponse
}
