package helpers

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

const (
	fixedSeparatorHex = "#bdae93"
	chipCornerRadius  = 10
	cardCornerRadius  = 12
	chipAlpha         = 200
	strokeAlpha       = 255
	strokeWidth       = 1
)

// CreateSpacer returns invisible fixed-size spacer
func CreateSpacer(width float32, height float32) fyne.CanvasObject {
	width = nonNegativeSize(width)
	height = nonNegativeSize(height)

	r := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})

	return container.NewGridWrap(fyne.NewSize(width, height), r)
}

// CreateFixedSeparator returns warm fixed-color separator
func CreateFixedSeparator() fyne.CanvasObject {
	rect := canvas.NewRectangle(Hex(fixedSeparatorHex))
	rect.SetMinSize(fyne.NewSize(1, 1))

	return rect
}

// CreateChipStyle wraps object into compact rounded chip
func CreateChipStyle(obj fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(ToNRGBA(theme.Color(theme.ColorNameHover)))
	bg.CornerRadius = chipCornerRadius

	// Чип делаем чуть плотнее, иначе на теме он теряется
	c := ToNRGBA(bg.FillColor)
	c.A = chipAlpha
	bg.FillColor = c

	applySeparatorStroke(bg)

	return container.NewMax(bg, container.NewPadded(obj))
}

// CreateCardStyle wraps object into padded rounded card
func CreateCardStyle(obj fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(ToNRGBA(theme.Color(theme.ColorNameInputBackground)))
	bg.CornerRadius = cardCornerRadius

	applySeparatorStroke(bg)

	return container.NewMax(bg, container.NewPadded(obj))
}

// applySeparatorStroke adds visible separator border
func applySeparatorStroke(rect *canvas.Rectangle) {
	sep := ToNRGBA(theme.Color(theme.ColorNameSeparator))
	sep.A = strokeAlpha

	rect.StrokeColor = sep
	rect.StrokeWidth = strokeWidth
}

// nonNegativeSize keeps layout sizes sane
func nonNegativeSize(size float32) float32 {
	if size < 0 {
		// Отрицательный spacer тут только ломает верстку
		return 0
	}

	return size
}
