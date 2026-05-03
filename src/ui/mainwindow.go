package ui

import (
	"fmt"
	"image/color"
	"time"

	assets "godo/resources"
	"godo/src/localization"
	"godo/src/models"
	"godo/src/persistence"
	"godo/src/ui/forms"
	"godo/src/ui/helpers"
	"godo/src/ui/widgets"
	"godo/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// mainWindowColors keeps resolved theme colors
type mainWindowColors struct {
	backgroundStart color.Color
	backgroundEnd   color.Color
	title           color.Color
	logo            fyne.Resource
	navBg           color.Color
	navFg           color.Color
	selectBg        color.Color
	themeBg         color.Color
	themeFg         color.Color
	pomodoroBg      color.Color
	pomodoroFg      color.Color
}

// MainWindow manages main app UI
type MainWindow struct {
	window        fyne.Window
	dataManager   persistence.TodoRepository
	configManager persistence.ConfigRepository
	config        *models.Config
	todoForm      *forms.TodoForm
	timeline      *Timeline

	titleLabel  *widget.Label
	addButton   *widget.Button
	prevButton  *widget.Button
	nextButton  *widget.Button
	viewModeBtn *widget.Button // legacy hidden control
	themeBtn    *widget.Button // legacy hidden control

	viewSelect      *widgets.CustomSelect
	prevRectBtn     *widgets.SimpleRectButton
	nextRectBtn     *widgets.SimpleRectButton
	pomodoroRectBtn *widgets.SimpleRectButton
	themeRectBtn    *widgets.SimpleRectButton

	currentDate    time.Time
	viewMode       models.ViewMode
	todos          []*models.TodoItem
	isGruvbox      bool
	pomodoroWindow *PomodoroWindow
	todoFormWindow fyne.Window

	pendingReorderYear  int
	pendingReorderMonth int
	pendingReorderTodos []*models.TodoItem
}

// NewMainWindow builds main application window
func NewMainWindow(window fyne.Window, dataManager persistence.TodoRepository, configManager persistence.ConfigRepository) *MainWindow {
	mw := &MainWindow{
		window:        window,
		dataManager:   dataManager,
		configManager: configManager,
		currentDate:   time.Now(),
		viewMode:      models.ViewIncomplete,
		isGruvbox:     false,
	}

	mw.loadConfig()

	if mw.config.GetCurrentDate().IsZero() {
		mw.findAndSetCurrentDateFromDataFile()
	}

	mw.todoForm = forms.NewTodoForm(window, mw.dataManager)

	mw.timeline = NewTimeline(mw.dataManager)
	mw.timeline.SetWindow(window)
	mw.timeline.SetOnTodoSelected(mw.onTodoSelected)
	mw.timeline.SetOnTodoReorder(mw.onTodoReorder)
	mw.timeline.SetOnReorderFinished(mw.onReorderFinished)
	mw.timeline.SetOnTodosChanged(func() {
		mw.loadTodos()
		mw.refreshView()
	})

	mw.setupUI()
	mw.loadTodos()
	mw.refreshView()

	return mw
}

// findAndSetCurrentDateFromDataFile selects latest todo date from storage
func (mw *MainWindow) findAndSetCurrentDateFromDataFile() {
	months, err := mw.dataManager.GetAllMonths()

	if err != nil {
		return
	}

	var latestTime time.Time
	for _, dateKey := range months {
		year, month := utils.ParseDateKey(dateKey)

		if year == 0 {
			continue
		}

		todos, err := mw.dataManager.GetTodosForMonth(year, month)
		if err != nil {
			continue
		}

		for _, todo := range todos {
			if todo.TodoTime.After(latestTime) {
				latestTime = todo.TodoTime
			}
		}
	}

	if !latestTime.IsZero() {
		mw.currentDate = latestTime
	}
}

// setupUI rebuilds main window content
func (mw *MainWindow) setupUI() {
	mw.setupWindow()
	mw.setupLegacyControls()

	colors := mainWindowThemeColors()
	header := createHeader(colors)
	controls := mw.createControls(colors)
	content := mw.createContent(header, controls, colors)
	background := NewGradientRect(colors.backgroundStart, colors.backgroundEnd, 0)

	mw.window.SetContent(container.NewMax(background, content))
}

// setupWindow applies fixed main window properties
func (mw *MainWindow) setupWindow() {
	mw.window.SetTitle(localization.GetString("window_title"))
	mw.window.Resize(fyne.NewSize(MainWindowWidth, MainWindowHeight))
	mw.window.SetFixedSize(true)
}

// setupLegacyControls keeps old hidden controls alive
func (mw *MainWindow) setupLegacyControls() {
	mw.titleLabel = widget.NewLabel(localization.GetString("window_title") + " - User")
	mw.titleLabel.Alignment = fyne.TextAlignCenter
	mw.titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	mw.titleLabel.Wrapping = fyne.TextTruncate

	mw.addButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), mw.onAddButtonClicked)
	mw.addButton.Importance = widget.HighImportance

	mw.prevButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), mw.onPrevDayClicked)
	mw.prevButton.Hide()

	mw.nextButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), mw.onNextDayClicked)
	mw.nextButton.Hide()

	mw.viewModeBtn = widget.NewButton(mw.viewMode.GetLabel(), mw.onViewModeClicked)
	mw.viewModeBtn.Hide()

	mw.themeBtn = widget.NewButton("Gruvbox", mw.onThemeToggleClicked)
	mw.themeBtn.Hide()
}

