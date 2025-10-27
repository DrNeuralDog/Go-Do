package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// GruvboxBlackTheme provides a near-black Gruvbox-inspired dark theme.
// Only color palette is customized; fonts, icons, and sizes are inherited
// from the default dark theme for compatibility.
type GruvboxBlackTheme struct{}

// NewGruvboxBlackTheme returns a new instance of the Gruvbox black theme.
func NewGruvboxBlackTheme() fyne.Theme {
	return &GruvboxBlackTheme{}
}

func (t *GruvboxBlackTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	// Core Gruvbox-like palette (black variant)
	// Backgrounds
	bg := hex("#0d0e0f")    // almost black
	bg2 := hex("#1d2021")   // darker gray for surfaces
	panel := hex("#32302f") // panel/button background
	input := hex("#1f1f1f") // inputs
	// Foregrounds
	fg := hex("#ebdbb2")    // primary text
	muted := hex("#a89984") // secondary text / placeholders
	// Accents
	primary := hex("#d79921")   // yellow accent
	focus := hex("#fabd2f")     // focus ring
	hover := hex("#3c3836")     // hover overlay
	selection := hex("#665c54") // selection background
	disabled := hex("#504945")  // disabled elements

	switch name {
	case theme.ColorNameBackground:
		return bg
	case theme.ColorNameForeground:
		return fg
	case theme.ColorNameButton:
		return panel
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
		return input
	case theme.ColorNamePlaceHolder:
		return muted
	case theme.ColorNameSeparator:
		return bg2
	case theme.ColorNameSelection:
		return selection
	default:
		return theme.DarkTheme().Color(name, theme.VariantDark)
	}
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

// hex parses a #RRGGBB hex color to color.NRGBA.
func hex(h string) color.NRGBA {
	var r, g, b uint8
	// assume valid #RRGGBB
	if len(h) == 7 && h[0] == '#' {
		// naive parse; avoid fmt/strconv to keep it small
		r = fromHex(h[1])<<4 | fromHex(h[2])
		g = fromHex(h[3])<<4 | fromHex(h[4])
		b = fromHex(h[5])<<4 | fromHex(h[6])
	}
	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

func fromHex(c byte) uint8 {
	if c >= '0' && c <= '9' {
		return uint8(c - '0')
	}
	if c >= 'a' && c <= 'f' {
		return uint8(10 + c - 'a')
	}
	if c >= 'A' && c <= 'F' {
		return uint8(10 + c - 'A')
	}
	return 0
}
