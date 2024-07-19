package ui

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"

	assets "todo-list-migration/doc"
	"todo-list-migration/src/localization"
	"todo-list-migration/src/models"
	"todo-list-migration/src/persistence"
	"todo-list-migration/src/ui/forms"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// MainWindow represents the main application window
type MainWindow struct {
	window      fyne.Window
	dataManager *persistence.MonthlyManager
	todoForm    *forms.TodoForm
	timeline    *Timeline

	// UI components
	titleLabel  *widget.Label
	addButton   *widget.Button
	prevButton  *widget.Button
	nextButton  *widget.Button
	viewModeBtn *widget.Button // legacy hidden
	themeBtn    *widget.Button // legacy hidden
	// Styled controls
	viewSelect      *widget.Select
	prevRectBtn     *SimpleRectButton
	nextRectBtn     *SimpleRectButton
	pomodoroRectBtn *SimpleRectButton
	themeRectBtn    *SimpleRectButton

	// State
	currentDate time.Time // Changed to time.Time for daily view
	viewMode    models.ViewMode
	todos       []*models.TodoItem
	isGruvbox   bool
}

// NewMainWindow creates a new main window
func NewMainWindow(window fyne.Window, dataDir string) *MainWindow {
	mw := &MainWindow{
		window:      window,
		dataManager: persistence.NewMonthlyManager(dataDir),
		currentDate: time.Now(), // Start with today
		viewMode:    models.ViewIncomplete,
		isGruvbox:   false,
	}

	// Find and set to latest day with data
	mw.findAndSetCurrentDateFromDataFile()

	// Initialize todo form
	mw.todoForm = forms.NewTodoForm(window, mw.dataManager)

	// Initialize timeline
	mw.timeline = NewTimeline(mw.dataManager)
	mw.timeline.SetOnTodoSelected(mw.onTodoSelected)

	mw.setupUI()
	mw.loadTodos()
	mw.refreshView()

	return mw
}

// findAndSetCurrentDateFromDataFile looks for existing data files and sets currentDate to the latest day with todos
func (mw *MainWindow) findAndSetCurrentDateFromDataFile() {
	dataDir := mw.dataManager.GetDataDir()

	// Look for data files
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return // Use current date if we can't read the directory
	}

	var latestTime time.Time
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".txt") {
			continue
		}
		filename := strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".txt")
		if len(filename) != 6 {
			continue
		}
		year, err1 := strconv.Atoi(filename[:4])
		month, err2 := strconv.Atoi(filename[4:])
		if err1 != nil || err2 != nil || month < 1 || month > 12 {
			continue
		}

		// Load todos for this month
		todos, err := mw.dataManager.GetTodosForMonth(year, month)
		if err != nil {
			continue
		}

		// Find the latest todo time in this month
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

