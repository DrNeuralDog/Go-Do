package forms

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"time"

	"godo/src/localization"
	"godo/src/models"
	"godo/src/persistence"
	"godo/src/ui/helpers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TodoForm represents a form for creating/editing todo items
type TodoForm struct {
	parentWindow fyne.Window // Main application window (for dialogs)
	formWindow   fyne.Window // Current form window (for standalone windows)
	dataManager  persistence.TodoRepository

	// Form fields
	nameEntry      *widget.Entry
	contentEntry   *widget.Entry
	placeEntry     *widget.Entry
	labelEntry     *widget.Entry
	dateTimeEntry  *widget.Entry
	dateTimeButton *widget.Button
	prioritySelect *widget.Select
	kindSelect     *widget.Select
	warnTimeSlider *ReminderSlider
	warnTimeLabel  *canvas.Text

	// Date/Time picker components
	selectedDateTime time.Time

	// Form container
	formContainer *fyne.Container

	// Current dialog reference for closing
	currentDialog dialog.Dialog

	// State
	isEditMode     bool
	originalTodo   *models.TodoItem
	originalTime   time.Time
	onSaveCallback func()
}

// NewTodoForm creates a new todo form dialog
func NewTodoForm(window fyne.Window, dataManager persistence.TodoRepository) *TodoForm {
	tf := &TodoForm{
		parentWindow: window,
		dataManager:  dataManager,
		isEditMode:   false,
	}

	tf.setupForm()
	return tf
}

// createFormItemWithWhiteLabel creates FormItem for dialog.NewForm
func createFormItemWithWhiteLabel(labelText string, w fyne.CanvasObject) *widget.FormItem {
	return &widget.FormItem{Text: labelText, Widget: w}
}

// ShowCreateDialog shows the form for creating a new todo
func (tf *TodoForm) ShowCreateDialog(onSave func()) {
	tf.isEditMode = false
	tf.onSaveCallback = onSave
	tf.resetForm()

	title := localization.GetString("form_title_add")

	// Use dialog.NewForm for proper form handling
	formItems := []*widget.FormItem{
		{Text: "Name:", Widget: tf.nameEntry},
		{Text: "Date/Time:", Widget: container.NewBorder(nil, nil, nil, tf.dateTimeButton, tf.dateTimeEntry)},
		{Text: "Location:", Widget: tf.placeEntry},
		{Text: "Label:", Widget: tf.labelEntry},
		{Text: "Type:", Widget: tf.kindSelect},
		{Text: "Priority:", Widget: tf.prioritySelect},
		{Text: "Reminder:", Widget: container.NewVBox(tf.warnTimeSlider, tf.warnTimeLabel)},
	}

	// Add content field as a separate form item
	contentFormItem := &widget.FormItem{
		Text:   "Content:",
		Widget: container.NewScroll(tf.contentEntry),
	}
	formItems = append(formItems, contentFormItem)

	dialog := dialog.NewForm(title, localization.GetString("form_button_add"), localization.GetString("form_button_cancel"), formItems, func(submitted bool) {
		if submitted {
			tf.onSubmit()
		}
	}, tf.parentWindow)
	// Make the Add dialog wider so the Date/Time row has enough space
	dialog.Resize(fyne.NewSize(700, 600))
	tf.currentDialog = dialog
	dialog.Show()
}

// ShowEditDialog shows the form for editing an existing todo
func (tf *TodoForm) ShowEditDialog(todo *models.TodoItem, originalTime time.Time, onSave func()) {
	tf.isEditMode = true
	tf.originalTodo = todo
	tf.originalTime = originalTime
	tf.onSaveCallback = onSave
	tf.populateForm(todo)

	title := localization.GetString("form_title_edit")

	// Use dialog.NewForm for proper form handling
	formItems := []*widget.FormItem{
		{Text: "Name:", Widget: tf.nameEntry},
		{Text: "Date/Time:", Widget: container.NewBorder(nil, nil, nil, tf.dateTimeButton, tf.dateTimeEntry)},
		{Text: "Location:", Widget: tf.placeEntry},
		{Text: "Label:", Widget: tf.labelEntry},
		{Text: "Type:", Widget: tf.kindSelect},
		{Text: "Priority:", Widget: tf.prioritySelect},
		{Text: "Reminder:", Widget: container.NewVBox(tf.warnTimeSlider, tf.warnTimeLabel)},
	}

	// Add content field as a separate form item
	contentFormItem := &widget.FormItem{
		Text:   "Content:",
		Widget: container.NewScroll(tf.contentEntry),
	}
	formItems = append(formItems, contentFormItem)

	dialog := dialog.NewForm(title, localization.GetString("form_button_save"), localization.GetString("form_button_cancel"), formItems, func(submitted bool) {
		if submitted {
			tf.onSubmit()
		}
	}, tf.parentWindow)
	// Keep edit dialog consistent with add dialog width
	dialog.Resize(fyne.NewSize(700, 600))
	tf.currentDialog = dialog
	dialog.Show()
}

