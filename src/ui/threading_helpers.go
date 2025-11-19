package ui

import "fyne.io/fyne/v2"

// runOnMainThread ensures the provided function executes on the Fyne UI thread.
// Fyne widgets must only be touched from the main thread, so goroutines should
// schedule their UI work through this helper to avoid Do/DoAndWait panics.
func runOnMainThread(fn func()) {
	if fn == nil {
		return
	}

	app := fyne.CurrentApp()
	if app == nil {
		fn()
		return
	}

	driver := app.Driver()
	if driver == nil {
		fn()
		return
	}

	// Different Fyne driver versions expose either RunOnMain or CallOnMain;
	// try both before falling back to executing inline.
	if runner, ok := driver.(interface{ RunOnMain(func()) }); ok {
		runner.RunOnMain(fn)
		return
	}
	if caller, ok := driver.(interface{ CallOnMain(func()) }); ok {
		caller.CallOnMain(fn)
		return
	}

	fn()
}
