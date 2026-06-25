package wallpaper

import (
	"fmt"
	"os"
	"runtime"
)

// Engine handles wallpaper setting for different terminal emulators
type Engine struct {
	OS string
}

// NewEngine creates a wallpaper engine for the current OS
func NewEngine() *Engine {
	return &Engine{OS: runtime.GOOS}
}

// Apply sets the wallpaper for the detected terminal
func (e *Engine) Apply(imagePath string, opacity float64, stretch string, alignment string) error {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image not found: %s", imagePath)
	}

	if opacity < 0 || opacity > 1 {
		return fmt.Errorf("opacity must be between 0.0 and 1.0")
	}

	switch e.OS {
	case "windows":
		return applyWindowsTerminal(imagePath, opacity, stretch, alignment)
	case "darwin":
		return applyITerm2(imagePath, opacity)
	case "linux":
		return applyKitty(imagePath, opacity)
	default:
		return fmt.Errorf("unsupported OS: %s", e.OS)
	}
}

// Remove clears the wallpaper from the terminal config
func (e *Engine) Remove() error {
	switch e.OS {
	case "windows":
		return removeWindowsTerminal()
	case "darwin":
		return removeITerm2()
	case "linux":
		return removeKitty()
	default:
		return fmt.Errorf("unsupported OS: %s", e.OS)
	}
}

// DetectedTerminal returns the terminal emulator name for the current OS
func (e *Engine) DetectedTerminal() string {
	switch e.OS {
	case "windows":
		return "Windows Terminal"
	case "darwin":
		return "iTerm2"
	case "linux":
		return "Kitty"
	default:
		return "Unknown"
	}
}