// ShowCreateWindow opens the form in a standalone window for creating a new todo
func (tf *TodoForm) ShowCreateWindow(onSave func(), onWindowCreated func(fyne.Window), onWindowClosed func()) {
	tf.isEditMode = false
	tf.onSaveCallback = onSave
	tf.resetForm()

	title := localization.GetString("form_title_add")

	// Compute window size based on main window (always use parentWindow, not old formWindow)
	parentSize := fyne.NewSize(700, 600)
	if tf.parentWindow != nil && tf.parentWindow.Canvas() != nil {
		parentSize = tf.parentWindow.Canvas().Size()
	}
	targetW := parentSize.Width - 40
	if targetW < 200 {
		targetW = parentSize.Width // fallback, just in case
	}
	// Make the standalone window noticeably more compact in height
	targetH := parentSize.Height * 0.45

	win := fyne.CurrentApp().NewWindow(title)
	tf.formWindow = win

	// Set close callback to notify parent
	if onWindowClosed != nil {
		win.SetOnClosed(onWindowClosed)
	}

	// Notify parent that window was created
	if onWindowCreated != nil {
		onWindowCreated(win)
	}

	// Build custom form content with styled labels
	rows := []fyne.CanvasObject{
		tf.makeRowLabel("Name:", tf.nameEntry),
		tf.makeRowLabel("Date/Time:", container.NewBorder(nil, nil, nil, tf.dateTimeButton, tf.dateTimeEntry)),
		tf.makeRowLabel("Location:", tf.placeEntry),
		tf.makeRowLabel("Label:", tf.labelEntry),
		tf.makeRowLabel("Type:", tf.kindSelect),
		tf.makeRowLabel("Priority:", tf.prioritySelect),
		tf.makeRowLabel("Content:", container.NewScroll(tf.contentEntry)),
		tf.makeRowLabel("Reminder:", container.NewVBox(tf.warnTimeSlider, tf.warnTimeLabel)),
	}
	// Add vertical spacing between rows so fields don't stick together
	spacedRows := make([]fyne.CanvasObject, 0, len(rows)*2-1)
	for i, r := range rows {
		if i > 0 {
			spacedRows = append(spacedRows, newVSpacer(5))
		}
		spacedRows = append(spacedRows, r)
	}
	formBox := container.NewVBox(spacedRows...)

	// Buttons
	addBtn := tf.makePrimaryButton(localization.GetString("form_button_add"), func() {
		if err := tf.trySubmit(); err != nil {
			dialog.ShowError(err, tf.formWindow)
			return
		}
		if tf.onSaveCallback != nil {
			tf.onSaveCallback()
		}
		win.Close()
	})
	cancelBtn := tf.makeCancelButton(localization.GetString("form_button_cancel"), func() { win.Close() })

	// Center buttons and set order: Cancel (left), Add (right)
	buttonRow := container.NewCenter(container.NewHBox(cancelBtn, addBtn))
	// Extra space between Reminder text and buttons, and 7px margin at the bottom
	bottomWithSpacer := container.NewVBox(newVSpacer(30), buttonRow, newVSpacer(5))
	// Wrap form in padding and add a bit of top margin so fields are not glued to the window edges
	paddedForm := container.NewPadded(formBox)
	content := container.NewBorder(newVSpacer(10), bottomWithSpacer, nil, nil, paddedForm)
	win.SetContent(content)
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(targetW, targetH))
	win.Show()
}

