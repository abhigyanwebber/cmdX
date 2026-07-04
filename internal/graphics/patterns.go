package graphics

import "strings"

// Pattern identifies a background fill pattern style usable in theme
// previews and wallpaper-like terminal decoration.
type Pattern string

// Available background patterns.
const (
	PatternDots     Pattern = "dots"
	PatternGrid     Pattern = "grid"
	PatternDiagonal Pattern = "diagonal"
	PatternBricks   Pattern = "bricks"
	PatternCircuit  Pattern = "circuit"
	PatternStars    Pattern = "stars"
)

// RenderPattern renders a width x height block of the given pattern,
// colored with colorHex, as a multi-line string with embedded ANSI
// escape codes. Falls back to plain spaces if colorHex is invalid.
func RenderPattern(p Pattern, width int, height int, colorHex string) string {
	c, err := ParseHex(colorHex)
	if err != nil {
		return strings.Repeat(" ", width*height)
	}
	rgb := ToRGB(c)

	var lines []string
	for row := 0; row < height; row++ {
		var line strings.Builder
		for col := 0; col < width; col++ {
			ch := patternChar(p, row, col)
			line.WriteString(AnsiColor(rgb, string(ch)))
		}
		lines = append(lines, line.String())
	}
	return strings.Join(lines, "\n")
}

func patternChar(p Pattern, row, col int) rune {
	switch p {
	case PatternDots:
		if row%2 == 0 && col%4 == 0 {
			return '·'
		}
		return ' '
	case PatternGrid:
		if row%2 == 0 || col%4 == 0 {
			return '+'
		}
		return ' '
	case PatternDiagonal:
		if (row+col)%4 == 0 {
			return '\\'
		}
		if (row-col+100)%4 == 0 {
			return '/'
		}
		return ' '
	case PatternBricks:
		if row%2 == 0 {
			if col%6 == 0 {
				return '│'
			}
		} else {
			if col%6 == 3 {
				return '│'
			}
		}
		if row%2 == 1 {
			return '─'
		}
		return ' '
	case PatternCircuit:
		chars := [][]rune{
			{' ', ' ', '─', '─', '┐', ' '},
			{' ', '┌', '─', '─', '┘', ' '},
			{'─', '┘', ' ', '┌', '─', '─'},
			{' ', ' ', ' ', '└', '─', '┐'},
		}
		return chars[row%len(chars)][col%6]
	case PatternStars:
		if (row*3+col)%7 == 0 {
			return '✦'
		}
		if (row*2+col)%11 == 0 {
			return '·'
		}
		return ' '
	default:
		return ' '
	}
}
