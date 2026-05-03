package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"godo/src/ui/helpers"
)

var gruvboxBlackPalette = struct {
	background      color.Color
	foreground      color.Color
	inputBackground color.Color
	primary         color.Color
	hover           color.Color
	selection       color.Color
	disabled        color.Color
	separator       color.Color
	placeholder     color.Color
	gradientStart   color.Color
	gradientEnd     color.Color
}{
	background:      helpers.Hex("#282828"),
	foreground:      helpers.Hex("#ebdbb2"),
	inputBackground: helpers.Hex("#3c3836"),
	primary:         helpers.Hex("#fabd2f"),
	hover:           helpers.Hex("#504945"),
	selection:       helpers.Hex("#665c54"),
	disabled:        helpers.Hex("#504945"),
	separator:       helpers.Hex("#3c3836"),
	placeholder:     helpers.Hex("#a89984"),
	gradientStart:   helpers.Hex("#282828"),
	gradientEnd:     helpers.Hex("#3c3836"),
}

type GruvboxBlackTheme struct{}

// NewGruvboxBlackTheme creates Gruvbox black app theme
func NewGruvboxBlackTheme() fyne.Theme {
	return &GruvboxBlackTheme{}
}

// IsLight reports dark theme styling
func (t *GruvboxBlackTheme) IsLight() bool { return false }

// Color returns Gruvbox palette color by Fyne color name
func (t *GruvboxBlackTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return gruvboxBlackPalette.background
	case theme.ColorNameOverlayBackground:
		return gruvboxBlackPalette.background
	case theme.ColorNameMenuBackground:
		return gruvboxBlackPalette.inputBackground
	case theme.ColorNameForeground:
		return gruvboxBlackPalette.foreground
	case theme.ColorNameButton:
		return gruvboxBlackPalette.primary
	case theme.ColorNameDisabled:
		return gruvboxBlackPalette.disabled
	case theme.ColorNameDisabledButton:
		return gruvboxBlackPalette.disabled
	case theme.ColorNamePrimary:
		return gruvboxBlackPalette.primary
	case theme.ColorNameFocus:
		return gruvboxBlackPalette.primary
	case theme.ColorNameHover:
		return gruvboxBlackPalette.hover
	case theme.ColorNameInputBackground:
		return gruvboxBlackPalette.inputBackground
	case theme.ColorNamePlaceHolder:
		return gruvboxBlackPalette.placeholder
	case theme.ColorNameSeparator:
		return gruvboxBlackPalette.separator
	case theme.ColorNameSelection:
		return gruvboxBlackPalette.selection
	default:
		return theme.DarkTheme().Color(name, theme.VariantDark)
	}
}

// GetHeaderGradientColors returns main window gradient colors
func (t *GruvboxBlackTheme) GetHeaderGradientColors() (color.Color, color.Color) {
	return gruvboxBlackPalette.gradientStart, gruvboxBlackPalette.gradientEnd
}

// Font returns default dark theme font
func (t *GruvboxBlackTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

// Icon returns default dark theme icon
func (t *GruvboxBlackTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

// Size returns default dark theme size
func (t *GruvboxBlackTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DarkTheme().Size(name)
}