// ShowEditWindow opens the form in a standalone window for editing an existing todo
func (tf *TodoForm) ShowEditWindow(todo *models.TodoItem, originalTime time.Time, onSave func(), onWindowCreated func(fyne.Window), onWindowClosed func()) {
	tf.isEditMode = true
	tf.originalTodo = todo
	tf.originalTime = originalTime
	tf.onSaveCallback = onSave
	tf.populateForm(todo)

	title := localization.GetString("form_title_edit")

	// Compute window size based on main window (always use parentWindow, not old formWindow)
	parentSize := fyne.NewSize(700, 600)
	if tf.parentWindow != nil && tf.parentWindow.Canvas() != nil {
		parentSize = tf.parentWindow.Canvas().Size()
	}
	targetW := parentSize.Width - 40
	if targetW < 200 {
		targetW = parentSize.Width
	}
	// Match compact height used in ShowCreateWindow
	targetH := parentSize.Height * 0.45

	win := fyne.CurrentApp().NewWindow(title)
	tf.formWindow = win

	// Set close callback to notify parent
	if onWindowClosed != nil {
		win.SetOnClosed(onWindowClosed)
	}

	// Notify parent that window was created
	if onWindowCreated != nil {
		onWindowCreated(win)
	}

	// Build custom form content with styled labels
	rows := []fyne.CanvasObject{
		tf.makeRowLabel("Name:", tf.nameEntry),
		tf.makeRowLabel("Date/Time:", container.NewBorder(nil, nil, nil, tf.dateTimeButton, tf.dateTimeEntry)),
		tf.makeRowLabel("Location:", tf.placeEntry),
		tf.makeRowLabel("Label:", tf.labelEntry),
		tf.makeRowLabel("Type:", tf.kindSelect),
		tf.makeRowLabel("Priority:", tf.prioritySelect),
		tf.makeRowLabel("Content:", container.NewScroll(tf.contentEntry)),
		tf.makeRowLabel("Reminder:", container.NewVBox(tf.warnTimeSlider, tf.warnTimeLabel)),
	}
	// Add vertical spacing between rows so fields don't stick together
	spacedRows := make([]fyne.CanvasObject, 0, len(rows)*2-1)
	for i, r := range rows {
		if i > 0 {
			spacedRows = append(spacedRows, newVSpacer(5))
		}
		spacedRows = append(spacedRows, r)
	}
	formBox := container.NewVBox(spacedRows...)

	// Buttons
	saveBtn := tf.makePrimaryButton(localization.GetString("form_button_save"), func() {
		if err := tf.trySubmit(); err != nil {
			dialog.ShowError(err, tf.formWindow)
			return
		}
		if tf.onSaveCallback != nil {
			tf.onSaveCallback()
		}
		win.Close()
	})
	cancelBtn := tf.makeCancelButton(localization.GetString("form_button_cancel"), func() { win.Close() })

	// Center buttons and set order: Cancel (left), Save (right)
	buttonRow := container.NewCenter(container.NewHBox(cancelBtn, saveBtn))
	// Extra space between Reminder text and buttons, and 7px margin at the bottom
	bottomWithSpacer := container.NewVBox(newVSpacer(30), buttonRow, newVSpacer(7))
	// Wrap form in padding and add a bit of top margin so fields are not glued to the window edges
	paddedForm := container.NewPadded(formBox)
	content := container.NewBorder(newVSpacer(10), bottomWithSpacer, nil, nil, paddedForm)
	win.SetContent(content)
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(targetW, targetH))
	win.Show()
}

