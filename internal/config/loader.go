package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadTheme reads a JSON theme file from the given path
// and returns a parsed Theme struct
func LoadTheme(path string) (*Theme, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("could not read theme file: %w", err)
	}

	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("invalid JSON in theme file: %w", err)
	}

	return &theme, nil
}

// LoadThemeFromBytes parses a Theme directly from raw JSON bytes
// useful for embedded/default themes
func LoadThemeFromBytes(data []byte) (*Theme, error) {
	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme: %w", err)
	}
	return &theme, nil
}

// GetThemesDir returns the path to the themes directory
// relative to the binary location
func GetThemesDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not find executable path: %w", err)
	}
	return filepath.Join(filepath.Dir(exe), "themes"), nil
}

// ListThemes returns all available .json theme files from the themes directory
func ListThemes() ([]string, error) {
	dir, err := GetThemesDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read themes directory: %w", err)
	}

	var themes []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			name := entry.Name()[:len(entry.Name())-5] // strip .json
			themes = append(themes, name)
		}
	}

	return themes, nil
}