// mainWindowThemeColors returns colors for current theme
func mainWindowThemeColors() mainWindowColors {
	colors := mainWindowColors{}
	currentTheme := fyne.CurrentApp().Settings().Theme()
	isLightTheme := helpers.IsLightTheme()

	if gradientTheme, ok := currentTheme.(interface {
		GetHeaderGradientColors() (color.Color, color.Color)
	}); ok {
		colors.backgroundStart, colors.backgroundEnd = gradientTheme.GetHeaderGradientColors()
	} else {
		colors.backgroundStart = theme.Color(theme.ColorNamePrimary)
		colors.backgroundEnd = colors.backgroundStart
	}

	if isLightTheme {
		colors.title = color.White
		colors.logo = assets.HeaderIconLight
		colors.navBg = helpers.Hex(ColorHexAccentLight)
		colors.navFg = color.White
		colors.selectBg = color.White
		colors.themeBg = helpers.Hex(ColorHexAccentLight)
		colors.themeFg = color.White
		colors.pomodoroBg = helpers.Hex(ColorHexAccentLight)
		colors.pomodoroFg = color.White
	} else {
		colors.title = helpers.Hex(ColorHexGruvboxPrimary)
		colors.logo = assets.HeaderIconDark
		colors.navBg = helpers.Hex(ColorHexAccentDark)
		colors.navFg = helpers.Hex(ColorHexGruvboxPrimary)
		colors.selectBg = helpers.Hex(ColorHexGruvboxSurface)
		colors.themeBg = helpers.Hex(ColorHexAccentDark)
		colors.themeFg = helpers.Hex(ColorHexGruvboxPrimary)
		colors.pomodoroBg = helpers.Hex(ColorHexAccentDark)
		colors.pomodoroFg = helpers.Hex(ColorHexGruvboxPrimary)
	}

	return colors
}

