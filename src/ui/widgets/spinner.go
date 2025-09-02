package widgets

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

	"godo/src/ui/helpers"
	"godo/src/ui/threading"
)

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
	threading.RunOnMainThread(func() {
		ns.Refresh()
	})
	if ns.OnChanged != nil {
		ns.OnChanged(v)
	}
}

func (ns *NumberSpinner) MinSize() fyne.Size {
	// увеличили высоту примерно на 10% (32 -> 35)
	return fyne.NewSize(160, 35)
}

func (ns *NumberSpinner) CreateRenderer() fyne.WidgetRenderer {
	// background (white by requirement)
	bg := canvas.NewRectangle(helpers.ToNRGBA(ns.BgColor))
	bg.CornerRadius = 8
	// subtle border to mimic input field look
	sep := helpers.ToNRGBA(theme.Color(theme.ColorNameSeparator))
	sep.A = 255
	bg.StrokeColor = sep
	bg.StrokeWidth = 1

	// display value text - STORE THIS IN RENDERER
	txt := canvas.NewText(fmt.Sprintf("%d", ns.Value), helpers.ToNRGBA(ns.TextColor))
	txt.TextSize = 16
	txt.Alignment = fyne.TextAlignLeading

	// up/down buttons (compact 16x16 icons, wrapped to 24x24 each for total 48px height)
	// Use Center to vertically align arrows in the middle
	// Note: ArrowUpIcon and ArrowDownIcon need to be imported from parent ui package
	upBtn := NewTinyIconButton(theme.MoveUpIcon(), func() { ns.increment() })
	downBtn := NewTinyIconButton(theme.MoveDownIcon(), func() { ns.decrement() })
	upDown := container.NewVBox(
		container.NewCenter(upBtn),
		container.NewCenter(downBtn),
	)

	// left padding for text
	textPadded := container.NewBorder(nil, nil, helpers.CreateSpacer(12, 1), nil, container.NewCenter(txt))

	content := container.NewBorder(nil, nil, nil, upDown, textPadded)
	stack := container.NewStack(bg, content)

	// Return custom renderer that can update text
	return &numberSpinnerRenderer{
		spinner: ns,
		bg:      bg,
		txt:     txt,
		upBtn:   upBtn,
		downBtn: downBtn,
		stack:   stack,
	}
}

// numberSpinnerRenderer is a custom renderer that properly updates text on value change
type numberSpinnerRenderer struct {
	spinner *NumberSpinner
	bg      *canvas.Rectangle
	txt     *canvas.Text
	upBtn   *TinyIconButton
	downBtn *TinyIconButton
	stack   *fyne.Container
}

func (r *numberSpinnerRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
}

func (r *numberSpinnerRenderer) MinSize() fyne.Size {
	return r.spinner.MinSize()
}

func (r *numberSpinnerRenderer) Refresh() {
	// Update text display with current value
	r.txt.Text = fmt.Sprintf("%d", r.spinner.Value)
	threading.RunOnMainThread(func() {
		r.txt.Refresh()
		r.bg.Refresh()
		r.stack.Refresh()
	})
}

func (r *numberSpinnerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.stack}
}

func (r *numberSpinnerRenderer) Destroy() {
	// Cleanup if needed
}

func (r *numberSpinnerRenderer) BackgroundColor() fyne.ThemeColorName {
	return ""
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
