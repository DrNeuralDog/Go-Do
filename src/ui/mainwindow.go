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
	viewModeBtn *widget.Button
	themeBtn    *widget.Button
	viewModeSel *widget.Select

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
		viewMode:    models.ViewAll,
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
	// Set window properties
	mw.window.SetTitle(localization.GetString("window_title"))
	mw.window.Resize(fyne.NewSize(420, 600))
	mw.window.SetFixedSize(true)

	// Create UI components
	mw.titleLabel = widget.NewLabel(localization.GetString("window_title") + " - User")
	mw.titleLabel.Alignment = fyne.TextAlignCenter
	mw.titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	// Disable wrapping to prevent the label from changing size
	mw.titleLabel.Wrapping = fyne.TextTruncate

	mw.addButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), mw.onAddButtonClicked)
	mw.addButton.Importance = widget.HighImportance

	mw.prevButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), mw.onPrevMonthClicked)
	mw.prevButton.Importance = widget.HighImportance
	mw.nextButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), mw.onNextMonthClicked)
	mw.nextButton.Importance = widget.HighImportance

	// Legacy button kept for compatibility with refreshView; hidden from UI
	mw.viewModeBtn = widget.NewButton(mw.viewMode.GetLabel(), mw.onViewModeClicked)
	mw.viewModeBtn.Hide()

	// View mode select (combobox)
	mw.viewModeSel = widget.NewSelect([]string{
		models.ViewAll.GetLabel(),
		models.ViewIncomplete.GetLabel(),
		models.ViewComplete.GetLabel(),
		models.ViewStarred.GetLabel(),
	}, func(value string) {
		switch value {
		case models.ViewAll.GetLabel():
			mw.viewMode = models.ViewAll
		case models.ViewIncomplete.GetLabel():
			mw.viewMode = models.ViewIncomplete
		case models.ViewComplete.GetLabel():
			mw.viewMode = models.ViewComplete
		case models.ViewStarred.GetLabel():
			mw.viewMode = models.ViewStarred
		}
		mw.loadTodos()
		mw.refreshView()
	})
	mw.viewModeSel.SetSelected(models.ViewAll.GetLabel())

	mw.themeBtn = widget.NewButton("Gruvbox", mw.onThemeToggleClicked)

	// Prepare fixed-width wrappers for symmetry
	// Use a larger fixed width to accommodate all view mode labels without resizing
	selSize := fyne.NewSize(150, mw.viewModeSel.MinSize().Height)
	btnSize := fyne.NewSize(140, mw.themeBtn.MinSize().Height)

	// Wrap Select in a Max container to prevent it from changing size
	selContainer := container.NewMax(mw.viewModeSel)
	selWrap := container.NewGridWrap(selSize, selContainer)
	btnWrap := container.NewGridWrap(btnSize, mw.themeBtn)
	spacer := func(w float32) *canvas.Rectangle {
		r := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		r.SetMinSize(fyne.NewSize(w, 1))
		return r
	}
	chip := func(obj fyne.CanvasObject) fyne.CanvasObject {
		bg := canvas.NewRectangle(theme.Color(theme.ColorNameHover))
		return container.NewMax(bg, container.NewPadded(obj))
	}

	// Order: [Select] [sp] [Prev] | [Title] | [Next] [sp] [Theme]
	leftGroup := container.NewHBox(chip(selWrap), spacer(16), chip(mw.prevButton))
	rightGroup := container.NewHBox(chip(mw.nextButton), spacer(16), chip(btnWrap))

	// Make theme toggle more prominent
	mw.themeBtn.Importance = widget.HighImportance

	// Create full-width top bar with centered title and controls left/right (no background)
	topBar := container.NewBorder(
		nil, nil,
		leftGroup,
		rightGroup,
		container.NewCenter(mw.titleLabel),
	)
	// Add separator under the top bar
	header := container.NewVBox(topBar, widget.NewSeparator())

	// Set up timeline with current date and view mode
	mw.timeline.SetDate(mw.currentDate)
	mw.timeline.SetViewMode(mw.viewMode)

	// Create main content (single scroll managed inside Timeline widget)
	// Thin border around the panel
	lineColor := theme.Color(theme.ColorNameSeparator)
	topLine := canvas.NewRectangle(lineColor)
	topLine.SetMinSize(fyne.NewSize(1, 1))
	bottomLine := canvas.NewRectangle(lineColor)
	bottomLine.SetMinSize(fyne.NewSize(1, 1))
	leftLine := canvas.NewRectangle(lineColor)
	leftLine.SetMinSize(fyne.NewSize(1, 1))
	rightLine := canvas.NewRectangle(lineColor)
	rightLine.SetMinSize(fyne.NewSize(1, 1))
	// No background here; background will be applied only under items within timeline
	panelWithBorder := container.NewBorder(topLine, bottomLine, leftLine, rightLine, mw.timeline)
	contentCenter := container.NewPadded(panelWithBorder)

	// Square + button at bottom (reduced by ~15%)
	addWrap := container.NewGridWrap(fyne.NewSize(34, 34), mw.addButton)

	content := container.NewBorder(
		header,                     // Top
		container.NewHBox(addWrap), // Bottom
		nil, nil,                   // Left, Right
		contentCenter, // Center
	)

	mw.window.SetContent(content)
}

func (mw *MainWindow) onThemeToggleClicked() {
	mw.isGruvbox = !mw.isGruvbox
	if mw.isGruvbox {
		fyne.CurrentApp().Settings().SetTheme(NewGruvboxBlackTheme())
		mw.themeBtn.SetText("Light")
	} else {
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
		mw.themeBtn.SetText("Gruvbox")
	}
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
