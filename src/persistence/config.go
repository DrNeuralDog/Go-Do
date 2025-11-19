package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"godo/src/models"
)

// ConfigManager manages application configuration persistence
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a new config manager
func NewConfigManager(dataDir string) *ConfigManager {
	return &ConfigManager{
		configPath: filepath.Join(dataDir, "config.json"),
	}
}

// LoadConfig loads the configuration from disk
// Returns default config if file doesn't exist
func (cm *ConfigManager) LoadConfig() (*models.Config, error) {
	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return models.NewDefaultConfig(), nil
	}

	// Read file
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to disk
func (cm *ConfigManager) SaveConfig(config *models.Config) error {
	// Ensure data directory exists
	dataDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to temporary file first (atomic write pattern)
	tmpPath := cm.configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Rename temporary file to actual config file
	if err := os.Rename(tmpPath, cm.configPath); err != nil {
		// Clean up temp file on error
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the path to the config file
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configPath
}
