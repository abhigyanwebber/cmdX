package plugin

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager handles plugin discovery and loading
type Manager struct {
	PluginsDir string
	Loaded     map[string]*Plugin
}

// NewManager creates a plugin manager for the given directory
func NewManager(pluginsDir string) (*Manager, error) {
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		// create plugins dir if it doesn't exist
		if err := os.MkdirAll(pluginsDir, 0755); err != nil {
			return nil, fmt.Errorf("could not create plugins directory: %w", err)
		}
	}

	return &Manager{
		PluginsDir: pluginsDir,
		Loaded:     make(map[string]*Plugin),
	}, nil
}

// Discover scans the plugins directory and loads all valid plugins
func (m *Manager) Discover() error {
	entries, err := os.ReadDir(m.PluginsDir)
	if err != nil {
		return fmt.Errorf("could not read plugins directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginDir := filepath.Join(m.PluginsDir, entry.Name())
		p, err := LoadPlugin(pluginDir)
		if err != nil {
			fmt.Printf("  ! Skipping plugin '%s': %s\n", entry.Name(), err)
			continue
		}

		m.Loaded[p.Meta.Name] = p
	}

	return nil
}

// Get returns a loaded plugin by name
func (m *Manager) Get(name string) (*Plugin, error) {
	p, ok := m.Loaded[name]
	if !ok {
		return nil, fmt.Errorf("plugin '%s' not found", name)
	}
	return p, nil
}

// List returns all loaded plugin names
func (m *Manager) List() []string {
	names := make([]string, 0, len(m.Loaded))
	for name := range m.Loaded {
		names = append(names, name)
	}
	return names
}

// GetSpinner finds a named spinner set across all loaded plugins
func (m *Manager) GetSpinner(name string) (*SpinnerSet, error) {
	for _, p := range m.Loaded {
		for _, s := range p.Spinners {
			if s.Name == name {
				return &s, nil
			}
		}
	}
	return nil, fmt.Errorf("spinner '%s' not found in any plugin", name)
}

// GetBanner finds a named banner template across all loaded plugins
func (m *Manager) GetBanner(name string) (*BannerTemplate, error) {
	for _, p := range m.Loaded {
		for _, b := range p.Banners {
			if b.Name == name {
				return &b, nil
			}
		}
	}
	return nil, fmt.Errorf("banner '%s' not found in any plugin", name)
}

// GetPrompt finds a named prompt template across all loaded plugins
func (m *Manager) GetPrompt(name string) (*PromptTemplate, error) {
	for _, p := range m.Loaded {
		for _, pr := range p.Prompts {
			if pr.Name == name {
				return &pr, nil
			}
		}
	}
	return nil, fmt.Errorf("prompt '%s' not found in any plugin", name)
}

// ValidatePlugin checks a plugin for required fields
func (m *Manager) ValidatePlugin(p *Plugin) error {
	if p.Meta.Name == "" {
		return fmt.Errorf("plugin name is required")
	}
	if p.Meta.Version == "" {
		return fmt.Errorf("plugin version is required")
	}
	for _, s := range p.Spinners {
		if len(s.Frames) == 0 {
			return fmt.Errorf("spinner '%s' must have at least one frame", s.Name)
		}
		if s.IntervalMs <= 0 {
			return fmt.Errorf("spinner '%s' interval_ms must be greater than 0", s.Name)
		}
	}
	return nil
}
