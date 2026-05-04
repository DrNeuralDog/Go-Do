package forms

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"godo/src/localization"
	"godo/src/models"
	"godo/src/persistence"
	"godo/src/ui/helpers"
	"godo/src/ui/widgets"
)

const (
	formDialogWidth             = 700
	formDialogHeight            = 600
	formStandaloneMinWidth      = 200
	formStandaloneWidthPadding  = 40
	formStandaloneHeightRatio   = 0.45
	formContentEntryWidth       = 380
	formContentEntryHeight      = 100
	formDateTimeEntryWidth      = 320
	formDateTimeEntryHeight     = 35
	formDateEntryWidth          = 300
	formDateEntryHeight         = 45
	formDateTimeDialogWidth     = 700
	formDateTimeDialogPadding   = 40
	formRowSpacerHeight         = 5
	formTopSpacerHeight         = 10
	formButtonTopSpacerHeight   = 30
	formCreateBottomSpacer      = 5
	formEditBottomSpacer        = 7
	formLabelWidth              = 100
	formLabelHeightPadding      = 8
	formButtonWidth             = 100
	formButtonHeight            = 44
	formButtonRadius            = 8
	formCancelButtonLightAmount = 0.20
	reminderSliderMin           = 0
	reminderSliderMax           = 864
	reminderSliderStep          = 5
	reminderSliderPadding       = 12
	reminderSliderTrackHeight   = 4
	reminderSliderKnobRadius    = 8
	reminderSliderMinWidth      = 100
	reminderSliderMinHeight     = 32
	softLightBackgroundValue    = 0x3c
	whiteColorChannel           = 0xff
)

const (
	dateTimeDisplayLayout = "02.01.2006 15:04"
	dateInputLayout       = "02.01.2006"
	timeInputLayout       = "15:04"
)

// TodoForm manages create and edit todo forms
type TodoForm struct {
	parentWindow fyne.Window
	formWindow   fyne.Window
	dataManager  persistence.TodoRepository

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

	selectedDateTime time.Time
	isEditMode       bool
	originalTodo     *models.TodoItem
	originalTime     time.Time
	onSaveCallback   func()
}

// NewTodoForm builds todo form controller
func NewTodoForm(window fyne.Window, dataManager persistence.TodoRepository) *TodoForm {
	tf := &TodoForm{
		parentWindow: window,
		dataManager:  dataManager,
	}

	tf.setupForm()

	return tf
}

// ShowCreateDialog shows create form as dialog
func (tf *TodoForm) ShowCreateDialog(onSave func()) {
	tf.beginCreate(onSave)

	formDialog := tf.newFormDialog(
		localization.GetString("form_title_add"),
		localization.GetString("form_button_add"),
	)
	formDialog.Show()
}

// ShowEditDialog shows edit form as dialog
func (tf *TodoForm) ShowEditDialog(todo *models.TodoItem, originalTime time.Time, onSave func()) {
	if !tf.beginEdit(todo, originalTime, onSave) {
		return
	}

	formDialog := tf.newFormDialog(
		localization.GetString("form_title_edit"),
		localization.GetString("form_button_save"),
	)
	formDialog.Show()
}

// ShowCreateWindow opens standalone create form
func (tf *TodoForm) ShowCreateWindow(onSave func(), onWindowCreated func(fyne.Window), onWindowClosed func()) {
	tf.beginCreate(onSave)

	tf.showStandaloneWindow(
		localization.GetString("form_title_add"),
		localization.GetString("form_button_add"),
		formCreateBottomSpacer,
		onWindowCreated,
		onWindowClosed,
	)
}

// ShowEditWindow opens standalone edit form
func (tf *TodoForm) ShowEditWindow(todo *models.TodoItem, originalTime time.Time, onSave func(), onWindowCreated func(fyne.Window), onWindowClosed func()) {
	if !tf.beginEdit(todo, originalTime, onSave) {
		return
	}

	tf.showStandaloneWindow(
		localization.GetString("form_title_edit"),
		localization.GetString("form_button_save"),
		formEditBottomSpacer,
		onWindowCreated,
		onWindowClosed,
	)
}

