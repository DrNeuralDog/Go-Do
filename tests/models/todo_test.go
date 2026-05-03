package models_test

import (
	"image/color"
	"testing"
	"time"

	"godo/src/models"
)

// TestNewTodoItem checks default todo item values
func TestNewTodoItem(t *testing.T) {
	todo := models.NewTodoItem()

	if todo == nil {
		t.Fatal("NewTodoItem() returned nil")
	}

	if todo.Done {
		t.Error("New todo item should not be done by default")
	}

	if todo.Kind != 0 {
		t.Errorf("Expected kind 0, got %d", todo.Kind)
	}

	if todo.Level != 0 {
		t.Errorf("Expected level 0, got %d", todo.Level)
	}

	if todo.WarnTime != 0 {
		t.Errorf("Expected warn time 0, got %d", todo.WarnTime)
	}

	if todo.Starred {
		t.Error("New todo item should not be starred by default")
	}

	if todo.Order != 0 {
		t.Errorf("Expected order 0, got %d", todo.Order)
	}
}

// TestTodoItemSetters checks field setters
func TestTodoItemSetters(t *testing.T) {
	todo := models.NewTodoItem()
	testTime := time.Date(2026, time.April, 29, 12, 30, 0, 0, time.UTC)

	todo.SetName("Test Todo")
	todo.SetContent("Test content")
	todo.SetPlace("Test location")
	todo.SetLabel("Test label")
	todo.SetKind(1)
	todo.SetLevel(2)
	todo.SetTime(testTime)
	todo.SetWarnTime(60)
	todo.MarkAsDone(true)
	todo.SetOrder(10)

	if todo.Name != "Test Todo" {
		t.Errorf("Expected name %s, got %s", "Test Todo", todo.Name)
	}

	if todo.Content != "Test content" {
		t.Errorf("Expected content %s, got %s", "Test content", todo.Content)
	}

	if todo.Place != "Test location" {
		t.Errorf("Expected place %s, got %s", "Test location", todo.Place)
	}

	if todo.Label != "Test label" {
		t.Errorf("Expected label %s, got %s", "Test label", todo.Label)
	}

	if todo.Kind != 1 {
		t.Errorf("Expected kind 1, got %d", todo.Kind)
	}

	if todo.Level != 2 {
		t.Errorf("Expected level 2, got %d", todo.Level)
	}

	if !todo.TodoTime.Equal(testTime) {
		t.Error("Time was not set correctly")
	}

	if todo.WarnTime != 60 {
		t.Errorf("Expected warn time 60, got %d", todo.WarnTime)
	}

	if !todo.Done {
		t.Error("Todo should be marked as done")
	}

	if todo.Order != 10 {
		t.Errorf("Expected order 10, got %d", todo.Order)
	}
}

// TestTodoItemSetLevelIgnoresInvalidValues checks level range guard
func TestTodoItemSetLevelIgnoresInvalidValues(t *testing.T) {
	todo := models.NewTodoItem()
	todo.SetLevel(2)

	todo.SetLevel(5)

	if todo.Level != 2 {
		t.Errorf("Expected invalid high level to be ignored, got %d", todo.Level)
	}

	todo.SetLevel(-1)

	if todo.Level != 2 {
		t.Errorf("Expected invalid low level to be ignored, got %d", todo.Level)
	}
}

// TestTodoItemGetters checks getter methods
func TestTodoItemGetters(t *testing.T) {
	testTime := time.Date(2026, time.April, 29, 18, 45, 0, 0, time.UTC)
	todo := &models.TodoItem{
		Name:     "Test Todo",
		Content:  "Test content",
		Place:    "Test location",
		Label:    "Test label",
		Kind:     1,
		Level:    2,
		TodoTime: testTime,
		Done:     true,
		WarnTime: 30,
		Order:    4,
	}

	if todo.GetName() != "Test Todo" {
		t.Errorf("GetName() failed: expected %s, got %s", "Test Todo", todo.GetName())
	}

	if todo.GetContent() != "Test content" {
		t.Errorf("GetContent() failed: expected %s, got %s", "Test content", todo.GetContent())
	}

	if todo.GetPlace() != "Test location" {
		t.Errorf("GetPlace() failed: expected %s, got %s", "Test location", todo.GetPlace())
	}

	if todo.GetLabel() != "Test label" {
		t.Errorf("GetLabel() failed: expected %s, got %s", "Test label", todo.GetLabel())
	}

	if todo.GetKind() != 1 {
		t.Errorf("GetKind() failed: expected 1, got %d", todo.GetKind())
	}

	if todo.GetLevel() != 2 {
		t.Errorf("GetLevel() failed: expected 2, got %d", todo.GetLevel())
	}

	if !todo.GetTime().Equal(testTime) {
		t.Error("GetTime() failed: time not equal")
	}

	if todo.GetWarnTime() != 30 {
		t.Errorf("GetWarnTime() failed: expected 30, got %d", todo.GetWarnTime())
	}

	if !todo.IsDone() {
		t.Error("IsDone() failed: should return true")
	}

	if !todo.HaveDone() {
		t.Error("HaveDone() failed: should return true")
	}

	if todo.GetOrder() != 4 {
		t.Errorf("GetOrder() failed: expected 4, got %d", todo.GetOrder())
	}
}

