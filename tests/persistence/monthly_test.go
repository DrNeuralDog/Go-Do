package persistence_test

import (
	"testing"
	"time"

	"godo/src/models"
	"godo/src/persistence"
)

// TestMonthlyManagerUpdateTodoAllowsRename checks edit save after name change
func TestMonthlyManagerUpdateTodoAllowsRename(t *testing.T) {
	manager := persistence.NewMonthlyManager(t.TempDir())
	todoTime := time.Date(2026, time.May, 4, 12, 30, 0, 0, time.UTC)

	original := &models.TodoItem{
		Name:     "Old name",
		TodoTime: todoTime,
		Done:     true,
		Starred:  true,
		Order:    3,
	}

	if err := manager.AddTodo(original); err != nil {
		t.Fatalf("AddTodo() failed: %v", err)
	}

	updated := *original
	updated.Name = "New name"

	if err := manager.UpdateTodo(&updated, todoTime); err != nil {
		t.Fatalf("UpdateTodo() failed after rename: %v", err)
	}

	todos, err := manager.GetTodosForMonth(todoTime.Year(), int(todoTime.Month()))
	if err != nil {
		t.Fatalf("GetTodosForMonth() failed: %v", err)
	}

	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}

	if todos[0].Name != "New name" {
		t.Errorf("expected renamed todo, got %q", todos[0].Name)
	}

	if !todos[0].Done || !todos[0].Starred || todos[0].Order != 3 {
		t.Errorf("todo state was not preserved: %#v", todos[0])
	}
}
