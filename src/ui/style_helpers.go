package ui

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
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
	hovered  bool
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
	bg := toNRGBA(r.btn.Bg)
	if r.btn.hovered {
		bg = lighten(bg, 0.12)
	}
	r.circle.FillColor = bg
	r.circle.Refresh()
	r.icon.Refresh()
}

func (b *RoundIconButton) MinSize() fyne.Size {
	// Default minimum; actual size is controlled by parent container (GridWrap)
	return fyne.NewSize(36, 36)
}

func (b *RoundIconButton) Tapped(*fyne.PointEvent) {
	// quick press flash by darkening background briefly
	orig := toNRGBA(b.Bg)
	b.Bg = darken(orig, 0.85)
	b.Refresh()
	go func() {
		time.Sleep(120 * time.Millisecond)
		b.Bg = orig
		b.Refresh()
	}()
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

// Hover handling (desktop only)
func (b *RoundIconButton) MouseIn(*desktop.MouseEvent)    { b.hovered = true; b.Refresh() }
func (b *RoundIconButton) MouseMoved(*desktop.MouseEvent) {}
func (b *RoundIconButton) MouseOut()                      { b.hovered = false; b.Refresh() }

// toNRGBA converts any color.Color to color.NRGBA for manipulation.
func toNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// darken returns a darker version of the color by the provided factor (0..1)
func darken(c color.NRGBA, factor float32) color.NRGBA {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}
	return color.NRGBA{
		R: uint8(float32(c.R) * factor),
		G: uint8(float32(c.G) * factor),
		B: uint8(float32(c.B) * factor),
		A: c.A,
	}
}

// lighten returns a lighter version of the color by mixing with white (0..1)
func lighten(c color.NRGBA, amount float32) color.NRGBA {
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	mix := func(v uint8) uint8 {
		return uint8(float32(v)*(1-amount) + 255*amount)
	}
	return color.NRGBA{R: mix(c.R), G: mix(c.G), B: mix(c.B), A: c.A}
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
	// Sync colors & text from button state
	btn := r.button
	bgCol := toNRGBA(btn.Bg)
	fgCol := toNRGBA(btn.Fg)
	if btn.Disabled {
		// dim colors
		bgCol = darken(bgCol, 0.6)
		fgCol.A = 160
	} else if btn.hovered {
		bgCol = lighten(bgCol, 0.10)
	}
	r.bg.FillColor = bgCol
	r.text.Color = fgCol
	r.text.Text = btn.Text
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
	if b.Disabled {
		return
	}
	// flash background on press
	orig := toNRGBA(b.Bg)
	b.Bg = darken(orig, 0.85)
	b.Refresh()
	go func() {
		time.Sleep(120 * time.Millisecond)
		b.Bg = orig
		b.Refresh()
	}()
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *SimpleRectButton) SetText(text string) {
	b.Text = text
	b.Refresh()
}

func (b *SimpleRectButton) Enable() {
	if !b.Disabled {
		return
	}
	b.Disabled = false
	b.Refresh()
}

func (b *SimpleRectButton) Disable() {
	if b.Disabled {
		return
	}
	b.Disabled = true
	b.Refresh()
}

// Hover handling (desktop only)
func (b *SimpleRectButton) MouseIn(*desktop.MouseEvent) {
	if !b.Disabled {
		b.hovered = true
		b.Refresh()
	}
}
func (b *SimpleRectButton) MouseMoved(*desktop.MouseEvent) {}
func (b *SimpleRectButton) MouseOut() {
	if !b.Disabled {
		b.hovered = false
		b.Refresh()
	}
}

// CustomSelect is a custom dropdown widget with proper text color and no press/focus highlighting
type CustomSelect struct {
	widget.BaseWidget
	Options   []string
	Selected  string
	OnChanged func(string)
	TextColor color.Color
	window    fyne.Window
	hovered   bool
	pressed   bool
}

func NewCustomSelect(options []string, onChanged func(string)) *CustomSelect {
	cs := &CustomSelect{
		Options:   options,
		OnChanged: onChanged,
	}

	// Set text color based on theme
	currentTheme := fyne.CurrentApp().Settings().Theme()
	if _, ok := currentTheme.(*LightSoftTheme); ok {
		// Light theme: use dark text for visibility on white background
		cs.TextColor = hex("#3c3836") // Gruvbox dark gray
	} else {
		// Dark theme: use light text
		cs.TextColor = hex("#ebdbb2") // Gruvbox light
	}

	cs.ExtendBaseWidget(cs)
	return cs
}

func (cs *CustomSelect) SetSelected(s string) {
	cs.Selected = s
	cs.Refresh()
}

func (cs *CustomSelect) MinSize() fyne.Size {
	return fyne.NewSize(180, 44)
}

func (cs *CustomSelect) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(cs.Selected, toNRGBA(cs.TextColor))
	text.Alignment = fyne.TextAlignLeading
	text.TextSize = 14

	icon := widget.NewIcon(theme.MenuDropDownIcon())

	base := container.NewBorder(nil, nil, nil, icon, container.NewPadded(text))
	overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	overlay.CornerRadius = 8
	cont := container.NewMax(base, overlay)

	return &customSelectRenderer{
		select_: cs,
		text:    text,
		icon:    icon,
		cont:    cont,
		overlay: overlay,
	}
}

type customSelectRenderer struct {
	select_ *CustomSelect
	text    *canvas.Text
	icon    *widget.Icon
	cont    *fyne.Container
	overlay *canvas.Rectangle
}

