package main

import (
	"log"
	"os"
	"path/filepath"

	assets "godo/doc"
	"godo/src/persistence"
	"godo/src/ui"
	"godo/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

// showAlreadyRunningDialog displays an error message when another instance is detected
func showAlreadyRunningDialog() {
	// Create a minimal app just to show the dialog
	tempApp := app.New()
	tempWindow := tempApp.NewWindow("Todo List - Error")
	// Increase window size to accommodate dialog content and button
	tempWindow.Resize(fyne.NewSize(480, 280))
	tempWindow.CenterOnScreen()

	// Show information dialog with clear message and callback to close window
	errorDialog := dialog.NewInformation(
		"Application Already Running",
		"Todo List is already running.\n\nOnly one instance of the application can run at a time to prevent data conflicts.\n\nPlease check your taskbar or system tray.",
		tempWindow,
	)

	// Set callback to quit the app when dialog is dismissed (OK button clicked)
	errorDialog.SetOnClosed(func() {
		tempApp.Quit()
	})

	errorDialog.Show()
	tempWindow.ShowAndRun()
}

func main() {
	// Check for single instance
	instanceLock := utils.NewSingleInstance("todo-list-app")
	locked, err := instanceLock.TryLock()
	if err != nil {
		log.Fatal("Failed to check instance lock:", err)
	}
	if !locked {
		// Another instance is already running
		showAlreadyRunningDialog()
		return
	}
	// Ensure lock is released on exit
	defer instanceLock.Unlock()

	// Create the application
	myApp := app.New()
	// Default theme: Soft Light (per mockup)
	myApp.Settings().SetTheme(ui.NewLightSoftTheme())
	// Set app icon (used for window title bar/taskbar)
	if assets.AppIcon != nil {
		myApp.SetIcon(assets.AppIcon)
	}
	myWindow := myApp.NewWindow("My Day - Todo List")
	if assets.AppIcon != nil {
		myWindow.SetIcon(assets.AppIcon)
	}
	myWindow.SetMaster()

	// Get the data directory (create if doesn't exist)
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path:", err)
	}

	dataDir := filepath.Join(filepath.Dir(execPath), "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// Run one-shot migration from TXT to YAML on startup (non-fatal)
	migrator := persistence.NewMonthlyManager(dataDir)
	_ = migrator.MigrateAllToYAML()

	// Create the main UI
	ui.NewMainWindow(myWindow, dataDir)

	// Show the window and run the application
	myWindow.ShowAndRun()
}
