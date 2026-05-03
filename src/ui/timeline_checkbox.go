package ui

import (
	"image/color"

	"godo/src/ui/helpers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// squareCheckbox renders compact todo checkbox
type squareCheckbox struct {
	widget.BaseWidget

	checked   bool
	onChanged func(bool)
	rect      *canvas.Rectangle
	tick      *canvas.Text
	overlay   *canvas.Rectangle
	cont      *fyne.Container
}

// newSquareCheckbox creates compact checkbox widget
func newSquareCheckbox(initial bool, onChanged func(bool)) *squareCheckbox {
	c := &squareCheckbox{
		checked:   initial,
		onChanged: onChanged,
	}

	c.ExtendBaseWidget(c)

	return c
}

// CreateRenderer builds themed checkbox visuals
func (c *squareCheckbox) CreateRenderer() fyne.WidgetRenderer {
	isLightTheme := helpers.IsLightTheme()

	var border color.Color
	var fillChecked color.Color

	if isLightTheme {
		border = color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF}
		fillChecked = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
	} else {
		border = color.NRGBA{R: 0x66, G: 0x5c, B: 0x54, A: 0xFF}
		fillChecked = color.NRGBA{R: 0x98, G: 0x97, B: 0x1a, A: 0xFF}
	}

	c.rect = canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	c.rect.StrokeColor = border
	c.rect.StrokeWidth = 2
	c.rect.CornerRadius = 3
	c.rect.SetMinSize(fyne.NewSize(20, 20))

	c.tick = canvas.NewText("", color.White)

	if c.checked {
		c.rect.FillColor = fillChecked
		c.tick.Text = "✓"
	}

	c.overlay = canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	c.cont = container.NewGridWrap(
		fyne.NewSize(20, 20),
		container.NewMax(c.rect, container.NewCenter(c.tick), c.overlay),
	)

	return widget.NewSimpleRenderer(c.cont)
}

func (c *squareCheckbox) MinSize() fyne.Size {
	return fyne.NewSize(20, 20)
}

// Tapped toggles state and calls change callback
func (c *squareCheckbox) Tapped(*fyne.PointEvent) {
	c.checked = !c.checked

	if c.onChanged != nil {
		c.onChanged(c.checked)
	}

	runOnMainThread(func() {
		c.Refresh()
	})
}

// Refresh syncs checkbox visuals with current state
func (c *squareCheckbox) Refresh() {
	isLightTheme := helpers.IsLightTheme()

	var fillChecked color.Color

	if isLightTheme {
		fillChecked = color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}
	} else {
		fillChecked = color.NRGBA{R: 0x98, G: 0x97, B: 0x1a, A: 0xFF}
	}

	runOnMainThread(func() {
		if c.rect != nil && c.tick != nil {
			if c.checked {
				c.rect.FillColor = fillChecked
				c.tick.Text = "✓"
			} else {
				c.rect.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
				c.tick.Text = ""
			}

			c.rect.Refresh()
			c.tick.Refresh()
		}

		c.BaseWidget.Refresh()
	})
}