// createHeader builds app logo and title
func createHeader(colors mainWindowColors) fyne.CanvasObject {
	logoImg := canvas.NewImageFromResource(colors.logo)
	logoImg.FillMode = canvas.ImageFillContain

	logoSize := float32(84)
	logoImg.SetMinSize(fyne.NewSize(logoSize, logoSize))

	titleTxt := canvas.NewText("GO DO", colors.title)
	titleTxt.TextSize = 64
	titleTxt.TextStyle = fyne.TextStyle{Bold: false}

	iconPad := titleTxt.TextSize - logoSize

	if iconPad < 0 {
		iconPad = 0
	}

	logoAligned := container.NewVBox(
		helpers.CreateSpacer(1, iconPad),
		container.NewMax(logoImg),
	)

	return container.NewHBox(logoAligned, titleTxt)
}

// createControls builds top navigation controls
func (mw *MainWindow) createControls(colors mainWindowColors) fyne.CanvasObject {
	viewOptions := []string{"All", "Incomplete", "Complete", "Important"}
	mw.viewSelect = NewCustomSelect(viewOptions, mw.onViewModeSelected)

	mw.viewSelect.SetSelected(mw.viewMode.GetLabel())

	selectWrapper := CreateStyledSelect(mw.viewSelect, colors.selectBg, fyne.NewSize(180, ButtonHeight), BorderRadius)

	mw.prevRectBtn = NewSimpleRectButton("←", colors.navBg, colors.navFg, fyne.NewSize(ButtonHeight, ButtonHeight), BorderRadius, mw.onPrevDayClicked)
	mw.nextRectBtn = NewSimpleRectButton("→", colors.navBg, colors.navFg, fyne.NewSize(ButtonHeight, ButtonHeight), BorderRadius, mw.onNextDayClicked)

	addButtonRounded := RoundedIconButton(theme.ContentAddIcon(), mw.onAddButtonClicked)
	addWrapTop := container.NewGridWrap(fyne.NewSize(ButtonHeight, ButtonHeight), addButtonRounded)

	return container.NewHBox(
		selectWrapper,
		helpers.CreateSpacer(10, 1),

		mw.prevRectBtn,
		helpers.CreateSpacer(2, 1),

		mw.nextRectBtn,
		helpers.CreateSpacer(10, 1),
		addWrapTop,
	)
}

// createContent lays out main content over background
func (mw *MainWindow) createContent(header, controls fyne.CanvasObject, colors mainWindowColors) fyne.CanvasObject {
	mw.timeline.SetDate(mw.currentDate)
	mw.timeline.SetViewMode(mw.viewMode)

	timelineCard := CreateTasksContainer(mw.timeline)
	headerPadded := withHorizontalPadding(header, 24)
	controlsPadded := withHorizontalPadding(controls, 24)
	timelinePadded := withHorizontalPadding(timelineCard, 24)

	headerArea := container.NewVBox(
		helpers.CreateSpacer(1, 15),
		headerPadded,
		helpers.CreateSpacer(1, 30),
		controlsPadded,
		helpers.CreateSpacer(1, 30),
	)

	appBody := container.NewBorder(
		headerArea,
		helpers.CreateSpacer(1, 24),
		nil, nil,
		timelinePadded,
	)

	return container.NewBorder(
		nil,
		mw.setupBottomButtons(colors),
		nil, nil,
		appBody,
	)
}

// withHorizontalPadding wraps object with equal side padding
func withHorizontalPadding(object fyne.CanvasObject, padding float32) fyne.CanvasObject {
	return container.NewBorder(nil, nil, helpers.CreateSpacer(padding, 1), helpers.CreateSpacer(padding, 1), object)
}

// onViewModeSelected applies selected filter mode
func (mw *MainWindow) onViewModeSelected(selected string) {
	mw.viewMode = viewModeFromLabel(selected)

	mw.loadTodos()
	mw.refreshView()
	mw.saveConfig()
}

// viewModeFromLabel maps select label to view mode
func viewModeFromLabel(label string) models.ViewMode {
	switch label {
	case "All":
		return models.ViewAll
	case "Incomplete":
		return models.ViewIncomplete
	case "Complete":
		return models.ViewComplete
	case "Important":
		return models.ViewStarred
	default:
		return models.ViewIncomplete
	}
}

