package service

import (
	"codeassign/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserDetails struct {
	Name   string `json:"name"`
	PAN    string `json:"pan" validate:"required,pan"`
	Mobile int    `json:"mobile" validate:"required,mobile"`
	Email  string `json:"email" validate:"required,email"`
}

type Handler struct {
	userService UserService
}

func NewHandler(userService UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) AddDetails(ctx *gin.Context) {
	body := UserDetails{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrHandler("Error in decoding json"))
		return
	}

	_, err := h.userService.AddUser(&body)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrors {
				ctx.JSON(http.StatusInternalServerError, models.ErrHandler(
					fmt.Sprintf("Validation failed for field '%s, expected pattern didn't match.\n", err.Field())))
				return
			}
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, models.SuccessResponse("Successfully created user"))
}
