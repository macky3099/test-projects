package repository

import (
	"context"
	"errors"
	"sort"
	"sync"
	"todolist/models"
)

var ErrTodoNotFound = errors.New("todo not found")

type TodoRepository interface {
	Create(ctx context.Context, todo models.Todo) (models.Todo, error)
	GetByID(ctx context.Context, id string) (models.Todo, error)
	List(ctx context.Context, includeCompleted bool) ([]models.Todo, error)
	Update(ctx context.Context, todo models.Todo) (models.Todo, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryTodoRepository struct {
	mu    sync.RWMutex
	todos map[string]models.Todo
}

func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[string]models.Todo),
	}
}

func (r *InMemoryTodoRepository) Create(ctx context.Context, todo models.Todo) (models.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.todos[todo.ID] = todo
	return todo, nil
}

func (r *InMemoryTodoRepository) GetByID(ctx context.Context, id string) (models.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, ok := r.todos[id]
	if !ok {
		return models.Todo{}, ErrTodoNotFound
	}

	return todo, nil
}

func (r *InMemoryTodoRepository) List(ctx context.Context, includeCompleted bool) ([]models.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]models.Todo, 0, len(r.todos))

	for _, todo := range r.todos {
		if !includeCompleted && todo.Completed {
			continue
		}

		todos = append(todos, todo)
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].DueDate.Before(todos[j].DueDate)
	})

	return todos, nil
}

func (r *InMemoryTodoRepository) Update(ctx context.Context, todo models.Todo) (models.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[todo.ID]; !ok {
		return models.Todo{}, ErrTodoNotFound
	}

	r.todos[todo.ID] = todo
	return todo, nil
}

func (r *InMemoryTodoRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[id]; !ok {
		return ErrTodoNotFound
	}

	delete(r.todos, id)
	return nil
}
