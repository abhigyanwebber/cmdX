package assets

import (
	"os"

	"github.com/muesli/termenv"
)

// TerminalCapabilities holds what the current terminal supports
type TerminalCapabilities struct {
	ColorLevel        termenv.Profile
	SupportsSixel     bool
	SupportsColor     bool
	Supports256       bool
	SupportsTrueColor bool
}

// DetectCapabilities checks what the current terminal can render
func DetectCapabilities() TerminalCapabilities {
	profile := termenv.ColorProfile()

	return TerminalCapabilities{
		ColorLevel:        profile,
		SupportsColor:     profile > termenv.Ascii,
		Supports256:       profile >= termenv.ANSI256,
		SupportsTrueColor: profile >= termenv.TrueColor,
		SupportsSixel:     detectSixel(),
	}
}

// BestRenderMode returns the best chafa render mode for this terminal
func BestRenderMode() RenderMode {
	caps := DetectCapabilities()

	if caps.SupportsSixel {
		return RenderModeSixel
	}
	if caps.SupportsTrueColor {
		return RenderModeBraille
	}
	if caps.Supports256 {
		return RenderModeBlocks
	}
	return RenderModeASCII
}

// BestColorMode returns the best chafa color mode for this terminal
func BestColorMode() ColorMode {
	caps := DetectCapabilities()

	if caps.SupportsTrueColor {
		return ColorModeTrueColor
	}
	if caps.Supports256 {
		return ColorMode256
	}
	if caps.SupportsColor {
		return ColorModeAnsi
	}
	return ColorModeNone
}

// detectSixel checks for known sixel-supporting terminals
func detectSixel() bool {
	// Windows Terminal
	if os.Getenv("WT_SESSION") != "" {
		return true
	}
	// iTerm2
	if os.Getenv("TERM_PROGRAM") == "iTerm.app" {
		return true
	}
	// Kitty
	if os.Getenv("TERM") == "xterm-kitty" {
		return true
	}
	return false
}
