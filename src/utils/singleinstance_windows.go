//go:build windows
// +build windows

package utils

import (
	"fmt"
	"os"
	"syscall"
)

// isProcessRunning reports whether PID is active on Windows
func isProcessRunning(pid int) bool {
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))

	if err != nil {
		return false
	}

	defer syscall.CloseHandle(handle)

	var exitCode uint32
	err = syscall.GetExitCodeProcess(handle, &exitCode)

	if err != nil {
		return false
	}

	// Exit code 259 значит, что процесс еще жив
	const STILL_ACTIVE = 259

	return exitCode == STILL_ACTIVE
}

// isLockStale reports whether lock file points to inactive process
func (si *SingleInstance) isLockStale() bool {
	data, err := os.ReadFile(si.lockPath)

	if err != nil {
		return true
	}

	var pid int
	_, err = fmt.Sscanf(string(data), "%d", &pid)

	if err != nil || pid <= 0 {
		return true
	}

	return !isProcessRunning(pid)
}
