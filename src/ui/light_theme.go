package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// LightSoftTheme customizes the default light theme to a softer palette
// inspired by the Light mockup (rounded, subtle borders, warm background).
type LightSoftTheme struct{}

func NewLightSoftTheme() fyne.Theme { return &LightSoftTheme{} }

func (t *LightSoftTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	// Define color palette
	bg := hex("#3c3c3c")       // main background - gray matching window gradient
	fg := hex("#3c3836")       // DARK text for main window (tasks, menus)
	inputBg := hex("#ffffff")  // white inputs for contrast
	primary := hex("#ff8c42")  // accent/add-button
	separator := hex("#d0d0d0") // borders
	placeholder := hex("#aaaaaa") // lighter muted text for visibility on gray
	hoverColor := hex("#ffd27a") // hover state

	switch name {
	case theme.ColorNameBackground:
		return bg // gray background for all dialogs/windows
	case theme.ColorNameOverlayBackground:
		return bg // dialog overlay background - same gray
	case theme.ColorNameMenuBackground:
		return inputBg // white background for popup menus
	case theme.ColorNameForeground:
		return fg // DARK text for main window
	case theme.ColorNameSeparator:
		return separator // borders
	case theme.ColorNameInputBackground:
		return inputBg // white inputs
	case theme.ColorNamePrimary:
		return primary // accent
	case theme.ColorNameButton:
		return primary // buttons
	case theme.ColorNamePlaceHolder:
		return placeholder // lighter muted text
	case theme.ColorNameHover:
		return hoverColor // hover
	case theme.ColorNameFocus:
		return primary // focus
	case theme.ColorNameSelection:
		return separator // selection
	case theme.ColorNameDisabled:
		return hex("#999999") // disabled state
	case theme.ColorNameDisabledButton:
		return hex("#cccccc") // disabled buttons
	default:
		return theme.LightTheme().Color(name, theme.VariantLight)
	}
}

// GetHeaderGradientColors returns the two colors for the header gradient.
// DARK background gradient for light theme (same as dark theme but slightly different shade)
// Dark gradient from nearly black to dark gray
func (t *LightSoftTheme) GetHeaderGradientColors() (color.Color, color.Color) {
	start := color.NRGBA{R: 0x2a, G: 0x2a, B: 0x2a, A: 0xFF} // #2a2a2a - Nearly black at top
	end := color.NRGBA{R: 0x3c, G: 0x3c, B: 0x3c, A: 0xFF}   // #3c3c3c - Dark gray at bottom
	return start, end
}

func (t *LightSoftTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.LightTheme().Font(style)
}

func (t *LightSoftTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.LightTheme().Icon(name)
}

func (t *LightSoftTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.LightTheme().Size(name)
}
