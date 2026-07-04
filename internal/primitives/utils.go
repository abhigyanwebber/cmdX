package primitives

import "github.com/abhigyanwebber/cmd-customizer/internal/config"

// resolveColor maps a color key to its hex value from the theme.
// Delegates to the centralized config.ResolveColor.
func resolveColor(t *config.Theme, key string) string {
	return config.ResolveColor(t.Colors, key)
}
