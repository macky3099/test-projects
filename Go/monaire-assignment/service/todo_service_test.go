package service

import (
	"context"
	"testing"
	"todolist/models"
	"todolist/repository"
)

func TestCreateTodoSuccess(t *testing.T) {
	repo := repository.NewInMemoryTodoRepository()
	svc := NewTodoService(repo)

	req := models.CreateTodoRequest{
		Text:    "Finish assignment",
		DueDate: "2026-05-03T10:00:00Z",
	}

	todo, err := svc.CreateTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if todo.ID == "" {
		t.Fatal("expected todo ID to be generated")
	}

	if todo.Text != "Finish assignment" {
		t.Fatalf("expected text to match, got %s", todo.Text)
	}

	if todo.Completed {
		t.Fatal("expected new todo to be incomplete")
	}
}

func TestCreateTodoEmptyText(t *testing.T) {
	repo := repository.NewInMemoryTodoRepository()
	svc := NewTodoService(repo)

	req := models.CreateTodoRequest{
		Text:    "   ",
		DueDate: "2026-05-03T10:00:00Z",
	}

	_, err := svc.CreateTodo(context.Background(), req)
	if err != ErrInvalidText {
		t.Fatalf("expected ErrInvalidText, got %v", err)
	}
}

func TestCreateTodoInvalidDueDate(t *testing.T) {
	repo := repository.NewInMemoryTodoRepository()
	svc := NewTodoService(repo)

	req := models.CreateTodoRequest{
		Text:    "Finish assignment",
		DueDate: "03-05-2026",
	}

	_, err := svc.CreateTodo(context.Background(), req)
	if err != ErrInvalidDueDate {
		t.Fatalf("expected ErrInvalidDueDate, got %v", err)
	}
}

func TestListTodosSortedByDueDate(t *testing.T) {
	repo := repository.NewInMemoryTodoRepository()
	svc := NewTodoService(repo)

	_, _ = svc.CreateTodo(context.Background(), models.CreateTodoRequest{
		Text:    "Second task",
		DueDate: "2026-05-05T10:00:00Z",
	})

	_, _ = svc.CreateTodo(context.Background(), models.CreateTodoRequest{
		Text:    "First task",
		DueDate: "2026-05-03T10:00:00Z",
	})

	todos, err := svc.ListTodos(context.Background(), false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(todos) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(todos))
	}

	if todos[0].Text != "First task" {
		t.Fatalf("expected earliest due date first, got %s", todos[0].Text)
	}
}

func TestListTodosExcludesCompletedByDefault(t *testing.T) {
	repo := repository.NewInMemoryTodoRepository()
	svc := NewTodoService(repo)

	todo, _ := svc.CreateTodo(context.Background(), models.CreateTodoRequest{
		Text:    "Completed task",
		DueDate: "2026-05-03T10:00:00Z",
	})

	completed := true
	_, _ = svc.UpdateTodo(context.Background(), todo.ID, models.UpdateTodoRequest{
		Completed: &completed,
	})

	todos, err := svc.ListTodos(context.Background(), false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(todos) != 0 {
		t.Fatalf("expected completed todos to be excluded, got %d", len(todos))
	}
}

func TestListTodosIncludeCompleted(t *testing.T) {
	repo := repository.NewInMemoryTodoRepository()
	svc := NewTodoService(repo)

	todo, _ := svc.CreateTodo(context.Background(), models.CreateTodoRequest{
		Text:    "Completed task",
		DueDate: "2026-05-03T10:00:00Z",
	})

	completed := true
	_, _ = svc.UpdateTodo(context.Background(), todo.ID, models.UpdateTodoRequest{
		Completed: &completed,
	})

	todos, err := svc.ListTodos(context.Background(), true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(todos) != 1 {
		t.Fatalf("expected completed todo to be included, got %d", len(todos))
	}
}