func (r *customSelectRenderer) Layout(size fyne.Size) {
	r.cont.Resize(size)
}

func (r *customSelectRenderer) MinSize() fyne.Size {
	return r.select_.MinSize()
}

func (r *customSelectRenderer) Refresh() {
	r.text.Text = r.select_.Selected

	// Dynamically set text color based on current theme
	currentTheme := fyne.CurrentApp().Settings().Theme()
	if _, ok := currentTheme.(*LightSoftTheme); ok {
		// Light theme: use dark text for visibility on white background
		r.text.Color = toNRGBA(hex("#3c3836")) // Dark gray text
	} else {
		// Dark theme: use light text
		r.text.Color = toNRGBA(hex("#ebdbb2")) // Light text
	}
	r.text.Refresh()

	// hover / press overlay (theme-aware)
	var col color.NRGBA
	if r.select_.pressed {
		// subtle dark press for both themes
		col = color.NRGBA{R: 0, G: 0, B: 0, A: 70}
	} else if r.select_.hovered {
		if _, ok := fyne.CurrentApp().Settings().Theme().(*LightSoftTheme); ok {
			// light theme: use dark translucent overlay
			col = color.NRGBA{R: 0, G: 0, B: 0, A: 26} // ~10% black
		} else {
			// dark theme: use light translucent overlay
			col = color.NRGBA{R: 255, G: 255, B: 255, A: 26} // ~10% white
		}
	} else {
		col = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}
	r.overlay.FillColor = col
	r.overlay.Refresh()
	r.cont.Refresh()
}

func (r *customSelectRenderer) BackgroundColor() fyne.ThemeColorName {
	return ""
}

func (r *customSelectRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.cont}
}

func (r *customSelectRenderer) Destroy() {}

func (cs *CustomSelect) Tapped(_ *fyne.PointEvent) {
	// brief press flash
	cs.pressed = true
	cs.Refresh()
	go func(s *CustomSelect) {
		time.Sleep(120 * time.Millisecond)
		s.pressed = false
		s.Refresh()
	}(cs)
	// Create a popup menu with options
	items := make([]*fyne.MenuItem, len(cs.Options))
	for i, opt := range cs.Options {
		option := opt // Capture loop variable
		items[i] = fyne.NewMenuItem(option, func() {
			cs.Selected = option
			cs.Refresh()
			if cs.OnChanged != nil {
				cs.OnChanged(option)
			}
		})
	}

	// Show popup menu at widget position with width matching control width
	m := fyne.NewMenu("", items...)
	cnv := fyne.CurrentApp().Driver().CanvasForObject(cs)
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(cs)
	popup := widget.NewPopUpMenu(m, cnv)
	popup.ShowAtPosition(pos)
	// enforce popup width = control width
	ctrlWidth := cs.Size().Width
	min := popup.MinSize()
	popup.Resize(fyne.NewSize(ctrlWidth, min.Height))
}

// FocusGained is overridden to do nothing (no focus highlight)
func (cs *CustomSelect) FocusGained() {
	// Don't show focus highlight
}

// FocusLost is overridden to do nothing
func (cs *CustomSelect) FocusLost() {
	// Don't show focus highlight
}

// Hover handling (desktop only)
func (cs *CustomSelect) MouseIn(*desktop.MouseEvent)    { cs.hovered = true; cs.Refresh() }
func (cs *CustomSelect) MouseMoved(*desktop.MouseEvent) {}
func (cs *CustomSelect) MouseOut()                      { cs.hovered = false; cs.Refresh() }

// CreateStyledSelect wraps a Select widget in a rounded container with custom background
// Accepts any CanvasObject (including widget.Select and CustomSelect)
func CreateStyledSelect(selectWidget fyne.CanvasObject, bgColor color.Color, size fyne.Size, radius float32) fyne.CanvasObject {
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
	Window    fyne.Window
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
		Window:    win,
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
	// Increased to ensure arrow buttons have enough space and do not overlap
	return fyne.NewSize(160, 48)
}

func (ns *NumberSpinner) CreateRenderer() fyne.WidgetRenderer {
	// background (white by requirement)
	bg := canvas.NewRectangle(toNRGBA(ns.BgColor))
	bg.CornerRadius = 8
	// subtle border to mimic input field look
	sep := toNRGBA(theme.Color(theme.ColorNameSeparator))
	sep.A = 255
	bg.StrokeColor = sep
	bg.StrokeWidth = 1

	// display value text
	txt := canvas.NewText(fmt.Sprintf("%d", ns.Value), toNRGBA(ns.TextColor))
	txt.TextSize = 16
	txt.Alignment = fyne.TextAlignLeading

	// up/down buttons (icon-based, fixed-size to prevent layout shifts)
	upBtn := widget.NewButtonWithIcon("", ArrowUpIcon, func() { ns.increment() })
	downBtn := widget.NewButtonWithIcon("", ArrowDownIcon, func() { ns.decrement() })
	upDown := container.NewVBox(
		container.NewGridWrap(fyne.NewSize(32, 24), upBtn),
		container.NewGridWrap(fyne.NewSize(32, 24), downBtn),
	)

	// left padding for text
	textPadded := container.NewBorder(nil, nil, CreateSpacer(12, 1), nil, container.NewCenter(txt))

	content := container.NewBorder(nil, nil, nil, upDown, textPadded)
	return widget.NewSimpleRenderer(container.NewStack(bg, content))
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
