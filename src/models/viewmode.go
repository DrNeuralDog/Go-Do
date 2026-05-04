package models

import (
	"strings"
)

// ViewMode filters visible todos
type ViewMode int

const (
	ViewAll ViewMode = iota
	ViewIncomplete
	ViewComplete
	ViewStarred
)

// GetLabel returns view mode label
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

// String returns persisted view mode value
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

// ViewModeFromString parses persisted view mode value
func ViewModeFromString(s string) ViewMode {
	switch strings.TrimSpace(strings.ToLower(s)) {
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

// FilterItems returns todos visible in this mode
func (v ViewMode) FilterItems(items []*TodoItem) []*TodoItem {
	filtered := make([]*TodoItem, 0, len(items))

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

// GetNextMode returns next view mode in cycle
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
