package controller

import (
	"errors"
	"net/http"
	"todolist/models"
	"todolist/repository"
	"todolist/service"

	"github.com/labstack/echo/v4"
)

type TodoController struct {
	service service.TodoService
}

func NewTodoController(service service.TodoService) *TodoController {
	return &TodoController{service: service}
}

func (tc *TodoController) CreateTodo(c echo.Context) error {
	var req models.CreateTodoRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	todo, err := tc.service.CreateTodo(c.Request().Context(), req)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, service.ErrInvalidText) || errors.Is(err, service.ErrInvalidDueDate) {
			status = http.StatusBadRequest
		}

		return c.JSON(status, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, todo)
}

func (tc *TodoController) GetTodoByID(c echo.Context) error {
	id := c.Param("id")

	todo, err := tc.service.GetTodoByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "todo not found",
		})
	}

	return c.JSON(http.StatusOK, todo)
}

func (tc *TodoController) ListTodos(c echo.Context) error {
	includeCompleted := c.QueryParam("include_completed") == "true"

	todos, err := tc.service.ListTodos(c.Request().Context(), includeCompleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to list todos",
		})
	}

	return c.JSON(http.StatusOK, todos)
}

func (tc *TodoController) UpdateTodo(c echo.Context) error {
	id := c.Param("id")

	var req models.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	todo, err := tc.service.UpdateTodo(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidText) || errors.Is(err, service.ErrInvalidDueDate) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		if errors.Is(err, repository.ErrTodoNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "todo not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update todo",
		})
	}

	return c.JSON(http.StatusOK, todo)
}

func (tc *TodoController) DeleteTodo(c echo.Context) error {
	id := c.Param("id")

	err := tc.service.DeleteTodo(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "todo not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete todo",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
