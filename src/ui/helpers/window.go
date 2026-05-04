package helpers

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"godo/src/ui/threading"
)

const (
	flashPulseCount    = 3
	flashPulseDuration = 100 * time.Millisecond
	flashOverlayAlpha  = 60
)

// FlashWindow briefly flashes already opened window
func FlashWindow(win fyne.Window) {
	threading.RunOnMainThread(func() {
		flashWindow(win)
	})
}

// flashWindow runs flash animation setup on UI thread
func flashWindow(win fyne.Window) {
	if win == nil {
		return
	}

	content := win.Content()

	if content == nil {
		return
	}

	overlay := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 0})
	flashContent := container.NewStack(content, overlay)

	win.SetContent(flashContent)
	win.RequestFocus()

	for i := 0; i < flashPulseCount; i++ {
		step := i

		time.AfterFunc(time.Duration(step*2)*flashPulseDuration, func() {
			setFlashAlpha(overlay, flashOverlayAlpha)
		})

		time.AfterFunc(time.Duration(step*2+1)*flashPulseDuration, func() {
			setFlashAlpha(overlay, 0)
		})
	}

	time.AfterFunc(time.Duration(flashPulseCount*2)*flashPulseDuration, func() {
		restoreFlashContent(win, flashContent, content)
	})
}

// setFlashAlpha updates flash overlay opacity
func setFlashAlpha(overlay *canvas.Rectangle, alpha uint8) {
	threading.RunOnMainThread(func() {
		overlay.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: alpha}

		overlay.Refresh()
	})
}

// restoreFlashContent removes flash overlay
func restoreFlashContent(win fyne.Window, flashContent, content fyne.CanvasObject) {
	threading.RunOnMainThread(func() {
		if win == nil {
			return
		}

		// Возвращаем content только если поверх него все еще текущий flash
		if win.Content() == flashContent {
			win.SetContent(content)
		}
	})
}
