package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// SingleInstance manages application instance locking
type SingleInstance struct {
	lockFile *os.File
	lockPath string
}

// NewSingleInstance creates a new single instance manager
func NewSingleInstance(appName string) *SingleInstance {
	// Use system temp directory for lock file
	tempDir := os.TempDir()
	lockPath := filepath.Join(tempDir, fmt.Sprintf("%s.lock", appName))

	return &SingleInstance{
		lockPath: lockPath,
	}
}

// TryLock attempts to acquire the instance lock
// Returns true if lock acquired successfully, false if another instance is running
func (si *SingleInstance) TryLock() (bool, error) {
	// Try to create lock file with exclusive access
	// O_CREATE | O_EXCL ensures atomic creation - fails if file exists
	// O_RDWR for read/write access
	file, err := os.OpenFile(si.lockPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)

	if err != nil {
		if os.IsExist(err) {
			// Lock file exists - check if it's stale
			if si.isLockStale() {
				// Remove stale lock and retry
				os.Remove(si.lockPath)
				return si.TryLock()
			}
			// Another instance is running
			return false, nil
		}
		// Other error occurred
		return false, fmt.Errorf("failed to create lock file: %w", err)
	}

	// Write PID to lock file for debugging
	pid := os.Getpid()
	_, err = file.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		file.Close()
		os.Remove(si.lockPath)
		return false, fmt.Errorf("failed to write PID to lock file: %w", err)
	}

	si.lockFile = file
	return true, nil
}

// isLockStale is implemented in platform-specific files:
// - singleinstance_windows.go for Windows
// - singleinstance_unix.go for Unix-like systems

// Unlock releases the instance lock
func (si *SingleInstance) Unlock() error {
	if si.lockFile != nil {
		// Close the file
		err := si.lockFile.Close()
		if err != nil {
			return fmt.Errorf("failed to close lock file: %w", err)
		}
		si.lockFile = nil
	}

	// Remove the lock file
	if _, err := os.Stat(si.lockPath); err == nil {
		err = os.Remove(si.lockPath)
		if err != nil {
			return fmt.Errorf("failed to remove lock file: %w", err)
		}
	}

	return nil
}

// GetLockPath returns the path to the lock file (for debugging)
func (si *SingleInstance) GetLockPath() string {
	return si.lockPath
}
