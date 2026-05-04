package app

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetDataDirectory returns data dir near executable
func GetDataDirectory() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Данные держу рядом с exe так проще таскать приложение целиком
	dataDir := filepath.Join(filepath.Dir(execPath), "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create data directory: %w", err)
	}

	return dataDir, nil
}
