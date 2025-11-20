package ui

import (
	"godo/src/ui/threading"
)

// runOnMainThread is a deprecated wrapper for threading.RunOnMainThread.
// This is kept for backward compatibility within the ui package.
// New code should use threading.RunOnMainThread directly.
func runOnMainThread(fn func()) {
	threading.RunOnMainThread(fn)
}
