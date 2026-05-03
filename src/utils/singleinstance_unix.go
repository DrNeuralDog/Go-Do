//go:build !windows
// +build !windows

package utils

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

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

	process, err := os.FindProcess(pid)

	if err != nil {
		return true
	}

	// Signal 0 только проверяет существование процесса
	err = process.Signal(syscall.Signal(0))

	if errors.Is(err, syscall.EPERM) {
		return false
	}

	if err != nil {
		return true
	}

	return false
}
