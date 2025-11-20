/*
Package helpers provides utility functions for UI operations in the Go Do application.

The helpers package contains helper functions for common UI tasks:

Color Utilities (color.go):
  - ToNRGBA: Converts color.Color to color.NRGBA
  - Lighten: Lightens a color by a given factor
  - Darken: Darkens a color by a given factor
  - Hex: Converts hex color string to color.Color

Theme Utilities (theme.go):
  - IsLightTheme: Checks if current theme is light
  - GetBackgroundColor: Returns theme-appropriate background color
  - GetForegroundColor: Returns theme-appropriate foreground color
  - GetCardColor: Returns theme-appropriate card color

Layout Utilities (layout.go):
  - CreateSpacer: Creates a fixed-size spacer widget
  - CreateCardStyle: Creates a styled card container

Window Utilities (window.go):
  - FlashWindow: Creates a visual flash effect on a window
*/
package helpers
