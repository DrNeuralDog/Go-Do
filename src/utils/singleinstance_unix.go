// go:build !windows
// +build !windows

package utils

import (
	"fmt"
	"os"
	"syscall"
)

// isLockStale checks if the lock file belongs to a dead process (Unix version)
func (si *SingleInstance) isLockStale() bool {
	// Read PID from lock file
	data, err := os.ReadFile(si.lockPath)
	if err != nil {
		return true // Can't read, consider it stale
	}

	var pid int
	_, err = fmt.Sscanf(string(data), "%d", &pid)
	if err != nil {
		return true // Invalid PID format, consider it stale
	}

	// Check if process exists using kill with signal 0
	// Signal 0 doesn't actually send a signal but checks if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		return true // Process doesn't exist
	}

	// Send signal 0 to check if process is alive
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return true // Process doesn't exist or we can't signal it
	}

	return false // Process exists, lock is valid
}
