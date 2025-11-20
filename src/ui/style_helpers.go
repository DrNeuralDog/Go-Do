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

// RoundedIconButton creates a circular accent button with a centered icon.
// Wrapper for backward compatibility.
func RoundedIconButton(icon fyne.Resource, tapped func()) fyne.CanvasObject {
	return widgets.NewRoundIconButton(icon, tapped)
}

// NewRoundIconButton creates a new RoundIconButton widget.
// Wrapper for backward compatibility.
func NewRoundIconButton(icon fyne.Resource, onTap func()) *widgets.RoundIconButton {
	return widgets.NewRoundIconButton(icon, onTap)
}

// NewSimpleRectButton creates a new SimpleRectButton widget.
// Wrapper for backward compatibility.
func NewSimpleRectButton(text string, bg, fg color.Color, size fyne.Size, radius float32, onTap func()) *widgets.SimpleRectButton {
	return widgets.NewSimpleRectButton(text, bg, fg, size, radius, onTap)
}

// NewGradientRect creates a new GradientRect widget.
// Wrapper for backward compatibility.
func NewGradientRect(start, end color.Color, radius float32) *widgets.GradientRect {
	return widgets.NewGradientRect(start, end, radius)
}

// NewCustomSelect creates a new CustomSelect widget.
// Wrapper for backward compatibility.
func NewCustomSelect(options []string, onChanged func(string)) *widgets.CustomSelect {
	return widgets.NewCustomSelect(options, onChanged)
}

// CreateStyledSelect wraps a Select widget in a rounded container with custom background
// Accepts any CanvasObject (including widget.Select and CustomSelect)
func CreateStyledSelect(selectWidget fyne.CanvasObject, bgColor color.Color, size fyne.Size, radius float32) fyne.CanvasObject {
	bg := canvas.NewRectangle(helpers.ToNRGBA(bgColor))
	bg.CornerRadius = radius

	// Wrap select in container with fixed size
	selectWrapper := container.NewGridWrap(size, selectWidget)

	return container.NewMax(
		container.NewGridWrap(size, bg),
		selectWrapper,
	)
}

// NewTinyIconButton creates a new TinyIconButton widget.
// Wrapper for backward compatibility.
func NewTinyIconButton(icon fyne.Resource, onTap func()) *widgets.TinyIconButton {
	return widgets.NewTinyIconButton(icon, onTap)
}

// CreateTasksContainer wraps timeline in a rounded container with theme-specific bg.
func CreateTasksContainer(content fyne.CanvasObject) fyne.CanvasObject {
	bgColor := theme.Color(theme.ColorNameInputBackground)
	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
		// Light theme: keep tasks window (card) white
		bgColor = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	} else if _, ok := fyne.CurrentApp().Settings().Theme().(*GruvboxBlackTheme); ok {
		// Dark theme: force solid background (no gradient) using theme base background
		bgColor = theme.Color(theme.ColorNameBackground)
	}
	rect := canvas.NewRectangle(helpers.ToNRGBA(bgColor))
	rect.CornerRadius = 12
	// Use custom padding to control exact spacing - 8px all around instead of default Fyne padding
	padded := container.NewPadded(content)
	return container.NewMax(rect, padded)
}

// NewNumberSpinner creates a new NumberSpinner widget.
// Wrapper for backward compatibility.
func NewNumberSpinner(win fyne.Window, initial, min, max, step int, textColor, bgColor color.Color, onChanged func(int)) *widgets.NumberSpinner {
	return widgets.NewNumberSpinner(win, initial, min, max, step, textColor, bgColor, onChanged)
}

// FlashWindow creates a visual flash effect on a window to indicate it's already open.
// Wrapper for backward compatibility.
func FlashWindow(win fyne.Window) {
	helpers.FlashWindow(win)
}
