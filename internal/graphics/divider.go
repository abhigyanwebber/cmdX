package graphics

import (
	"fmt"
	"strings"
)

// DividerStyle identifies a horizontal divider line style usable in
// theme previews and section separators.
type DividerStyle string

// Available divider styles.
const (
	StyleLine   DividerStyle = "line"
	StyleWave   DividerStyle = "wave"
	StyleDots   DividerStyle = "dots"
	StyleStars  DividerStyle = "stars"
	StyleDouble DividerStyle = "double"
	StyleArrow  DividerStyle = "arrow"
	StyleZigzag DividerStyle = "zigzag"
)

func dividerPattern(style DividerStyle, width int) string {
	switch style {
	case StyleLine:
		return strings.Repeat("─", width)
	case StyleWave:
		unit := "~-~"
		return strings.Repeat(unit, width/len(unit))
	case StyleDots:
		unit := "· "
		return strings.Repeat(unit, width/len(unit))
	case StyleStars:
		unit := "✦ "
		return strings.Repeat(unit, width/len(unit))
	case StyleDouble:
		return strings.Repeat("═", width)
	case StyleArrow:
		return strings.Repeat("»", width)
	case StyleZigzag:
		unit := "/\\"
		return strings.Repeat(unit, width/len(unit))
	default:
		return strings.Repeat("─", width)
	}
}

// Divider renders a styled horizontal divider
func Divider(style DividerStyle, width int, colorHex string) string {
	c, err := ParseHex(colorHex)
	if err != nil {
		return dividerPattern(style, width)
	}
	rgb := ToRGB(c)
	return AnsiColor(rgb, dividerPattern(style, width))
}

// GradientDivider renders a divider with a LAB gradient
func GradientDivider(style DividerStyle, width int, fromHex string, toHex string) (string, error) {
	var char string
	switch style {
	case StyleLine:
		char = "─"
	case StyleDouble:
		char = "═"
	case StyleDots:
		char = "·"
	case StyleStars:
		char = "✦"
	default:
		char = "─"
	}
	line := strings.Repeat(char, width)
	return GradientText(line, fromHex, toHex)
}

// SectionHeader renders a titled divider
func SectionHeader(title string, width int, colorHex string) string {
	c, err := ParseHex(colorHex)
	if err != nil {
		return title
	}
	rgb := ToRGB(c)

	titleLen := len([]rune(title)) + 2
	remaining := width - titleLen
	if remaining < 0 {
		remaining = 0
	}
	left := remaining / 2
	right := remaining - left

	result := fmt.Sprintf("%s %s %s",
		strings.Repeat("─", left),
		title,
		strings.Repeat("─", right),
	)
	return AnsiColor(rgb, result)
}
