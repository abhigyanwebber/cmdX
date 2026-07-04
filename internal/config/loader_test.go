package config

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleThemeJSON = `{
	"meta": {"name": "sample", "version": "1.0.0", "author": "tester", "description": "a test theme"},
	"colors": {
		"primary": "#FF00FF", "secondary": "#00FFFF", "background": "#0D0D0D",
		"foreground": "#FFFFFF", "accent": "#FFD700", "error": "#FF4444",
		"success": "#00FF88", "warning": "#FFA500", "muted": "#444444"
	},
	"prompt": {"symbol": "▶", "format": "{user}@{dir}", "style": "single"},
	"loader": {"frames": ["◐", "◓"], "interval_ms": 100},
	"progress_bar": {"filled": "█", "empty": "░", "width": 20},
	"cursor": {"style": "block"}
}`

func writeTempTheme(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("could not write temp theme file: %v", err)
	}
	return path
}

func TestLoadTheme_ValidFile(t *testing.T) {
	path := writeTempTheme(t, sampleThemeJSON)

	th, err := LoadTheme(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if th.Meta.Name != "sample" {
		t.Errorf("expected name 'sample', got %q", th.Meta.Name)
	}
	if th.Colors.Primary != "#FF00FF" {
		t.Errorf("expected primary color #FF00FF, got %q", th.Colors.Primary)
	}
}

func TestLoadTheme_FileNotFound(t *testing.T) {
	_, err := LoadTheme("/nonexistent/path/theme.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestLoadTheme_InvalidJSON(t *testing.T) {
	path := writeTempTheme(t, `{not valid json,,,`)

	_, err := LoadTheme(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoadTheme_RelativePath(t *testing.T) {
	path := writeTempTheme(t, sampleThemeJSON)
	abs, _ := filepath.Abs(path)

	th, err := LoadTheme(abs)
	if err != nil {
		t.Fatalf("expected no error loading absolute path, got: %v", err)
	}
	if th.Meta.Name != "sample" {
		t.Errorf("expected name 'sample', got %q", th.Meta.Name)
	}
}

func TestLoadThemeFromBytes_Valid(t *testing.T) {
	th, err := LoadThemeFromBytes([]byte(sampleThemeJSON))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if th.Meta.Author != "tester" {
		t.Errorf("expected author 'tester', got %q", th.Meta.Author)
	}
}

func TestLoadThemeFromBytes_InvalidJSON(t *testing.T) {
	_, err := LoadThemeFromBytes([]byte(`{broken`))
	if err == nil {
		t.Fatal("expected error for invalid JSON bytes, got nil")
	}
}

func TestLoadThemeFromBytes_EmptyBytes(t *testing.T) {
	_, err := LoadThemeFromBytes([]byte(``))
	if err == nil {
		t.Fatal("expected error for empty bytes, got nil")
	}
}

func TestListThemes_FindsJSONFiles(t *testing.T) {
	dir := t.TempDir()

	// write two theme files and one non-json file
	os.WriteFile(filepath.Join(dir, "alpha.json"), []byte(sampleThemeJSON), 0644)
	os.WriteFile(filepath.Join(dir, "beta.json"), []byte(sampleThemeJSON), 0644)
	os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("not a theme"), 0644)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not read temp dir: %v", err)
	}

	var jsonCount int
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			jsonCount++
		}
	}

	if jsonCount != 2 {
		t.Errorf("expected 2 json theme files, found %d", jsonCount)
	}
}

func TestGetThemesDir_ReturnsExecutableRelativePath(t *testing.T) {
	dir, err := GetThemesDir()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if filepath.Base(dir) != "themes" {
		t.Errorf("expected dir to end in 'themes', got: %s", dir)
	}
}
