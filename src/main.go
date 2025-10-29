package main

import (
	"log"
	"os"
	"path/filepath"

	"todo-list-migration/src/persistence"
	"todo-list-migration/src/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	// Create the application
	myApp := app.New()
	// Default theme: Soft Light (per mockup)
	myApp.Settings().SetTheme(ui.NewLightSoftTheme())
	myWindow := myApp.NewWindow("My Day - Todo List")
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
