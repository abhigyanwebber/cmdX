package graphics

import (
	"strings"
)

// Pattern defines a repeating visual texture
type Pattern string

const (
	PatternDots     Pattern = "dots"
	PatternGrid     Pattern = "grid"
	PatternDiagonal Pattern = "diagonal"
	PatternBricks   Pattern = "bricks"
	PatternCircuit  Pattern = "circuit"
	PatternStars    Pattern = "stars"
)

// RenderPattern renders a pattern block of given width and height
func RenderPattern(p Pattern, width int, height int, colorHex string) string {
	color, err := ParseHex(colorHex)
	if err != nil {
		color = RGB{80, 80, 80}
	}

	var lines []string

	for row := 0; row < height; row++ {
		var line strings.Builder

		for col := 0; col < width; col++ {
			ch := patternChar(p, row, col)
			line.WriteString(AnsiColor(color, string(ch)))
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