// onThemeToggleClicked switches app theme
func (mw *MainWindow) onThemeToggleClicked() {
	mw.isGruvbox = !mw.isGruvbox

	if mw.isGruvbox {
		fyne.CurrentApp().Settings().SetTheme(NewGruvboxBlackTheme())
		mw.themeBtn.SetText("Light")
	} else {
		fyne.CurrentApp().Settings().SetTheme(NewLightSoftTheme())
		mw.themeBtn.SetText("Gruvbox")
	}

	mw.setupUI()
	mw.loadTodos()
	mw.refreshView()

	if mw.pomodoroWindow != nil {
		mw.pomodoroWindow.UpdateTheme(mw.isGruvbox)
	}

	mw.saveConfig()
}

// loadTodos loads visible todos for selected day
func (mw *MainWindow) loadTodos() {
	year, month := mw.currentDate.Year(), int(mw.currentDate.Month())
	monthlyTodos, err := mw.dataManager.GetTodosForMonth(year, month)

	if err != nil {
		fmt.Println(localization.GetStringWithArgs("error_load_failed", err.Error()))
		mw.todos = []*models.TodoItem{}

		return
	}

	currentTime := time.Now()
	dailyTodos := mw.currentDayTodos(monthlyTodos)
	mw.todos = mw.viewMode.FilterItems(dailyTodos, currentTime)

	models.SortTodosByOrder(mw.todos)
}

// currentDayTodos returns todos inside selected day
func (mw *MainWindow) currentDayTodos(todos []*models.TodoItem) []*models.TodoItem {
	startOfDay, endOfDay := mw.currentDayRange()
	dailyTodos := make([]*models.TodoItem, 0)

	for _, todo := range todos {
		if todo == nil {
			continue
		}

		if todoInDay(todo, startOfDay, endOfDay) {
			dailyTodos = append(dailyTodos, todo)
		}
	}

	return dailyTodos
}

// currentDayRange returns inclusive start and exclusive end
func (mw *MainWindow) currentDayRange() (time.Time, time.Time) {
	startOfDay := time.Date(mw.currentDate.Year(), mw.currentDate.Month(), mw.currentDate.Day(), 0, 0, 0, 0, mw.currentDate.Location())

	return startOfDay, startOfDay.Add(24 * time.Hour)
}

// todoInDay checks selected day bounds
func todoInDay(todo *models.TodoItem, startOfDay, endOfDay time.Time) bool {
	return !todo.TodoTime.Before(startOfDay) && todo.TodoTime.Before(endOfDay)
}

// refreshView syncs timeline state
func (mw *MainWindow) refreshView() {
	mw.timeline.SetDate(mw.currentDate)
	mw.timeline.SetViewMode(mw.viewMode)
	mw.timeline.SetTodos(mw.todos)

	mw.timeline.Refresh()
}

// onAddButtonClicked opens create todo form
func (mw *MainWindow) onAddButtonClicked() {
	if mw.todoFormWindow != nil {
		FlashWindow(mw.todoFormWindow)

		return
	}

	mw.todoForm.ShowCreateWindow(
		func() {
			mw.loadTodos()
			mw.refreshView()
		},
		func(win fyne.Window) {
			mw.todoFormWindow = win
		},
		func() {
			mw.todoFormWindow = nil
		},
	)
}

// onPrevDayClicked moves selection to previous day
func (mw *MainWindow) onPrevDayClicked() {
	mw.currentDate = mw.currentDate.AddDate(0, 0, -1)

	mw.loadTodos()
	mw.refreshView()
	mw.saveConfig()
}

// onNextDayClicked moves selection to next day
func (mw *MainWindow) onNextDayClicked() {
	mw.currentDate = mw.currentDate.AddDate(0, 0, 1)

	mw.loadTodos()
	mw.refreshView()
	mw.saveConfig()
}