// beginCreate prepares form state for new todo
func (tf *TodoForm) beginCreate(onSave func()) {
	tf.isEditMode = false
	tf.originalTodo = nil
	tf.originalTime = time.Time{}
	tf.onSaveCallback = onSave
	tf.resetForm()
}

// beginEdit prepares form state for existing todo
func (tf *TodoForm) beginEdit(todo *models.TodoItem, originalTime time.Time, onSave func()) bool {
	if todo == nil {
		return false
	}

	tf.isEditMode = true
	tf.originalTodo = todo
	tf.originalTime = originalTime
	tf.onSaveCallback = onSave
	tf.populateForm(todo)

	return true
}

// newFormDialog creates classic Fyne form dialog
func (tf *TodoForm) newFormDialog(title, submitText string) dialog.Dialog {
	formDialog := dialog.NewForm(
		title,
		submitText,
		localization.GetString("form_button_cancel"),
		tf.dialogFormItems(),
		func(submitted bool) {
			if submitted {
				tf.onSubmit()
			}
		},
		tf.parentWindow,
	)
	formDialog.Resize(fyne.NewSize(formDialogWidth, formDialogHeight))

	return formDialog
}

// dialogFormItems builds dialog field list
func (tf *TodoForm) dialogFormItems() []*widget.FormItem {
	return []*widget.FormItem{
		{Text: "Name:", Widget: tf.nameEntry},
		{Text: "Date/Time:", Widget: container.NewBorder(nil, nil, nil, tf.dateTimeButton, tf.dateTimeEntry)},
		{Text: "Location:", Widget: tf.placeEntry},
		{Text: "Label:", Widget: tf.labelEntry},
		{Text: "Type:", Widget: tf.kindSelect},
		{Text: "Priority:", Widget: tf.prioritySelect},
		{Text: "Reminder:", Widget: container.NewVBox(tf.warnTimeSlider, tf.warnTimeLabel)},
		{Text: "Content:", Widget: container.NewScroll(tf.contentEntry)},
	}
}

// showStandaloneWindow creates separate todo form window
func (tf *TodoForm) showStandaloneWindow(title, submitText string, bottomSpacerHeight float32, onWindowCreated func(fyne.Window), onWindowClosed func()) {
	app := fyne.CurrentApp()
	if app == nil {
		return
	}

	win := app.NewWindow(title)
	tf.formWindow = win

	if onWindowClosed != nil {
		win.SetOnClosed(onWindowClosed)
	}

	if onWindowCreated != nil {
		onWindowCreated(win)
	}

	submitBtn := tf.makePrimaryButton(submitText, func() {
		tf.submitStandalone(win)
	})
	cancelBtn := tf.makeCancelButton(localization.GetString("form_button_cancel"), func() {
		win.Close()
	})

	win.SetContent(tf.standaloneContent(cancelBtn, submitBtn, bottomSpacerHeight))
	win.SetFixedSize(true)
	win.Resize(tf.standaloneWindowSize())
	win.Show()
}

// standaloneWindowSize returns compact form window size
func (tf *TodoForm) standaloneWindowSize() fyne.Size {
	parentSize := fyne.NewSize(formDialogWidth, formDialogHeight)

	if tf.parentWindow != nil && tf.parentWindow.Canvas() != nil {
		parentSize = tf.parentWindow.Canvas().Size()
	}

	targetW := parentSize.Width - formStandaloneWidthPadding
	if targetW < formStandaloneMinWidth {
		targetW = parentSize.Width
	}

	return fyne.NewSize(targetW, parentSize.Height*formStandaloneHeightRatio)
}

