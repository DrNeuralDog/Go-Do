package models

import "image/color"

// PriorityLevel represents the priority levels matching the original C++ implementation
type PriorityLevel int

const (
	PriorityLow    PriorityLevel = 0 // Not Important - Not Urgent (Green)
	PriorityMedium PriorityLevel = 1 // Not Important - Urgent (Blue)
	PriorityHigh   PriorityLevel = 2 // Important - Not Urgent (Orange)
	PriorityUrgent PriorityLevel = 3 // Important - Urgent (Red)
)

// GetColor returns the color for each priority level compatible with Fyne
func (p PriorityLevel) GetColor() color.RGBA {
	switch p {
	case PriorityLow:
		return color.RGBA{R: 184, G: 187, B: 38, A: 255} // Gruvbox Green (#b8bb26)
	case PriorityMedium:
		return color.RGBA{R: 131, G: 165, B: 152, A: 255} // Gruvbox Blue (#83a598)
	case PriorityHigh:
		return color.RGBA{R: 254, G: 128, B: 25, A: 255} // Gruvbox Orange (#fe8019)
	case PriorityUrgent:
		return color.RGBA{R: 251, G: 73, B: 52, A: 255} // Gruvbox Red (#fb4934)
	default:
		return color.RGBA{R: 184, G: 187, B: 38, A: 255}
	}
}

// GetLabel returns the English label for each priority level
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

// GetShortLabel returns a short version of the priority label
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
