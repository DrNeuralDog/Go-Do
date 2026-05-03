package ui

import (
	"image/color"

	"godo/src/models"
	"godo/src/ui/helpers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// statusIndicator renders done or starred state
type statusIndicator struct {
	widget.BaseWidget

	todo     *models.TodoItem
	onToggle func(toggleStar bool)
	overlay  *canvas.Rectangle
	cont     *fyne.Container
}

// newStatusIndicator creates done/star marker
func newStatusIndicator(todo *models.TodoItem, onToggle func(bool)) *statusIndicator {
	s := &statusIndicator{
		todo:     todo,
		onToggle: onToggle,
	}

	s.ExtendBaseWidget(s)

	return s
}

func (s *statusIndicator) CreateRenderer() fyne.WidgetRenderer {
	txt, col := s.visualState()
	text := canvas.NewText(txt, col)
	text.TextSize = 20

	s.overlay = canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	s.cont = container.NewMax(container.NewCenter(text), s.overlay)

	return widget.NewSimpleRenderer(s.cont)
}

func (s *statusIndicator) Tapped(*fyne.PointEvent) {
	if !s.todo.Done && s.onToggle != nil {
		s.onToggle(true)
	}
}

func (s *statusIndicator) MinSize() fyne.Size {
	return fyne.NewSize(20, 20)
}

// visualState returns icon and color for todo status
func (s *statusIndicator) visualState() (string, color.Color) {
	if s.todo.Done {
		return "✓", doneStatusColor()
	}

	if s.todo.Starred {
		return "★", starredStatusColor()
	}

	return "★", color.NRGBA{R: 0x9E, G: 0x9E, B: 0x9E, A: 0xFF}
}

// doneStatusColor returns theme-aware done color
func doneStatusColor() color.Color {
	if helpers.IsLightTheme() {
		return color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
	}

	return color.NRGBA{R: 0x98, G: 0x97, B: 0x1a, A: 0xFF}
}

// starredStatusColor returns theme-aware star color
func starredStatusColor() color.Color {
	if helpers.IsLightTheme() {
		return color.NRGBA{R: 0x3C, G: 0x82, B: 0xFF, A: 0xFF}
	}

	return color.NRGBA{R: 0xFE, G: 0x80, B: 0x19, A: 0xFF}
}
