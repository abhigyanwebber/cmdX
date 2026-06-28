package graphics

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
)

// GlitchText applies a glitch effect to text
func GlitchText(text string, intensity float64) string {
	glitchChars := []rune("!@#$%^&*<>[]{}|~`")
	runes := []rune(text)
	result := make([]rune, len(runes))
	src := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i, ch := range runes {
		if ch != ' ' && src.Float64() < intensity {
			result[i] = glitchChars[src.Intn(len(glitchChars))]
		} else {
			result[i] = ch
		}
	}
	return string(result)
}

// GlitchFrames generates multiple frames of glitch animation
func GlitchFrames(text string, frameCount int) []string {
	frames := make([]string, frameCount)
	intensities := []float64{0.3, 0.6, 0.2, 0.5, 0.1, 0.0}
	for i := 0; i < frameCount; i++ {
		frames[i] = GlitchText(text, intensities[i%len(intensities)])
	}
	return frames
}

// PulseText wraps text with a brightness pulse effect
func PulseText(text string, frame int) string {
	styles := []string{"\033[2m", "\033[0m", "\033[1m", "\033[0m", "\033[2m"}
	style := styles[frame%len(styles)]
	return fmt.Sprintf("%s%s\033[0m", style, text)
}

// NeonText applies a neon glow simulation
func NeonText(text string, colorHex string) (string, error) {
	c, err := ParseHex(colorHex)
	if err != nil {
		return text, err
	}

	// create dimmed glow version
	h, s, v := c.Hsv()
	glow := colorful.Hsv(h, s*0.5, v*0.5)

	runes := []rune(text)
	var result strings.Builder

	for i, ch := range runes {
		var col colorful.Color
		if i == 0 || i == len(runes)-1 {
			col = glow
		} else {
			col = c
		}
		rgb := ToRGB(col)
		result.WriteString(AnsiColor(rgb, string(ch)))
	}

	return "\033[1m" + result.String() + "\033[0m", nil
}

// TypewriterPrint prints text character by character with a delay
func TypewriterPrint(text string, delayMs int) {
	for _, ch := range text {
		fmt.Print(string(ch))
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
	fmt.Println()
}