// standaloneContent builds custom form layout
func (tf *TodoForm) standaloneContent(cancelBtn, submitBtn fyne.CanvasObject, bottomSpacerHeight float32) fyne.CanvasObject {
	formBox := container.NewVBox(spacedRows(tf.formRows())...)
	buttonRow := container.NewCenter(container.NewHBox(cancelBtn, submitBtn))
	bottom := container.NewVBox(newVSpacer(formButtonTopSpacerHeight), buttonRow, newVSpacer(bottomSpacerHeight))
	paddedForm := container.NewPadded(formBox)

	return container.NewBorder(newVSpacer(formTopSpacerHeight), bottom, nil, nil, paddedForm)
}

// formRows returns standalone form rows
func (tf *TodoForm) formRows() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		tf.makeRowLabel("Name:", tf.nameEntry),
		tf.makeRowLabel("Date/Time:", container.NewBorder(nil, nil, nil, tf.dateTimeButton, tf.dateTimeEntry)),
		tf.makeRowLabel("Location:", tf.placeEntry),
		tf.makeRowLabel("Label:", tf.labelEntry),
		tf.makeRowLabel("Type:", tf.kindSelect),
		tf.makeRowLabel("Priority:", tf.prioritySelect),
		tf.makeRowLabel("Content:", container.NewScroll(tf.contentEntry)),
		tf.makeRowLabel("Reminder:", container.NewVBox(tf.warnTimeSlider, tf.warnTimeLabel)),
	}
}

// spacedRows inserts small gaps between form rows
func spacedRows(rows []fyne.CanvasObject) []fyne.CanvasObject {
	if len(rows) == 0 {
		return nil
	}

	spaced := make([]fyne.CanvasObject, 0, len(rows)*2-1)
	for i, row := range rows {
		if i > 0 {
			spaced = append(spaced, newVSpacer(formRowSpacerHeight))
		}

		spaced = append(spaced, row)
	}

	return spaced
}

// setupForm initializes form widgets
func (tf *TodoForm) setupForm() {
	tf.nameEntry = widget.NewEntry()
	tf.nameEntry.SetPlaceHolder(localization.GetString("field_name_placeholder"))
	tf.nameEntry.TextStyle = fyne.TextStyle{Bold: true}

	tf.contentEntry = widget.NewMultiLineEntry()
	tf.contentEntry.SetPlaceHolder(localization.GetString("field_content_placeholder"))
	tf.contentEntry.Resize(fyne.NewSize(formContentEntryWidth, formContentEntryHeight))

	tf.placeEntry = widget.NewEntry()
	tf.placeEntry.SetPlaceHolder(localization.GetString("field_location_placeholder"))

	tf.labelEntry = widget.NewEntry()
	tf.labelEntry.SetPlaceHolder(localization.GetString("field_label_placeholder"))

	tf.dateTimeEntry = widget.NewEntry()
	tf.dateTimeEntry.SetPlaceHolder(localization.GetString("field_datetime_placeholder"))
	tf.dateTimeEntry.Disable()
	tf.dateTimeEntry.Resize(fyne.NewSize(formDateTimeEntryWidth, formDateTimeEntryHeight))

	tf.dateTimeButton = widget.NewButton(localization.GetString("select_datetime"), func() {
		tf.showDateTimePicker()
	})

	tf.selectedDateTime = time.Now()

	priorityOptions := []string{
		localization.GetString("priority_0"),
		localization.GetString("priority_1"),
		localization.GetString("priority_2"),
		localization.GetString("priority_3"),
	}
	tf.prioritySelect = widget.NewSelect(priorityOptions, nil)
	tf.prioritySelect.SetSelectedIndex(0)

	kindOptions := []string{localization.GetString("type_event"), localization.GetString("type_task")}
	tf.kindSelect = widget.NewSelect(kindOptions, nil)
	tf.kindSelect.SetSelectedIndex(0)

	// Старый лимит: до 14.4 часа до задачи
	tf.warnTimeSlider = NewReminderSlider(reminderSliderMin, reminderSliderMax)
	tf.warnTimeSlider.Step = reminderSliderStep
	tf.warnTimeSlider.Value = reminderSliderMin
	tf.warnTimeSlider.OnChanged = tf.onWarnTimeChanged

	tf.warnTimeLabel = canvas.NewText(localization.GetString("reminder_none"), tf.reminderLabelColor())
	tf.warnTimeLabel.Alignment = fyne.TextAlignCenter
}