// setupForm initializes the form fields
func (tf *TodoForm) setupForm() {
	// Name entry (large text input)
	tf.nameEntry = widget.NewEntry()
	tf.nameEntry.SetPlaceHolder(localization.GetString("field_name_placeholder"))
	tf.nameEntry.TextStyle = fyne.TextStyle{Bold: true}

	// Content entry (multi-line)
	tf.contentEntry = widget.NewMultiLineEntry()
	tf.contentEntry.SetPlaceHolder(localization.GetString("field_content_placeholder"))
	tf.contentEntry.Resize(fyne.NewSize(380, 100))

	// Place entry
	tf.placeEntry = widget.NewEntry()
	tf.placeEntry.SetPlaceHolder(localization.GetString("field_location_placeholder"))

	// Label entry
	tf.labelEntry = widget.NewEntry()
	tf.labelEntry.SetPlaceHolder(localization.GetString("field_label_placeholder"))

	// Date/Time entry and picker
	tf.dateTimeEntry = widget.NewEntry()
	tf.dateTimeEntry.SetPlaceHolder(localization.GetString("field_datetime_placeholder"))
	tf.dateTimeEntry.Disable()                     // Make it read-only, use button for editing
	tf.dateTimeEntry.Resize(fyne.NewSize(320, 35)) // Increase width to fit full date/time

	tf.dateTimeButton = widget.NewButton(localization.GetString("select_datetime"), func() {
		tf.showDateTimePicker()
	})

	// Initialize selected date/time
	tf.selectedDateTime = time.Now()

	// Priority selection
	priorityOptions := []string{
		localization.GetString("priority_0"),
		localization.GetString("priority_1"),
		localization.GetString("priority_2"),
		localization.GetString("priority_3"),
	}
	tf.prioritySelect = widget.NewSelect(priorityOptions, nil)
	tf.prioritySelect.SetSelectedIndex(0)

	// Kind selection (Event/Task)
	kindOptions := []string{localization.GetString("type_event"), localization.GetString("type_task")}
	tf.kindSelect = widget.NewSelect(kindOptions, nil)
	tf.kindSelect.SetSelectedIndex(0)

	// Warning time slider (0-864 minutes = 0-14.4 hours)
	tf.warnTimeSlider = NewReminderSlider(0, 864)
	tf.warnTimeSlider.Step = 5
	tf.warnTimeSlider.Value = 0
	tf.warnTimeSlider.OnChanged = tf.onWarnTimeChanged

	// Reminder label: use custom-colored canvas text so it stays visible
	// on the dark "light" theme background.
	tf.warnTimeLabel = canvas.NewText(localization.GetString("reminder_none"), tf.reminderLabelColor())
	tf.warnTimeLabel.Alignment = fyne.TextAlignCenter
}

// Note: createFormContent is no longer needed as we use dialog.NewForm directly

// resetForm clears all form fields for new todo creation
func (tf *TodoForm) resetForm() {
	tf.nameEntry.SetText("")
	tf.contentEntry.SetText("")
	tf.placeEntry.SetText("")
	tf.labelEntry.SetText("")

	// Set current date/time in DD.MM.YYYY HH:MM format
	now := time.Now()
	tf.selectedDateTime = now
	currentDateTime := now.Format("02.01.2006 15:04")
	tf.dateTimeEntry.SetText(currentDateTime)

	tf.prioritySelect.SetSelectedIndex(0)
	tf.kindSelect.SetSelectedIndex(0)
	tf.warnTimeSlider.SetValue(0)
	tf.onWarnTimeChanged(0)
}

// populateForm fills form fields with existing todo data
func (tf *TodoForm) populateForm(todo *models.TodoItem) {
	tf.nameEntry.SetText(todo.Name)
	tf.contentEntry.SetText(todo.Content)
	tf.placeEntry.SetText(todo.Place)
	tf.labelEntry.SetText(todo.Label)

	// Format date/time for display in DD.MM.YYYY HH:MM format
	tf.selectedDateTime = todo.TodoTime
	dateTimeStr := todo.TodoTime.Format("02.01.2006 15:04")
	tf.dateTimeEntry.SetText(dateTimeStr)

	tf.prioritySelect.SetSelectedIndex(todo.Level)
	tf.kindSelect.SetSelectedIndex(todo.Kind)
	tf.warnTimeSlider.SetValue(float64(todo.WarnTime))
	tf.onWarnTimeChanged(float64(todo.WarnTime))
}

