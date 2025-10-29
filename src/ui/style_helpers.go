package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateSpacer returns an invisible spacer with fixed size.
func CreateSpacer(width float32, height float32) fyne.CanvasObject {
	r := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	return container.NewGridWrap(fyne.NewSize(width, height), r)
}

// CreateChipStyle wraps the given object in a rounded, lightly tinted background
// to create a compact "chip" appearance.
func CreateChipStyle(obj fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(toNRGBA(theme.Color(theme.ColorNameHover)))
	bg.CornerRadius = 10
	// Increase opacity slightly for visibility
	c := toNRGBA(bg.FillColor)
	c.A = 200
	bg.FillColor = c

	// Optional subtle border using separator color
	sep := toNRGBA(theme.Color(theme.ColorNameSeparator))
	sep.A = 140
	bg.StrokeColor = sep
	bg.StrokeWidth = 1

	return container.NewMax(bg, container.NewPadded(obj))
}

// CreateCardStyle wraps content in a rounded card with padding and subtle border.
func CreateCardStyle(obj fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(toNRGBA(theme.Color(theme.ColorNameInputBackground)))
	bg.CornerRadius = 12

	sep := toNRGBA(theme.Color(theme.ColorNameSeparator))
	sep.A = 160
	bg.StrokeColor = sep
	bg.StrokeWidth = 1

	return container.NewMax(bg, container.NewPadded(obj))
}

// RoundedIconButton creates a circular accent button with a centered icon.
func RoundedIconButton(icon fyne.Resource, tapped func()) fyne.CanvasObject {
	// Background circle
	circle := canvas.NewCircle(toNRGBA(theme.Color(theme.ColorNamePrimary)))
	circle.StrokeColor = toNRGBA(theme.Color(theme.ColorNameFocus))
	circle.StrokeWidth = 0

	// Foreground button (for accessibility and focus handling)
	btn := widget.NewButtonWithIcon("", icon, tapped)
	btn.Importance = widget.HighImportance

	// Center the button over the circle
	layered := container.NewMax(circle, container.NewCenter(btn))
	return layered
}

// toNRGBA converts any color.Color to color.NRGBA for manipulation.
func toNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// (removed old CreateGradientHeader helper; using full-screen gradient instead)

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

	start := toNRGBA(r.gradient.StartColor)
	end := toNRGBA(r.gradient.EndColor)

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
	return fyne.NewSize(100, 100)
}

func (r *gradientRenderer) Refresh() {
	// Trigger re-layout
	r.Layout(r.gradient.Size())
}

func (r *gradientRenderer) BackgroundColor() fyne.ThemeColorName {
	return theme.ColorNameBackground
}

func (r *gradientRenderer) Objects() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, len(r.rects))
	for i, rect := range r.rects {
		objs[i] = rect
	}
	return objs
}

func (r *gradientRenderer) Destroy() {}

// SimpleRectButton is a minimal custom button with rounded rectangle background.
type SimpleRectButton struct {
	widget.BaseWidget
	Text      string
	Bg        color.Color
	Fg        color.Color
	SizeFixed fyne.Size
	Radius    float32
	OnTapped  func()
}

func NewSimpleRectButton(text string, bg, fg color.Color, size fyne.Size, radius float32, onTap func()) *SimpleRectButton {
	b := &SimpleRectButton{
		Text:      text,
		Bg:        bg,
		Fg:        fg,
		SizeFixed: size,
		Radius:    radius,
		OnTapped:  onTap,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *SimpleRectButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(toNRGBA(b.Bg))
	bg.CornerRadius = b.Radius
	txt := canvas.NewText(b.Text, toNRGBA(b.Fg))
	txt.Alignment = fyne.TextAlignCenter
	txt.TextStyle = fyne.TextStyle{Bold: false}
	cont := container.NewMax(bg, container.NewCenter(txt))
	return widget.NewSimpleRenderer(cont)
}

func (b *SimpleRectButton) MinSize() fyne.Size {
	if b.SizeFixed.Width > 0 && b.SizeFixed.Height > 0 {
		return b.SizeFixed
	}
	return fyne.NewSize(80, 36)
}

func (b *SimpleRectButton) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *SimpleRectButton) SetText(text string) {
	b.Text = text
	b.Refresh()
}

// CreateTasksContainer wraps timeline in a rounded container with theme-specific bg.
func CreateTasksContainer(content fyne.CanvasObject) fyne.CanvasObject {
	var bg color.Color
	if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
		bg = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // #ffffff
	} else {
		bg = color.NRGBA{R: 0x28, G: 0x28, B: 0x28, A: 0xFF} // #282828
	}
	rect := canvas.NewRectangle(toNRGBA(bg))
	rect.CornerRadius = 12
	// Optional subtle border akin to mockup
	sep := toNRGBA(theme.Color(theme.ColorNameSeparator))
	sep.A = 80
	rect.StrokeColor = sep
	rect.StrokeWidth = 1
	return container.NewMax(rect, container.NewPadded(content))
}
