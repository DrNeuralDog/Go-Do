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
	if helpers.IsLightTheme() {
		// Light theme: use dark text for visibility on white background
		cs.TextColor = helpers.Hex("#3c3836") // Gruvbox dark gray
	} else {
		// Dark theme: use light text
		cs.TextColor = helpers.Hex("#ebdbb2") // Gruvbox light
	}

	cs.ExtendBaseWidget(cs)
	return cs
}

func (cs *CustomSelect) SetSelected(s string) {
	cs.Selected = s
	threading.RunOnMainThread(func() {
		cs.Refresh()
	})
}

func (cs *CustomSelect) MinSize() fyne.Size {
	return fyne.NewSize(180, 44)
}

func (cs *CustomSelect) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(cs.Selected, helpers.ToNRGBA(cs.TextColor))
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
	if helpers.IsLightTheme() {
		// Light theme: use dark text for visibility on white background
		r.text.Color = helpers.ToNRGBA(helpers.Hex("#3c3836")) // Dark gray text
	} else {
		// Dark theme: use light text
		r.text.Color = helpers.ToNRGBA(helpers.Hex("#ebdbb2")) // Light text
	}

	// hover / press overlay (theme-aware)
	var col color.NRGBA
	if r.select_.pressed {
		// subtle dark press for both themes
		col = color.NRGBA{R: 0, G: 0, B: 0, A: 70}
	} else if r.select_.hovered {
		if helpers.IsLightTheme() {
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

func (cs *CustomSelect) Tapped(_ *fyne.PointEvent) {
	// brief press flash
	cs.pressed = true
	threading.RunOnMainThread(func() {
		cs.Refresh()
	})
	go func(s *CustomSelect) {
		time.Sleep(120 * time.Millisecond)
		threading.RunOnMainThread(func() {
			s.pressed = false
			s.Refresh()
		})
	}(cs)
	// Create a popup menu with options
	items := make([]*fyne.MenuItem, len(cs.Options))
	for i, opt := range cs.Options {
		option := opt // Capture loop variable
		items[i] = fyne.NewMenuItem(option, func() {
			cs.Selected = option
			threading.RunOnMainThread(func() {
				cs.Refresh()
			})
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
func (cs *CustomSelect) MouseIn(*desktop.MouseEvent) {
	cs.hovered = true
	threading.RunOnMainThread(func() { cs.Refresh() })
}
func (cs *CustomSelect) MouseMoved(*desktop.MouseEvent) {}
func (cs *CustomSelect) MouseOut() {
	cs.hovered = false
	threading.RunOnMainThread(func() { cs.Refresh() })
}