// onViewModeClicked cycles legacy hidden view mode
func (mw *MainWindow) onViewModeClicked() {
	mw.viewMode = mw.viewMode.GetNextMode()

	if mw.viewSelect != nil {
		mw.viewSelect.SetSelected(mw.viewMode.GetLabel())
	}

	mw.loadTodos()
	mw.refreshView()
	mw.saveConfig()
}

// onTodoSelected opens edit todo form
func (mw *MainWindow) onTodoSelected(todo *models.TodoItem, todoTime time.Time) {
	if mw.todoFormWindow != nil {
		FlashWindow(mw.todoFormWindow)

		return
	}

	mw.todoForm.ShowEditWindow(
		todo,
		todoTime,
		func() {
			mw.loadTodos()
			mw.refreshView()
		},
		func(win fyne.Window) {
			mw.todoFormWindow = win
		},
		func() {
			mw.todoFormWindow = nil
		},
	)
}

// onTodoReorder moves todo inside selected day
func (mw *MainWindow) onTodoReorder(todo *models.TodoItem, delta int) {
	if delta == 0 || todo == nil {
		return
	}

	year, month := mw.currentDate.Year(), int(mw.currentDate.Month())
	monthlyTodos, err := mw.dataManager.GetTodosForMonth(year, month)

	if err != nil {
		return
	}

	dayTodos := mw.currentDayTodos(monthlyTodos)
	models.SortTodosByOrder(dayTodos)

	idx := findTodoIndex(dayTodos, todo)
	if idx == -1 {
		return
	}

	newIdx := boundedIndex(idx+delta, len(dayTodos))
	if newIdx == idx {
		return
	}

	moveTodo(dayTodos, idx, newIdx)
	assignTodoOrders(dayTodos)
	mw.syncVisibleOrders(dayTodos)
	mw.rememberPendingReorder(year, month, monthlyTodos)

	models.SortTodosByOrder(mw.todos)
	mw.timeline.SetTodos(mw.todos)

	mw.timeline.Refresh()
}

// findTodoIndex returns todo position in day list
func findTodoIndex(todos []*models.TodoItem, target *models.TodoItem) int {
	for i, todo := range todos {
		if todo == target || sameTodoIdentity(todo, target) {
			return i
		}
	}

	return -1
}

// sameTodoIdentity compares persisted todo identity
func sameTodoIdentity(left, right *models.TodoItem) bool {
	if left == nil || right == nil {
		return false
	}

	return left.TodoTime.Equal(right.TodoTime) && left.Name == right.Name
}

// boundedIndex keeps index inside slice bounds
func boundedIndex(index, length int) int {
	if index < 0 {
		return 0
	}

	if index >= length {
		return length - 1
	}

	return index
}

// moveTodo moves item inside slice
func moveTodo(todos []*models.TodoItem, oldIndex, newIndex int) {
	item := todos[oldIndex]

	if newIndex > oldIndex {
		copy(todos[oldIndex:], todos[oldIndex+1:newIndex+1])
	} else {
		copy(todos[newIndex+1:], todos[newIndex:oldIndex])
	}

	todos[newIndex] = item
}

// assignTodoOrders writes sequential order values
func assignTodoOrders(todos []*models.TodoItem) {
	for i, todo := range todos {
		todo.Order = i + 1
	}
}

// syncVisibleOrders updates filtered visible list order
func (mw *MainWindow) syncVisibleOrders(dayTodos []*models.TodoItem) {
	orderMap := make(map[string]int, len(dayTodos))

	for _, todo := range dayTodos {
		orderMap[todoIdentityKey(todo)] = todo.Order
	}

	for _, todo := range mw.todos {
		if order, ok := orderMap[todoIdentityKey(todo)]; ok {
			todo.Order = order
		}
	}
}

