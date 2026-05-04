package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"godo/src/ui/helpers"
)

const gradientBandCount = 10

// GradientRect draws a vertical linear gradient
type GradientRect struct {
	widget.BaseWidget
	StartColor color.Color
	EndColor   color.Color
	Radius     float32
}

// NewGradientRect builds vertical gradient rectangle
func NewGradientRect(start, end color.Color, radius float32) *GradientRect {
	gradient := &GradientRect{
		StartColor: start,
		EndColor:   end,
		Radius:     radius,
	}

	gradient.ExtendBaseWidget(gradient)

	return gradient
}

func (gradient *GradientRect) CreateRenderer() fyne.WidgetRenderer {
	return newGradientRenderer(gradient)
}

// gradientRenderer keeps gradient bands stable for Fyne renderer
type gradientRenderer struct {
	gradient *GradientRect
	rects    []*canvas.Rectangle
	objects  []fyne.CanvasObject
}

// newGradientRenderer creates gradient bands once
func newGradientRenderer(gradient *GradientRect) *gradientRenderer {
	renderer := &gradientRenderer{
		gradient: gradient,
		rects:    make([]*canvas.Rectangle, gradientBandCount),
		objects:  make([]fyne.CanvasObject, gradientBandCount),
	}

	for i := 0; i < gradientBandCount; i++ {
		rect := canvas.NewRectangle(renderer.bandColor(i))
		rect.CornerRadius = gradient.Radius

		renderer.rects[i] = rect
		renderer.objects[i] = rect
	}

	return renderer
}

// Layout positions gradient bands
func (renderer *gradientRenderer) Layout(size fyne.Size) {
	bandHeight := size.Height / gradientBandCount

	for i, rect := range renderer.rects {
		rect.CornerRadius = renderer.gradient.Radius
		rect.Resize(fyne.NewSize(size.Width, bandHeight+1))
		rect.Move(fyne.NewPos(0, float32(i)*bandHeight))
	}
}

// MinSize allows gradient to fill any parent size
func (renderer *gradientRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

// Refresh updates gradient colors
func (renderer *gradientRenderer) Refresh() {
	renderer.Layout(renderer.gradient.Size())

	for i, rect := range renderer.rects {
		rect.FillColor = renderer.bandColor(i)
		rect.CornerRadius = renderer.gradient.Radius

		rect.Refresh()
	}
}

func (renderer *gradientRenderer) BackgroundColor() fyne.ThemeColorName {
	return ""
}

// Objects returns gradient band objects
func (renderer *gradientRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

// Destroy releases renderer resources
func (renderer *gradientRenderer) Destroy() {}

// bandColor returns interpolated color for band index
func (renderer *gradientRenderer) bandColor(index int) color.NRGBA {
	start := helpers.ToNRGBA(renderer.gradient.StartColor)
	end := helpers.ToNRGBA(renderer.gradient.EndColor)
	t := float32(index) / float32(gradientBandCount-1)

	return color.NRGBA{
		R: blendColorByte(start.R, end.R, t),
		G: blendColorByte(start.G, end.G, t),
		B: blendColorByte(start.B, end.B, t),
		A: blendColorByte(start.A, end.A, t),
	}
}

// blendColorByte mixes one color channel
func blendColorByte(start, end uint8, t float32) uint8 {
	return uint8(float32(start)*(1-t) + float32(end)*t)
}