// resetForm clears form fields for new todo
func (tf *TodoForm) resetForm() {
	tf.nameEntry.SetText("")
	tf.contentEntry.SetText("")
	tf.placeEntry.SetText("")
	tf.labelEntry.SetText("")

	tf.selectedDateTime = time.Now()
	tf.updateDateTimeDisplay()

	tf.prioritySelect.SetSelectedIndex(0)
	tf.kindSelect.SetSelectedIndex(0)
	tf.warnTimeSlider.SetValue(0)
	tf.onWarnTimeChanged(0)
}

// populateForm fills fields from existing todo
func (tf *TodoForm) populateForm(todo *models.TodoItem) {
	tf.nameEntry.SetText(todo.Name)
	tf.contentEntry.SetText(todo.Content)
	tf.placeEntry.SetText(todo.Place)
	tf.labelEntry.SetText(todo.Label)

	tf.selectedDateTime = todo.TodoTime
	tf.updateDateTimeDisplay()

	setSelectIndex(tf.prioritySelect, todo.Level, 0)
	setSelectIndex(tf.kindSelect, todo.Kind, 0)
	tf.warnTimeSlider.SetValue(float64(todo.WarnTime))
	tf.onWarnTimeChanged(float64(todo.WarnTime))
}

// onWarnTimeChanged updates reminder label
func (tf *TodoForm) onWarnTimeChanged(value float64) {
	warnTime := int(value)
	tf.warnTimeLabel.Text = reminderText(warnTime)
	tf.warnTimeLabel.Color = tf.reminderLabelColor()
	tf.warnTimeLabel.Refresh()
}

// reminderText returns localized reminder text
func reminderText(warnTime int) string {
	if warnTime <= 0 {
		return localization.GetString("reminder_none")
	}

	days := warnTime / (24 * 60)
	hours := (warnTime % (24 * 60)) / 60
	minutes := warnTime % 60

	parts := make([]string, 0, 3)
	parts = appendTimePart(parts, days, "time_day", "time_days")
	parts = appendTimePart(parts, hours, "time_hour", "time_hours")
	parts = appendTimePart(parts, minutes, "time_minute", "time_minutes")

	if len(parts) == 0 {
		return localization.GetString("reminder_none")
	}

	return fmt.Sprintf(localization.GetString("reminder_format"), strings.Join(parts, " "))
}

// appendTimePart appends non-zero localized time part
func appendTimePart(parts []string, value int, oneKey, manyKey string) []string {
	if value <= 0 {
		return parts
	}

	label := localization.GetString(oneKey)
	if value > 1 {
		label = localization.GetString(manyKey)
	}

	return append(parts, fmt.Sprintf("%d %s", value, label))
}

// onSubmit handles dialog form submit
func (tf *TodoForm) onSubmit() {
	if err := tf.trySubmit(); err != nil {
		dialog.ShowError(err, tf.parentWindow)
		return
	}

	tf.callOnSave()
}

// submitStandalone saves todo and closes window
func (tf *TodoForm) submitStandalone(win fyne.Window) {
	if err := tf.trySubmit(); err != nil {
		dialog.ShowError(err, win)
		return
	}

	tf.callOnSave()
	win.Close()
}

// callOnSave runs save callback
func (tf *TodoForm) callOnSave() {
	if tf.onSaveCallback != nil {
		tf.onSaveCallback()
	}
}

// trySubmit validates and saves todo
func (tf *TodoForm) trySubmit() error {
	if tf.nameEntry.Text == "" {
		return errors.New(localization.GetString("error_name_required"))
	}

	todoTime := tf.selectedDateTime
	if todoTime.IsZero() {
		todoTime = time.Now()
	}

	todo := tf.todoFromForm(todoTime)

	if tf.isEditMode {
		return tf.dataManager.UpdateTodo(todo, tf.originalTime)
	}

	return tf.dataManager.AddTodo(todo)
}

