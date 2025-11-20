package models

import (
	"strings"
	"time"
)

// ViewMode represents different filtering modes for the todo list
type ViewMode int

const (
	ViewAll        ViewMode = 0 // Show all items
	ViewIncomplete ViewMode = 1 // Show only incomplete items
	ViewComplete   ViewMode = 2 // Show only completed items
	ViewStarred    ViewMode = 3 // Show only starred items
)

// GetLabel returns the English label for each view mode
func (v ViewMode) GetLabel() string {
	switch v {
	case ViewAll:
		return "All"
	case ViewIncomplete:
		return "Incomplete"
	case ViewComplete:
		return "Complete"
	case ViewStarred:
		return "Important"
	default:
		return "All"
	}
}

// String converts a ViewMode to its persisted string value.
func (v ViewMode) String() string {
	switch v {
	case ViewAll:
		return "all"
	case ViewIncomplete:
		return "incomplete"
	case ViewComplete:
		return "complete"
	case ViewStarred:
		return "starred"
	default:
		return "incomplete"
	}
}

// ViewModeFromString returns a ViewMode from its string representation.
func ViewModeFromString(s string) ViewMode {
	switch strings.ToLower(s) {
	case "all":
		return ViewAll
	case "incomplete":
		return ViewIncomplete
	case "complete":
		return ViewComplete
	case "starred":
		return ViewStarred
	default:
		return ViewIncomplete
	}
}

// FilterItems filters a slice of todo items based on the current view mode
func (v ViewMode) FilterItems(items []*TodoItem, currentTime time.Time) []*TodoItem {
	var filtered []*TodoItem

	for _, item := range items {
		switch v {
		case ViewAll:
			filtered = append(filtered, item)
		case ViewIncomplete:
			if !item.Done {
				filtered = append(filtered, item)
			}
		case ViewComplete:
			if item.Done {
				filtered = append(filtered, item)
			}
		case ViewStarred:
			if item.Starred {
				filtered = append(filtered, item)
			}
		}
	}

	return filtered
}

// GetNextMode returns the next view mode in cycle
func (v ViewMode) GetNextMode() ViewMode {
	switch v {
	case ViewAll:
		return ViewIncomplete
	case ViewIncomplete:
		return ViewComplete
	case ViewComplete:
		return ViewStarred
	case ViewStarred:
		return ViewAll
	default:
		return ViewAll
	}
}
