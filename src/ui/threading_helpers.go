package ui

import (
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
)

// runOnMainThread ensures the provided function executes on the Fyne UI thread.
// Fyne widgets must only be touched from the main thread, so goroutines should
// schedule their UI work through this helper to avoid Do/DoAndWait panics.
func runOnMainThread(fn func()) {
	if fn == nil {
		return
	}

	// Check if we're already on the main thread by examining the call stack
	// If we're on the main goroutine, just execute directly
	if isMainThread() {
		fn()
		return
	}

	// Otherwise, schedule on main thread without waiting (to avoid deadlock)
	fyne.Do(fn)
}

// isMainThread checks if we're currently on the main UI thread
func isMainThread() bool {
	// Get stack trace
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	// If stack contains "main.main" or we're in goroutine 1, we're on main thread
	// This is a heuristic but works for most cases
	return strings.Contains(stack, "main.main") || strings.HasPrefix(stack, "goroutine 1 ")
}
