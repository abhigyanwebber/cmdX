package config

// ResolveColor maps a semantic color key ("primary", "accent", etc.) to its
// hex value from the given Colors set. If the key is not a recognized
// semantic name, it is returned unchanged — this lets callers pass through
// raw hex strings (e.g. "#FF00FF") transparently alongside named keys.
//
// This is the single source of truth for color-key resolution. Do not
// duplicate this map elsewhere; import and call this function instead.
func ResolveColor(c Colors, key string) string {
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
