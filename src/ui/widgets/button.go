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

const (
	roundIconButtonSize           = 36
	simpleRectButtonDefaultWidth  = 80
	simpleRectButtonDefaultHeight = 36
	tinyIconButtonSize            = 16
	tinyIconButtonCornerRadius    = 4
	buttonPressDelay              = 120 * time.Millisecond
	buttonPressDarkenFactor       = 0.85
	roundIconButtonHoverAmount    = 0.12
	simpleRectButtonHoverAmount   = 0.10
	simpleRectButtonDisabledAlpha = 160
	simpleRectButtonDisabledDim   = 0.6
	tinyIconButtonPressAlpha      = 80
	tinyIconButtonHoverAlpha      = 60
)

// RoundIconButton is a circular icon button without rectangular hover overlay
type RoundIconButton struct {
	widget.BaseWidget
	Icon fyne.Resource
	Bg   color.Color

	OnTapped func()
	hovered  bool
}

// NewRoundIconButton builds circular icon button
func NewRoundIconButton(icon fyne.Resource, onTap func()) *RoundIconButton {
	b := &RoundIconButton{
		Icon:     icon,
		Bg:       helpers.ToNRGBA(theme.Color(theme.ColorNamePrimary)),
		OnTapped: onTap,
	}

	b.ExtendBaseWidget(b)

	return b
}

// CreateRenderer builds round icon button renderer
func (b *RoundIconButton) CreateRenderer() fyne.WidgetRenderer {
	circle := canvas.NewCircle(helpers.ToNRGBA(b.Bg))
	icon := widget.NewIcon(b.Icon)
	cont := container.NewMax(circle, container.NewCenter(icon))

	return &roundIconButtonRenderer{
		btn:     b,
		circle:  circle,
		icon:    icon,
		cont:    cont,
		objects: []fyne.CanvasObject{cont},
	}
}

// roundIconButtonRenderer keeps round button visuals in sync
type roundIconButtonRenderer struct {
	btn     *RoundIconButton
	circle  *canvas.Circle
	icon    *widget.Icon
	cont    *fyne.Container
	objects []fyne.CanvasObject
}

// Layout resizes round button content
func (r *roundIconButtonRenderer) Layout(size fyne.Size) { r.cont.Resize(size) }

// MinSize returns round button minimum size
func (r *roundIconButtonRenderer) MinSize() fyne.Size { return r.btn.MinSize() }

// BackgroundColor keeps renderer transparent
func (r *roundIconButtonRenderer) BackgroundColor() fyne.ThemeColorName { return "" }

// Objects returns round button objects
func (r *roundIconButtonRenderer) Objects() []fyne.CanvasObject { return r.objects }

// Destroy releases renderer resources
func (r *roundIconButtonRenderer) Destroy() {}

// Refresh redraws round button state
func (r *roundIconButtonRenderer) Refresh() {
	r.circle.FillColor = roundIconButtonBg(r.btn)

	threading.RunOnMainThread(func() {
		r.circle.Refresh()
		r.icon.Refresh()
	})
}

// MinSize returns fixed round button size
func (b *RoundIconButton) MinSize() fyne.Size {
	return fyne.NewSize(roundIconButtonSize, roundIconButtonSize)
}

