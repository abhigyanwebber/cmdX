package graphics

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GlitchText applies a glitch effect to text
// randomly replaces some characters with glitch chars
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
		intensity := intensities[i%len(intensities)]
		frames[i] = GlitchText(text, intensity)
	}
	return frames
}

// PulseText wraps text with a brightness pulse effect using bold/dim
func PulseText(text string, frame int) string {
	styles := []string{
		"\033[2m", // dim
		"\033[0m", // normal
		"\033[1m", // bold
		"\033[0m", // normal
		"\033[2m", // dim
	}
	style := styles[frame%len(styles)]
	return fmt.Sprintf("%s%s\033[0m", style, text)
}

// NeonText applies a neon glow simulation using color layering
func NeonText(text string, colorHex string) (string, error) {
	color, err := ParseHex(colorHex)
	if err != nil {
		return text, err
	}

	// create a slightly dimmed version for the glow effect
	glow := RGB{
		R: uint8(float64(color.R) * 0.5),
		G: uint8(float64(color.G) * 0.5),
		B: uint8(float64(color.B) * 0.5),
	}

	var result strings.Builder
	runes := []rune(text)

	for i, ch := range runes {
		if i == 0 || i == len(runes)-1 {
			result.WriteString(AnsiColor(glow, string(ch)))
		} else {
			result.WriteString(AnsiColor(color, string(ch)))
		}
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