// todoIdentityKey returns stable key for current storage model
func todoIdentityKey(todo *models.TodoItem) string {
	if todo == nil {
		return ""
	}

	return fmt.Sprintf("%d|%s", todo.TodoTime.UnixNano(), todo.Name)
}

// rememberPendingReorder stores changed monthly list for final save
func (mw *MainWindow) rememberPendingReorder(year, month int, todos []*models.TodoItem) {
	mw.pendingReorderYear = year
	mw.pendingReorderMonth = month
	mw.pendingReorderTodos = todos
}

// onReorderFinished persists reordered list once
func (mw *MainWindow) onReorderFinished() {
	if mw.pendingReorderTodos == nil {
		return
	}

	if err := mw.dataManager.SaveTodosForMonth(mw.pendingReorderYear, mw.pendingReorderMonth, mw.pendingReorderTodos); err != nil {
		fmt.Printf("Failed to save todo order: %v\n", err)
	}

	mw.pendingReorderTodos = nil
}

// setupBottomButtons creates theme and Pomodoro controls
func (mw *MainWindow) setupBottomButtons(colors mainWindowColors) fyne.CanvasObject {
	themeLabel := "Dark"

	if mw.isGruvbox {
		themeLabel = "Light"
	}

	mw.pomodoroRectBtn = NewSimpleRectButton("Pomodoro", colors.pomodoroBg, colors.pomodoroFg, fyne.NewSize(100, ButtonHeight), BorderRadius, mw.onPomodoroTopClicked)
	mw.themeRectBtn = NewSimpleRectButton(themeLabel, colors.themeBg, colors.themeFg, fyne.NewSize(100, ButtonHeight), BorderRadius, mw.onThemeToggleClicked)

	bottomButtons := container.NewBorder(
		nil, nil,
		container.NewBorder(nil, nil, helpers.CreateSpacer(25, 1), nil, mw.themeRectBtn),
		container.NewBorder(nil, nil, nil, helpers.CreateSpacer(25, 1), mw.pomodoroRectBtn),
		canvas.NewRectangle(color.Transparent),
	)

	return container.NewVBox(
		bottomButtons,
		helpers.CreateSpacer(1, 10),
	)
}

// onPomodoroTopClicked opens Pomodoro window
func (mw *MainWindow) onPomodoroTopClicked() {
	if mw.pomodoroWindow != nil {
		FlashWindow(mw.pomodoroWindow.window)

		return
	}

	mw.pomodoroWindow = NewPomodoroWindow(fyne.CurrentApp(), mw.isGruvbox)

	mw.pomodoroWindow.SetOnClosed(func() {
		mw.pomodoroWindow = nil
	})

	mw.pomodoroWindow.Show()
}

// loadConfig loads persisted UI state
func (mw *MainWindow) loadConfig() {
	config, err := mw.configManager.LoadConfig()

	if err != nil {
		fmt.Printf("Failed to load config: %v, using defaults\n", err)

		config = models.NewDefaultConfig()
	}

	mw.config = config

	if config.GetTheme() == "dark" {
		mw.isGruvbox = true
		fyne.CurrentApp().Settings().SetTheme(NewGruvboxBlackTheme())
	} else {
		mw.isGruvbox = false
		fyne.CurrentApp().Settings().SetTheme(NewLightSoftTheme())
	}

	mw.viewMode = models.ViewModeFromString(config.GetViewMode())

	if !config.GetCurrentDate().IsZero() {
		mw.currentDate = config.GetCurrentDate()
	}
}

// saveConfig persists current UI state
func (mw *MainWindow) saveConfig() {
	if mw.isGruvbox {
		mw.config.SetTheme("dark")
	} else {
		mw.config.SetTheme("light")
	}

	mw.config.SetViewMode(mw.viewMode.String())

	mw.config.SetCurrentDate(mw.currentDate)

	if err := mw.configManager.SaveConfig(mw.config); err != nil {
		fmt.Printf("Failed to save config: %v\n", err)
	}
}
