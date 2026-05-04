package models

import "time"

// Config stores app settings
type Config struct {
	Version string   `json:"version"`
	UI      UIConfig `json:"ui"`
}

// UIConfig stores saved UI state
type UIConfig struct {
	Theme        string    `json:"theme"`
	ViewMode     string    `json:"viewMode"`
	CurrentDate  time.Time `json:"currentDate"`
	WindowWidth  float32   `json:"windowWidth"`
	WindowHeight float32   `json:"windowHeight"`
}

// NewDefaultConfig creates default config
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

func (c *Config) GetTheme() string {
	return c.UI.Theme
}

func (c *Config) SetTheme(theme string) {
	c.UI.Theme = theme
}

func (c *Config) GetViewMode() string {
	return c.UI.ViewMode
}

func (c *Config) SetViewMode(viewMode string) {
	c.UI.ViewMode = viewMode
}

func (c *Config) GetCurrentDate() time.Time {
	return c.UI.CurrentDate
}

func (c *Config) SetCurrentDate(date time.Time) {
	c.UI.CurrentDate = date
}
