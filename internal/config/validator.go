package config

import (
	"fmt"
	"regexp"
)

var hexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)

// ValidateTheme checks a parsed Theme for missing or invalid fields
func ValidateTheme(t *Theme) error {
	if err := validateMeta(t.Meta); err != nil {
		return fmt.Errorf("meta: %w", err)
	}
	if err := validateColors(t.Colors); err != nil {
		return fmt.Errorf("colors: %w", err)
	}
	if err := validatePrompt(t.Prompt); err != nil {
		return fmt.Errorf("prompt: %w", err)
	}
	if err := validateLoader(t.Loader); err != nil {
		return fmt.Errorf("loader: %w", err)
	}
	if err := validateProgressBar(t.ProgressBar); err != nil {
		return fmt.Errorf("progress_bar: %w", err)
	}
	if err := validateCursor(t.Cursor); err != nil {
		return fmt.Errorf("cursor: %w", err)
	}
	return nil
}

func validateMeta(m Meta) error {
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	if m.Version == "" {
		return fmt.Errorf("version is required")
	}
	return nil
}

func validateColors(c Colors) error {
	fields := map[string]string{
		"primary":    c.Primary,
		"secondary":  c.Secondary,
		"background": c.Background,
		"foreground": c.Foreground,
		"accent":     c.Accent,
		"error":      c.Error,
		"success":    c.Success,
		"warning":    c.Warning,
		"muted":      c.Muted,
	}

	for name, value := range fields {
		if value == "" {
			return fmt.Errorf("%s color is required", name)
		}
		if !hexColorRegex.MatchString(value) {
			return fmt.Errorf("%s must be a valid hex color (e.g. #FF00FF), got: %s", name, value)
		}
	}
	return nil
}

func validatePrompt(p Prompt) error {
	if p.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if p.Format == "" {
		return fmt.Errorf("format is required")
	}
	validStyles := map[string]bool{"single": true, "multiline": true}
	if !validStyles[p.Style] {
		return fmt.Errorf("style must be 'single' or 'multiline', got: %s", p.Style)
	}
	return nil
}

func validateLoader(l Loader) error {
	if len(l.Frames) == 0 {
		return fmt.Errorf("frames cannot be empty")
	}
	if l.IntervalMs <= 0 {
		return fmt.Errorf("interval_ms must be greater than 0")
	}
	return nil
}

func validateProgressBar(p ProgressBar) error {
	if p.Filled == "" {
		return fmt.Errorf("filled character is required")
	}
	if p.Empty == "" {
		return fmt.Errorf("empty character is required")
	}
	if p.Width <= 0 {
		return fmt.Errorf("width must be greater than 0")
	}
	return nil
}

func validateCursor(c Cursor) error {
	validStyles := map[string]bool{"block": true, "bar": true, "underline": true}
	if !validStyles[c.Style] {
		return fmt.Errorf("cursor style must be 'block', 'bar', or 'underline', got: %s", c.Style)
	}
	return nil
}