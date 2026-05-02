package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"todolist/models"
	"todolist/repository"

	"github.com/google/uuid"
)

var (
	ErrInvalidText    = errors.New("todo text cannot be empty")
	ErrInvalidDueDate = errors.New("invalid due_date, expected RFC3339 format")
)

type TodoService interface {
	CreateTodo(ctx context.Context, req models.CreateTodoRequest) (models.Todo, error)
	GetTodoByID(ctx context.Context, id string) (models.Todo, error)
	ListTodos(ctx context.Context, includeCompleted bool) ([]models.Todo, error)
	UpdateTodo(ctx context.Context, id string, req models.UpdateTodoRequest) (models.Todo, error)
	DeleteTodo(ctx context.Context, id string) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(ctx context.Context, req models.CreateTodoRequest) (models.Todo, error) {
	if strings.TrimSpace(req.Text) == "" {
		return models.Todo{}, ErrInvalidText
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		return models.Todo{}, ErrInvalidDueDate
	}

	now := time.Now().UTC()

	todo := models.Todo{
		ID:        uuid.NewString(),
		Text:      strings.TrimSpace(req.Text),
		DueDate:   dueDate,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(ctx, todo)
}

func (s *todoService) GetTodoByID(ctx context.Context, id string) (models.Todo, error) {
	if strings.TrimSpace(id) == "" {
		return models.Todo{}, repository.ErrTodoNotFound
	}

	return s.repo.GetByID(ctx, id)
}

func (s *todoService) ListTodos(ctx context.Context, includeCompleted bool) ([]models.Todo, error) {
	return s.repo.List(ctx, includeCompleted)
}

func (s *todoService) UpdateTodo(ctx context.Context, id string, req models.UpdateTodoRequest) (models.Todo, error) {
	if strings.TrimSpace(id) == "" {
		return models.Todo{}, repository.ErrTodoNotFound
	}

	existingTodo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.Todo{}, err
	}

	if req.Text != nil {
		trimmedText := strings.TrimSpace(*req.Text)
		if trimmedText == "" {
			return models.Todo{}, ErrInvalidText
		}

		existingTodo.Text = trimmedText
	}

	if req.DueDate != nil {
		dueDate, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return models.Todo{}, ErrInvalidDueDate
		}

		existingTodo.DueDate = dueDate
	}

	if req.Completed != nil {
		existingTodo.Completed = *req.Completed
	}

	existingTodo.UpdatedAt = time.Now().UTC()

	return s.repo.Update(ctx, existingTodo)
}

func (s *todoService) DeleteTodo(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return repository.ErrTodoNotFound
	}

	return s.repo.Delete(ctx, id)
}
