package helpers_test

import (
	"image/color"
	"testing"

	"godo/src/ui/helpers"
)

// TestToNRGBAConvertsColors checks NRGBA conversion
func TestToNRGBAConvertsColors(t *testing.T) {
	tests := []struct {
		name     string
		input    color.Color
		expected color.NRGBA
	}{
		{"nil", nil, color.NRGBA{}},
		{"rgba", color.RGBA{R: 10, G: 20, B: 30, A: 255}, color.NRGBA{R: 10, G: 20, B: 30, A: 255}},
		{"transparent nrgba", color.NRGBA{R: 120, G: 80, B: 40, A: 128}, color.NRGBA{R: 120, G: 80, B: 40, A: 128}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := helpers.ToNRGBA(test.input)

			if result != test.expected {
				t.Errorf("ToNRGBA() expected %#v, got %#v", test.expected, result)
			}
		})
	}
}

// TestDarkenClampsFactor checks darker color bounds
func TestDarkenClampsFactor(t *testing.T) {
	base := color.NRGBA{R: 100, G: 120, B: 140, A: 200}

	tests := []struct {
		name     string
		factor   float32
		expected color.NRGBA
	}{
		{"below zero", -1, color.NRGBA{R: 0, G: 0, B: 0, A: 200}},
		{"inside range", 0.5, color.NRGBA{R: 50, G: 60, B: 70, A: 200}},
		{"above one", 2, base},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := helpers.Darken(base, test.factor)

			if result != test.expected {
				t.Errorf("Darken() expected %#v, got %#v", test.expected, result)
			}
		})
	}
}

// TestLightenClampsAmount checks white mix bounds
func TestLightenClampsAmount(t *testing.T) {
	base := color.NRGBA{R: 100, G: 120, B: 140, A: 200}

	tests := []struct {
		name     string
		amount   float32
		expected color.NRGBA
	}{
		{"below zero", -1, base},
		{"inside range", 0.5, color.NRGBA{R: 177, G: 187, B: 197, A: 200}},
		{"above one", 2, color.NRGBA{R: 255, G: 255, B: 255, A: 200}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := helpers.Lighten(base, test.amount)

			if result != test.expected {
				t.Errorf("Lighten() expected %#v, got %#v", test.expected, result)
			}
		})
	}
}

// TestHexParsesColors checks #RRGGBB parsing
func TestHexParsesColors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected color.NRGBA
	}{
		{"lowercase", "#fabd2f", color.NRGBA{R: 250, G: 189, B: 47, A: 255}},
		{"uppercase", "#FABD2F", color.NRGBA{R: 250, G: 189, B: 47, A: 255}},
		{"bad length", "#fff", color.NRGBA{}},
		{"missing hash", "fabd2f", color.NRGBA{}},
		{"bad symbol", "#fagd2f", color.NRGBA{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := helpers.Hex(test.input)

			if result != test.expected {
				t.Errorf("Hex(%q) expected %#v, got %#v", test.input, test.expected, result)
			}
		})
	}
}
