package helpers

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type lightThemeDetector interface {
	IsLight() bool
}

// IsLightTheme reports whether active theme is light
func IsLightTheme() bool {
	app := fyne.CurrentApp()
	if app == nil || app.Settings() == nil {
		return false
	}

	settings := app.Settings()
	if detector, ok := settings.Theme().(lightThemeDetector); ok {
		return detector.IsLight()
	}

	return settings.ThemeVariant() == theme.VariantLight
}

// GetBackgroundColor returns active theme background
func GetBackgroundColor() color.Color {
	return colorFromTheme(theme.ColorNameBackground)
}

// GetForegroundColor returns active theme foreground
func GetForegroundColor() color.Color {
	return colorFromTheme(theme.ColorNameForeground)
}

// GetCardColor returns active theme input background
func GetCardColor() color.Color {
	return colorFromTheme(theme.ColorNameInputBackground)
}

// colorFromTheme returns active theme color by name
func colorFromTheme(name fyne.ThemeColorName) color.Color {
	app := fyne.CurrentApp()
	if app == nil || app.Settings() == nil {
		return theme.Color(name)
	}

	settings := app.Settings()
	currentTheme := settings.Theme()

	if currentTheme == nil {
		return theme.Color(name)
	}

	return currentTheme.Color(name, settings.ThemeVariant())
}
