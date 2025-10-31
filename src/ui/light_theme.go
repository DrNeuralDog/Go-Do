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
	switch name {
	case theme.ColorNameBackground:
		// Force dark app background to avoid light bleed outside cards
		return color.NRGBA{R: 0x28, G: 0x28, B: 0x28, A: 0xFF} // #282828
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0x3C, G: 0x38, B: 0x36, A: 0xFF} // #3c3836 main text
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF} // #d0d0d0 task-item border (darker for visibility)
	case theme.ColorNameInputBackground:
		// Use dark panel for containers/cards to match overall dark background
		return color.NRGBA{R: 0x3C, G: 0x38, B: 0x36, A: 0xFF} // #3c3836
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0xFF, G: 0x8C, B: 0x42, A: 0xFF} // #ff8c42 accent/add-button
	case theme.ColorNameButton:
		return color.NRGBA{R: 0xFF, G: 0x8C, B: 0x42, A: 0xFF} // #ff8c42 nav/theme bg
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xFF} // #666 task-time
	case theme.ColorNameHover:
		return color.NRGBA{R: 0xFF, G: 0xD2, B: 0x7A, A: 0xFF} // Approx hover
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xFF, G: 0x8C, B: 0x42, A: 0xFF} // #ff8c42 focus
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0xD0, G: 0xD0, B: 0xD0, A: 0xFF} // #d0d0d0 checkbox border
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
