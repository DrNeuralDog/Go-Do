package models

import (
	"image/color"
	"time"
)

// TodoItem represents a single todo item with all its properties
// This struct matches the original C++ TodoItem class structure
type TodoItem struct {
	Name     string    `json:"name"`                                   // Todo item name
	Content  string    `json:"content"`                                // Detailed content/description
	Place    string    `json:"place"`                                  // Location information
	Label    string    `json:"label"`                                  // Custom label/tag
	Kind     int       `json:"kind"`                                   // Type: 0=Event, 1=Task
	Level    int       `json:"level"`                                  // Priority level: 0=Low, 1=Medium, 2=High, 3=Urgent
	TodoTime time.Time `json:"todoTime"`                               // Due date and time
	Done     bool      `json:"done"`                                   // Completion status
	WarnTime int       `json:"warnTime"`                               // Reminder time in minutes before due time
	Starred  bool      `json:"starred"`                                // Mark as important
	Order    int       `json:"order,omitempty" yaml:"order,omitempty"` // Implicit UI order within a day (0 = unset)
}

// NewTodoItem creates a new TodoItem with default values
func NewTodoItem() *TodoItem {
	return &TodoItem{
		Kind:    0, // Default to Event
		Level:   0, // Default to lowest priority
		Done:    false,
		Starred: false,
	}
}

// Setters
func (t *TodoItem) SetName(name string) {
	t.Name = name
}

func (t *TodoItem) SetContent(content string) {
	t.Content = content
}

func (t *TodoItem) SetPlace(place string) {
	t.Place = place
}

func (t *TodoItem) SetLabel(label string) {
	t.Label = label
}

func (t *TodoItem) SetKind(kind int) {
	t.Kind = kind
}

func (t *TodoItem) SetLevel(level int) {
	if level >= 0 && level <= 3 {
		t.Level = level
	}
}

func (t *TodoItem) SetTime(todoTime time.Time) {
	t.TodoTime = todoTime
}

func (t *TodoItem) SetWarnTime(warnTime int) {
	t.WarnTime = warnTime
}

func (t *TodoItem) MarkAsDone(done bool) {
	t.Done = done
}

// SetOrder sets the implicit order value for UI sorting within the same day
func (t *TodoItem) SetOrder(order int) {
	t.Order = order
}

// GetOrder returns the implicit UI order value
func (t *TodoItem) GetOrder() int {
	return t.Order
}

// Getters
func (t *TodoItem) GetName() string {
	return t.Name
}

func (t *TodoItem) GetContent() string {
	return t.Content
}

func (t *TodoItem) GetPlace() string {
	return t.Place
}

func (t *TodoItem) GetLabel() string {
	return t.Label
}

func (t *TodoItem) GetKind() int {
	return t.Kind
}

func (t *TodoItem) GetLevel() int {
	return t.Level
}

func (t *TodoItem) GetTime() time.Time {
	return t.TodoTime
}

func (t *TodoItem) GetWarnTime() int {
	return t.WarnTime
}

func (t *TodoItem) IsDone() bool {
	return t.Done
}

// HaveDone is deprecated, use IsDone() instead
func (t *TodoItem) HaveDone() bool {
	return t.IsDone()
}

// IsBefore returns true if this todo item comes before the other item chronologically
func (t *TodoItem) IsBefore(other *TodoItem) bool {
	return t.TodoTime.Before(other.TodoTime)
}

// GetKindString returns string representation of the kind
func (t *TodoItem) GetKindString() string {
	if t.Kind == 0 {
		return "Event"
	}
	return "Task"
}

// GetLevelString returns string representation of the priority level
func (t *TodoItem) GetLevelString() string {
	switch t.Level {
	case 0:
		return "Not Important - Not Urgent"
	case 1:
		return "Not Important - Urgent"
	case 2:
		return "Important - Not Urgent"
	case 3:
		return "Important - Urgent"
	default:
		return "Unknown"
	}
}

// GetLevelColor returns the color for the priority level compatible with Fyne
func (t *TodoItem) GetLevelColor() color.RGBA {
	switch t.Level {
	case 0:
		return color.RGBA{R: 184, G: 187, B: 38, A: 255} // Gruvbox Green (#b8bb26)
	case 1:
		return color.RGBA{R: 131, G: 165, B: 152, A: 255} // Gruvbox Blue (#83a598)
	case 2:
		return color.RGBA{R: 254, G: 128, B: 25, A: 255} // Gruvbox Orange (#fe8019)
	case 3:
		return color.RGBA{R: 251, G: 73, B: 52, A: 255} // Gruvbox Red (#fb4934)
	default:
		return color.RGBA{R: 184, G: 187, B: 38, A: 255} // Default Gruvbox Green
	}
}

// ShouldRemind checks if this item should trigger a reminder
func (t *TodoItem) ShouldRemind(currentTime time.Time) bool {
	if t.WarnTime == 0 || t.Done {
		return false
	}

	remindTime := t.TodoTime.Add(-time.Duration(t.WarnTime) * time.Minute)
	return currentTime.After(remindTime) && currentTime.Before(t.TodoTime)
}
