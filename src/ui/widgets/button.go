package widgets

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"godo/src/ui/helpers"
	"godo/src/ui/threading"
)

// RoundIconButton is a circular icon button that is fully clickable without rectangular hover overlay.
type RoundIconButton struct {
	widget.BaseWidget
	Icon     fyne.Resource
	Bg       color.Color
	Fg       color.Color
	OnTapped func()
	hovered  bool
}

func NewRoundIconButton(icon fyne.Resource, onTap func()) *RoundIconButton {
	b := &RoundIconButton{
		Icon:     icon,
		Bg:       helpers.ToNRGBA(theme.Color(theme.ColorNamePrimary)),
		Fg:       helpers.ToNRGBA(theme.Color(theme.ColorNameForeground)),
		OnTapped: onTap,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *RoundIconButton) CreateRenderer() fyne.WidgetRenderer {
	circle := canvas.NewCircle(helpers.ToNRGBA(b.Bg))
	icon := widget.NewIcon(b.Icon)
	cont := container.NewMax(circle, container.NewCenter(icon))
	return &roundIconButtonRenderer{
		btn:    b,
		circle: circle,
		icon:   icon,
		cont:   cont,
	}
}

type roundIconButtonRenderer struct {
	btn    *RoundIconButton
	circle *canvas.Circle
	icon   *widget.Icon
	cont   *fyne.Container
}

func (r *roundIconButtonRenderer) Layout(size fyne.Size)                { r.cont.Resize(size) }
func (r *roundIconButtonRenderer) MinSize() fyne.Size                   { return r.btn.MinSize() }
func (r *roundIconButtonRenderer) BackgroundColor() fyne.ThemeColorName { return "" }
func (r *roundIconButtonRenderer) Objects() []fyne.CanvasObject         { return []fyne.CanvasObject{r.cont} }
func (r *roundIconButtonRenderer) Destroy()                             {}
func (r *roundIconButtonRenderer) Refresh() {
	bg := helpers.ToNRGBA(r.btn.Bg)
	if r.btn.hovered {
		bg = helpers.Lighten(bg, 0.12)
	}
	r.circle.FillColor = bg
	threading.RunOnMainThread(func() {
		r.circle.Refresh()
		r.icon.Refresh()
	})
}

func (b *RoundIconButton) MinSize() fyne.Size {
	// Default minimum; actual size is controlled by parent container (GridWrap)
	return fyne.NewSize(36, 36)
}

func (b *RoundIconButton) Tapped(*fyne.PointEvent) {
	// quick press flash by darkening background briefly
	orig := helpers.ToNRGBA(b.Bg)
	b.Bg = helpers.Darken(orig, 0.85)
	threading.RunOnMainThread(func() {
		b.Refresh()
	})
	go func() {
		time.Sleep(120 * time.Millisecond)
		threading.RunOnMainThread(func() {
			b.Bg = orig
			b.Refresh()
		})
	}()
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// Hover handling (desktop only)
func (b *RoundIconButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	threading.RunOnMainThread(func() { b.Refresh() })
}
func (b *RoundIconButton) MouseMoved(*desktop.MouseEvent) {}
func (b *RoundIconButton) MouseOut() {
	b.hovered = false
	threading.RunOnMainThread(func() { b.Refresh() })
}

// SimpleRectButton is a minimal custom button with rounded rectangle background.
type SimpleRectButton struct {
	widget.BaseWidget
	Text      string
	Bg        color.Color
	Fg        color.Color
	SizeFixed fyne.Size
	Radius    float32
	OnTapped  func()
	Disabled  bool
	hovered   bool
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
	bg := canvas.NewRectangle(helpers.ToNRGBA(b.Bg))
	bg.CornerRadius = b.Radius
	txt := canvas.NewText(b.Text, helpers.ToNRGBA(b.Fg))
	txt.Alignment = fyne.TextAlignCenter
	// Make button text bold for better emphasis
	txt.TextStyle = fyne.TextStyle{Bold: true}
	cont := container.NewMax(bg, container.NewCenter(txt))
	return &simpleRectButtonRenderer{
		button: b,
		bg:     bg,
		text:   txt,
		cont:   cont,
	}
}

// simpleRectButtonRenderer is a custom renderer for SimpleRectButton
type simpleRectButtonRenderer struct {
	button *SimpleRectButton
	bg     *canvas.Rectangle
	text   *canvas.Text
	cont   *fyne.Container
}

func (r *simpleRectButtonRenderer) Layout(size fyne.Size) {
	r.cont.Resize(size)
}

func (r *simpleRectButtonRenderer) MinSize() fyne.Size {
	return r.button.MinSize()
}

func (r *simpleRectButtonRenderer) Refresh() {
	// Sync colors & text from button state
	btn := r.button
	bgCol := helpers.ToNRGBA(btn.Bg)
	fgCol := helpers.ToNRGBA(btn.Fg)
	if btn.Disabled {
		// dim colors
		bgCol = helpers.Darken(bgCol, 0.6)
		fgCol.A = 160
	} else if btn.hovered {
		bgCol = helpers.Lighten(bgCol, 0.10)
	}
	r.bg.FillColor = bgCol
	r.text.Color = fgCol
	r.text.Text = btn.Text
	threading.RunOnMainThread(func() {
		r.text.Refresh()
		r.bg.Refresh()
	})
}

func (r *simpleRectButtonRenderer) BackgroundColor() fyne.ThemeColorName {
	return ""
}

func (r *simpleRectButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.cont}
}

func (r *simpleRectButtonRenderer) Destroy() {}

func (b *SimpleRectButton) MinSize() fyne.Size {
	if b.SizeFixed.Width > 0 && b.SizeFixed.Height > 0 {
		return b.SizeFixed
	}
	return fyne.NewSize(80, 36)
}

func (b *SimpleRectButton) Tapped(*fyne.PointEvent) {
	if b.Disabled {
		return
	}
	// flash background on press
	orig := helpers.ToNRGBA(b.Bg)
	b.Bg = helpers.Darken(orig, 0.85)
	threading.RunOnMainThread(func() {
		b.Refresh()
	})
	go func() {
		time.Sleep(120 * time.Millisecond)
		threading.RunOnMainThread(func() {
			b.Bg = orig
			b.Refresh()
		})
	}()
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *SimpleRectButton) SetText(text string) {
	b.Text = text
	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

func (b *SimpleRectButton) Enable() {
	if !b.Disabled {
		return
	}
	b.Disabled = false
	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

func (b *SimpleRectButton) Disable() {
	if b.Disabled {
		return
	}
	b.Disabled = true
	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

// Hover handling (desktop only)
func (b *SimpleRectButton) MouseIn(*desktop.MouseEvent) {
	if !b.Disabled {
		b.hovered = true
		threading.RunOnMainThread(func() {
			b.Refresh()
		})
	}
}
func (b *SimpleRectButton) MouseMoved(*desktop.MouseEvent) {}
func (b *SimpleRectButton) MouseOut() {
	if !b.Disabled {
		b.hovered = false
		threading.RunOnMainThread(func() {
			b.Refresh()
		})
	}
}

// TinyIconButton is a minimal icon button (16x16) without padding or background for compact layouts.
type TinyIconButton struct {
	widget.BaseWidget
	Icon     fyne.Resource
	OnTapped func()
	pressed  bool // track press state for animation
	hovered  bool // track hover state
}

func NewTinyIconButton(icon fyne.Resource, onTap func()) *TinyIconButton {
	b := &TinyIconButton{
		Icon:     icon,
		OnTapped: onTap,
		pressed:  false,
		hovered:  false,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *TinyIconButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	bg.CornerRadius = 4
	return &tinyIconButtonRenderer{
		button: b,
		icon:   widget.NewIcon(b.Icon),
		bg:     bg,
	}
}

func (b *TinyIconButton) MinSize() fyne.Size {
	return fyne.NewSize(16, 16)
}

func (b *TinyIconButton) Tapped(*fyne.PointEvent) {
	// Trigger press animation
	b.pressed = true
	threading.RunOnMainThread(func() {
		b.Refresh()
	})

	// Execute callback
	if b.OnTapped != nil {
		b.OnTapped()
	}

	// Animate back to normal after 120ms
	go func() {
		time.Sleep(120 * time.Millisecond)
		threading.RunOnMainThread(func() {
			b.pressed = false
			b.Refresh()
		})
	}()
}

func (b *TinyIconButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

func (b *TinyIconButton) MouseMoved(*desktop.MouseEvent) {
	// Keep hovered state while mouse is in
}

func (b *TinyIconButton) MouseOut() {
	b.hovered = false
	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

// tinyIconButtonRenderer renders a TinyIconButton with visual feedback
type tinyIconButtonRenderer struct {
	button *TinyIconButton
	icon   *widget.Icon
	bg     *canvas.Rectangle
}

func (r *tinyIconButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))
	r.icon.Resize(size)
	r.icon.Move(fyne.NewPos(0, 0))
}

func (r *tinyIconButtonRenderer) MinSize() fyne.Size {
	return r.button.MinSize()
}

func (r *tinyIconButtonRenderer) Refresh() {
	// Adjust background color based on state for visual feedback
	if r.button.pressed {
		// Darken when pressed - use a dark semi-transparent background
		r.bg.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 80}
	} else if r.button.hovered {
		// Brighten when hovered - use a light semi-transparent background
		r.bg.FillColor = color.NRGBA{R: 200, G: 200, B: 200, A: 60}
	} else {
		// Normal state - transparent
		r.bg.FillColor = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}
	threading.RunOnMainThread(func() {
		r.bg.Refresh()
		r.icon.Refresh()
	})
}

func (r *tinyIconButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.icon}
}

func (r *tinyIconButtonRenderer) Destroy() {
	// Cleanup if needed
}

func (r *tinyIconButtonRenderer) BackgroundColor() fyne.ThemeColorName {
	return ""
}
