package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"godo/src/ui/helpers"
)

var lightSoftPalette = struct {
	background      color.Color
	foreground      color.Color
	inputBackground color.Color
	primary         color.Color
	separator       color.Color
	placeholder     color.Color
	hover           color.Color
	disabled        color.Color
	disabledButton  color.Color
	gradientStart   color.Color
	gradientEnd     color.Color
}{
	background:      helpers.Hex("#3c3c3c"),
	foreground:      helpers.Hex("#3c3836"),
	inputBackground: helpers.Hex("#ffffff"),
	primary:         helpers.Hex("#ff8c42"),
	separator:       helpers.Hex("#d0d0d0"),
	placeholder:     helpers.Hex("#aaaaaa"),
	hover:           helpers.Hex("#ffd27a"),
	disabled:        helpers.Hex("#999999"),
	disabledButton:  helpers.Hex("#cccccc"),
	gradientStart:   helpers.Hex("#2a2a2a"),
	gradientEnd:     helpers.Hex("#3c3c3c"),
}

// LightSoftTheme provides light controls on dark window background
type LightSoftTheme struct{}

func NewLightSoftTheme() fyne.Theme {
	return &LightSoftTheme{}
}

func (t *LightSoftTheme) IsLight() bool { return true }

// Color returns palette color by Fyne color name
func (t *LightSoftTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return lightSoftPalette.background
	case theme.ColorNameOverlayBackground:
		return lightSoftPalette.background
	case theme.ColorNameMenuBackground:
		return lightSoftPalette.inputBackground
	case theme.ColorNameForeground:
		return lightSoftPalette.foreground
	case theme.ColorNameSeparator:
		return lightSoftPalette.separator
	case theme.ColorNameInputBackground:
		return lightSoftPalette.inputBackground
	case theme.ColorNamePrimary:
		return lightSoftPalette.primary
	case theme.ColorNameButton:
		return lightSoftPalette.primary
	case theme.ColorNamePlaceHolder:
		return lightSoftPalette.placeholder
	case theme.ColorNameHover:
		return lightSoftPalette.hover
	case theme.ColorNameFocus:
		return lightSoftPalette.primary
	case theme.ColorNameSelection:
		return lightSoftPalette.separator
	case theme.ColorNameDisabled:
		return lightSoftPalette.disabled
	case theme.ColorNameDisabledButton:
		return lightSoftPalette.disabledButton
	default:
		return theme.LightTheme().Color(name, theme.VariantLight)
	}
}

// GetHeaderGradientColors returns main window gradient colors
func (t *LightSoftTheme) GetHeaderGradientColors() (color.Color, color.Color) {
	return lightSoftPalette.gradientStart, lightSoftPalette.gradientEnd
}

// Font returns default light theme font
func (t *LightSoftTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.LightTheme().Font(style)
}

// Icon returns default light theme icon
func (t *LightSoftTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.LightTheme().Icon(name)
}

// Size returns default light theme size
func (t *LightSoftTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.LightTheme().Size(name)
}