// todoFromForm builds todo from current fields
func (tf *TodoForm) todoFromForm(todoTime time.Time) *models.TodoItem {
	todo := models.NewTodoItem()

	if tf.isEditMode && tf.originalTodo != nil {
		// Редактирование не должно сбрасывать done/star/order
		copyTodo := *tf.originalTodo
		todo = &copyTodo
	}

	todo.Name = tf.nameEntry.Text
	todo.Content = tf.contentEntry.Text
	todo.Place = tf.placeEntry.Text
	todo.Label = tf.labelEntry.Text
	todo.Kind = selectedIndexOrDefault(tf.kindSelect, 0)
	todo.Level = selectedIndexOrDefault(tf.prioritySelect, 0)
	todo.TodoTime = todoTime
	todo.WarnTime = int(tf.warnTimeSlider.Value)

	return todo
}

// makeRowLabel creates fixed label and stretchable field row
func (tf *TodoForm) makeRowLabel(label string, w fyne.CanvasObject) fyne.CanvasObject {
	lbl := tf.makeStyledLabel(label)

	return container.NewBorder(nil, nil, lbl, nil, w)
}

// makeStyledLabel builds form row label
func (tf *TodoForm) makeStyledLabel(text string) fyne.CanvasObject {
	n := helpers.ToNRGBA(theme.Color(theme.ColorNameForeground))
	if tf.useLightFormText() {
		n = color.NRGBA{R: whiteColorChannel, G: whiteColorChannel, B: whiteColorChannel, A: whiteColorChannel}
	}

	t := canvas.NewText(text, n)
	if tf.useLightFormText() {
		t.TextStyle = fyne.TextStyle{Bold: true}
	}

	labelHeight := t.TextSize + formLabelHeightPadding

	return container.NewGridWrap(fyne.NewSize(formLabelWidth, labelHeight), container.NewCenter(t))
}

// useLightFormText reports whether labels need white color
func (tf *TodoForm) useLightFormText() bool {
	bg := helpers.ToNRGBA(theme.Color(theme.ColorNameBackground))

	if bg.R == softLightBackgroundValue && bg.G == softLightBackgroundValue && bg.B == softLightBackgroundValue {
		// LightSoftTheme на деле темный фон, так что текст нужен белый
		return true
	}

	return false
}

// reminderLabelColor returns visible reminder label color
func (tf *TodoForm) reminderLabelColor() color.Color {
	if tf.useLightFormText() {
		return color.NRGBA{R: whiteColorChannel, G: whiteColorChannel, B: whiteColorChannel, A: whiteColorChannel}
	}

	return helpers.ToNRGBA(theme.Color(theme.ColorNameForeground))
}

// ReminderSlider is a flat slider for reminder time
type ReminderSlider struct {
	widget.BaseWidget
	Min, Max  float64
	Value     float64
	Step      float64
	OnChanged func(float64)
}

// NewReminderSlider builds reminder slider
func NewReminderSlider(min, max float64) *ReminderSlider {
	s := &ReminderSlider{
		Min:   min,
		Max:   max,
		Value: min,
	}

	s.ExtendBaseWidget(s)

	return s
}

// SetValue applies clamped and snapped slider value
func (s *ReminderSlider) SetValue(v float64) {
	if s.Max <= s.Min {
		return
	}

	v = s.clampedValue(v)
	if s.Value == v {
		return
	}

	s.Value = v

	if s.OnChanged != nil {
		s.OnChanged(v)
	}

	s.Refresh()
}

// Tapped moves slider knob to tap position
func (s *ReminderSlider) Tapped(ev *fyne.PointEvent) {
	s.updateFromPos(ev.Position, s.Size())
}

// Dragged updates slider value during drag
func (s *ReminderSlider) Dragged(ev *fyne.DragEvent) {
	s.updateFromPos(ev.Position, s.Size())
}

// DragEnd completes slider drag
func (s *ReminderSlider) DragEnd() {}

