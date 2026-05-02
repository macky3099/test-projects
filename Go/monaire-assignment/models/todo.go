package models

import "time"

type Todo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	DueDate   time.Time `json:"due_date"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTodoRequest struct {
	Text    string `json:"text"`
	DueDate string `json:"due_date"`
}

type UpdateTodoRequest struct {
	Text      *string `json:"text,omitempty"`
	DueDate   *string `json:"due_date,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}
