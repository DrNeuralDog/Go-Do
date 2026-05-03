package models_test

import (
	"testing"
	"time"

	"godo/src/models"
)

// TestViewModeLabels checks visible labels for all view modes
func TestViewModeLabels(t *testing.T) {
	tests := []struct {
		mode     models.ViewMode
		expected string
	}{
		{models.ViewAll, "All"},
		{models.ViewIncomplete, "Incomplete"},
		{models.ViewComplete, "Complete"},
		{models.ViewStarred, "Important"},
	}

	for _, test := range tests {
		result := test.mode.GetLabel()

		if result != test.expected {
			t.Errorf("GetLabel() for %v: expected %s, got %s", test.mode, test.expected, result)
		}
	}
}

// TestViewModeCycle checks view mode switching order
func TestViewModeCycle(t *testing.T) {
	tests := []struct {
		name     string
		mode     models.ViewMode
		expected models.ViewMode
	}{
		{"all to incomplete", models.ViewAll, models.ViewIncomplete},
		{"incomplete to complete", models.ViewIncomplete, models.ViewComplete},
		{"complete to starred", models.ViewComplete, models.ViewStarred},
		{"starred to all", models.ViewStarred, models.ViewAll},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.mode.GetNextMode()

			if result != test.expected {
				t.Errorf("GetNextMode(): expected %v, got %v", test.expected, result)
			}
		})
	}
}

// TestViewModeString checks persisted values for all view modes
func TestViewModeString(t *testing.T) {
	tests := []struct {
		mode     models.ViewMode
		expected string
	}{
		{models.ViewAll, "all"},
		{models.ViewIncomplete, "incomplete"},
		{models.ViewComplete, "complete"},
		{models.ViewStarred, "starred"},
	}

	for _, test := range tests {
		result := test.mode.String()

		if result != test.expected {
			t.Errorf("String() for %v: expected %s, got %s", test.mode, test.expected, result)
		}
	}
}

// TestViewModeFromString checks parsing of persisted values
func TestViewModeFromString(t *testing.T) {
	tests := []struct {
		value    string
		expected models.ViewMode
	}{
		{"all", models.ViewAll},
		{"incomplete", models.ViewIncomplete},
		{"complete", models.ViewComplete},
		{"starred", models.ViewStarred},
		{"STARRED", models.ViewStarred},
		{"unknown", models.ViewIncomplete},
	}

	for _, test := range tests {
		result := models.ViewModeFromString(test.value)

		if result != test.expected {
			t.Errorf("ViewModeFromString(%q): expected %v, got %v", test.value, test.expected, result)
		}
	}
}

// TestViewModeFiltering checks todo filtering rules
func TestViewModeFiltering(t *testing.T) {
	now := time.Now()

	// Набор специально смешанный- завершенные незавершенные и важные задачи
	todos := []*models.TodoItem{
		{Name: "Completed task", Done: true, TodoTime: now.Add(-time.Hour)},
		{Name: "Incomplete task", Done: false, TodoTime: now.Add(time.Hour)},
		{Name: "Starred task", Done: false, Starred: true, TodoTime: now.Add(2 * time.Hour)},
		{Name: "Completed starred task", Done: true, Starred: true, TodoTime: now.Add(3 * time.Hour)},
	}

	tests := []struct {
		name     string
		mode     models.ViewMode
		expected []string
	}{
		{"all", models.ViewAll, []string{"Completed task", "Incomplete task", "Starred task", "Completed starred task"}},
		{"incomplete", models.ViewIncomplete, []string{"Incomplete task", "Starred task"}},
		{"complete", models.ViewComplete, []string{"Completed task", "Completed starred task"}},
		{"starred", models.ViewStarred, []string{"Starred task", "Completed starred task"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filtered := test.mode.FilterItems(todos, now)

			assertTodoNames(t, filtered, test.expected)
		})
	}
}

func assertTodoNames(t *testing.T, items []*models.TodoItem, expected []string) {
	t.Helper()

	if len(items) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(items))
	}

	for i, item := range items {
		if item.Name != expected[i] {
			t.Errorf("Item %d: expected %q, got %q", i, expected[i], item.Name)
		}
	}
}