// onWarnTimeChanged updates the warning time label
func (tf *TodoForm) onWarnTimeChanged(value float64) {
	warnTime := int(value)
	if warnTime == 0 {
		tf.warnTimeLabel.Text = localization.GetString("reminder_none")
		tf.warnTimeLabel.Color = tf.reminderLabelColor()
		tf.warnTimeLabel.Refresh()
		return
	}

	minutes := warnTime
	hours := minutes / 60
	days := hours / 24

	var parts []string
	if days > 0 {
		dayStr := localization.GetString("time_day")
		if days > 1 {
			dayStr = localization.GetString("time_days")
		}
		parts = append(parts, fmt.Sprintf("%d %s", days, dayStr))
		minutes %= 60
	}
	if hours > 0 {
		hourStr := localization.GetString("time_hour")
		if hours%24 > 1 {
			hourStr = localization.GetString("time_hours")
		}
		parts = append(parts, fmt.Sprintf("%d %s", hours%24, hourStr))
		minutes %= 60
	}
	if minutes > 0 {
		minuteStr := localization.GetString("time_minute")
		if minutes > 1 {
			minuteStr = localization.GetString("time_minutes")
		}
		parts = append(parts, fmt.Sprintf("%d %s", minutes, minuteStr))
	}

	if len(parts) == 0 {
		tf.warnTimeLabel.Text = localization.GetString("reminder_none")
	} else {
		reminderFormat := localization.GetString("reminder_format")
		tf.warnTimeLabel.Text = fmt.Sprintf(reminderFormat, joinStrings(parts, " "))
	}
	tf.warnTimeLabel.Color = tf.reminderLabelColor()
	tf.warnTimeLabel.Refresh()
}

// onSubmit handles form submission
func (tf *TodoForm) onSubmit() {
	if err := tf.trySubmit(); err != nil {
		dialog.ShowError(err, tf.parentWindow)
		return
	}
	if tf.onSaveCallback != nil {
		tf.onSaveCallback()
	}
}

// trySubmit validates and saves the todo, returning error on failure
func (tf *TodoForm) trySubmit() error {
	if tf.nameEntry.Text == "" {
		return errors.New(localization.GetString("error_name_required"))
	}

	// Use the selected date/time
	todoTime := tf.selectedDateTime
	if todoTime.IsZero() {
		todoTime = time.Now()
	}

	// Create todo item
	todo := models.NewTodoItem()
	todo.Name = tf.nameEntry.Text
	todo.Content = tf.contentEntry.Text
	todo.Place = tf.placeEntry.Text
	todo.Label = tf.labelEntry.Text
	todo.Kind = tf.kindSelect.SelectedIndex()
	todo.Level = tf.prioritySelect.SelectedIndex()
	todo.TodoTime = todoTime
	todo.WarnTime = int(tf.warnTimeSlider.Value)

	// Save todo
	var err error
	if tf.isEditMode {
		err = tf.dataManager.UpdateTodo(todo, tf.originalTime)
	} else {
		err = tf.dataManager.AddTodo(todo)
	}
	if err != nil {
		return err
	}
	return nil
}

// makeRowLabel creates a two-column row with a styled label and a widget
func (tf *TodoForm) makeRowLabel(label string, w fyne.CanvasObject) fyne.CanvasObject {
	lbl := tf.makeStyledLabel(label)
	// Put label on the left and let the widget take all remaining width,
	// so entries are noticeably wider than labels.
	return container.NewBorder(nil, nil, lbl, nil, w)
}

// makeStyledLabel builds a label that is bold and white in light themes
func (tf *TodoForm) makeStyledLabel(text string) fyne.CanvasObject {
	col := theme.Color(theme.ColorNameForeground)
	n := color.NRGBAModel.Convert(col).(color.NRGBA)
	if tf.isLightTheme() {
		n = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	}
	t := canvas.NewText(text, n)
	if tf.isLightTheme() {
		t.TextStyle = fyne.TextStyle{Bold: true}
	}
	// Fix label width so that all input fields start at the same X position.
	labelWidth := float32(100)
	labelHeight := t.TextSize + 8
	return container.NewGridWrap(fyne.NewSize(labelWidth, labelHeight), container.NewCenter(t))
}

