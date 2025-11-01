package ui

import (
    "fmt"
    "image/color"
    "strconv"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
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
	sep.A = 255 // Full opacity for visibility
	bg.StrokeColor = sep
	bg.StrokeWidth = 1

	return container.NewMax(bg, container.NewPadded(obj))
}

// CreateCardStyle wraps content in a rounded card with padding and subtle border.
func CreateCardStyle(obj fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(toNRGBA(theme.Color(theme.ColorNameInputBackground)))
	bg.CornerRadius = 12

	sep := toNRGBA(theme.Color(theme.ColorNameSeparator))
	sep.A = 255 // Full opacity for visibility
	bg.StrokeColor = sep
	bg.StrokeWidth = 1

	return container.NewMax(bg, container.NewPadded(obj))
}

// RoundedIconButton creates a circular accent button with a centered icon.
func RoundedIconButton(icon fyne.Resource, tapped func()) fyne.CanvasObject {
	return NewRoundIconButton(icon, tapped)
}

// RoundIconButton is a circular icon button that is fully clickable without rectangular hover overlay.
type RoundIconButton struct {
	widget.BaseWidget
	Icon     fyne.Resource
	Bg       color.Color
	Fg       color.Color
	OnTapped func()
}

func NewRoundIconButton(icon fyne.Resource, onTap func()) *RoundIconButton {
	b := &RoundIconButton{
		Icon:     icon,
		Bg:       toNRGBA(theme.Color(theme.ColorNamePrimary)),
		Fg:       toNRGBA(theme.Color(theme.ColorNameForeground)),
		OnTapped: onTap,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *RoundIconButton) CreateRenderer() fyne.WidgetRenderer {
	circle := canvas.NewCircle(toNRGBA(b.Bg))
	icon := widget.NewIcon(b.Icon)
	cont := container.NewMax(circle, container.NewCenter(icon))
	return widget.NewSimpleRenderer(cont)
}

func (b *RoundIconButton) MinSize() fyne.Size {
	// Default minimum; actual size is controlled by parent container (GridWrap)
	return fyne.NewSize(36, 36)
}

func (b *RoundIconButton) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
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
	// Update text content from button
	r.text.Text = r.button.Text
	r.text.Refresh()
	r.bg.Refresh()
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
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *SimpleRectButton) SetText(text string) {
	b.Text = text
	b.Refresh()
}

// CreateStyledSelect wraps a Select widget in a rounded container with custom background
func CreateStyledSelect(selectWidget *widget.Select, bgColor color.Color, size fyne.Size, radius float32) fyne.CanvasObject {
	bg := canvas.NewRectangle(toNRGBA(bgColor))
	bg.CornerRadius = radius

	// Wrap select in container with fixed size
	selectWrapper := container.NewGridWrap(size, selectWidget)

	return container.NewMax(
		container.NewGridWrap(size, bg),
		selectWrapper,
	)
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
	rect := canvas.NewRectangle(toNRGBA(bgColor))
	rect.CornerRadius = 12
	// Use custom padding to control exact spacing - 8px all around instead of default Fyne padding
	padded := container.NewPadded(content)
	return container.NewMax(rect, padded)
}

// NumberSpinner is a simple numeric stepper with up/down arrows on the right and a configurable background.
type NumberSpinner struct {
    widget.BaseWidget
    Window   fyne.Window
    Value     int
    Min       int
    Max       int
    Step      int
    TextColor color.Color
    BgColor   color.Color
    OnChanged func(int)
}

func NewNumberSpinner(win fyne.Window, initial, min, max, step int, textColor, bgColor color.Color, onChanged func(int)) *NumberSpinner {
    ns := &NumberSpinner{
        Window:   win,
        Value:     initial,
        Min:       min,
        Max:       max,
        Step:      step,
        TextColor: textColor,
        BgColor:   bgColor,
        OnChanged: onChanged,
    }
    ns.ExtendBaseWidget(ns)
    return ns
}

func (ns *NumberSpinner) increment() {
    next := ns.Value + ns.Step
    if ns.Max != 0 && next > ns.Max {
        next = ns.Max
    }
    ns.SetValue(next)
}

func (ns *NumberSpinner) decrement() {
    next := ns.Value - ns.Step
    if next < ns.Min {
        next = ns.Min
    }
    ns.SetValue(next)
}

func (ns *NumberSpinner) SetValue(v int) {
    if ns.Max != 0 && v > ns.Max {
        v = ns.Max
    }
    if v < ns.Min {
        v = ns.Min
    }
    if ns.Value == v {
        return
    }
    ns.Value = v
    ns.Refresh()
    if ns.OnChanged != nil {
        ns.OnChanged(v)
    }
}

func (ns *NumberSpinner) MinSize() fyne.Size {
    return fyne.NewSize(140, 36)
}

func (ns *NumberSpinner) CreateRenderer() fyne.WidgetRenderer {
    // background (white by requirement)
    bg := canvas.NewRectangle(toNRGBA(ns.BgColor))
    bg.CornerRadius = 8

    // display value text
    txt := canvas.NewText(fmt.Sprintf("%d", ns.Value), toNRGBA(ns.TextColor))
    txt.TextSize = 16
    txt.Alignment = fyne.TextAlignLeading

    // up/down buttons (compact)
    upBtn := widget.NewButton("▲", func() { ns.increment() })
    downBtn := widget.NewButton("▼", func() { ns.decrement() })
    upDown := container.NewVBox(
        container.NewGridWrap(fyne.NewSize(24, 18), upBtn),
        container.NewGridWrap(fyne.NewSize(24, 18), downBtn),
    )

    // left padding for text
    textPadded := container.NewBorder(nil, nil, CreateSpacer(10, 1), nil, container.NewCenter(txt))

    content := container.NewBorder(nil, nil, nil, upDown, textPadded)
    return widget.NewSimpleRenderer(container.NewMax(bg, content))
}

// Tapped opens a small dialog to allow manual numeric input
func (ns *NumberSpinner) Tapped(*fyne.PointEvent) {
    if ns.Window == nil {
        return
    }
    entry := widget.NewEntry()
    entry.SetText(fmt.Sprintf("%d", ns.Value))
    entry.PlaceHolder = "minutes"
    // simple numeric filter (allow empty while editing)
    entry.OnChanged = func(s string) {
        if s == "" {
            return
        }
        if _, err := strconv.Atoi(s); err != nil {
            // strip non-digits
            digits := make([]rune, 0, len(s))
            for _, r := range s {
                if r >= '0' && r <= '9' {
                    digits = append(digits, r)
                }
            }
            entry.SetText(string(digits))
        }
    }
    content := container.NewVBox(
        widget.NewLabel("Enter minutes:"),
        entry,
    )
    dialog.NewCustomConfirm("Set value", "OK", "Cancel", content, func(ok bool) {
        if !ok {
            return
        }
        v, err := strconv.Atoi(entry.Text)
        if err != nil {
            return
        }
        if ns.Max != 0 && v > ns.Max {
            v = ns.Max
        }
        if v < ns.Min {
            v = ns.Min
        }
        ns.SetValue(v)
    }, ns.Window).Show()
}
