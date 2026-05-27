package graphics

import (
	"fmt"
	"math"
	"strings"
)

// RGB holds red, green, blue values
type RGB struct {
	R, G, B uint8
}

// ParseHex converts a hex color string to RGB
func ParseHex(hex string) (RGB, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return RGB{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return RGB{}, fmt.Errorf("could not parse hex color: %w", err)
	}

	return RGB{r, g, b}, nil
}

// Interpolate blends two colors at position t (0.0 to 1.0)
func Interpolate(from, to RGB, t float64) RGB {
	return RGB{
		R: uint8(math.Round(float64(from.R) + t*float64(int(to.R)-int(from.R)))),
		G: uint8(math.Round(float64(from.G) + t*float64(int(to.G)-int(from.G)))),
		B: uint8(math.Round(float64(from.B) + t*float64(int(to.B)-int(from.B)))),
	}
}

// AnsiColor wraps text in a 24-bit ANSI foreground color
func AnsiColor(c RGB, text string) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", c.R, c.G, c.B, text)
}

// AnsiBgColor wraps text in a 24-bit ANSI background color
func AnsiBgColor(c RGB, text string) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", c.R, c.G, c.B, text)
}

// GradientText applies a horizontal gradient across each character
func GradientText(text string, fromHex string, toHex string) (string, error) {
	from, err := ParseHex(fromHex)
	if err != nil {
		return text, err
	}

	to, err := ParseHex(toHex)
	if err != nil {
		return text, err
	}

	runes := []rune(text)
	total := len(runes)
	if total == 0 {
		return "", nil
	}

	var result strings.Builder
	for i, ch := range runes {
		t := float64(i) / float64(total-1)
		if total == 1 {
			t = 0
		}
		color := Interpolate(from, to, t)
		result.WriteString(AnsiColor(color, string(ch)))
	}

	return result.String(), nil
}

// GradientLine renders a full line of a character with a gradient
func GradientLine(char string, width int, fromHex string, toHex string) (string, error) {
	line := strings.Repeat(char, width)
	return GradientText(line, fromHex, toHex)
}

// RainbowText applies rainbow colors across text
func RainbowText(text string) string {
	colors := []RGB{
		{255, 0, 0},
		{255, 127, 0},
		{255, 255, 0},
		{0, 255, 0},
		{0, 0, 255},
		{75, 0, 130},
		{148, 0, 211},
	}

	runes := []rune(text)
	var result strings.Builder
	for i, ch := range runes {
		color := colors[i%len(colors)]
		result.WriteString(AnsiColor(color, string(ch)))
	}
	return result.String()
}
