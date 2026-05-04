package helpers

import "image/color"

const (
	hexColorLength = 7
	hexMaxAlpha    = 255
)

// ToNRGBA converts color.Color to color.NRGBA
func ToNRGBA(c color.Color) color.NRGBA {
	if c == nil {
		return color.NRGBA{}
	}

	return color.NRGBAModel.Convert(c).(color.NRGBA)
}

// Darken returns darker color by factor
func Darken(c color.NRGBA, factor float32) color.NRGBA {
	factor = clampColorFactor(factor)

	return color.NRGBA{
		R: uint8(float32(c.R) * factor),
		G: uint8(float32(c.G) * factor),
		B: uint8(float32(c.B) * factor),

		A: c.A,
	}
}

// Lighten returns color mixed with white
func Lighten(c color.NRGBA, amount float32) color.NRGBA {
	amount = clampColorFactor(amount)

	return color.NRGBA{
		R: mixColorByte(c.R, amount),
		G: mixColorByte(c.G, amount),
		B: mixColorByte(c.B, amount),

		A: c.A,
	}
}

// Hex parses #RRGGBB into color.NRGBA
func Hex(h string) color.NRGBA {
	if len(h) != hexColorLength || h[0] != '#' {
		return color.NRGBA{}
	}

	r, ok := hexColorByte(h[1], h[2])
	if !ok {
		return color.NRGBA{}
	}

	g, ok := hexColorByte(h[3], h[4])
	if !ok {
		return color.NRGBA{}
	}

	b, ok := hexColorByte(h[5], h[6])
	if !ok {
		return color.NRGBA{}
	}

	return color.NRGBA{R: r, G: g, B: b, A: hexMaxAlpha}
}

// clampColorFactor keeps color factor inside 0..1
func clampColorFactor(value float32) float32 {
	if value < 0 {
		return 0
	}

	if value > 1 {
		return 1
	}

	return value
}

// mixColorByte mixes one channel with white
func mixColorByte(value uint8, amount float32) uint8 {
	return uint8(float32(value)*(1-amount) + hexMaxAlpha*amount)
}

func hexColorByte(hi, lo byte) (uint8, bool) {
	hiValue, ok := hexNibble(hi)
	if !ok {
		return 0, false
	}

	loValue, ok := hexNibble(lo)
	if !ok {
		return 0, false
	}

	return hiValue<<4 | loValue, true
}

// hexNibble parses one hex symbol
func hexNibble(c byte) (uint8, bool) {
	switch {
	case c >= '0' && c <= '9':
		return uint8(c - '0'), true
	case c >= 'a' && c <= 'f':
		return uint8(10 + c - 'a'), true
	case c >= 'A' && c <= 'F':
		return uint8(10 + c - 'A'), true
	default:
		return 0, false
	}
}
