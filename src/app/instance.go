package app

import (
	"godo/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

// CheckSingleInstance checks if another instance of the application is running.
// Returns the instance lock and a boolean indicating if the lock was acquired.
// If the lock cannot be acquired, it displays an error dialog and returns nil, false.
func CheckSingleInstance() (*utils.SingleInstance, bool) {
	instanceLock := utils.NewSingleInstance("todo-list-app")
	locked, err := instanceLock.TryLock()
	if err != nil {
		// If we can't check the lock, fail fast
		return nil, false
	}
	if !locked {
		// Another instance is already running - show dialog
		showAlreadyRunningDialog()
		return nil, false
	}
	return instanceLock, true
}

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
