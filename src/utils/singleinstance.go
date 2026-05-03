package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// SingleInstance owns process lock file
type SingleInstance struct {
	lockFile *os.File

	lockPath string
}

// NewSingleInstance creates lock holder for app name
func NewSingleInstance(appName string) *SingleInstance {
	tempDir := os.TempDir()
	lockPath := filepath.Join(tempDir, fmt.Sprintf("%s.lock", appName))

	return &SingleInstance{
		lockPath: lockPath,
	}
}

// TryLock tries to acquire process lock
func (si *SingleInstance) TryLock() (bool, error) {
	// O_EXCL дает атомарный захват lock-файла
	file, err := os.OpenFile(si.lockPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)

	if err != nil {
		if os.IsExist(err) {
			// Протух - lock можно удалить и попробовать снова
			if si.isLockStale() {
				if err := os.Remove(si.lockPath); err != nil {
					return false, fmt.Errorf("failed to remove stale lock file: %w", err)
				}

				return si.TryLock()
			}

			return false, nil
		}

		return false, fmt.Errorf("failed to create lock file: %w", err)
	}

	// PID в lock-файле помогает понять, кто держит запуск
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

// Unlock releases process lock
func (si *SingleInstance) Unlock() error {
	if si.lockFile != nil {
		err := si.lockFile.Close()

		if err != nil {
			return fmt.Errorf("failed to close lock file: %w", err)
		}

		si.lockFile = nil
	}

	if _, err := os.Stat(si.lockPath); err == nil {
		err = os.Remove(si.lockPath)

		if err != nil {
			return fmt.Errorf("failed to remove lock file: %w", err)
		}
	}

	return nil
}

// GetLockPath returns lock file path
func (si *SingleInstance) GetLockPath() string {
	return si.lockPath
}