// isLightTheme returns true for the custom "light" theme and any truly light background,
// so that labels can switch to bold white for better contrast.
func (tf *TodoForm) isLightTheme() bool {
	bg := theme.Color(theme.ColorNameBackground)
	n := color.NRGBAModel.Convert(bg).(color.NRGBA)

	// Our custom LightSoftTheme uses a dark gray background (#3c3c3c).
	// Detect it explicitly so labels become white on that background.
	if n.R == 0x3c && n.G == 0x3c && n.B == 0x3c {
		return true
	}

	// relative luminance
	l := 0.2126*float64(n.R)/255.0 + 0.7152*float64(n.G)/255.0 + 0.0722*float64(n.B)/255.0
	return l > 0.5
}

// reminderLabelColor returns a color for the reminder text
// that is always visible on the current theme background.
func (tf *TodoForm) reminderLabelColor() color.Color {
	col := theme.Color(theme.ColorNameForeground)
	n := color.NRGBAModel.Convert(col).(color.NRGBA)
	if tf.isLightTheme() {
		// Force white text on dark "light" background.
		return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	}
	return n
}

// ReminderSlider is a minimal custom slider for the reminder time,
// with a flat line and a clearly visible knob that adapts to the theme.
type ReminderSlider struct {
	widget.BaseWidget
	Min, Max  float64
	Value     float64
	Step      float64
	OnChanged func(float64)
}

// NewReminderSlider creates a new ReminderSlider with the given range.
func NewReminderSlider(min, max float64) *ReminderSlider {
	s := &ReminderSlider{
		Min:   min,
		Max:   max,
		Value: min,
		Step:  0,
	}
	s.ExtendBaseWidget(s)
	return s
}

// SetValue sets the slider value, clamped to range and snapped to Step.
func (s *ReminderSlider) SetValue(v float64) {
	if s.Max <= s.Min {
		return
	}
	if v < s.Min {
		v = s.Min
	}
	if v > s.Max {
		v = s.Max
	}
	if s.Step > 0 {
		v = math.Round(v/s.Step) * s.Step
	}
	if v == s.Value {
		return
	}
	s.Value = v
	if s.OnChanged != nil {
		s.OnChanged(v)
	}
	s.Refresh()
}

// Tapped moves the knob to the tap position.
func (s *ReminderSlider) Tapped(ev *fyne.PointEvent) {
	s.updateFromPos(ev.Position, s.Size())
}

// Dragged updates the value during drag.
func (s *ReminderSlider) Dragged(ev *fyne.DragEvent) {
	s.updateFromPos(ev.Position, s.Size())
}

func (s *ReminderSlider) DragEnd() {}

// updateFromPos converts a point inside the widget to a value.
func (s *ReminderSlider) updateFromPos(pos fyne.Position, size fyne.Size) {
	if s.Max <= s.Min {
		return
	}
	padding := float32(12)
	width := size.Width - 2*padding
	if width <= 0 {
		return
	}
	x := pos.X
	if x < padding {
		x = padding
	}
	if x > size.Width-padding {
		x = size.Width - padding
	}
	ratio := float64((x - padding) / width)
	v := s.Min + ratio*(s.Max-s.Min)
	s.SetValue(v)
}

// CreateRenderer implements fyne.Widget.
func (s *ReminderSlider) CreateRenderer() fyne.WidgetRenderer {
	track := canvas.NewRectangle(theme.Color(theme.ColorNameSeparator))
	knob := canvas.NewCircle(theme.Color(theme.ColorNamePrimary))
	objs := []fyne.CanvasObject{track, knob}
	return &reminderSliderRenderer{
		slider: s,
		track:  track,
		knob:   knob,
		objs:   objs,
	}
}

type reminderSliderRenderer struct {
	slider *ReminderSlider
	track  *canvas.Rectangle
	knob   *canvas.Circle
	objs   []fyne.CanvasObject
}

func (r *reminderSliderRenderer) Layout(size fyne.Size) {
	padding := float32(12)
	centerY := size.Height / 2
	trackHeight := float32(4)

	// Track spans the full width used by other inputs.
	trackWidth := size.Width - 2*padding
	if trackWidth < 0 {
		trackWidth = 0
	}
	r.track.Resize(fyne.NewSize(trackWidth, trackHeight))
	r.track.Move(fyne.NewPos(padding, centerY-trackHeight/2))

	// Knob position based on value.
	if r.slider.Max > r.slider.Min {
		ratio := float32((r.slider.Value - r.slider.Min) / (r.slider.Max - r.slider.Min))
		if ratio < 0 {
			ratio = 0
		}
		if ratio > 1 {
			ratio = 1
		}
		x := padding + ratio*trackWidth
		knobRadius := float32(8)
		r.knob.Resize(fyne.NewSize(knobRadius*2, knobRadius*2))
		r.knob.Move(fyne.NewPos(x-knobRadius, centerY-knobRadius))
	}
}

