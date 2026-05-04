package ui

import (
	"fmt"
	"image/color"
	"time"

	"godo/src/localization"
	"godo/src/models"
	"godo/src/ui/helpers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

// confirmDelete shows delete confirmation
func (t *Timeline) confirmDelete(todo *models.TodoItem) {
	if todo == nil {
		return
	}

	if t.window == nil {
		t.deleteTodo(todo)

		return
	}

	var restoreTheme func()
	currentTheme := fyne.CurrentApp().Settings().Theme()
	isLightTheme := helpers.IsLightTheme()
	messageColor := helpers.ToNRGBA(theme.Color(theme.ColorNameForeground))
	overrideTheme := &foregroundOverrideTheme{
		base:         currentTheme,
		headingDelta: 5,
		textDelta:    5,
	}

	if isLightTheme {
		messageColor = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		overrideTheme.overrideForeground = true
		overrideTheme.overrideColor = messageColor
	}

	fyne.CurrentApp().Settings().SetTheme(overrideTheme)

	restoreTheme = func() {
		fyne.CurrentApp().Settings().SetTheme(currentTheme)
	}

	message := canvas.NewText(localization.GetString("confirm_delete_message"), messageColor)
	message.Alignment = fyne.TextAlignLeading
	message.TextSize = 16

	content := container.NewVBox(
		container.NewHBox(message, layout.NewSpacer()),
		helpers.CreateSpacer(1, 20),
	)

	conf := dialog.NewCustomConfirm(
		localization.GetString("confirm_delete_title"),
		localization.GetString("shortcut_delete"),
		localization.GetString("form_button_cancel"),
		content,
		func(confirm bool) {
			if confirm {
				t.deleteTodo(todo)
			}
		},
		t.window,
	)

	conf.SetOnClosed(func() {
		if restoreTheme != nil {
			restoreTheme()
		}
	})

	conf.Show()
}

// deleteTodo removes todo and refreshes timeline
func (t *Timeline) deleteTodo(todo *models.TodoItem) {
	if todo == nil {
		return
	}

	if err := t.dataManager.RemoveTodo(todo.TodoTime); err != nil {
		t.showError(err)

		return
	}

	t.notifyTodosChanged()
}

// notifyTodosChanged delegates refresh to parent callback
func (t *Timeline) notifyTodosChanged() {
	if t.onTodosChanged != nil {
		t.onTodosChanged()

		return
	}

	t.reloadVisibleTodos()
}

// reloadVisibleTodos reloads selected day after inline changes
func (t *Timeline) reloadVisibleTodos() {
	year, month := t.currentDate.Year(), int(t.currentDate.Month())
	todos, err := t.dataManager.GetTodosForMonth(year, month)

	if err != nil {
		t.showError(err)

		return
	}

	dailyTodos := t.filterCurrentDayTodos(todos)
	filtered := t.viewMode.FilterItems(dailyTodos)
	models.SortTodosByOrder(filtered)

	t.SetTodos(filtered)
	t.Refresh()
}

// filterCurrentDayTodos keeps monthly data scoped to selected day
func (t *Timeline) filterCurrentDayTodos(todos []*models.TodoItem) []*models.TodoItem {
	startOfDay := time.Date(t.currentDate.Year(), t.currentDate.Month(), t.currentDate.Day(), 0, 0, 0, 0, t.currentDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	dailyTodos := make([]*models.TodoItem, 0, len(todos))

	for _, todo := range todos {
		if todo == nil {
			continue
		}

		if !todo.TodoTime.Before(startOfDay) && todo.TodoTime.Before(endOfDay) {
			dailyTodos = append(dailyTodos, todo)
		}
	}

	return dailyTodos
}

// showError shows dialog error or logs fallback
func (t *Timeline) showError(err error) {
	if err == nil {
		return
	}

	if t.window != nil {
		dialog.NewError(err, t.window).Show()

		return
	}

	fmt.Println("timeline error:", err)
}

// foregroundOverrideTheme adjusts dialog text color
type foregroundOverrideTheme struct {
	base               fyne.Theme
	overrideColor      color.Color
	overrideForeground bool
	headingDelta       float32
	textDelta          float32
}

func (t *foregroundOverrideTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameForeground && t.overrideForeground {
		return t.overrideColor
	}

	if t.base != nil {
		return t.base.Color(name, variant)
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *foregroundOverrideTheme) Font(style fyne.TextStyle) fyne.Resource {
	if t.base != nil {
		return t.base.Font(style)
	}

	return theme.DefaultTheme().Font(style)
}

func (t *foregroundOverrideTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if t.base != nil {
		return t.base.Icon(name)
	}

	return theme.DefaultTheme().Icon(name)
}

func (t *foregroundOverrideTheme) Size(name fyne.ThemeSizeName) float32 {
	var size float32

	if t.base != nil {
		size = t.base.Size(name)
	} else {
		size = theme.DefaultTheme().Size(name)
	}

	if name == theme.SizeNameHeadingText {
		size += t.headingDelta
	}

	if name == theme.SizeNameText {
		size += t.textDelta
	}

	return size
}
