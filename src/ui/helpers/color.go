package helpers

import "image/color"

// ToNRGBA converts any color.Color to color.NRGBA for manipulation.
func ToNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

// Darken returns a darker version of the color by the provided factor (0..1).
func Darken(c color.NRGBA, factor float32) color.NRGBA {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}
	return color.NRGBA{
		R: uint8(float32(c.R) * factor),
		G: uint8(float32(c.G) * factor),
		B: uint8(float32(c.B) * factor),
		A: c.A,
	}
}

// Lighten returns a lighter version of the color by mixing with white (0..1).
func Lighten(c color.NRGBA, amount float32) color.NRGBA {
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	mix := func(v uint8) uint8 {
		return uint8(float32(v)*(1-amount) + 255*amount)
	}
	return color.NRGBA{R: mix(c.R), G: mix(c.G), B: mix(c.B), A: c.A}
}

// Hex parses a #RRGGBB hex color string into color.NRGBA.
func Hex(h string) color.NRGBA {
	var r, g, b uint8
	if len(h) == 7 && h[0] == '#' {
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
