package app

import (
	"fmt"

	assets "godo/resources"
	"godo/src/persistence"
	"godo/src/ui"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
)

const mainWindowTitle = "My Day - Todo List"

// Application wires app startup
type Application struct {
	fyneApp fyne.App
	window  fyne.Window
	dataDir string
}

// New creates app shell
func New() (*Application, error) {
	dataDir, err := GetDataDirectory()

	if err != nil {
		return nil, fmt.Errorf("failed to initialize data directory: %w", err)
	}

	return &Application{
		dataDir: dataDir,
	}, nil
}

// Initialize prepares Fyne app and main window
func (a *Application) Initialize() error {
	a.fyneApp = fyneApp.New()

	a.fyneApp.Settings().SetTheme(ui.NewLightSoftTheme())

	a.setAppIcon()

	a.window = a.fyneApp.NewWindow(mainWindowTitle)
	a.setWindowIcon()
	a.window.SetMaster()

	return nil
}

// setAppIcon applies app icon
func (a *Application) setAppIcon() {
	if assets.AppIcon != nil {
		a.fyneApp.SetIcon(assets.AppIcon)
	}
}

// setWindowIcon applies window icon
func (a *Application) setWindowIcon() {
	if assets.AppIcon != nil {
		a.window.SetIcon(assets.AppIcon)
	}
}

// RunMigration migrates old todo files
func (a *Application) RunMigration() error {
	migrator := persistence.NewMonthlyManager(a.dataDir)
	return migrator.MigrateAllToYAML()
}

// CreateMainUI builds main UI
func (a *Application) CreateMainUI() {
	dataManager := persistence.NewMonthlyManager(a.dataDir)
	configManager := persistence.NewConfigManager(a.dataDir)

	ui.NewMainWindow(a.window, dataManager, configManager)
}

// Run starts app loop
func (a *Application) Run() {
	a.window.ShowAndRun()
}
