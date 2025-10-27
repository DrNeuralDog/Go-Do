package models

import (
	"testing"
	"time"
)

func TestViewModeLabels(t *testing.T) {
	tests := []struct {
		mode     ViewMode
		expected string
	}{
		{ViewAll, "All"},
		{ViewIncomplete, "Incomplete"},
		{ViewReminders, "Reminders"},
	}

	for _, test := range tests {
		result := test.mode.GetLabel()
		if result != test.expected {
			t.Errorf("GetLabel() for %v: expected %s, got %s", test.mode, test.expected, result)
		}
	}
}

func TestViewModeCycle(t *testing.T) {
	// Test cycling through view modes
	mode := ViewAll
	expected := ViewIncomplete

	result := mode.GetNextMode()
	if result != expected {
		t.Errorf("GetNextMode() from ViewAll: expected %v, got %v", expected, result)
	}

	mode = ViewIncomplete
	expected = ViewReminders

	result = mode.GetNextMode()
	if result != expected {
		t.Errorf("GetNextMode() from ViewIncomplete: expected %v, got %v", expected, result)
	}

	mode = ViewReminders
	expected = ViewAll

	result = mode.GetNextMode()
	if result != expected {
		t.Errorf("GetNextMode() from ViewReminders: expected %v, got %v", expected, result)
	}
}

func TestViewModeFiltering(t *testing.T) {
	// Create test todos
	now := time.Now()

	todos := []*TodoItem{
		{Name: "Completed task", Done: true, TodoTime: now.Add(-time.Hour)},
		{Name: "Incomplete task", Done: false, TodoTime: now.Add(time.Hour)},
		{Name: "Reminder task", Done: false, TodoTime: now.Add(30 * time.Minute), WarnTime: 60},
		{Name: "No reminder task", Done: false, TodoTime: now.Add(2 * time.Hour), WarnTime: 0},
	}

	// Test ViewAll
	allMode := ViewAll
	filtered := allMode.FilterItems(todos, now)
	if len(filtered) != 4 {
		t.Errorf("ViewAll should return all 4 items, got %d", len(filtered))
	}

	// Test ViewIncomplete
	incompleteMode := ViewIncomplete
	filtered = incompleteMode.FilterItems(todos, now)
	if len(filtered) != 3 {
		t.Errorf("ViewIncomplete should return 3 items, got %d", len(filtered))
	}

	// Check that only incomplete items are returned
	for _, item := range filtered {
		if item.Done {
			t.Error("ViewIncomplete should not return completed items")
		}
	}

	// Test ViewReminders
	reminderMode := ViewReminders
	filtered = reminderMode.FilterItems(todos, now)
	if len(filtered) != 1 {
		t.Errorf("ViewReminders should return 1 item, got %d", len(filtered))
	}

	// Check that only the reminder item is returned
	if filtered[0].Name != "Reminder task" {
		t.Error("ViewReminders should return only the reminder task")
	}
}

func TestViewModeFilteringWithCurrentTime(t *testing.T) {
	// Create test todos with different reminder times
	now := time.Now()

	todos := []*TodoItem{
		{
			Name:     "Past reminder",
			Done:     false,
			TodoTime: now.Add(2 * time.Hour),
			WarnTime: 60, // Should remind 1 hour before
		},
		{
			Name:     "Future reminder",
			Done:     false,
			TodoTime: now.Add(2 * time.Hour),
			WarnTime: 30, // Should remind 30 minutes before
		},
		{
			Name:     "No reminder",
			Done:     false,
			TodoTime: now.Add(2 * time.Hour),
			WarnTime: 0,
		},
	}

	// Test that past reminder time doesn't trigger reminder
	reminderMode := ViewReminders
	filtered := reminderMode.FilterItems(todos, now)

	// The past reminder (60 min) should not trigger since it's more than 1 hour before todo time
	// The future reminder (30 min) should trigger since it's within the reminder window
	// The no reminder should not appear

	expectedCount := 1 // Only the future reminder should trigger
	if len(filtered) != expectedCount {
		t.Errorf("ViewReminders should return %d items, got %d", expectedCount, len(filtered))
		for _, item := range filtered {
			t.Logf("Returned item: %s", item.Name)
		}
	}
}
