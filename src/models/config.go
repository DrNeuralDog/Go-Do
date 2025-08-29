package models

import "time"

// Config represents the application configuration
type Config struct {
	Version string    `json:"version"`
	UI      UIConfig  `json:"ui"`
}

// UIConfig stores UI state preferences
type UIConfig struct {
	Theme           string    `json:"theme"`           // "light" or "dark"
	ViewMode        string    `json:"viewMode"`        // "all", "incomplete", "complete", "starred"
	CurrentDate     time.Time `json:"currentDate"`     // Last viewed date
	WindowWidth     float32   `json:"windowWidth"`     // Window dimensions (for future)
	WindowHeight    float32   `json:"windowHeight"`    // Window dimensions (for future)
}

// NewDefaultConfig creates a default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		UI: UIConfig{
			Theme:        "light",
			ViewMode:     "incomplete",
			CurrentDate:  time.Now(),
			WindowWidth:  420,
			WindowHeight: 800,
		},
	}
}

// GetTheme returns the current theme
func (c *Config) GetTheme() string {
	return c.UI.Theme
}

// SetTheme sets the current theme
func (c *Config) SetTheme(theme string) {
	c.UI.Theme = theme
}

// GetViewMode returns the current view mode
func (c *Config) GetViewMode() string {
	return c.UI.ViewMode
}

// SetViewMode sets the current view mode
func (c *Config) SetViewMode(viewMode string) {
	c.UI.ViewMode = viewMode
}

// GetCurrentDate returns the last viewed date
func (c *Config) GetCurrentDate() time.Time {
	return c.UI.CurrentDate
}

// SetCurrentDate sets the current date
func (c *Config) SetCurrentDate(date time.Time) {
	c.UI.CurrentDate = date
}