// TestTodoItemLevelStrings checks priority text values
func TestTodoItemLevelStrings(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{0, "Not Important - Not Urgent"},
		{1, "Not Important - Urgent"},
		{2, "Important - Not Urgent"},
		{3, "Important - Urgent"},
		{5, "Unknown"},
	}

	for _, test := range tests {
		todo := &models.TodoItem{Level: test.level}
		result := todo.GetLevelString()

		if result != test.expected {
			t.Errorf("GetLevelString() for level %d: expected %s, got %s", test.level, test.expected, result)
		}
	}
}

// TestTodoItemLevelColors checks priority color values
func TestTodoItemLevelColors(t *testing.T) {
	tests := []struct {
		level    int
		expected color.RGBA
	}{
		{0, color.RGBA{R: 184, G: 187, B: 38, A: 255}},
		{1, color.RGBA{R: 131, G: 165, B: 152, A: 255}},
		{2, color.RGBA{R: 254, G: 128, B: 25, A: 255}},
		{3, color.RGBA{R: 251, G: 73, B: 52, A: 255}},
		{5, color.RGBA{R: 184, G: 187, B: 38, A: 255}},
	}

	for _, test := range tests {
		todo := &models.TodoItem{Level: test.level}
		result := todo.GetLevelColor()

		if result != test.expected {
			t.Errorf("GetLevelColor() for level %d: expected %#v, got %#v", test.level, test.expected, result)
		}
	}
}

// TestTodoItemKindString checks kind labels
func TestTodoItemKindString(t *testing.T) {
	tests := []struct {
		kind     int
		expected string
	}{
		{0, "Event"},
		{1, "Task"},
		{9, "Task"},
	}

	for _, test := range tests {
		todo := &models.TodoItem{Kind: test.kind}
		result := todo.GetKindString()

		if result != test.expected {
			t.Errorf("GetKindString() for kind %d: expected %s, got %s", test.kind, test.expected, result)
		}
	}
}

// TestTodoItemIsBefore checks chronological comparison
func TestTodoItemIsBefore(t *testing.T) {
	now := time.Date(2026, time.April, 29, 10, 0, 0, 0, time.UTC)
	earlier := &models.TodoItem{TodoTime: now}
	later := &models.TodoItem{TodoTime: now.Add(time.Hour)}

	if !earlier.IsBefore(later) {
		t.Error("Earlier todo should be before later todo")
	}

	if later.IsBefore(earlier) {
		t.Error("Later todo should not be before earlier todo")
	}
}

// TestTodoItemShouldRemind checks reminder window logic
func TestTodoItemShouldRemind(t *testing.T) {
	now := time.Date(2026, time.April, 29, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		todo     *models.TodoItem
		at       time.Time
		expected bool
	}{
		{
			name:     "without warn time",
			todo:     &models.TodoItem{TodoTime: now.Add(time.Hour)},
			at:       now,
			expected: false,
		},
		{
			name:     "done todo",
			todo:     &models.TodoItem{Done: true, WarnTime: 60, TodoTime: now.Add(time.Hour)},
			at:       now,
			expected: false,
		},
		{
			name:     "inside reminder window",
			todo:     &models.TodoItem{WarnTime: 60, TodoTime: now.Add(30 * time.Minute)},
			at:       now,
			expected: true,
		},
		{
			name:     "before reminder window",
			todo:     &models.TodoItem{WarnTime: 30, TodoTime: now.Add(time.Hour)},
			at:       now,
			expected: false,
		},
		{
			name:     "after todo time",
			todo:     &models.TodoItem{WarnTime: 60, TodoTime: now.Add(-time.Minute)},
			at:       now,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.todo.ShouldRemind(test.at)

			if result != test.expected {
				t.Errorf("ShouldRemind(): expected %v, got %v", test.expected, result)
			}
		})
	}
}
