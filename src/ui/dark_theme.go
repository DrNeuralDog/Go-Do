package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"godo/src/ui/helpers"
)

// GruvboxBlackTheme provides a near-black Gruvbox-inspired dark theme.
// Only color palette is customized; fonts, icons, and sizes are inherited
// from the default dark theme for compatibility.
type GruvboxBlackTheme struct{}

// NewGruvboxBlackTheme returns a new instance of the Gruvbox black theme.
func NewGruvboxBlackTheme() fyne.Theme {
	return &GruvboxBlackTheme{}
}

// IsLight indicates this theme should use dark styling.
func (t *GruvboxBlackTheme) IsLight() bool { return false }

func (t *GruvboxBlackTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	// Updated Gruvbox palette strictly matching the mockup CSS
	// Backgrounds
	bg := helpers.Hex("#282828")    // tasks-container bg
	panel := helpers.Hex("#3c3836") // dropdown bg
	// Foregrounds
	fg := helpers.Hex("#ebdbb2")    //#ebdbb2 main text
	muted := helpers.Hex("#a89984") // task-time
	// Accents
	primary := helpers.Hex("#fabd2f")   // logo, add-button bg, theme-selector text
	focus := helpers.Hex("#fabd2f")     // focus same as primary
	hover := helpers.Hex("#504945")     // nav-btn bg, theme-selector bg
	selection := helpers.Hex("#665c54") // checkbox border
	disabled := helpers.Hex("#504945")  // disabled
	border := helpers.Hex("#3c3836")    // task-item border

	switch name {
	case theme.ColorNameBackground:
		return bg
	case theme.ColorNameForeground:
		return fg
	case theme.ColorNameButton:
		return primary // Use accent color for buttons
	case theme.ColorNameDisabled:
		return disabled
	case theme.ColorNameDisabledButton:
		return disabled
	case theme.ColorNamePrimary:
		return primary
	case theme.ColorNameFocus:
		return focus
	case theme.ColorNameHover:
		return hover
	case theme.ColorNameInputBackground:
		return panel
	case theme.ColorNamePlaceHolder:
		return muted
	case theme.ColorNameSeparator:
		return border
	case theme.ColorNameSelection:
		return selection
	default:
		return theme.DarkTheme().Color(name, theme.VariantDark)
	}
}

// GetHeaderGradientColors returns the two colors for the header gradient in Gruvbox theme.
// DARK background gradient - nearly black to dark gray
func (t *GruvboxBlackTheme) GetHeaderGradientColors() (color.Color, color.Color) {
	start := helpers.Hex("#282828") // Nearly black at top
	end := helpers.Hex("#3c3836")   // Dark gray at bottom
	return start, end
}

func (t *GruvboxBlackTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

func (t *GruvboxBlackTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

func (t *GruvboxBlackTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DarkTheme().Size(name)
}