// setupUI initializes the user interface
func (mw *MainWindow) setupUI() {
	// Set window properties - matching mockup dimensions
	mw.window.SetTitle(localization.GetString("window_title"))
	mw.window.Resize(fyne.NewSize(420, 800))
	mw.window.SetFixedSize(true)

	// Create UI components
	mw.titleLabel = widget.NewLabel(localization.GetString("window_title") + " - User")
	mw.titleLabel.Alignment = fyne.TextAlignCenter
	mw.titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	// Disable wrapping to prevent the label from changing size
	mw.titleLabel.Wrapping = fyne.TextTruncate

	// Create buttons with icons (for dialogs etc.)
	mw.addButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), mw.onAddButtonClicked)
	mw.addButton.Importance = widget.HighImportance

	// Legacy hidden controls for compatibility
	mw.prevButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), mw.onPrevDayClicked)
	mw.prevButton.Hide()
	mw.nextButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), mw.onNextDayClicked)
	mw.nextButton.Hide()
	mw.viewModeBtn = widget.NewButton(mw.viewMode.GetLabel(), mw.onViewModeClicked)
	mw.viewModeBtn.Hide()
	mw.themeBtn = widget.NewButton("Gruvbox", mw.onThemeToggleClicked)
	mw.themeBtn.Hide()

	// --- Get gradient colors from current theme (used as full-window background) ---
	// Get gradient colors from current theme
	var startColor, endColor color.Color
	currentTheme := fyne.CurrentApp().Settings().Theme()
	if lightTheme, ok := currentTheme.(*LightSoftTheme); ok {
		startColor, endColor = lightTheme.GetHeaderGradientColors()
	} else if gruvboxTheme, ok := currentTheme.(*GruvboxBlackTheme); ok {
		startColor, endColor = gruvboxTheme.GetHeaderGradientColors()
	} else {
		// Fallback to primary color if theme doesn't support gradients
		startColor = theme.Color(theme.ColorNamePrimary)
		endColor = startColor
	}

	// Themed header icon from embedded assets and bottom-aligned with title
	var titleColor color.Color
	var logoRes fyne.Resource
	if _, ok := currentTheme.(*LightSoftTheme); ok {
		titleColor = color.White
		logoRes = assets.HeaderIconLight
	} else {
		titleColor = hex("#fabd2f")
		logoRes = assets.HeaderIconDark
	}
	logoImg := canvas.NewImageFromResource(logoRes)
	logoImg.FillMode = canvas.ImageFillContain
	// Icon size reduced to 54px (was 162px - 3x smaller)
	logoSize := float32(80)
	logoImg.SetMinSize(fyne.NewSize(logoSize, logoSize))

	// Title text - 56px from mockup
	titleTxt := canvas.NewText("GO DO", titleColor)
	titleTxt.TextSize = 64                           // From mockup
	titleTxt.TextStyle = fyne.TextStyle{Bold: false} // Weight 300 in mockup = light, use normal

	// Bottom-align both icon and title using VBox with spacers
	// Calculate how much padding needed to push icon to bottom
	iconPad := titleTxt.TextSize - logoSize
	if iconPad < 0 {
		iconPad = 0
	}
	logoAligned := container.NewVBox(
		CreateSpacer(1, iconPad),
		container.NewMax(logoImg),
	)

	header := container.NewHBox(logoAligned, CreateSpacer(5, 1), titleTxt)

	// --- Controls row: [Select] [←] [→] [Pomodoro] ---
	var navBg, navFg, pomodoroBg, pomodoroFg color.Color
	if _, ok := currentTheme.(*LightSoftTheme); ok {
		navBg = hex("#ff8c42")
		navFg = color.White
		pomodoroBg = hex("#ff8c42")
		pomodoroFg = color.White
	} else {
		navBg = hex("#504945")
		navFg = hex("#fabd2f")
		pomodoroBg = hex("#504945")
		pomodoroFg = hex("#fabd2f")
	}

	// Create Select widget with all view modes
	viewOptions := []string{"All", "Incomplete", "Complete", "Important"}
	mw.viewSelect = widget.NewSelect(viewOptions, func(selected string) {
		// Map selected string to ViewMode
		switch selected {
		case "All":
			mw.viewMode = models.ViewAll
		case "Incomplete":
			mw.viewMode = models.ViewIncomplete
		case "Complete":
			mw.viewMode = models.ViewComplete
		case "Important":
			mw.viewMode = models.ViewStarred
		}
		mw.loadTodos()
		mw.refreshView()
	})
	mw.viewSelect.SetSelected(mw.viewMode.GetLabel())

	// Wrap Select in styled container with white background for light theme
	var selectBg color.Color
	if _, ok := currentTheme.(*LightSoftTheme); ok {
		selectBg = color.White
	} else {
		selectBg = hex("#3c3836")
	}
	selectWrapper := CreateStyledSelect(mw.viewSelect, selectBg, fyne.NewSize(180, 44), 8)

	mw.prevRectBtn = NewSimpleRectButton("←", navBg, navFg, fyne.NewSize(44, 44), 8, mw.onPrevDayClicked)
	mw.nextRectBtn = NewSimpleRectButton("→", navBg, navFg, fyne.NewSize(44, 44), 8, mw.onNextDayClicked)

	// Pomodoro button
	mw.pomodoroRectBtn = NewSimpleRectButton("Pomodoro", pomodoroBg, pomodoroFg, fyne.NewSize(100, 44), 8, mw.onPomodoroTopClicked)

	controls := container.NewHBox(
		selectWrapper,
		CreateSpacer(10, 1),
		mw.prevRectBtn,
		CreateSpacer(10, 1),
		mw.nextRectBtn,
		CreateSpacer(10, 1),
		mw.pomodoroRectBtn,
	)

	// Set up timeline with current date and view mode
	mw.timeline.SetDate(mw.currentDate)
	mw.timeline.SetViewMode(mw.viewMode)

	// Create main content tasks container strictly per mockup
	timelineCard := CreateTasksContainer(mw.timeline)
	// Horizontal padding 24px for header and controls, top padding 30px
	headerPadded := container.NewBorder(nil, nil, CreateSpacer(24, 1), CreateSpacer(24, 1), header)
	controlsPadded := container.NewBorder(nil, nil, CreateSpacer(24, 1), CreateSpacer(24, 1), controls)
	timelinePadded := container.NewBorder(nil, nil, CreateSpacer(24, 1), CreateSpacer(24, 1), timelineCard)

	// Build header section (fixed at top)
	headerArea := container.NewVBox(
		CreateSpacer(1, 15), // Reduced from 30px to 15px (2x smaller)
		headerPadded,
		CreateSpacer(1, 30),
		controlsPadded,
		CreateSpacer(1, 30),
	)
	topSection := headerArea

	// Use Border layout to make timeline fill remaining space, with bottom margin for add button
	appBody := container.NewBorder(
		topSection,          // top: header + controls
		CreateSpacer(1, 24), // bottom: 24px margin (space for add button which floats)
		nil, nil,            // left, right
		timelinePadded, // center: timeline fills remaining vertical space
	)

	// Bottom buttons (add button in center, pomodoro on right)
	addButtonWithMargin := mw.setupBottomButtons()

	// Content stack - NO PADDING to avoid gaps
	content := container.NewBorder(
		nil,
		addButtonWithMargin, // Add button with margin at bottom
		nil, nil,
		appBody,
	)

	// Full-window background gradient
	background := NewGradientRect(startColor, endColor, 0)
	finalContent := container.NewMax(background, content)
	mw.window.SetContent(finalContent)
}

