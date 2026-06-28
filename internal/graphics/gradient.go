package graphics

import (
	"fmt"
	"math"
	"strings"

	colorful "github.com/lucasb-eyer/go-colorful"
)

// ParseHex converts a hex color string to a colorful.Color
func ParseHex(hex string) (colorful.Color, error) {
	c, err := colorful.Hex(hex)
	if err != nil {
		return colorful.Color{}, fmt.Errorf("invalid hex color '%s': %w", hex, err)
	}
	return c, nil
}

// RGB holds red, green, blue values (0-255) for ANSI output
type RGB struct {
	R, G, B uint8
}

// ToRGB converts a colorful.Color to RGB struct
func ToRGB(c colorful.Color) RGB {
	r, g, b := c.Clamped().RGB255()
	return RGB{r, g, b}
}

// AnsiColor wraps text in a 24-bit ANSI foreground color
func AnsiColor(c RGB, text string) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", c.R, c.G, c.B, text)
}

// AnsiBgColor wraps text in a 24-bit ANSI background color
func AnsiBgColor(c RGB, text string) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", c.R, c.G, c.B, text)
}

// AnsiColorHex wraps text using a hex color string directly
func AnsiColorHex(hex string, text string) string {
	c, err := ParseHex(hex)
	if err != nil {
		return text
	}
	rgb := ToRGB(c)
	return AnsiColor(rgb, text)
}

// GradientText applies a smooth LAB color space gradient across each character
// LAB blending produces natural transitions without muddy midpoints
func GradientText(text string, fromHex string, toHex string) (string, error) {
	from, err := ParseHex(fromHex)
	if err != nil {
		return text, fmt.Errorf("gradient from color: %w", err)
	}

	to, err := ParseHex(toHex)
	if err != nil {
		return text, fmt.Errorf("gradient to color: %w", err)
	}

	runes := []rune(text)
	total := len(runes)
	if total == 0 {
		return "", nil
	}

	var result strings.Builder
	for i, ch := range runes {
		t := 0.0
		if total > 1 {
			t = float64(i) / float64(total-1)
		}
		// BlendLab gives perceptually uniform color transitions
		blended := from.BlendLab(to, t).Clamped()
		rgb := ToRGB(blended)
		result.WriteString(AnsiColor(rgb, string(ch)))
	}

	return result.String(), nil
}

// GradientLine renders a full line of a character with a LAB gradient
func GradientLine(char string, width int, fromHex string, toHex string) (string, error) {
	line := strings.Repeat(char, width)
	return GradientText(line, fromHex, toHex)
}

// RainbowText applies perceptually uniform rainbow colors across text
// Uses LCH color space for vivid, consistent brightness
func RainbowText(text string) string {
	runes := []rune(text)
	var result strings.Builder

	for i, ch := range runes {
		// rotate hue evenly across 360 degrees
		hue := float64(i%7) / 7.0 * 360.0
		c := colorful.Hcl(hue, 0.8, 0.6)
		rgb := ToRGB(c)
		result.WriteString(AnsiColor(rgb, string(ch)))
	}
	return result.String()
}

// ComplementaryColor returns the complementary hex color
func ComplementaryColor(hex string) (string, error) {
	c, err := ParseHex(hex)
	if err != nil {
		return hex, err
	}
	h, s, v := c.Hsv()
	complementH := math.Mod(h+180, 360)
	comp := colorful.Hsv(complementH, s, v)
	return comp.Hex(), nil
}

// AnalogousColors returns two analogous colors (30 degrees apart)
func AnalogousColors(hex string) (string, string, error) {
	c, err := ParseHex(hex)
	if err != nil {
		return hex, hex, err
	}
	h, s, v := c.Hsv()
	c1 := colorful.Hsv(math.Mod(h+30, 360), s, v)
	c2 := colorful.Hsv(math.Mod(h-30+360, 360), s, v)
	return c1.Hex(), c2.Hex(), nil
}

// TriadicColors returns two triadic colors (120 degrees apart)
func TriadicColors(hex string) (string, string, error) {
	c, err := ParseHex(hex)
	if err != nil {
		return hex, hex, err
	}
	h, s, v := c.Hsv()
	c1 := colorful.Hsv(math.Mod(h+120, 360), s, v)
	c2 := colorful.Hsv(math.Mod(h+240, 360), s, v)
	return c1.Hex(), c2.Hex(), nil
}

// WarmColor shifts a color slightly warmer
func WarmColor(hex string, amount float64) (string, error) {
	c, err := ParseHex(hex)
	if err != nil {
		return hex, err
	}
	h, s, v := c.Hsv()
	// shift hue toward red/orange (warmer)
	h = math.Mod(h-amount*30+360, 360)
	return colorful.Hsv(h, s, v).Hex(), nil
}

// CoolColor shifts a color slightly cooler
func CoolColor(hex string, amount float64) (string, error) {
	c, err := ParseHex(hex)
	if err != nil {
		return hex, err
	}
	h, s, v := c.Hsv()
	// shift hue toward blue (cooler)
	h = math.Mod(h+amount*30, 360)
	return colorful.Hsv(h, s, v).Hex(), nil
}

// LightenColor increases the lightness of a color
func LightenColor(hex string, amount float64) (string, error) {
	c, err := ParseHex(hex)
	if err != nil {
		return hex, err
	}
	h, s, v := c.Hsv()
	v = math.Min(1.0, v+amount)
	return colorful.Hsv(h, s, v).Hex(), nil
}

// DarkenColor decreases the lightness of a color
func DarkenColor(hex string, amount float64) (string, error) {
	c, err := ParseHex(hex)
	if err != nil {
		return hex, err
	}
	h, s, v := c.Hsv()
	v = math.Max(0.0, v-amount)
	return colorful.Hsv(h, s, v).Hex(), nil
}
