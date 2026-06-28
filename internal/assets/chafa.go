package assets

import (
	"fmt"
	"os/exec"
	"strings"
)

// ChafaOptions holds all options for a chafa render call
type ChafaOptions struct {
	Mode      RenderMode
	ColorMode ColorMode
	SymbolSet SymbolSet
	Width     int
	Height    int
	Dither    bool
	Stretch   bool
	Threshold float64
	Animate   bool
	FPS       int
}

// ChafaAvailable checks if chafa is installed on the system
func ChafaAvailable() bool {
	_, err := exec.LookPath("chafa")
	return err == nil
}

// ChafaVersion returns the installed chafa version string
func ChafaVersion() (string, error) {
	out, err := exec.Command("chafa", "--version").Output()
	if err != nil {
		return "", fmt.Errorf("chafa not found: %w", err)
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "", nil
}

// Render converts an image file to terminal output using chafa
func Render(imagePath string, opts ChafaOptions) (string, error) {
	if !ChafaAvailable() {
		return "", fmt.Errorf("chafa is not installed. Install it with: winget install hpjansson.chafa")
	}

	args := buildChafaArgs(imagePath, opts)
	out, err := exec.Command("chafa", args...).Output()
	if err != nil {
		return "", fmt.Errorf("chafa render failed: %w", err)
	}

	return string(out), nil
}

// RenderFrames converts multiple image files to terminal frames
func RenderFrames(imagePaths []string, opts ChafaOptions) ([]string, error) {
	if !ChafaAvailable() {
		return nil, fmt.Errorf("chafa is not installed")
	}

	frames := make([]string, len(imagePaths))
	for i, path := range imagePaths {
		frame, err := Render(path, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to render frame %d (%s): %w", i+1, path, err)
		}
		frames[i] = frame
	}

	return frames, nil
}

// buildChafaArgs constructs the chafa command line arguments
func buildChafaArgs(imagePath string, opts ChafaOptions) []string {
	args := []string{}

	// output mode
	switch opts.Mode {
	case RenderModeBraille:
		args = append(args, "--symbols", "braille")
	case RenderModeBlocks:
		args = append(args, "--symbols", "block")
	case RenderModeASCII:
		args = append(args, "--symbols", "ascii")
	case RenderModeSixel:
		args = append(args, "--format", "sixels")
	case RenderModeColor:
		args = append(args, "--symbols", "block")
	default:
		if opts.SymbolSet != "" {
			args = append(args, "--symbols", string(opts.SymbolSet))
		}
	}

	// color mode
	switch opts.ColorMode {
	case ColorModeNone:
		args = append(args, "--colors", "none")
	case ColorModeAnsi:
		args = append(args, "--colors", "16")
	case ColorMode256:
		args = append(args, "--colors", "256")
	case ColorModeTrueColor:
		args = append(args, "--colors", "full")
	default:
		args = append(args, "--colors", "full")
	}

	// dimensions
	if opts.Width > 0 && opts.Height > 0 {
		args = append(args, "--size", fmt.Sprintf("%dx%d", opts.Width, opts.Height))
	} else if opts.Width > 0 {
		args = append(args, "--size", fmt.Sprintf("%d", opts.Width))
	}

	// dithering
	if opts.Dither {
		args = append(args, "--dither", "ordered")
	}

	// stretch
	if opts.Stretch {
		args = append(args, "--stretch")
	}

	// threshold
	if opts.Threshold > 0 {
		args = append(args, "--threshold", fmt.Sprintf("%.2f", opts.Threshold))
	}

	// animate GIF
	if opts.Animate {
		args = append(args, "--animate", "on")
		if opts.FPS > 0 {
			args = append(args, "--speed", fmt.Sprintf("%d", opts.FPS))
		}
	}

	args = append(args, imagePath)
	return args
}

// DefaultRenderConfig returns sensible defaults for a given render mode
// DefaultRenderConfig returns the best render config for the current terminal
func DefaultRenderConfig(mode RenderMode) RenderConfig {
	// if mode is not specified use best available
	if mode == "" {
		mode = BestRenderMode()
	}

	colorMode := BestColorMode()

	base := RenderConfig{
		Mode:      mode,
		ColorMode: colorMode,
		Dither:    false,
		Stretch:   false,
		Threshold: 0.5,
	}

	switch mode {
	case RenderModeBraille:
		base.Width = 8
		base.Height = 4
		base.SymbolSet = SymbolSetBraille
	case RenderModeBlocks:
		base.Width = 10
		base.Height = 4
		base.SymbolSet = SymbolSetBlock
	case RenderModeASCII:
		base.Width = 16
		base.Height = 4
		base.SymbolSet = SymbolSetAscii
		base.ColorMode = ColorModeNone
	case RenderModeSixel:
		base.Width = 32
		base.Height = 16
	case RenderModeColor:
		base.Width = 12
		base.Height = 6
		base.SymbolSet = SymbolSetBlock
	}

	return base
}