// updateFromPos converts widget position to slider value
func (s *ReminderSlider) updateFromPos(pos fyne.Position, size fyne.Size) {
	if s.Max <= s.Min {
		return
	}

	width := size.Width - 2*reminderSliderPadding
	if width <= 0 {
		return
	}

	x := clampFloat32(pos.X, reminderSliderPadding, size.Width-reminderSliderPadding)
	ratio := float64((x - reminderSliderPadding) / width)

	s.SetValue(s.Min + ratio*(s.Max-s.Min))
}

// clampedValue keeps slider value inside range
func (s *ReminderSlider) clampedValue(v float64) float64 {
	v = clampFloat64(v, s.Min, s.Max)

	if s.Step > 0 {
		v = math.Round((v-s.Min)/s.Step)*s.Step + s.Min
		v = clampFloat64(v, s.Min, s.Max)
	}

	return v
}

// CreateRenderer builds reminder slider renderer
func (s *ReminderSlider) CreateRenderer() fyne.WidgetRenderer {
	track := canvas.NewRectangle(theme.Color(theme.ColorNameSeparator))
	knob := canvas.NewCircle(theme.Color(theme.ColorNamePrimary))

	return &reminderSliderRenderer{
		slider: s,
		track:  track,
		knob:   knob,
		objs:   []fyne.CanvasObject{track, knob},
	}
}

// reminderSliderRenderer keeps slider visuals in sync
type reminderSliderRenderer struct {
	slider *ReminderSlider
	track  *canvas.Rectangle
	knob   *canvas.Circle
	objs   []fyne.CanvasObject
}

// Layout positions slider track and knob
func (r *reminderSliderRenderer) Layout(size fyne.Size) {
	centerY := size.Height / 2
	trackWidth := size.Width - 2*reminderSliderPadding
	if trackWidth < 0 {
		trackWidth = 0
	}

	r.track.Resize(fyne.NewSize(trackWidth, reminderSliderTrackHeight))
	r.track.Move(fyne.NewPos(reminderSliderPadding, centerY-reminderSliderTrackHeight/2))

	if r.slider.Max <= r.slider.Min {
		return
	}

	ratio := float32((r.slider.Value - r.slider.Min) / (r.slider.Max - r.slider.Min))
	ratio = clampFloat32(ratio, 0, 1)

	x := reminderSliderPadding + ratio*trackWidth
	r.knob.Resize(fyne.NewSize(reminderSliderKnobRadius*2, reminderSliderKnobRadius*2))
	r.knob.Move(fyne.NewPos(x-reminderSliderKnobRadius, centerY-reminderSliderKnobRadius))
}

// MinSize returns slider minimum size
func (r *reminderSliderRenderer) MinSize() fyne.Size {
	return fyne.NewSize(reminderSliderMinWidth, reminderSliderMinHeight)
}

// Refresh redraws slider colors and position
func (r *reminderSliderRenderer) Refresh() {
	r.Layout(r.slider.Size())

	r.track.FillColor = theme.Color(theme.ColorNameSeparator)
	r.track.Refresh()

	r.knob.FillColor = theme.Color(theme.ColorNamePrimary)
	r.knob.Refresh()
}

// BackgroundColor keeps renderer transparent
func (r *reminderSliderRenderer) BackgroundColor() fyne.ThemeColorName { return "" }

// Objects returns slider objects
func (r *reminderSliderRenderer) Objects() []fyne.CanvasObject { return r.objs }

// Destroy releases renderer resources
func (r *reminderSliderRenderer) Destroy() {}

// newVSpacer creates fixed vertical spacer
func newVSpacer(height float32) fyne.CanvasObject {
	r := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	r.SetMinSize(fyne.NewSize(1, height))

	return r
}

// makeCancelButton builds secondary form action
func (tf *TodoForm) makeCancelButton(text string, onTap func()) fyne.CanvasObject {
	bg := helpers.Lighten(helpers.ToNRGBA(theme.Color(theme.ColorNameBackground)), formCancelButtonLightAmount)
	fg := helpers.ToNRGBA(theme.Color(theme.ColorNameForeground))

	return widgets.NewSimpleRectButton(text, bg, fg, fyne.NewSize(formButtonWidth, formButtonHeight), formButtonRadius, onTap)
}

