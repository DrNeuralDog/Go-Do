package helpers

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type lightAware interface {
	IsLight() bool
}

// IsLightTheme reports whether the current application theme should be treated as light.
func IsLightTheme() bool {
	app := fyne.CurrentApp()
	if app == nil {
		return false
	}

	th := app.Settings().Theme()
	if th == nil {
		return false
	}

	if detector, ok := th.(lightAware); ok {
		return detector.IsLight()
	}

	return app.Settings().ThemeVariant() == theme.VariantLight
}

// GetBackgroundColor returns the background color from the active theme.
func GetBackgroundColor() color.Color {
	return colorFromTheme(theme.ColorNameBackground)
}

// GetForegroundColor returns the foreground color from the active theme.
func GetForegroundColor() color.Color {
	return colorFromTheme(theme.ColorNameForeground)
}

// GetCardColor returns the input background color useful for cards and panels.
func GetCardColor() color.Color {
	return colorFromTheme(theme.ColorNameInputBackground)
}

func colorFromTheme(name fyne.ThemeColorName) color.Color {
	app := fyne.CurrentApp()
	if app == nil || app.Settings() == nil || app.Settings().Theme() == nil {
		return theme.Color(name)
	}

	variant := app.Settings().ThemeVariant()
	return app.Settings().Theme().Color(name, variant)
}
