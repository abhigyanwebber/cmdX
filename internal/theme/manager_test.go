package theme

import (
	"os"
	"path/filepath"
	"testing"
)

const testThemeJSON = `{
	"meta": {"name": "%s", "version": "1.0.0", "author": "tester", "description": "test theme"},
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

func setupThemesDir(t *testing.T, names ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range names {
		content := []byte(strReplace(testThemeJSON, "%s", name))
		path := filepath.Join(dir, name+".json")
		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("could not write theme file: %v", err)
		}
	}
	return dir
}

func strReplace(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

func TestNewManager_ValidDir(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, err := NewManager(dir)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if m.ThemesDir != dir {
		t.Errorf("expected ThemesDir %s, got %s", dir, m.ThemesDir)
	}
}

func TestNewManager_NonexistentDir(t *testing.T) {
	_, err := NewManager("/nonexistent/themes/dir")
	if err == nil {
		t.Fatal("expected error for nonexistent themes directory, got nil")
	}
}

func TestManager_Load_ValidTheme(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)

	th, err := m.Load("cyberpunk")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if th.Meta.Name != "cyberpunk" {
		t.Errorf("expected theme name 'cyberpunk', got %q", th.Meta.Name)
	}
}

func TestManager_Load_NotFound(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)

	_, err := m.Load("nonexistent-theme")
	if err == nil {
		t.Fatal("expected error for missing theme, got nil")
	}
}

func TestManager_Load_InvalidTheme(t *testing.T) {
	dir := t.TempDir()
	// missing required fields -> should fail validation
	os.WriteFile(filepath.Join(dir, "broken.json"), []byte(`{"meta":{"name":"broken"}}`), 0644)
	m, _ := NewManager(dir)

	_, err := m.Load("broken")
	if err == nil {
		t.Fatal("expected validation error for incomplete theme, got nil")
	}
}

func TestManager_Apply_SetsActiveTheme(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)

	if err := m.Apply("cyberpunk"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if m.ActiveTheme == nil {
		t.Fatal("expected ActiveTheme to be set after Apply")
	}
	if m.ActiveTheme.Meta.Name != "cyberpunk" {
		t.Errorf("expected active theme 'cyberpunk', got %q", m.ActiveTheme.Meta.Name)
	}
}

func TestManager_Apply_NonexistentTheme(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)

	if err := m.Apply("ghost-theme"); err == nil {
		t.Fatal("expected error applying nonexistent theme, got nil")
	}
}

func TestManager_List_ReturnsAllThemes(t *testing.T) {
	dir := setupThemesDir(t, "alpha", "beta", "gamma")
	m, _ := NewManager(dir)

	names, err := m.List()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 themes, got %d: %v", len(names), names)
	}
}

func TestManager_List_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	m, _ := NewManager(dir)

	_, err := m.List()
	if err == nil {
		t.Fatal("expected error for empty themes directory, got nil")
	}
}

func TestManager_GetActive_NoneApplied(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)

	_, err := m.GetActive()
	if err == nil {
		t.Fatal("expected error when no theme has been applied, got nil")
	}
}

func TestManager_GetActive_AfterApply(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)
	m.Apply("cyberpunk")

	th, err := m.GetActive()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if th.Meta.Name != "cyberpunk" {
		t.Errorf("expected 'cyberpunk', got %q", th.Meta.Name)
	}
}

func TestManager_Exists(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)

	if !m.Exists("cyberpunk") {
		t.Error("expected Exists to return true for existing theme")
	}
	if m.Exists("ghost") {
		t.Error("expected Exists to return false for nonexistent theme")
	}
}

func TestManager_ResolveColor_ValidKey(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)
	m.Apply("cyberpunk")

	val, err := m.ResolveColor("primary")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if val != "#FF00FF" {
		t.Errorf("expected #FF00FF, got %s", val)
	}
}

func TestManager_ResolveColor_InvalidKey(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)
	m.Apply("cyberpunk")

	_, err := m.ResolveColor("not-a-real-color")
	if err == nil {
		t.Fatal("expected error for unknown color key, got nil")
	}
}

func TestManager_ResolveColor_NoActiveTheme(t *testing.T) {
	dir := setupThemesDir(t, "cyberpunk")
	m, _ := NewManager(dir)

	_, err := m.ResolveColor("primary")
	if err == nil {
		t.Fatal("expected error when no theme applied, got nil")
	}
}