// makePrimaryButton builds primary form action
func (tf *TodoForm) makePrimaryButton(text string, onTap func()) fyne.CanvasObject {
	bg := helpers.ToNRGBA(theme.Color(theme.ColorNamePrimary))
	fg := helpers.ToNRGBA(theme.Color(theme.ColorNameForeground))

	return widgets.NewSimpleRectButton(text, bg, fg, fyne.NewSize(formButtonWidth, formButtonHeight), formButtonRadius, onTap)
}

// showDateTimePicker displays date and time picker
func (tf *TodoForm) showDateTimePicker() {
	dateEntry := widget.NewEntry()
	dateEntry.SetText(tf.selectedDateTime.Format(dateInputLayout))
	dateEntry.SetPlaceHolder("DD.MM.YYYY")
	dateEntry.Resize(fyne.NewSize(formDateEntryWidth, formDateEntryHeight))

	timeEntry := widget.NewEntry()
	timeEntry.SetText(tf.selectedDateTime.Format(timeInputLayout))
	timeEntry.SetPlaceHolder("HH:MM")
	timeEntry.Resize(fyne.NewSize(formDateEntryWidth, formDateEntryHeight))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Date (DD.MM.YYYY):", Widget: dateEntry},
			{Text: "Time (HH:MM):", Widget: timeEntry},
		},
	}

	dialogParent := tf.dialogParent()
	var dateTimeDialog dialog.Dialog

	form.OnSubmit = func() {
		selected, err := tf.parseDateTimeInput(dateEntry.Text, timeEntry.Text)
		if err != nil {
			dialog.ShowError(err, dialogParent)
			return
		}

		tf.selectedDateTime = selected
		tf.updateDateTimeDisplay()
		dateTimeDialog.Hide()
	}

	form.OnCancel = func() {
		dateTimeDialog.Hide()
	}

	dateTimeDialog = dialog.NewCustomWithoutButtons("Select Date and Time", container.NewVBox(form), dialogParent)
	dateTimeDialog.Resize(fyne.NewSize(formDateTimeDialogWidth, form.MinSize().Height+formDateTimeDialogPadding))
	dateTimeDialog.Show()
}

// dialogParent returns active parent window for nested dialogs
func (tf *TodoForm) dialogParent() fyne.Window {
	if tf.formWindow != nil {
		return tf.formWindow
	}

	return tf.parentWindow
}

// parseDateTimeInput parses date picker values
func (tf *TodoForm) parseDateTimeInput(dateText, timeText string) (time.Time, error) {
	date, err := time.Parse(dateInputLayout, dateText)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %w", err)
	}

	clock, err := time.Parse(timeInputLayout, timeText)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time: %w", err)
	}

	year, month, day := date.Date()
	location := tf.selectedDateTime.Location()

	return time.Date(year, month, day, clock.Hour(), clock.Minute(), 0, 0, location), nil
}

// updateDateTimeDisplay updates date/time entry text
func (tf *TodoForm) updateDateTimeDisplay() {
	tf.dateTimeEntry.SetText(tf.selectedDateTime.Format(dateTimeDisplayLayout))
}

// setSelectIndex applies selected index with fallback
func setSelectIndex(selectWidget *widget.Select, index, fallback int) {
	if selectWidget == nil {
		return
	}

	if index < 0 || index >= len(selectWidget.Options) {
		selectWidget.SetSelectedIndex(fallback)
		return
	}

	selectWidget.SetSelectedIndex(index)
}

// selectedIndexOrDefault returns selected index or fallback
func selectedIndexOrDefault(selectWidget *widget.Select, fallback int) int {
	if selectWidget == nil {
		return fallback
	}

	index := selectWidget.SelectedIndex()
	if index < 0 {
		return fallback
	}

	return index
}

// clampFloat64 keeps value inside range
func clampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

// clampFloat32 keeps value inside range
func clampFloat32(value, min, max float32) float32 {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
