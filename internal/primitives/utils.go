package primitives

import "github.com/abhigyanwebber/cmd-customizer/internal/config"

// resolveColor maps a color key to its hex value from the theme
func resolveColor(t *config.Theme, key string) string {
	c := t.Colors
	colorMap := map[string]string{
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
	if val, ok := colorMap[key]; ok {
		return val
	}
	return key
}