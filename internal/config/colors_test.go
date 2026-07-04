package config

import "testing"

func sampleColors() Colors {
	return Colors{
		Primary:    "#FF00FF",
		Secondary:  "#00FFFF",
		Background: "#0D0D0D",
		Foreground: "#FFFFFF",
		Accent:     "#FFD700",
		Error:      "#FF4444",
		Success:    "#00FF88",
		Warning:    "#FFA500",
		Muted:      "#444444",
	}
}

func TestResolveColor_KnownKeys(t *testing.T) {
	c := sampleColors()
	cases := map[string]string{
		"primary":    "#FF00FF",
		"secondary":  "#00FFFF",
		"background": "#0D0D0D",
		"foreground": "#FFFFFF",
		"accent":     "#FFD700",
		"error":      "#FF4444",
		"success":    "#00FF88",
		"warning":    "#FFA500",
		"muted":      "#444444",
	}
	for key, expected := range cases {
		if got := ResolveColor(c, key); got != expected {
			t.Errorf("ResolveColor(%q) = %q, want %q", key, got, expected)
		}
	}
}

func TestResolveColor_UnknownKeyPassesThrough(t *testing.T) {
	c := sampleColors()
	// unrecognized keys (e.g. raw hex strings) should pass through unchanged
	if got := ResolveColor(c, "#123ABC"); got != "#123ABC" {
		t.Errorf("expected unknown key to pass through unchanged, got %q", got)
	}
}

func TestResolveColor_EmptyKey(t *testing.T) {
	c := sampleColors()
	if got := ResolveColor(c, ""); got != "" {
		t.Errorf("expected empty key to pass through as empty string, got %q", got)
	}
}
