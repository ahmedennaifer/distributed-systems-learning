package internal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	Status     string         `json:"status"`
	Type       string         `json:"type"`
	Params     map[string]any `json:"body"`
	CreatedAt  time.Time      `json:"created_at"`
	FinishedAt time.Time      `json:"finished_at"`
}

func NewTask() *Task {
	return &Task{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Status:    "Pending",
	}
}

func fieldEmptyError(field string) error {
	return fmt.Errorf("error: field %v cannot be empty", field)
}

func (t *Task) ValidateFields() error {
	if t.Name == "" {
		return fieldEmptyError("name")
	}
	if t.Type == "" {
		return fieldEmptyError("type")
	}
	if len(t.Params) == 0 {
		return fieldEmptyError("body")
	}
	return nil
}
