package app

import (
	"godo/src/utils"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

const (
	appLockName = "todo-list-app"

	alreadyRunningTitle   = "Application Already Running"
	alreadyRunningMessage = "Todo List is already running.\n\nOnly one instance of the application can run at a time to prevent data conflicts.\n\nPlease check your taskbar or system tray."
)

// CheckSingleInstance checks app process lock
func CheckSingleInstance() (*utils.SingleInstance, bool) {
	instanceLock := utils.NewSingleInstance(appLockName)
	locked, err := instanceLock.TryLock()

	if err != nil {
		return nil, false
	}

	if !locked {
		showAlreadyRunningDialog()

		return nil, false
	}

	return instanceLock, true
}

// showAlreadyRunningDialog shows duplicate launch warning
func showAlreadyRunningDialog() {
	// Основное приложение еще не создано, так что поднимаем маленькое окно под диалог
	tempApp := fyneApp.New()
	tempWindow := tempApp.NewWindow("Todo List - Error")

	tempWindow.Resize(fyne.NewSize(480, 280))
	tempWindow.CenterOnScreen()

	errorDialog := dialog.NewInformation(
		alreadyRunningTitle,
		alreadyRunningMessage,
		tempWindow,
	)

	errorDialog.SetOnClosed(func() {
		tempApp.Quit()
	})

	errorDialog.Show()
	tempWindow.ShowAndRun()
}
