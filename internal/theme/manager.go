package theme

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
)

// Manager handles all theme operations
type Manager struct {
	ThemesDir   string
	ActiveTheme *config.Theme
}

// NewManager creates a new theme manager
// it looks for themes in the given directory
func NewManager(themesDir string) (*Manager, error) {
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("themes directory not found: %s", themesDir)
	}
	return &Manager{ThemesDir: themesDir}, nil
}

// Load reads and validates a theme by name from the themes directory
func (m *Manager) Load(name string) (*config.Theme, error) {
	path := filepath.Join(m.ThemesDir, name+".json")

	theme, err := config.LoadTheme(path)
	if err != nil {
		return nil, fmt.Errorf("could not load theme '%s': %w", name, err)
	}

	if err := config.ValidateTheme(theme); err != nil {
		return nil, fmt.Errorf("theme '%s' is invalid: %w", name, err)
	}

	return theme, nil
}

// Apply sets a theme as the active theme
func (m *Manager) Apply(name string) error {
	theme, err := m.Load(name)
	if err != nil {
		return err
	}
	m.ActiveTheme = theme
	fmt.Printf("✓ Theme '%s' applied successfully\n", name)
	return nil
}

// List returns all available theme names in the themes directory
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.ThemesDir)
	if err != nil {
		return nil, fmt.Errorf("could not read themes directory: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			name := entry.Name()[:len(entry.Name())-5]
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("no themes found in %s", m.ThemesDir)
	}

	return names, nil
}

// GetActive returns the currently active theme
// returns an error if no theme has been applied yet
func (m *Manager) GetActive() (*config.Theme, error) {
	if m.ActiveTheme == nil {
		return nil, fmt.Errorf("no active theme — run 'cmdx theme apply <name>' first")
	}
	return m.ActiveTheme, nil
}

// Exists checks if a theme with the given name exists
func (m *Manager) Exists(name string) bool {
	path := filepath.Join(m.ThemesDir, name+".json")
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ResolveColor maps a color key like "primary" to its hex value
// from the active theme
func (m *Manager) ResolveColor(key string) (string, error) {
	theme, err := m.GetActive()
	if err != nil {
		return "", err
	}

	c := theme.Colors
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

	val, ok := colorMap[key]
	if !ok {
		return "", fmt.Errorf("unknown color key: '%s'", key)
	}
	return val, nil
}