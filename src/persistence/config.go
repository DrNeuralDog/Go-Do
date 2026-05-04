package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"godo/src/models"
)

const configFileName = "config.json"

// ConfigManager handles app config file
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates config storage
func NewConfigManager(dataDir string) *ConfigManager {
	return &ConfigManager{
		configPath: filepath.Join(dataDir, configFileName),
	}
}

// LoadConfig loads config or returns defaults
func (cm *ConfigManager) LoadConfig() (*models.Config, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return models.NewDefaultConfig(), nil
		}

		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if len(data) == 0 {
		return models.NewDefaultConfig(), nil
	}

	config := models.NewDefaultConfig()

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves config to disk
func (cm *ConfigManager) SaveConfig(config *models.Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	dataDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dataDir, dataDirPermission); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return writeFileAtomic(cm.configPath, data)
}

// GetConfigPath returns config file path
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configPath
}
