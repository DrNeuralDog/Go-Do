package models

import (
	"testing"
	"time"
)

func TestNewTodoItem(t *testing.T) {
	todo := NewTodoItem()

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
}

func TestTodoItemSetters(t *testing.T) {
	todo := NewTodoItem()

	// Test name setter
	testName := "Test Todo"
	todo.SetName(testName)
	if todo.Name != testName {
		t.Errorf("Expected name %s, got %s", testName, todo.Name)
	}

	// Test content setter
	testContent := "Test content"
	todo.SetContent(testContent)
	if todo.Content != testContent {
		t.Errorf("Expected content %s, got %s", testContent, todo.Content)
	}

	// Test place setter
	testPlace := "Test location"
	todo.SetPlace(testPlace)
	if todo.Place != testPlace {
		t.Errorf("Expected place %s, got %s", testPlace, todo.Place)
	}

	// Test label setter
	testLabel := "Test label"
	todo.SetLabel(testLabel)
	if todo.Label != testLabel {
		t.Errorf("Expected label %s, got %s", testLabel, todo.Label)
	}

	// Test kind setter
	todo.SetKind(1)
	if todo.Kind != 1 {
		t.Errorf("Expected kind 1, got %d", todo.Kind)
	}

	// Test level setter (should be clamped to 0-3)
	todo.SetLevel(5)
	if todo.Level != 0 {
		t.Errorf("Expected level 0 (clamped), got %d", todo.Level)
	}

	todo.SetLevel(2)
	if todo.Level != 2 {
		t.Errorf("Expected level 2, got %d", todo.Level)
	}

	// Test time setter
	testTime := time.Now()
	todo.SetTime(testTime)
	if !todo.TodoTime.Equal(testTime) {
		t.Error("Time was not set correctly")
	}

	// Test warn time setter
	testWarnTime := 60
	todo.SetWarnTime(testWarnTime)
	if todo.WarnTime != testWarnTime {
		t.Errorf("Expected warn time %d, got %d", testWarnTime, todo.WarnTime)
	}

	// Test mark as done
	todo.MarkAsDone(true)
	if !todo.Done {
		t.Error("Todo should be marked as done")
	}
}

func TestTodoItemGetters(t *testing.T) {
	todo := NewTodoItem()

	// Set values
	testName := "Test Todo"
	testContent := "Test content"
	testPlace := "Test location"
	testLabel := "Test label"
	testTime := time.Now()
	testWarnTime := 30

	todo.SetName(testName)
	todo.SetContent(testContent)
	todo.SetPlace(testPlace)
	todo.SetLabel(testLabel)
	todo.SetTime(testTime)
	todo.SetWarnTime(testWarnTime)
	todo.SetKind(1)
	todo.SetLevel(2)
	todo.MarkAsDone(true)

	// Test getters
	if todo.GetName() != testName {
		t.Errorf("GetName() failed: expected %s, got %s", testName, todo.GetName())
	}

	if todo.GetContent() != testContent {
		t.Errorf("GetContent() failed: expected %s, got %s", testContent, todo.GetContent())
	}

	if todo.GetPlace() != testPlace {
		t.Errorf("GetPlace() failed: expected %s, got %s", testPlace, todo.GetPlace())
	}

	if todo.GetLabel() != testLabel {
		t.Errorf("GetLabel() failed: expected %s, got %s", testLabel, todo.GetLabel())
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

	if todo.GetWarnTime() != testWarnTime {
		t.Errorf("GetWarnTime() failed: expected %d, got %d", testWarnTime, todo.GetWarnTime())
	}

	if !todo.HaveDone() {
		t.Error("HaveDone() failed: should return true")
	}
}

func TestTodoItemLevelStrings(t *testing.T) {
	todo := NewTodoItem()

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
		todo.SetLevel(test.level)
		result := todo.GetLevelString()
		if result != test.expected {
			t.Errorf("GetLevelString() for level %d: expected %s, got %s", test.level, test.expected, result)
		}
	}
}

func TestTodoItemLevelColors(t *testing.T) {
	todo := NewTodoItem()

	tests := []struct {
		level    int
		expected string
	}{
		{0, "#4CAF50"},
		{1, "#2196F3"},
		{2, "#FF9800"},
		{3, "#F44336"},
	}

	for _, test := range tests {
		todo.SetLevel(test.level)
		result := todo.GetLevelColor()
		if result != test.expected {
			t.Errorf("GetLevelColor() for level %d: expected %s, got %s", test.level, test.expected, result)
		}
	}
}

func TestTodoItemKindString(t *testing.T) {
	todo := NewTodoItem()

	todo.SetKind(0)
	if todo.GetKindString() != "Event" {
		t.Errorf("Expected 'Event', got %s", todo.GetKindString())
	}

	todo.SetKind(1)
	if todo.GetKindString() != "Task" {
		t.Errorf("Expected 'Task', got %s", todo.GetKindString())
	}
}

func TestTodoItemIsBefore(t *testing.T) {
	todo1 := NewTodoItem()
	todo2 := NewTodoItem()

	// Set different times
	now := time.Now()
	todo1.SetTime(now)
	todo2.SetTime(now.Add(time.Hour))

	if !todo2.IsBefore(todo1) {
		t.Error("todo2 should be before todo1")
	}

	if todo1.IsBefore(todo2) {
		t.Error("todo1 should not be before todo2")
	}
}

func TestTodoItemShouldRemind(t *testing.T) {
	todo := NewTodoItem()
	currentTime := time.Now()

	// Test with no warning time
	todo.SetWarnTime(0)
	if todo.ShouldRemind(currentTime) {
		t.Error("Should not remind when warn time is 0")
	}

	// Test with warning time but item is done
	todo.SetWarnTime(60)
	todo.MarkAsDone(true)
	if todo.ShouldRemind(currentTime) {
		t.Error("Should not remind when item is done")
	}

	// Test with warning time and item not done
	todo.MarkAsDone(false)
	futureTime := currentTime.Add(30 * time.Minute)
	todo.SetTime(futureTime)

	// Should remind if current time is after remind time but before todo time
	if !todo.ShouldRemind(currentTime) {
		t.Error("Should remind when current time is after remind time")
	}

	// Should not remind if current time is before remind time
	pastTime := currentTime.Add(-2 * time.Hour)
	if todo.ShouldRemind(pastTime) {
		t.Error("Should not remind when current time is before remind time")
	}
}
