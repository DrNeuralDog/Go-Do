package models

import (
	"testing"
)

func TestPriorityLevelColors(t *testing.T) {
	tests := []struct {
		level    PriorityLevel
		expected string
	}{
		{PriorityLow, "#4CAF50"},
		{PriorityMedium, "#2196F3"},
		{PriorityHigh, "#FF9800"},
		{PriorityUrgent, "#F44336"},
	}

	for _, test := range tests {
		result := test.level.GetColor()
		if result != test.expected {
			t.Errorf("GetColor() for %v: expected %s, got %s", test.level, test.expected, result)
		}
	}
}

func TestPriorityLevelLabels(t *testing.T) {
	tests := []struct {
		level    PriorityLevel
		expected string
	}{
		{PriorityLow, "Not Important - Not Urgent"},
		{PriorityMedium, "Not Important - Urgent"},
		{PriorityHigh, "Important - Not Urgent"},
		{PriorityUrgent, "Important - Urgent"},
	}

	for _, test := range tests {
		result := test.level.GetLabel()
		if result != test.expected {
			t.Errorf("GetLabel() for %v: expected %s, got %s", test.level, test.expected, result)
		}
	}
}

func TestPriorityLevelShortLabels(t *testing.T) {
	tests := []struct {
		level    PriorityLevel
		expected string
	}{
		{PriorityLow, "Low"},
		{PriorityMedium, "Medium"},
		{PriorityHigh, "High"},
		{PriorityUrgent, "Urgent"},
	}

	for _, test := range tests {
		result := test.level.GetShortLabel()
		if result != test.expected {
			t.Errorf("GetShortLabel() for %v: expected %s, got %s", test.level, test.expected, result)
		}
	}
}

func TestPriorityLevelBoundaries(t *testing.T) {
	// Test that priority levels are properly defined
	if PriorityLow != 0 {
		t.Errorf("PriorityLow should be 0, got %d", PriorityLow)
	}

	if PriorityMedium != 1 {
		t.Errorf("PriorityMedium should be 1, got %d", PriorityMedium)
	}

	if PriorityHigh != 2 {
		t.Errorf("PriorityHigh should be 2, got %d", PriorityHigh)
	}

	if PriorityUrgent != 3 {
		t.Errorf("PriorityUrgent should be 3, got %d", PriorityUrgent)
	}
}
