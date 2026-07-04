package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writePluginManifest(t *testing.T, dir string, p Plugin) {
	t.Helper()
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("could not marshal plugin: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "plugin.json"), data, 0644); err != nil {
		t.Fatalf("could not write plugin.json: %v", err)
	}
}

func validPlugin(name string) Plugin {
	return Plugin{
		Meta: PluginMeta{Name: name, Version: "1.0.0", Author: "tester"},
		Spinners: []SpinnerSet{
			{Name: "dots", Frames: []string{"◐", "◓", "◑", "◒"}, IntervalMs: 100},
		},
	}
}

func TestNewManager_CreatesPluginsDir(t *testing.T) {
	base := t.TempDir()
	pluginsDir := filepath.Join(base, "plugins")

	m, err := NewManager(pluginsDir)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		t.Error("expected plugins directory to be created")
	}
	if m.Loaded == nil {
		t.Error("expected Loaded map to be initialized")
	}
}

func TestNewManager_ExistingDir(t *testing.T) {
	dir := t.TempDir()
	m, err := NewManager(dir)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if m.PluginsDir != dir {
		t.Errorf("expected PluginsDir %s, got %s", dir, m.PluginsDir)
	}
}

func TestDiscover_LoadsValidPlugins(t *testing.T) {
	base := t.TempDir()
	pluginDir := filepath.Join(base, "example-plugin")
	os.MkdirAll(pluginDir, 0755)
	writePluginManifest(t, pluginDir, validPlugin("example-plugin"))

	m, _ := NewManager(base)
	if err := m.Discover(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(m.Loaded) != 1 {
		t.Fatalf("expected 1 loaded plugin, got %d", len(m.Loaded))
	}
	if _, ok := m.Loaded["example-plugin"]; !ok {
		t.Error("expected plugin 'example-plugin' to be loaded")
	}
}

func TestDiscover_SkipsInvalidPlugins(t *testing.T) {
	base := t.TempDir()

	validDir := filepath.Join(base, "valid-plugin")
	os.MkdirAll(validDir, 0755)
	writePluginManifest(t, validDir, validPlugin("valid-plugin"))

	brokenDir := filepath.Join(base, "broken-plugin")
	os.MkdirAll(brokenDir, 0755)
	os.WriteFile(filepath.Join(brokenDir, "plugin.json"), []byte("{not json"), 0644)

	m, _ := NewManager(base)
	if err := m.Discover(); err != nil {
		t.Fatalf("expected no top-level error even with broken plugin, got: %v", err)
	}

	if len(m.Loaded) != 1 {
		t.Errorf("expected only the valid plugin to load, got %d loaded", len(m.Loaded))
	}
}

func TestDiscover_IgnoresNonDirEntries(t *testing.T) {
	base := t.TempDir()
	os.WriteFile(filepath.Join(base, "readme.txt"), []byte("not a plugin"), 0644)

	m, _ := NewManager(base)
	if err := m.Discover(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(m.Loaded) != 0 {
		t.Errorf("expected 0 loaded plugins, got %d", len(m.Loaded))
	}
}

func TestGet_FoundPlugin(t *testing.T) {
	base := t.TempDir()
	pluginDir := filepath.Join(base, "myplugin")
	os.MkdirAll(pluginDir, 0755)
	writePluginManifest(t, pluginDir, validPlugin("myplugin"))

	m, _ := NewManager(base)
	m.Discover()

	p, err := m.Get("myplugin")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if p.Meta.Name != "myplugin" {
		t.Errorf("expected plugin name 'myplugin', got %q", p.Meta.Name)
	}
}

func TestGet_NotFound(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	_, err := m.Get("ghost-plugin")
	if err == nil {
		t.Fatal("expected error for missing plugin, got nil")
	}
}

func TestList_ReturnsLoadedNames(t *testing.T) {
	base := t.TempDir()
	for _, name := range []string{"alpha", "beta"} {
		dir := filepath.Join(base, name)
		os.MkdirAll(dir, 0755)
		writePluginManifest(t, dir, validPlugin(name))
	}

	m, _ := NewManager(base)
	m.Discover()

	names := m.List()
	if len(names) != 2 {
		t.Errorf("expected 2 plugin names, got %d", len(names))
	}
}

func TestGetSpinner_Found(t *testing.T) {
	base := t.TempDir()
	pluginDir := filepath.Join(base, "spinplug")
	os.MkdirAll(pluginDir, 0755)
	writePluginManifest(t, pluginDir, validPlugin("spinplug"))

	m, _ := NewManager(base)
	m.Discover()

	s, err := m.GetSpinner("dots")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if s.Name != "dots" {
		t.Errorf("expected spinner 'dots', got %q", s.Name)
	}
}

func TestGetSpinner_NotFound(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	_, err := m.GetSpinner("nonexistent-spinner")
	if err == nil {
		t.Fatal("expected error for missing spinner, got nil")
	}
}

func TestGetBanner_NotFound(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	_, err := m.GetBanner("nonexistent-banner")
	if err == nil {
		t.Fatal("expected error for missing banner, got nil")
	}
}

func TestGetPrompt_NotFound(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	_, err := m.GetPrompt("nonexistent-prompt")
	if err == nil {
		t.Fatal("expected error for missing prompt, got nil")
	}
}

func TestValidatePlugin_Valid(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	p := validPlugin("ok-plugin")
	if err := m.ValidatePlugin(&p); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidatePlugin_MissingName(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	p := validPlugin("")
	if err := m.ValidatePlugin(&p); err == nil {
		t.Fatal("expected error for missing plugin name, got nil")
	}
}

func TestValidatePlugin_MissingVersion(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	p := validPlugin("noversion")
	p.Meta.Version = ""
	if err := m.ValidatePlugin(&p); err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
}

func TestValidatePlugin_SpinnerEmptyFrames(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	p := validPlugin("badspinner")
	p.Spinners[0].Frames = []string{}
	if err := m.ValidatePlugin(&p); err == nil {
		t.Fatal("expected error for spinner with no frames, got nil")
	}
}

func TestValidatePlugin_SpinnerZeroInterval(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	p := validPlugin("slowspinner")
	p.Spinners[0].IntervalMs = 0
	if err := m.ValidatePlugin(&p); err == nil {
		t.Fatal("expected error for zero interval_ms, got nil")
	}
}

func TestLoadPlugin_MissingManifest(t *testing.T) {
	_, err := LoadPlugin(t.TempDir())
	if err == nil {
		t.Fatal("expected error for missing plugin.json, got nil")
	}
}

func TestLoadSpinners_FileNotExist(t *testing.T) {
	spinners, err := LoadSpinners(t.TempDir())
	if err != nil {
		t.Fatalf("expected no error when spinners.json absent, got: %v", err)
	}
	if spinners != nil {
		t.Errorf("expected nil spinners slice, got %v", spinners)
	}
}

func TestLoadBanners_FileNotExist(t *testing.T) {
	banners, err := LoadBanners(t.TempDir())
	if err != nil {
		t.Fatalf("expected no error when banners.json absent, got: %v", err)
	}
	if banners != nil {
		t.Errorf("expected nil banners slice, got %v", banners)
	}
}

func TestLoadPrompts_FileNotExist(t *testing.T) {
	prompts, err := LoadPrompts(t.TempDir())
	if err != nil {
		t.Fatalf("expected no error when prompts.json absent, got: %v", err)
	}
	if prompts != nil {
		t.Errorf("expected nil prompts slice, got %v", prompts)
	}
}
