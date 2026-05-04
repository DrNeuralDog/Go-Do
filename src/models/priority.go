package models

import "image/color"

// PriorityLevel stores todo priority
type PriorityLevel int

const (
	PriorityLow    PriorityLevel = 0
	PriorityMedium PriorityLevel = 1
	PriorityHigh   PriorityLevel = 2
	PriorityUrgent PriorityLevel = 3
)

// GetColor returns priority color
func (p PriorityLevel) GetColor() color.RGBA {
	switch p {
	case PriorityLow:
		return color.RGBA{R: 184, G: 187, B: 38, A: 255}
	case PriorityMedium:
		return color.RGBA{R: 131, G: 165, B: 152, A: 255}
	case PriorityHigh:
		return color.RGBA{R: 254, G: 128, B: 25, A: 255}
	case PriorityUrgent:
		return color.RGBA{R: 251, G: 73, B: 52, A: 255}
	default:
		return color.RGBA{R: 184, G: 187, B: 38, A: 255}
	}
}

// GetLabel returns full priority label
func (p PriorityLevel) GetLabel() string {
	switch p {
	case PriorityLow:
		return "Not Important - Not Urgent"
	case PriorityMedium:
		return "Not Important - Urgent"
	case PriorityHigh:
		return "Important - Not Urgent"
	case PriorityUrgent:
		return "Important - Urgent"
	default:
		return "Unknown"
	}
}

// GetShortLabel returns short priority label
func (p PriorityLevel) GetShortLabel() string {
	switch p {
	case PriorityLow:
		return "Low"
	case PriorityMedium:
		return "Medium"
	case PriorityHigh:
		return "High"
	case PriorityUrgent:
		return "Urgent"
	default:
		return "Unknown"
	}
}
