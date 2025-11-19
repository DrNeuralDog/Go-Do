// go:build windows
// +build windows

package utils

import (
	"fmt"
	"os"
	"syscall"
)

// isProcessRunning checks if a process with the given PID is running on Windows
func isProcessRunning(pid int) bool {
	// Try to open the process handle
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		// Process doesn't exist or we don't have permission
		return false
	}
	defer syscall.CloseHandle(handle)

	// Get the exit code
	var exitCode uint32
	err = syscall.GetExitCodeProcess(handle, &exitCode)
	if err != nil {
		return false
	}

	// STILL_ACTIVE (259) means the process is still running
	const STILL_ACTIVE = 259
	return exitCode == STILL_ACTIVE
}

// isLockStale checks if the lock file belongs to a dead process (Windows version)
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

	// Check if process is running using Windows-specific method
	return !isProcessRunning(pid)
}
