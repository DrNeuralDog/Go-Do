package models_test

import (
	"image/color"
	"testing"

	"godo/src/models"
)

// TestPriorityLevelColors checks priority color values
func TestPriorityLevelColors(t *testing.T) {
	tests := []struct {
		name     string
		level    models.PriorityLevel
		expected color.RGBA
	}{
		{"low", models.PriorityLow, color.RGBA{R: 184, G: 187, B: 38, A: 255}},
		{"medium", models.PriorityMedium, color.RGBA{R: 131, G: 165, B: 152, A: 255}},
		{"high", models.PriorityHigh, color.RGBA{R: 254, G: 128, B: 25, A: 255}},
		{"urgent", models.PriorityUrgent, color.RGBA{R: 251, G: 73, B: 52, A: 255}},
		{"unknown", models.PriorityLevel(99), color.RGBA{R: 184, G: 187, B: 38, A: 255}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.level.GetColor()

			if result != test.expected {
				t.Errorf("GetColor() for %v: expected %#v, got %#v", test.level, test.expected, result)
			}
		})
	}
}

// TestPriorityLevelLabels checks full priority labels
func TestPriorityLevelLabels(t *testing.T) {
	tests := []struct {
		name     string
		level    models.PriorityLevel
		expected string
	}{
		{"low", models.PriorityLow, "Not Important - Not Urgent"},
		{"medium", models.PriorityMedium, "Not Important - Urgent"},
		{"high", models.PriorityHigh, "Important - Not Urgent"},
		{"urgent", models.PriorityUrgent, "Important - Urgent"},
		{"unknown", models.PriorityLevel(99), "Unknown"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.level.GetLabel()

			if result != test.expected {
				t.Errorf("GetLabel() for %v: expected %s, got %s", test.level, test.expected, result)
			}
		})
	}
}

// TestPriorityLevelShortLabels checks compact priority labels
func TestPriorityLevelShortLabels(t *testing.T) {
	tests := []struct {
		name     string
		level    models.PriorityLevel
		expected string
	}{
		{"low", models.PriorityLow, "Low"},
		{"medium", models.PriorityMedium, "Medium"},
		{"high", models.PriorityHigh, "High"},
		{"urgent", models.PriorityUrgent, "Urgent"},
		{"unknown", models.PriorityLevel(99), "Unknown"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.level.GetShortLabel()

			if result != test.expected {
				t.Errorf("GetShortLabel() for %v: expected %s, got %s", test.level, test.expected, result)
			}
		})
	}
}

// TestPriorityLevelBoundaries checks stored enum values
func TestPriorityLevelBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		level    models.PriorityLevel
		expected models.PriorityLevel
	}{
		{"low", models.PriorityLow, 0},
		{"medium", models.PriorityMedium, 1},
		{"high", models.PriorityHigh, 2},
		{"urgent", models.PriorityUrgent, 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.level != test.expected {
				t.Errorf("%s should be %d, got %d", test.name, test.expected, test.level)
			}
		})
	}
}
