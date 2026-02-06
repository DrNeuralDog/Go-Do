package app

import (
	"fmt"

	assets "godo/resources"
	"godo/src/persistence"
	"godo/src/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// Application represents the main todo list application
type Application struct {
	fyneApp fyne.App
	window  fyne.Window
	dataDir string
}

// New creates a new Application instance
func New() (*Application, error) {
	// Get the data directory
	dataDir, err := GetDataDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize data directory: %w", err)
	}

	return &Application{
		dataDir: dataDir,
	}, nil
}

// Initialize sets up the Fyne application and window
func (a *Application) Initialize() error {
	// Create the Fyne application
	a.fyneApp = app.New()

	// Set default theme: Soft Light (per mockup)
	a.fyneApp.Settings().SetTheme(ui.NewLightSoftTheme())

	// Set app icon (used for window title bar/taskbar)
	a.setAppIcon()

	// Create main window
	a.window = a.fyneApp.NewWindow("My Day - Todo List")
	a.setWindowIcon()
	a.window.SetMaster()

	return nil
}

// setAppIcon sets the application icon if available
func (a *Application) setAppIcon() {
	if assets.AppIcon != nil {
		a.fyneApp.SetIcon(assets.AppIcon)
	}
}

// setWindowIcon sets the window icon if available
func (a *Application) setWindowIcon() {
	if assets.AppIcon != nil {
		a.window.SetIcon(assets.AppIcon)
	}
}

// RunMigration runs the one-shot migration from TXT to YAML format
func (a *Application) RunMigration() error {
	migrator := persistence.NewMonthlyManager(a.dataDir)
	if err := migrator.MigrateAllToYAML(); err != nil {
		// Migration is non-fatal, just log the error
		fmt.Printf("Warning: migration failed: %v\n", err)
	}
	return nil
}

// CreateMainUI creates and initializes the main user interface
func (a *Application) CreateMainUI() {
	dataManager := persistence.NewMonthlyManager(a.dataDir)
	configManager := persistence.NewConfigManager(a.dataDir)
	ui.NewMainWindow(a.window, dataManager, configManager)
}

// Run starts the application event loop
func (a *Application) Run() {
	a.window.ShowAndRun()
}
