package helpers

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// CreateSpacer returns an invisible spacer with fixed size.
func CreateSpacer(width float32, height float32) fyne.CanvasObject {
	r := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	return container.NewGridWrap(fyne.NewSize(width, height), r)
}

// CreateFixedSeparator returns a horizontal line with fixed color,
// so it remains visible and does not change when the app theme changes.
func CreateFixedSeparator() fyne.CanvasObject {
	rect := canvas.NewRectangle(Hex("#bdae93")) // warm light line, readable on dark background
	rect.SetMinSize(fyne.NewSize(1, 1))
	return rect
}

// CreateChipStyle wraps the given object in a rounded, lightly tinted background
// to create a compact "chip" appearance.
func CreateChipStyle(obj fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(ToNRGBA(theme.Color(theme.ColorNameHover)))
	bg.CornerRadius = 10
	// Increase opacity slightly for visibility
	c := ToNRGBA(bg.FillColor)
	c.A = 200
	bg.FillColor = c

	// Optional subtle border using separator color
	sep := ToNRGBA(theme.Color(theme.ColorNameSeparator))
	sep.A = 255 // Full opacity for visibility
	bg.StrokeColor = sep
	bg.StrokeWidth = 1

	return container.NewMax(bg, container.NewPadded(obj))
}

// CreateCardStyle wraps content in a rounded card with padding and subtle border.
func CreateCardStyle(obj fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(ToNRGBA(theme.Color(theme.ColorNameInputBackground)))
	bg.CornerRadius = 12

	sep := ToNRGBA(theme.Color(theme.ColorNameSeparator))
	sep.A = 255 // Full opacity for visibility
	bg.StrokeColor = sep
	bg.StrokeWidth = 1

	return container.NewMax(bg, container.NewPadded(obj))
}
