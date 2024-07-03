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
	// Gruvbox palette matching the mockup
	// Backgrounds
	bg := hex("#2a2723")    // main background --bg
	panel := hex("#3c3836") // panel/card background --panel-bg / --card-bg / input background
	// Foregrounds
	fg := hex("#ebdbb2")    // primary text --text
	muted := hex("#bdae93") // secondary text / placeholders --muted
	// Accents
	primary := hex("#fabd2f")   // yellow accent per mockup
	focus := hex("#fabd2f")     // focus ring equals accent
	hover := hex("#504945")     // hover overlay
	selection := hex("#665c54") // selection background
	disabled := hex("#504945")  // disabled elements
	border := hex("#3c3836")    // subtle border rgba(255,255,255,0.06)

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
// From mockup: Start: #3c2414, End: #6e4a30
func (t *GruvboxBlackTheme) GetHeaderGradientColors() (color.Color, color.Color) {
	start := hex("#3c2414") // --header-grad-a (from mockup)
	end := hex("#6e4a30")   // --header-grad-b (from mockup)
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
