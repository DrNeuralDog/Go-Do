package main

import (
	"log"

	"godo/src/app"
)

func main() {
	// Check for single instance
	instanceLock, locked := app.CheckSingleInstance()
	if !locked {
		// Another instance is already running or check failed
		return
	}
	// Ensure lock is released on exit
	defer instanceLock.Unlock()

	// Create and initialize the application
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Initialize(); err != nil {
		log.Fatal(err)
	}

	// Run one-shot migration from TXT to YAML on startup (non-fatal)
	if err := application.RunMigration(); err != nil {
		log.Printf("Warning: Migration failed: %v", err)
	}

	// Create the main UI
	application.CreateMainUI()

	// Show the window and run the application
	application.Run()
}