// Tapped runs button callback with press flash
func (b *RoundIconButton) Tapped(*fyne.PointEvent) {
	orig := helpers.ToNRGBA(b.Bg)
	b.Bg = helpers.Darken(orig, buttonPressDarkenFactor)

	threading.RunOnMainThread(func() {
		b.Refresh()
	})

	go func() {
		time.Sleep(buttonPressDelay)
		threading.RunOnMainThread(func() {
			b.Bg = orig
			b.Refresh()
		})
	}()

	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// MouseIn marks round button as hovered
func (b *RoundIconButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	threading.RunOnMainThread(func() { b.Refresh() })
}

// MouseMoved satisfies desktop hover interface
func (b *RoundIconButton) MouseMoved(*desktop.MouseEvent) {}

// MouseOut clears round button hover state
func (b *RoundIconButton) MouseOut() {
	b.hovered = false
	threading.RunOnMainThread(func() { b.Refresh() })
}

// SimpleRectButton is a minimal rounded rectangle button
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

// NewSimpleRectButton builds rounded rectangle button
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

// CreateRenderer builds rectangle button renderer
func (b *SimpleRectButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(helpers.ToNRGBA(b.Bg))
	bg.CornerRadius = b.Radius

	txt := canvas.NewText(b.Text, helpers.ToNRGBA(b.Fg))
	txt.Alignment = fyne.TextAlignCenter
	txt.TextStyle = fyne.TextStyle{Bold: true}

	cont := container.NewMax(bg, container.NewCenter(txt))

	return &simpleRectButtonRenderer{
		button:  b,
		bg:      bg,
		text:    txt,
		cont:    cont,
		objects: []fyne.CanvasObject{cont},
	}
}

// simpleRectButtonRenderer keeps rectangle button visuals in sync
type simpleRectButtonRenderer struct {
	button  *SimpleRectButton
	bg      *canvas.Rectangle
	text    *canvas.Text
	cont    *fyne.Container
	objects []fyne.CanvasObject
}

// Layout resizes rectangle button content
func (r *simpleRectButtonRenderer) Layout(size fyne.Size) {
	r.cont.Resize(size)
}

// MinSize returns rectangle button minimum size
func (r *simpleRectButtonRenderer) MinSize() fyne.Size {
	return r.button.MinSize()
}

// Refresh redraws rectangle button state
func (r *simpleRectButtonRenderer) Refresh() {
	bgCol, fgCol := simpleRectButtonColors(r.button)

	r.bg.FillColor = bgCol
	r.bg.CornerRadius = r.button.Radius
	r.text.Color = fgCol
	r.text.Text = r.button.Text

	threading.RunOnMainThread(func() {
		r.text.Refresh()
		r.bg.Refresh()
	})
}

// BackgroundColor keeps renderer transparent
func (r *simpleRectButtonRenderer) BackgroundColor() fyne.ThemeColorName {
	return ""
}

// Objects returns rectangle button objects
func (r *simpleRectButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy releases renderer resources
func (r *simpleRectButtonRenderer) Destroy() {}

// MinSize returns fixed or default rectangle button size
func (b *SimpleRectButton) MinSize() fyne.Size {
	if b.SizeFixed.Width > 0 && b.SizeFixed.Height > 0 {
		return b.SizeFixed
	}

	return fyne.NewSize(simpleRectButtonDefaultWidth, simpleRectButtonDefaultHeight)
}

// Tapped runs button callback with press flash
func (b *SimpleRectButton) Tapped(*fyne.PointEvent) {
	if b.Disabled {
		return
	}

	orig := helpers.ToNRGBA(b.Bg)
	b.Bg = helpers.Darken(orig, buttonPressDarkenFactor)

	threading.RunOnMainThread(func() {
		b.Refresh()
	})

	go func() {
		time.Sleep(buttonPressDelay)

		threading.RunOnMainThread(func() {
			b.Bg = orig
			b.Refresh()
		})
	}()

	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// SetText changes rectangle button text
func (b *SimpleRectButton) SetText(text string) {
	b.Text = text

	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

// Enable restores rectangle button interactions
func (b *SimpleRectButton) Enable() {
	if !b.Disabled {
		return
	}

	b.Disabled = false

	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

// Disable blocks rectangle button interactions
func (b *SimpleRectButton) Disable() {
	if b.Disabled {
		return
	}

	b.Disabled = true
	b.hovered = false

	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

// MouseIn marks rectangle button as hovered
func (b *SimpleRectButton) MouseIn(*desktop.MouseEvent) {
	if !b.Disabled {
		b.hovered = true

		threading.RunOnMainThread(func() {
			b.Refresh()
		})
	}
}

// MouseMoved satisfies desktop hover interface
func (b *SimpleRectButton) MouseMoved(*desktop.MouseEvent) {}

// MouseOut clears rectangle button hover state
func (b *SimpleRectButton) MouseOut() {
	if !b.Disabled {
		b.hovered = false

		threading.RunOnMainThread(func() {
			b.Refresh()
		})
	}
}

// TinyIconButton is a minimal icon button for compact layouts
type TinyIconButton struct {
	widget.BaseWidget
	Icon     fyne.Resource
	OnTapped func()

	pressed bool
	hovered bool
}

// NewTinyIconButton builds compact icon button
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

// CreateRenderer builds compact icon button renderer
func (b *TinyIconButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	bg.CornerRadius = tinyIconButtonCornerRadius

	icon := widget.NewIcon(b.Icon)

	return &tinyIconButtonRenderer{
		button:  b,
		icon:    icon,
		bg:      bg,
		objects: []fyne.CanvasObject{bg, icon},
	}
}

func (b *TinyIconButton) MinSize() fyne.Size {
	return fyne.NewSize(tinyIconButtonSize, tinyIconButtonSize)
}

// Tapped runs icon button callback with press flash
func (b *TinyIconButton) Tapped(*fyne.PointEvent) {
	b.pressed = true
	threading.RunOnMainThread(func() {
		b.Refresh()
	})

	if b.OnTapped != nil {
		b.OnTapped()
	}

	go func() {
		time.Sleep(buttonPressDelay)

		threading.RunOnMainThread(func() {
			b.pressed = false
			b.Refresh()
		})
	}()
}

// MouseIn marks compact icon as hovered
func (b *TinyIconButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true

	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

func (b *TinyIconButton) MouseMoved(*desktop.MouseEvent) {}

// MouseOut clears compact icon hover state
func (b *TinyIconButton) MouseOut() {
	b.hovered = false

	threading.RunOnMainThread(func() {
		b.Refresh()
	})
}

// tinyIconButtonRenderer keeps compact icon visuals in sync
type tinyIconButtonRenderer struct {
	button  *TinyIconButton
	icon    *widget.Icon
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

// Layout resizes compact icon content
func (r *tinyIconButtonRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.bg.Move(fyne.NewPos(0, 0))
	r.icon.Resize(size)
	r.icon.Move(fyne.NewPos(0, 0))
}

func (r *tinyIconButtonRenderer) MinSize() fyne.Size {
	return r.button.MinSize()
}

// Refresh redraws compact icon state
func (r *tinyIconButtonRenderer) Refresh() {
	r.bg.FillColor = tinyIconButtonBg(r.button)

	threading.RunOnMainThread(func() {
		r.bg.Refresh()
		r.icon.Refresh()
	})
}

// Objects returns compact icon button objects
func (r *tinyIconButtonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy releases renderer resources
func (r *tinyIconButtonRenderer) Destroy() {}

// BackgroundColor keeps renderer transparent
func (r *tinyIconButtonRenderer) BackgroundColor() fyne.ThemeColorName {
	return ""
}

// roundIconButtonBg returns current round button background
func roundIconButtonBg(button *RoundIconButton) color.NRGBA {
	bg := helpers.ToNRGBA(button.Bg)

	if button.hovered {
		return helpers.Lighten(bg, roundIconButtonHoverAmount)
	}

	return bg
}

// simpleRectButtonColors returns current rectangle button colors
func simpleRectButtonColors(button *SimpleRectButton) (color.NRGBA, color.NRGBA) {
	bgCol := helpers.ToNRGBA(button.Bg)
	fgCol := helpers.ToNRGBA(button.Fg)

	if button.Disabled {
		bgCol = helpers.Darken(bgCol, simpleRectButtonDisabledDim)
		fgCol.A = simpleRectButtonDisabledAlpha
	} else if button.hovered {
		bgCol = helpers.Lighten(bgCol, simpleRectButtonHoverAmount)
	}

	return bgCol, fgCol
}

// tinyIconButtonBg returns current compact icon background
func tinyIconButtonBg(button *TinyIconButton) color.NRGBA {
	if button.pressed {
		return color.NRGBA{R: 100, G: 100, B: 100, A: tinyIconButtonPressAlpha}
	}

	if button.hovered {
		return color.NRGBA{R: 200, G: 200, B: 200, A: tinyIconButtonHoverAlpha}
	}

	return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
}
