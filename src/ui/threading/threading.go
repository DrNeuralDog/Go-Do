package threading

import (
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
)

const mainThreadStackSize = 1024

// RunOnMainThread runs fn on Fyne UI thread
func RunOnMainThread(fn func()) {
	if fn == nil {
		return
	}

	if fyne.CurrentApp() == nil || isMainThread() {
		fn()

		return
	}

	fyne.Do(fn)
}

// isMainThread guesses whether current goroutine is main UI goroutine
func isMainThread() bool {
	buf := make([]byte, mainThreadStackSize)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	// Это не стопроцентно - но для текущей версии реализации UI refresh хватает
	hasMainFunc := strings.Contains(stack, "main.main")
	isFirstGoroutine := strings.HasPrefix(stack, "goroutine 1 ")

	return hasMainFunc || isFirstGoroutine
}
