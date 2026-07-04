package preview

import "github.com/abhigyanwebber/cmd-customizer/internal/config"

// resolve maps a color key to its hex value from the theme.
// Delegates to the centralized config.ResolveColor.
func (p *Preview) resolve(key string) string {
	return config.ResolveColor(p.theme.Colors, key)
}
