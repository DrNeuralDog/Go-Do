package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"godo/src/ui/helpers"
	"godo/src/ui/widgets"
)

// RoundedIconButton wraps widgets.NewRoundIconButton
func RoundedIconButton(icon fyne.Resource, tapped func()) fyne.CanvasObject {
	return widgets.NewRoundIconButton(icon, tapped)
}

// NewRoundIconButton keeps unused legacy round button wrapper just in case
func NewRoundIconButton(icon fyne.Resource, onTap func()) *widgets.RoundIconButton {
	return widgets.NewRoundIconButton(icon, onTap)
}

// NewSimpleRectButton wraps widgets.NewSimpleRectButton
func NewSimpleRectButton(text string, bg, fg color.Color, size fyne.Size, radius float32, onTap func()) *widgets.SimpleRectButton {
	return widgets.NewSimpleRectButton(text, bg, fg, size, radius, onTap)
}

// NewGradientRect wraps widgets.NewGradientRect
func NewGradientRect(start, end color.Color, radius float32) *widgets.GradientRect {
	return widgets.NewGradientRect(start, end, radius)
}

// NewCustomSelect wraps widgets.NewCustomSelect
func NewCustomSelect(options []string, onChanged func(string)) *widgets.CustomSelect {
	return widgets.NewCustomSelect(options, onChanged)
}

// CreateStyledSelect builds fixed-size select background
func CreateStyledSelect(selectWidget fyne.CanvasObject, bgColor color.Color, size fyne.Size, radius float32) fyne.CanvasObject {
	bg := canvas.NewRectangle(helpers.ToNRGBA(bgColor))
	bg.CornerRadius = radius

	selectWrapper := container.NewGridWrap(size, selectWidget)

	return container.NewMax(
		container.NewGridWrap(size, bg),
		selectWrapper,
	)
}

// NewTinyIconButton keeps unused legacy tiny button wrapper just in case
func NewTinyIconButton(icon fyne.Resource, onTap func()) *widgets.TinyIconButton {
	return widgets.NewTinyIconButton(icon, onTap)
}

// CreateTasksContainer builds themed task panel
func CreateTasksContainer(content fyne.CanvasObject) fyne.CanvasObject {
	bgColor := theme.Color(theme.ColorNameInputBackground)

	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
		// Светлая тема держит список на белой панели
		bgColor = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	} else if _, ok := fyne.CurrentApp().Settings().Theme().(*GruvboxBlackTheme); ok {
		// Темная тема убирает градиент под задачами
		bgColor = theme.Color(theme.ColorNameBackground)
	}

	rect := canvas.NewRectangle(helpers.ToNRGBA(bgColor))
	rect.CornerRadius = 12

	padded := container.NewPadded(content)

	return container.NewMax(rect, padded)
}

// NewNumberSpinner wraps widgets.NewNumberSpinner
func NewNumberSpinner(win fyne.Window, initial, min, max, step int, textColor, bgColor color.Color, onChanged func(int)) *widgets.NumberSpinner {
	return widgets.NewNumberSpinner(win, initial, min, max, step, textColor, bgColor, onChanged)
}

// FlashWindow wraps helpers.FlashWindow
func FlashWindow(win fyne.Window) {
	helpers.FlashWindow(win)
}
