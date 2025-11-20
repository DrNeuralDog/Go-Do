package helpers

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"godo/src/ui/threading"
)

// FlashWindow creates a visual flash effect on a window to indicate it's already open
// The window will flash 3 times over 600ms total
func FlashWindow(win fyne.Window) {
	if win == nil {
		return
	}

	// Get the window's current content
	content := win.Content()
	if content == nil {
		return
	}

	// Create a semi-transparent white overlay for the flash effect
	overlay := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 0})

	// Stack the overlay on top of existing content
	flashContent := container.NewStack(content, overlay)
	win.SetContent(flashContent)

	// Animate the flash: 3 quick pulses
	flashCount := 3
	flashDuration := 100 * time.Millisecond

	for i := 0; i < flashCount; i++ {
		i := i // capture for closure

		// Flash on
		time.AfterFunc(time.Duration(i*2)*flashDuration, func() {
			threading.RunOnMainThread(func() {
				overlay.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: 60}
				overlay.Refresh()
			})
		})

		// Flash off
		time.AfterFunc(time.Duration(i*2+1)*flashDuration, func() {
			threading.RunOnMainThread(func() {
				overlay.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: 0}
				overlay.Refresh()
			})
		})
	}

	// Remove overlay after animation completes
	time.AfterFunc(time.Duration(flashCount*2)*flashDuration, func() {
		threading.RunOnMainThread(func() {
			win.SetContent(content)
		})
	})

	// Also bring window to front
	win.RequestFocus()
}
