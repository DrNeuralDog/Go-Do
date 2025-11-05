package forms

import (
	"errors"
	"fmt"
	"time"

	"todo-list-migration/src/localization"
	"todo-list-migration/src/models"
	"todo-list-migration/src/persistence"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// TodoForm represents a form for creating/editing todo items
type TodoForm struct {
	window      fyne.Window
	dataManager *persistence.MonthlyManager

	// Form fields
	nameEntry      *widget.Entry
	contentEntry   *widget.Entry
	placeEntry     *widget.Entry
	labelEntry     *widget.Entry
	dateTimeEntry  *widget.Entry
	dateTimeButton *widget.Button
	prioritySelect *widget.Select
	kindSelect     *widget.Select
	warnTimeSlider *widget.Slider
	warnTimeLabel  *widget.Label

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
func NewTodoForm(window fyne.Window, dataManager *persistence.MonthlyManager) *TodoForm {
	tf := &TodoForm{
		window:      window,
		dataManager: dataManager,
		isEditMode:  false,
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
	}, tf.window)
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
	}, tf.window)
	// Keep edit dialog consistent with add dialog width
	dialog.Resize(fyne.NewSize(700, 600))
	tf.currentDialog = dialog
	dialog.Show()
}

// ShowCreateWindow opens the form in a standalone window for creating a new todo
func (tf *TodoForm) ShowCreateWindow(onSave func()) {
    tf.isEditMode = false
    tf.onSaveCallback = onSave
    tf.resetForm()

    title := localization.GetString("form_title_add")
    win := fyne.CurrentApp().NewWindow(title)
    tf.window = win

    // Build form content
    form := &widget.Form{
        Items: []*widget.FormItem{
            {Text: "Name:", Widget: tf.nameEntry},
            {Text: "Date/Time:", Widget: container.NewBorder(nil, nil, nil, tf.dateTimeButton, tf.dateTimeEntry)},
            {Text: "Location:", Widget: tf.placeEntry},
            {Text: "Label:", Widget: tf.labelEntry},
            {Text: "Type:", Widget: tf.kindSelect},
            {Text: "Priority:", Widget: tf.prioritySelect},
            {Text: "Content:", Widget: container.NewScroll(tf.contentEntry)},
            {Text: "Reminder:", Widget: container.NewVBox(tf.warnTimeSlider, tf.warnTimeLabel)},
        },
    }

    // Buttons
    addBtn := widget.NewButton(localization.GetString("form_button_add"), func() {
        if err := tf.trySubmit(); err != nil {
            dialog.ShowError(err, tf.window)
            return
        }
        if tf.onSaveCallback != nil {
            tf.onSaveCallback()
        }
        win.Close()
    })
    cancelBtn := widget.NewButton(localization.GetString("form_button_cancel"), func() {
        win.Close()
    })

    buttons := container.NewHBox(addBtn, cancelBtn)
    content := container.NewBorder(nil, buttons, nil, nil, form)
    win.SetContent(content)
    win.Resize(fyne.NewSize(700, 600))
    win.Show()
}

// ShowEditWindow opens the form in a standalone window for editing an existing todo
func (tf *TodoForm) ShowEditWindow(todo *models.TodoItem, originalTime time.Time, onSave func()) {
    tf.isEditMode = true
    tf.originalTodo = todo
    tf.originalTime = originalTime
    tf.onSaveCallback = onSave
    tf.populateForm(todo)

    title := localization.GetString("form_title_edit")
    win := fyne.CurrentApp().NewWindow(title)
    tf.window = win

    // Build form content
    form := &widget.Form{
        Items: []*widget.FormItem{
            {Text: "Name:", Widget: tf.nameEntry},
            {Text: "Date/Time:", Widget: container.NewBorder(nil, nil, nil, tf.dateTimeButton, tf.dateTimeEntry)},
            {Text: "Location:", Widget: tf.placeEntry},
            {Text: "Label:", Widget: tf.labelEntry},
            {Text: "Type:", Widget: tf.kindSelect},
            {Text: "Priority:", Widget: tf.prioritySelect},
            {Text: "Content:", Widget: container.NewScroll(tf.contentEntry)},
            {Text: "Reminder:", Widget: container.NewVBox(tf.warnTimeSlider, tf.warnTimeLabel)},
        },
    }

    // Buttons
    saveBtn := widget.NewButton(localization.GetString("form_button_save"), func() {
        if err := tf.trySubmit(); err != nil {
            dialog.ShowError(err, tf.window)
            return
        }
        if tf.onSaveCallback != nil {
            tf.onSaveCallback()
        }
        win.Close()
    })
    cancelBtn := widget.NewButton(localization.GetString("form_button_cancel"), func() {
        win.Close()
    })

    buttons := container.NewHBox(saveBtn, cancelBtn)
    content := container.NewBorder(nil, buttons, nil, nil, form)
    win.SetContent(content)
    win.Resize(fyne.NewSize(700, 600))
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
	tf.warnTimeSlider = widget.NewSlider(0, 864)
	tf.warnTimeSlider.Step = 5
	tf.warnTimeSlider.Value = 0
	tf.warnTimeSlider.OnChanged = tf.onWarnTimeChanged

	tf.warnTimeLabel = widget.NewLabel(localization.GetString("reminder_none"))
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
	tf.warnTimeSlider.Value = 0
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
	tf.warnTimeSlider.Value = float64(todo.WarnTime)
	tf.onWarnTimeChanged(float64(todo.WarnTime))
}

// onWarnTimeChanged updates the warning time label
func (tf *TodoForm) onWarnTimeChanged(value float64) {
	warnTime := int(value)
	if warnTime == 0 {
		tf.warnTimeLabel.SetText(localization.GetString("reminder_none"))
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
		tf.warnTimeLabel.SetText(localization.GetString("reminder_none"))
	} else {
		reminderFormat := localization.GetString("reminder_format")
		tf.warnTimeLabel.SetText(fmt.Sprintf(reminderFormat, joinStrings(parts, " ")))
	}
}

// onSubmit handles form submission
func (tf *TodoForm) onSubmit() {
    if err := tf.trySubmit(); err != nil {
        dialog.ShowError(err, tf.window)
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
	dateTimeDialog := dialog.NewCustomWithoutButtons("Select Date and Time", formContainer, tf.window)
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
