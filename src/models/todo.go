package models

import (
	"image/color"
	"time"
)

// Todo kind values match old storage
const (
	TodoKindEvent = iota
	TodoKindTask
)

// TodoItem stores one task or event
type TodoItem struct {
	Name     string    `json:"name"`
	Content  string    `json:"content"`
	Place    string    `json:"place"`
	Label    string    `json:"label"`
	Kind     int       `json:"kind"`
	Level    int       `json:"level"`
	TodoTime time.Time `json:"todoTime"`
	Done     bool      `json:"done"`
	WarnTime int       `json:"warnTime"`
	Starred  bool      `json:"starred"`
	Order    int       `json:"order,omitempty" yaml:"order,omitempty"`
}

// NewTodoItem creates default todo
func NewTodoItem() *TodoItem {
	return &TodoItem{
		Kind:    TodoKindEvent,
		Level:   int(PriorityLow),
		Done:    false,
		Starred: false,
	}
}

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
	if level >= int(PriorityLow) && level <= int(PriorityUrgent) {
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

// SetOrder sets UI order inside one day
func (t *TodoItem) SetOrder(order int) {
	t.Order = order
}

// GetOrder returns UI order
func (t *TodoItem) GetOrder() int {
	return t.Order
}

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

// HaveDone keeps old done getter
func (t *TodoItem) HaveDone() bool {
	return t.IsDone()
}

// IsBefore reports whether todo comes before another
func (t *TodoItem) IsBefore(other *TodoItem) bool {
	return t.TodoTime.Before(other.TodoTime)
}

// GetKindString returns kind label
func (t *TodoItem) GetKindString() string {
	if t.Kind == TodoKindEvent {
		return "Event"
	}

	return "Task"
}

// GetLevelString returns priority label
func (t *TodoItem) GetLevelString() string {
	return PriorityLevel(t.Level).GetLabel()
}

// GetLevelColor returns priority color
func (t *TodoItem) GetLevelColor() color.RGBA {
	return PriorityLevel(t.Level).GetColor()
}

// ShouldRemind reports whether reminder should fire now
func (t *TodoItem) ShouldRemind(currentTime time.Time) bool {
	if t.WarnTime <= 0 || t.Done {
		return false
	}

	remindTime := t.TodoTime.Add(-time.Duration(t.WarnTime) * time.Minute)

	return !currentTime.Before(remindTime) && currentTime.Before(t.TodoTime)
}