func (r *reminderSliderRenderer) MinSize() fyne.Size {
	// Height chosen to give comfortable click/drag target.
	return fyne.NewSize(100, 32)
}

func (r *reminderSliderRenderer) Refresh() {
	// Update colors according to theme.
	r.track.FillColor = theme.Color(theme.ColorNameSeparator)
	r.track.Refresh()
	r.knob.FillColor = theme.Color(theme.ColorNamePrimary)
	r.knob.Refresh()

	// Re-layout in case size or value changed.
	r.Layout(r.slider.Size())
}

func (r *reminderSliderRenderer) BackgroundColor() fyne.ThemeColorName { return "" }
func (r *reminderSliderRenderer) Objects() []fyne.CanvasObject         { return r.objs }
func (r *reminderSliderRenderer) Destroy()                             {}

// newVSpacer creates a transparent vertical spacer with fixed height.
func newVSpacer(height float32) fyne.CanvasObject {
	r := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	r.SetMinSize(fyne.NewSize(1, height))
	return r
}

// makeCancelButton returns a rectangular button 20% lighter than the current background
func (tf *TodoForm) makeCancelButton(text string, onTap func()) fyne.CanvasObject {
	bg := theme.Color(theme.ColorNameBackground)
	n := color.NRGBAModel.Convert(bg).(color.NRGBA)
	light := helpers.Lighten(n, 0.20)
	fg := theme.Color(theme.ColorNameForeground)
	fn := color.NRGBAModel.Convert(fg).(color.NRGBA)
	btn := &rectButton{
		Text:     text,
		Bg:       light,
		Fg:       fn,
		SizeHint: fyne.NewSize(100, 44),
		OnTapped: onTap,
	}
	btn.ExtendBaseWidget(btn)
	return btn
}

// makePrimaryButton returns a rectangular button using the primary accent color.
func (tf *TodoForm) makePrimaryButton(text string, onTap func()) fyne.CanvasObject {
	bg := theme.Color(theme.ColorNamePrimary)
	n := color.NRGBAModel.Convert(bg).(color.NRGBA)
	fg := theme.Color(theme.ColorNameForeground)
	fn := color.NRGBAModel.Convert(fg).(color.NRGBA)
	btn := &rectButton{
		Text:     text,
		Bg:       n,
		Fg:       fn,
		SizeHint: fyne.NewSize(100, 44),
		OnTapped: onTap,
	}
	btn.ExtendBaseWidget(btn)
	return btn
}

// rectButton is a minimal custom button used for Cancel styling
type rectButton struct {
	widget.BaseWidget
	Text     string
	Bg       color.NRGBA
	Fg       color.NRGBA
	SizeHint fyne.Size
	OnTapped func()
	hovered  bool
}

func (b *rectButton) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(b.Bg)
	bg.CornerRadius = 8
	txt := canvas.NewText(b.Text, b.Fg)
	txt.TextStyle = fyne.TextStyle{Bold: true}
	txt.Alignment = fyne.TextAlignCenter
	cont := container.NewMax(bg, container.NewCenter(txt))
	return &rectButtonRenderer{btn: b, bg: bg, txt: txt, cont: cont}
}

type rectButtonRenderer struct {
	btn  *rectButton
	bg   *canvas.Rectangle
	txt  *canvas.Text
	cont *fyne.Container
}

