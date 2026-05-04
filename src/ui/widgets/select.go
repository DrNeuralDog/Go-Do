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
	customSelectWidth        = 180
	customSelectHeight       = 44
	customSelectTextSize     = 14
	customSelectCornerRadius = 8
	customSelectPressAlpha   = 70
	customSelectHoverAlpha   = 26
	customSelectPressDelay   = 120 * time.Millisecond
)

// CustomSelect draws a dropdown without native focus highlight
type CustomSelect struct {
	widget.BaseWidget
	Options   []string
	Selected  string
	OnChanged func(string)
	hovered   bool
	pressed   bool
}

// NewCustomSelect builds custom dropdown selector
func NewCustomSelect(options []string, onChanged func(string)) *CustomSelect {
	cs := &CustomSelect{
		Options:   options,
		OnChanged: onChanged,
	}

	cs.ExtendBaseWidget(cs)

	return cs
}

// SetSelected updates visible selected option
func (cs *CustomSelect) SetSelected(value string) {
	cs.setSelected(value, false)
}

// MinSize returns fixed select size
func (cs *CustomSelect) MinSize() fyne.Size {
	return fyne.NewSize(customSelectWidth, customSelectHeight)
}

// CreateRenderer builds custom select renderer
func (cs *CustomSelect) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(cs.Selected, selectTextColor())
	text.Alignment = fyne.TextAlignLeading
	text.TextSize = customSelectTextSize

	icon := widget.NewIcon(theme.MenuDropDownIcon())

	base := container.NewBorder(nil, nil, nil, icon, container.NewPadded(text))
	overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	overlay.CornerRadius = customSelectCornerRadius
	cont := container.NewMax(base, overlay)

	return &customSelectRenderer{
		selectWidget: cs,
		text:         text,
		cont:         cont,
		overlay:      overlay,
	}
}

// customSelectRenderer keeps select visuals in sync
type customSelectRenderer struct {
	selectWidget *CustomSelect
	text         *canvas.Text
	cont         *fyne.Container
	overlay      *canvas.Rectangle
}

func (r *customSelectRenderer) Layout(size fyne.Size) {
	r.cont.Resize(size)
}

func (r *customSelectRenderer) MinSize() fyne.Size {
	return r.selectWidget.MinSize()
}

func (r *customSelectRenderer) Refresh() {
	r.text.Text = r.selectWidget.Selected
	r.text.Color = selectTextColor()
	r.overlay.FillColor = selectOverlayColor(r.selectWidget.hovered, r.selectWidget.pressed)

	threading.RunOnMainThread(func() {
		r.text.Refresh()
		r.overlay.Refresh()
		r.cont.Refresh()
	})
}

func (r *customSelectRenderer) BackgroundColor() fyne.ThemeColorName {
	return ""
}

func (r *customSelectRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.cont}
}

func (r *customSelectRenderer) Destroy() {}

// Tapped opens options popup
func (cs *CustomSelect) Tapped(_ *fyne.PointEvent) {
	cs.pressed = true
	threading.RunOnMainThread(func() {
		cs.Refresh()
	})

	go func(selectWidget *CustomSelect) {
		time.Sleep(customSelectPressDelay)
		threading.RunOnMainThread(func() {
			selectWidget.pressed = false
			selectWidget.Refresh()
		})
	}(cs)

	items := make([]*fyne.MenuItem, len(cs.Options))

	for i, option := range cs.Options {
		option := option

		items[i] = fyne.NewMenuItem(option, func() {
			cs.setSelected(option, true)
		})
	}

	menu := fyne.NewMenu("", items...)
	canvas := fyne.CurrentApp().Driver().CanvasForObject(cs)
	position := fyne.CurrentApp().Driver().AbsolutePositionForObject(cs)
	popup := widget.NewPopUpMenu(menu, canvas)
	popup.ShowAtPosition(position)

	// Ширину держим как у самого селекта
	ctrlWidth := cs.Size().Width
	min := popup.MinSize()
	popup.Resize(fyne.NewSize(ctrlWidth, min.Height))
}

// setSelected updates value and optionally fires callback
func (cs *CustomSelect) setSelected(value string, notify bool) {
	if cs.Selected == value {
		return
	}

	cs.Selected = value

	threading.RunOnMainThread(func() {
		cs.Refresh()
	})

	if notify && cs.OnChanged != nil {
		cs.OnChanged(value)
	}
}

// selectTextColor returns theme-aware select text color
func selectTextColor() color.NRGBA {
	if helpers.IsLightTheme() {
		return helpers.ToNRGBA(helpers.Hex("#3c3836"))
	}

	return helpers.ToNRGBA(helpers.Hex("#ebdbb2"))
}

// selectOverlayColor returns hover or press overlay color
func selectOverlayColor(hovered, pressed bool) color.NRGBA {
	if pressed {
		return color.NRGBA{R: 0, G: 0, B: 0, A: customSelectPressAlpha}
	}

	if !hovered {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}

	if helpers.IsLightTheme() {
		return color.NRGBA{R: 0, G: 0, B: 0, A: customSelectHoverAlpha}
	}

	return color.NRGBA{R: 255, G: 255, B: 255, A: customSelectHoverAlpha}
}

// FocusGained keeps native focus highlight disabled
func (cs *CustomSelect) FocusGained() {}

// FocusLost keeps native focus highlight disabled
func (cs *CustomSelect) FocusLost() {}

// MouseIn marks select as hovered
func (cs *CustomSelect) MouseIn(*desktop.MouseEvent) {
	cs.hovered = true

	threading.RunOnMainThread(func() { cs.Refresh() })
}

// MouseMoved satisfies desktop hover interface
func (cs *CustomSelect) MouseMoved(*desktop.MouseEvent) {}

// MouseOut clears hover state
func (cs *CustomSelect) MouseOut() {
	cs.hovered = false

	threading.RunOnMainThread(func() { cs.Refresh() })
}
