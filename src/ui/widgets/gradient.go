package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"godo/src/ui/helpers"
)

// GradientRect is a custom widget that draws a vertical linear gradient.
type GradientRect struct {
	widget.BaseWidget
	StartColor color.Color
	EndColor   color.Color
	Radius     float32
}

func NewGradientRect(start, end color.Color, radius float32) *GradientRect {
	gr := &GradientRect{
		StartColor: start,
		EndColor:   end,
		Radius:     radius,
	}
	gr.ExtendBaseWidget(gr)
	return gr
}

func (gr *GradientRect) CreateRenderer() fyne.WidgetRenderer {
	return &gradientRenderer{gradient: gr}
}

type gradientRenderer struct {
	gradient *GradientRect
	rects    []*canvas.Rectangle
}

func (r *gradientRenderer) Layout(size fyne.Size) {
	// Create gradient bands - 10 bands for smooth transition
	bandCount := 10
	bandHeight := size.Height / float32(bandCount)

	// Clear old rects if resizing
	r.rects = make([]*canvas.Rectangle, 0, bandCount)

	start := helpers.ToNRGBA(r.gradient.StartColor)
	end := helpers.ToNRGBA(r.gradient.EndColor)

	for i := 0; i < bandCount; i++ {
		// Interpolate color
		t := float32(i) / float32(bandCount-1)
		bandColor := color.NRGBA{
			R: uint8(float32(start.R)*(1-t) + float32(end.R)*t),
			G: uint8(float32(start.G)*(1-t) + float32(end.G)*t),
			B: uint8(float32(start.B)*(1-t) + float32(end.B)*t),
			A: 255,
		}

		rect := canvas.NewRectangle(bandColor)
		rect.CornerRadius = r.gradient.Radius
		rect.Resize(fyne.NewSize(size.Width, bandHeight+1)) // +1 to avoid gaps
		rect.Move(fyne.NewPos(0, float32(i)*bandHeight))
		r.rects = append(r.rects, rect)
	}
}

func (r *gradientRenderer) MinSize() fyne.Size {
	// Return zero size to allow gradient to be any size (no minimum constraint)
	return fyne.NewSize(0, 0)
}

func (r *gradientRenderer) Refresh() {
	// Trigger re-layout with current size
	if size := r.gradient.Size(); size.Width > 0 && size.Height > 0 {
		r.Layout(size)
	}
}

func (r *gradientRenderer) BackgroundColor() fyne.ThemeColorName {
	// Return empty/transparent - gradient fills everything
	return ""
}

func (r *gradientRenderer) Objects() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, len(r.rects))
	for i, rect := range r.rects {
		objs[i] = rect
	}
	return objs
}

func (r *gradientRenderer) Destroy() {}