func (mw *MainWindow) onThemeToggleClicked() {
	mw.isGruvbox = !mw.isGruvbox
	if mw.isGruvbox {
		fyne.CurrentApp().Settings().SetTheme(NewGruvboxBlackTheme())
		mw.themeBtn.SetText("Light")
	} else {
		fyne.CurrentApp().Settings().SetTheme(NewLightSoftTheme())
		mw.themeBtn.SetText("Gruvbox")
	}
	// Force refresh the entire window to update header gradient
	mw.setupUI()
	mw.loadTodos()
	mw.refreshView()
}

// loadTodos loads todos for the current day
func (mw *MainWindow) loadTodos() {
	year, month := mw.currentDate.Year(), int(mw.currentDate.Month())

	// Load all todos for the month
	monthlyTodos, err := mw.dataManager.GetTodosForMonth(year, month)
	if err != nil {
		fmt.Println(localization.GetStringWithArgs("error_load_failed", err.Error()))
		mw.todos = []*models.TodoItem{}
		return
	}

	// Filter for the current day and view mode
	currentTime := time.Now()
	var dailyTodos []*models.TodoItem
	startOfDay := time.Date(mw.currentDate.Year(), mw.currentDate.Month(), mw.currentDate.Day(), 0, 0, 0, 0, mw.currentDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	for _, todo := range monthlyTodos {
		if todo.TodoTime.After(startOfDay) && todo.TodoTime.Before(endOfDay) {
			dailyTodos = append(dailyTodos, todo)
		}
	}
	mw.todos = mw.viewMode.FilterItems(dailyTodos, currentTime)
}

// refreshView updates the UI display
func (mw *MainWindow) refreshView() {
	// Update timeline data
	mw.timeline.SetDate(mw.currentDate) // Now passes full time.Time
	mw.timeline.SetViewMode(mw.viewMode)
	mw.timeline.SetTodos(mw.todos)
	mw.timeline.Refresh()
}

// Event handlers

func (mw *MainWindow) onAddButtonClicked() {
	mw.todoForm.ShowCreateDialog(func() {
		// Refresh the todo list after saving
		mw.loadTodos()
		mw.refreshView()
	})
}

func (mw *MainWindow) onPrevDayClicked() {
	mw.currentDate = mw.currentDate.AddDate(0, 0, -1)
	mw.loadTodos()
	mw.refreshView()
}

func (mw *MainWindow) onNextDayClicked() {
	mw.currentDate = mw.currentDate.AddDate(0, 0, 1)
	mw.loadTodos()
	mw.refreshView()
}

func (mw *MainWindow) onViewModeClicked() {
	// This method is now handled by the Select widget callback
	// Kept for backward compatibility with legacy viewModeBtn
	mw.viewMode = mw.viewMode.GetNextMode()
	if mw.viewSelect != nil {
		mw.viewSelect.SetSelected(mw.viewMode.GetLabel())
	}
	mw.loadTodos()
	mw.refreshView()
}

func (mw *MainWindow) onTodoSelected(todo *models.TodoItem, todoTime time.Time) {
	fmt.Printf("Selected todo: %s\n", todo.Name)

	// Open edit dialog
	mw.todoForm.ShowEditDialog(todo, todoTime, func() {
		// Refresh the todo list after saving
		mw.loadTodos()
		mw.refreshView()
	})
}

// setupBottomButtons creates the bottom button layout with add and theme buttons
func (mw *MainWindow) setupBottomButtons() fyne.CanvasObject {
	// Get theme colors
	currentTheme := fyne.CurrentApp().Settings().Theme()
	var themeBg, themeFg color.Color
	if _, ok := currentTheme.(*LightSoftTheme); ok {
		themeBg = hex("#ff8c42")
		themeFg = color.White
	} else {
		themeBg = hex("#504945")
		themeFg = hex("#fabd2f")
	}

	// Add button on the left-center (60x60px)
	addButtonRounded := RoundedIconButton(theme.ContentAddIcon(), mw.onAddButtonClicked)
	addWrap := container.NewGridWrap(fyne.NewSize(55, 55), addButtonRounded)

	// Theme toggle button on the right (with text)
	themeLabel := "Dark"
	if mw.isGruvbox {
		themeLabel = "Light"
	} else {
		themeLabel = "Dark"
	}

	// Create theme button as SimpleRectButton
	mw.themeRectBtn = NewSimpleRectButton(themeLabel, themeBg, themeFg, fyne.NewSize(100, 44), 8, mw.onThemeToggleClicked)
	themeWrap := container.NewVBox(
		CreateSpacer(1, 8), // Add small spacer to align with add button
		mw.themeRectBtn,
	)

	// Create bottom button layout: add on left-center, theme on right with padding
	bottomButtons := container.NewBorder(
		nil, nil,
		container.NewBorder(nil, nil, CreateSpacer(25, 1), nil, addWrap),   // 20px left margin
		container.NewBorder(nil, nil, nil, CreateSpacer(25, 1), themeWrap), // 20px right margin
		canvas.NewRectangle(color.Transparent),                             // center placeholder
	)

	// Place spacer BELOW the buttons to lift them up from the bottom edge
	return container.NewVBox(
		bottomButtons,
		CreateSpacer(1, 10),
	)
}

// onPomodoroTopClicked handles the top pomodoro button click
func (mw *MainWindow) onPomodoroTopClicked() {
	// Create and show pomodoro window
	pomodoroWindow := NewPomodoroWindow(fyne.CurrentApp(), mw.isGruvbox)
	pomodoroWindow.Show()
}