func (r *rectButtonRenderer) Layout(size fyne.Size)                { r.cont.Resize(size) }
func (r *rectButtonRenderer) MinSize() fyne.Size                   { return r.btn.MinSize() }
func (r *rectButtonRenderer) BackgroundColor() fyne.ThemeColorName { return "" }
func (r *rectButtonRenderer) Objects() []fyne.CanvasObject         { return []fyne.CanvasObject{r.cont} }
func (r *rectButtonRenderer) Destroy()                             {}
func (r *rectButtonRenderer) Refresh() {
	// slight hover lightening
	bg := r.btn.Bg
	if r.btn.hovered {
		bg = color.NRGBA{R: uint8(float32(bg.R)*0.92 + 255*0.08), G: uint8(float32(bg.G)*0.92 + 255*0.08), B: uint8(float32(bg.B)*0.92 + 255*0.08), A: bg.A}
	}
	r.bg.FillColor = bg
	r.bg.Refresh()
	r.txt.Text = r.btn.Text
	r.txt.Color = r.btn.Fg
	r.txt.Refresh()
}

func (b *rectButton) MinSize() fyne.Size {
	if b.SizeHint.Width > 0 && b.SizeHint.Height > 0 {
		return b.SizeHint
	}
	return fyne.NewSize(100, 44)
}

func (b *rectButton) Tapped(*fyne.PointEvent) {
	if b.OnTapped != nil {
		b.OnTapped()
	}
}

func (b *rectButton) MouseIn(*desktop.MouseEvent)    { b.hovered = true; b.Refresh() }
func (b *rectButton) MouseMoved(*desktop.MouseEvent) {}
func (b *rectButton) MouseOut()                      { b.hovered = false; b.Refresh() }

// onCancel handles form cancellation
func (tf *TodoForm) onCancel() {
	// Dialog will be closed automatically by the form buttons
}

// showDateTimePicker displays a date and time picker dialog
func (tf *TodoForm) showDateTimePicker() {
	// Create a combined date/time picker dialog
	dateEntry := widget.NewEntry()
	dateEntry.SetText(tf.selectedDateTime.Format("02.01.2006"))
	dateEntry.SetPlaceHolder("DD.MM.YYYY")
	dateEntry.Resize(fyne.NewSize(300, 45))

	timeEntry := widget.NewEntry()
	timeEntry.SetText(tf.selectedDateTime.Format("15:04"))
	timeEntry.SetPlaceHolder("HH:MM")
	timeEntry.Resize(fyne.NewSize(300, 45))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Date (DD.MM.YYYY):", Widget: dateEntry},
			{Text: "Time (HH:MM):", Widget: timeEntry},
		},
	}

	// Create container for the form
	formContainer := container.NewVBox(form)

	// Remove bottom buttons by using custom dialog without buttons
	// Use formWindow if available (standalone mode), otherwise use parentWindow (dialog mode)
	dialogParent := tf.formWindow
	if dialogParent == nil {
		dialogParent = tf.parentWindow
	}
	dateTimeDialog := dialog.NewCustomWithoutButtons("Select Date and Time", formContainer, dialogParent)
	// Make the dialog wider and compact in height to avoid extra space
	dateTimeDialog.Resize(fyne.NewSize(700, form.MinSize().Height+40))

	// Handle date/time selection
	form.OnSubmit = func() {
		dateStr := dateEntry.Text
		timeStr := timeEntry.Text

		// Parse date
		if date, err := time.Parse("02.01.2006", dateStr); err == nil {
			// Parse time
			if parsedTime, err := time.Parse("15:04", timeStr); err == nil {
				// Combine date and time
				year, month, day := date.Date()
				hour, min := parsedTime.Hour(), parsedTime.Minute()
				location := tf.selectedDateTime.Location()
				tf.selectedDateTime = time.Date(year, month, day, hour, min, 0, 0, location)

				tf.updateDateTimeDisplay()
			}
		}
		dateTimeDialog.Hide()
	}

	form.OnCancel = func() {
		dateTimeDialog.Hide()
	}

	// Show dialog
	dateTimeDialog.Show()
}

// updateDateTimeDisplay updates the date/time entry display
func (tf *TodoForm) updateDateTimeDisplay() {
	displayText := tf.selectedDateTime.Format("02.01.2006 15:04")
	tf.dateTimeEntry.SetText(displayText)
}

// Helper function to join strings
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}
	if len(strings) == 1 {
		return strings[0]
	}

	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += " " + strings[i]
	}
	return result
}
