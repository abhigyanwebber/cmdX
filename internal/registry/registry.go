package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	RegistryBase  = "https://raw.githubusercontent.com/abhigyanwebber/cmdX-themes/main"
	RegistryIndex = RegistryBase + "/index.json"
)

// ThemeEntry represents a theme in the registry index
type ThemeEntry struct {
	Name        string   `json:"name"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
	URL         string   `json:"url"`
}

// Index is the full registry index
type Index struct {
	UpdatedAt string       `json:"updated_at"`
	Themes    []ThemeEntry `json:"themes"`
}

// client with timeout
var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

// FetchIndex downloads and parses the registry index
func FetchIndex() (*Index, error) {
	resp, err := httpClient.Get(RegistryIndex)
	if err != nil {
		return nil, fmt.Errorf("could not reach registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("registry returned status %d", resp.StatusCode)
	}

	var index Index
	if err := json.NewDecoder(resp.Body).Decode(&index); err != nil {
		return nil, fmt.Errorf("could not parse registry index: %w", err)
	}

	return &index, nil
}

// FetchTheme downloads a theme JSON from the registry and saves it
func FetchTheme(name string, themesDir string) error {
	url := fmt.Sprintf("%s/themes/%s.json", RegistryBase, name)

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("could not download theme: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("theme '%s' not found in registry", name)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	outPath := filepath.Join(themesDir, name+".json")
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return fmt.Errorf("could not save theme: %w", err)
	}

	return nil
}

// Search filters the index by name or tag
func Search(index *Index, query string) []ThemeEntry {
	var results []ThemeEntry
	for _, t := range index.Themes {
		if contains(t.Name, query) || contains(t.Description, query) {
			results = append(results, t)
			continue
		}
		for _, tag := range t.Tags {
			if contains(tag, query) {
				results = append(results, t)
				break
			}
		}
	}
	return results
}

func contains(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	return len(s) >= len(sub) && (s == sub ||
		len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
