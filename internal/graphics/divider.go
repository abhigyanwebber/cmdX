package graphics

import (
	"fmt"
	"strings"
)

// DividerStyle defines available divider styles
type DividerStyle string

const (
	StyleLine   DividerStyle = "line"
	StyleWave   DividerStyle = "wave"
	StyleDots   DividerStyle = "dots"
	StyleStars  DividerStyle = "stars"
	StyleDouble DividerStyle = "double"
	StyleArrow  DividerStyle = "arrow"
	StyleZigzag DividerStyle = "zigzag"
)

// Divider renders a styled horizontal divider
func Divider(style DividerStyle, width int, colorHex string) string {
	color, err := ParseHex(colorHex)
	if err != nil {
		color = RGB{100, 100, 100}
	}

	var pattern string

	switch style {
	case StyleLine:
		pattern = strings.Repeat("─", width)
	case StyleWave:
		unit := "~-~"
		reps := width / len(unit)
		pattern = strings.Repeat(unit, reps)
	case StyleDots:
		unit := "· "
		reps := width / len(unit)
		pattern = strings.Repeat(unit, reps)
	case StyleStars:
		unit := "✦ "
		reps := width / len(unit)
		pattern = strings.Repeat(unit, reps)
	case StyleDouble:
		pattern = strings.Repeat("═", width)
	case StyleArrow:
		unit := "»"
		reps := width / len(unit)
		pattern = strings.Repeat(unit, reps)
	case StyleZigzag:
		unit := "/\\"
		reps := width / len(unit)
		pattern = strings.Repeat(unit, reps)
	default:
		pattern = strings.Repeat("─", width)
	}

	return AnsiColor(color, pattern)
}

// GradientDivider renders a divider with a gradient color
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
	color, err := ParseHex(colorHex)
	if err != nil {
		color = RGB{100, 100, 100}
	}

	titleLen := len([]rune(title)) + 2
	remaining := width - titleLen
	if remaining < 0 {
		remaining = 0
	}

	left := remaining / 2
	right := remaining - left

	leftLine := strings.Repeat("─", left)
	rightLine := strings.Repeat("─", right)

	result := fmt.Sprintf("%s %s %s", leftLine, title, rightLine)
	return AnsiColor(color, result)
}
