package ui

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"

	"todo-list-migration/src/localization"
	"todo-list-migration/src/models"
	"todo-list-migration/src/persistence"
	"todo-list-migration/src/ui/forms"
	"todo-list-migration/src/utils"

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
	dropdownBtn  *SimpleRectButton
	prevRectBtn  *SimpleRectButton
	nextRectBtn  *SimpleRectButton
	themeRectBtn *SimpleRectButton

	// State
	currentDate *utils.CustomDate
	viewMode    models.ViewMode
	todos       []*models.TodoItem
	isGruvbox   bool
}

// NewMainWindow creates a new main window
func NewMainWindow(window fyne.Window, dataDir string) *MainWindow {
	mw := &MainWindow{
		window:      window,
		dataManager: persistence.NewMonthlyManager(dataDir),
		currentDate: utils.GetCurrentDate(),
		viewMode:    models.ViewIncomplete,
		isGruvbox:   false,
	}

	// Try to find existing data file and set currentDate accordingly
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

// findAndSetCurrentDateFromDataFile looks for existing data files and sets currentDate to match
func (mw *MainWindow) findAndSetCurrentDateFromDataFile() {
	dataDir := mw.dataManager.GetDataDir()

	// Look for data files in the format YYYYMM.yaml (preferred) or YYYYMM.txt
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return // Use current date if we can't read the directory
	}

	var latestFile string
	var latestYear, latestMonth int
	// first pass gather YAML months
	yamlMonths := make(map[string]struct{})
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".yaml") {
			filename := strings.TrimSuffix(entry.Name(), ".yaml")
			if len(filename) == 6 {
				if _, err := strconv.Atoi(filename); err == nil {
					yamlMonths[filename] = struct{}{}
				}
			}
		}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		isYAML := strings.HasSuffix(name, ".yaml")
		isTXT := strings.HasSuffix(name, ".txt")
		if !isYAML && !isTXT {
			continue
		}
		var filename string
		if isYAML {
			filename = strings.TrimSuffix(name, ".yaml")
		} else {
			filename = strings.TrimSuffix(name, ".txt")
		}
		if len(filename) != 6 {
			continue
		}
		yearStr := filename[:4]
		monthStr := filename[4:]
		year, err1 := strconv.Atoi(yearStr)
		month, err2 := strconv.Atoi(monthStr)
		if err1 != nil || err2 != nil || month < 1 || month > 12 {
			continue
		}
		// prefer YAML when both exist
		if _, hasYAML := yamlMonths[filename]; !isYAML && hasYAML {
			continue
		}
		if latestFile == "" || year > latestYear || (year == latestYear && month > latestMonth) {
			latestFile = name
			latestYear = year
			latestMonth = month
		}
	}

	// If we found a valid data file, set currentDate to match
	if latestFile != "" {
		mw.currentDate = utils.NewCustomDateFromValues(latestYear, latestMonth, 1, 0, 0, 0)
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
	mw.prevButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), mw.onPrevMonthClicked)
	mw.prevButton.Hide()
	mw.nextButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), mw.onNextMonthClicked)
	mw.nextButton.Hide()
	mw.viewModeBtn = widget.NewButton(mw.viewMode.GetLabel(), mw.onViewModeClicked)
	mw.viewModeBtn.Hide()
	mw.themeBtn = widget.NewButton("Gruvbox", mw.onThemeToggleClicked)
	mw.themeBtn.Hide()

	// --- Background gradient for the whole app area + header/logo/title ---
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

	// Logo: circular outline with checkmark - 54x54px from mockup
	logoCircle := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	// stroke/text color depends on theme
	var logoStroke color.Color
	var logoTextColor color.Color
	if _, ok := currentTheme.(*LightSoftTheme); ok {
		logoStroke = color.White
		logoTextColor = color.White
	} else {
		// Gruvbox: use bright yellow accent for both
		logoStroke = hex("#fabd2f")
		logoTextColor = hex("#fabd2f")
	}
	logoCircle.StrokeColor = logoStroke
	logoCircle.StrokeWidth = 3
	logoTxt := canvas.NewText("✓", logoTextColor)
	logoTxt.TextSize = 26 // Increased for 54px circle
	logoTxt.TextStyle = fyne.TextStyle{Bold: true}
	logo := container.NewGridWrap(fyne.NewSize(54, 54), container.NewMax(logoCircle, container.NewCenter(logoTxt)))

	// Title text - 56px from mockup with letter-spacing (simulated via text size)
	titleTxt := canvas.NewText("GO DO", color.White)
	titleTxt.TextSize = 56                           // From mockup
	titleTxt.TextStyle = fyne.TextStyle{Bold: false} // Weight 300 in mockup = light, use normal

	header := container.NewHBox(logo, CreateSpacer(20, 1), titleTxt)

	// --- Controls row: [Dropdown (cycle)] [←] [→] [Theme] --- strictly per mockup
	var dropdownBg, dropdownFg, navBg, navFg, themeBg, themeFg color.Color
	if _, ok := currentTheme.(*LightSoftTheme); ok {
		dropdownBg = hex("#ffffff")
		dropdownFg = hex("#3c3836")
		navBg = hex("#ff8c42")
		navFg = color.White
		themeBg = hex("#ff8c42")
		themeFg = color.White
	} else {
		dropdownBg = hex("#3c3836")
		dropdownFg = hex("#ebdbb2")
		navBg = hex("#504945")
		navFg = hex("#fabd2f")
		themeBg = hex("#504945")
		themeFg = hex("#fabd2f")
	}

	// Dropdown cycles through modes on tap
	mw.dropdownBtn = NewSimpleRectButton(mw.viewMode.GetLabel()+" ▼", dropdownBg, dropdownFg, fyne.NewSize(180, 44), 8, func() {
		mw.onViewModeClicked()
		mw.dropdownBtn.SetText(mw.viewMode.GetLabel() + " ▼")
	})

	mw.prevRectBtn = NewSimpleRectButton("←", navBg, navFg, fyne.NewSize(44, 44), 8, mw.onPrevMonthClicked)
	mw.nextRectBtn = NewSimpleRectButton("→", navBg, navFg, fyne.NewSize(44, 44), 8, mw.onNextMonthClicked)

	// Theme toggle shows current target
	themeLabel := "Dark"
	if mw.isGruvbox {
		themeLabel = "Light"
	} else {
		// current is light by default, button shows Dark
		themeLabel = "Dark"
	}
	mw.themeRectBtn = NewSimpleRectButton(themeLabel, themeBg, themeFg, fyne.NewSize(100, 44), 8, mw.onThemeToggleClicked)

	controls := container.NewHBox(
		mw.dropdownBtn,
		CreateSpacer(10, 1),
		mw.prevRectBtn,
		CreateSpacer(10, 1),
		mw.nextRectBtn,
		CreateSpacer(10, 1),
		mw.themeRectBtn,
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

	appBody := container.NewVBox(
		CreateSpacer(1, 30),
		headerPadded,
		CreateSpacer(1, 14),
		controlsPadded,
		CreateSpacer(1, 20),
		timelinePadded,
	)

	// Floating add button at bottom-center (60x60px from mockup)
	addButtonRounded := RoundedIconButton(theme.ContentAddIcon(), mw.onAddButtonClicked)
	addWrap := container.NewGridWrap(fyne.NewSize(60, 60), addButtonRounded)

	// Full-window background gradient per mockup
	background := NewGradientRect(startColor, endColor, 0)
	// Stack gradient behind content
	content := container.NewBorder(
		nil,
		container.NewCenter(addWrap), // Centered at bottom per mockup
		nil, nil,
		container.NewPadded(appBody),
	)
	mw.window.SetContent(container.NewMax(background, content))
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

// loadTodos loads todos for the current month
func (mw *MainWindow) loadTodos() {
	todos, err := mw.dataManager.GetTodosForMonth(mw.currentDate.Year, mw.currentDate.Month)
	if err != nil {
		// For now, just log the error. In a real app, you'd show a user-friendly message
		msg := localization.GetStringWithArgs("error_load_failed", err.Error())
		fmt.Println(msg)
		todos = []*models.TodoItem{}
	}

	// Filter based on current view mode
	currentTime := time.Now()
	mw.todos = mw.viewMode.FilterItems(todos, currentTime)
}

// refreshView updates the UI display
func (mw *MainWindow) refreshView() {
	// Update title
	titleText := fmt.Sprintf("%s - %d/%02d - %s",
		localization.GetString("window_title"), mw.currentDate.Year, mw.currentDate.Month, mw.viewMode.GetLabel())
	mw.titleLabel.SetText(titleText)

	// Update view mode button
	mw.viewModeBtn.SetText(mw.viewMode.GetLabel())

	// Update timeline data (without triggering refresh yet)
	mw.timeline.SetDate(mw.currentDate)
	mw.timeline.SetViewMode(mw.viewMode)
	mw.timeline.SetTodos(mw.todos)

	// Single refresh at the end to avoid multiple layout recalculations
	mw.timeline.Refresh()

	// DO NOT call content.Refresh() here! It triggers async layout recalculation
	// that causes the UI to shrink. Timeline already refreshed itself above.
}

// Event handlers

func (mw *MainWindow) onAddButtonClicked() {
	mw.todoForm.ShowCreateDialog(func() {
		// Refresh the todo list after saving
		mw.loadTodos()
		mw.refreshView()
	})
}

func (mw *MainWindow) onPrevMonthClicked() {
	mw.currentDate.ToLastMonth()
	mw.loadTodos()
	mw.refreshView()
}

func (mw *MainWindow) onNextMonthClicked() {
	mw.currentDate.ToNextMonth()
	mw.loadTodos()
	mw.refreshView()
}

func (mw *MainWindow) onViewModeClicked() {
	mw.viewMode = mw.viewMode.GetNextMode()
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
