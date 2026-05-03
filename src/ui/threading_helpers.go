package ui

import (
	"godo/src/ui/threading"
)

// runOnMainThread keeps legacy ui calls on Fyne main thread
func runOnMainThread(fn func()) {
	threading.RunOnMainThread(fn)
}
