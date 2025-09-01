package app

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetDataDirectory returns the data directory path, creating it if it doesn't exist.
// The data directory is located in the same directory as the executable.
func GetDataDirectory() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	dataDir := filepath.Join(filepath.Dir(execPath), "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create data directory: %w", err)
	}

	return dataDir, nil
}
