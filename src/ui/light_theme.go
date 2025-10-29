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
		// #FFF9F2 - warm background
		return color.NRGBA{R: 0xFF, G: 0xF9, B: 0xF2, A: 0xFF}
	case theme.ColorNameForeground:
		// #3c3836 - main text per mockup
		return color.NRGBA{R: 0x3C, G: 0x38, B: 0x36, A: 0xFF}
	case theme.ColorNameSeparator:
		// rgba(0,0,0,0.06) - subtle borders
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x0F}
	case theme.ColorNameInputBackground:
		// #fff - card background
		return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	case theme.ColorNamePrimary:
		// #FF8C42 - accent for buttons/FAB
		return color.NRGBA{R: 0xFF, G: 0x8C, B: 0x42, A: 0xFF}
	case theme.ColorNameButton:
		// #FFD27A - button/control background
		return color.NRGBA{R: 0xFF, G: 0xD2, B: 0x7A, A: 0xFF}
	case theme.ColorNamePlaceHolder:
		// #6b6b6b - muted text
		return color.NRGBA{R: 0x6B, G: 0x6B, B: 0x6B, A: 0xFF}
	case theme.ColorNameHover:
		// rgba(255,255,255,0.9) - control hover background
		return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xE6}
	case theme.ColorNameFocus:
		// #F0C23D - accent color for focus
		return color.NRGBA{R: 0xF0, G: 0xC2, B: 0x3D, A: 0xFF}
	case theme.ColorNameSelection:
		// Very subtle selection
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x18}
	default:
		return theme.LightTheme().Color(name, theme.VariantLight)
	}
}

// GetHeaderGradientColors returns the two colors for the header gradient.
// From mockup: Start: #ff9a4d, End: #fff6e6
func (t *LightSoftTheme) GetHeaderGradientColors() (color.Color, color.Color) {
	start := color.NRGBA{R: 0xFF, G: 0x9A, B: 0x4D, A: 0xFF} // #ff9a4d (from mockup)
	end := color.NRGBA{R: 0xFF, G: 0xF6, B: 0xE6, A: 0xFF}   // #fff6e6 (from mockup)
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
