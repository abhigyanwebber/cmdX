package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetThemesDir_EnvVarOverride(t *testing.T) {
	custom := filepath.Join(t.TempDir(), "my-custom-themes")
	t.Setenv("CMDX_THEMES_DIR", custom)

	got := getThemesDir()
	if got != custom {
		t.Errorf("expected CMDX_THEMES_DIR override %q, got %q", custom, got)
	}
}

func TestGetThemesDir_FallsBackWithoutEnvVar(t *testing.T) {
	t.Setenv("CMDX_THEMES_DIR", "")
	os.Unsetenv("CMDX_THEMES_DIR")

	got := getThemesDir()
	if got == "" {
		t.Error("expected a non-empty fallback path when env var is unset")
	}
}

func TestGetAssetsDir_EnvVarOverride(t *testing.T) {
	custom := filepath.Join(t.TempDir(), "my-custom-assets")
	t.Setenv("CMDX_ASSETS_DIR", custom)

	got := getAssetsDir()
	if got != custom {
		t.Errorf("expected CMDX_ASSETS_DIR override %q, got %q", custom, got)
	}
}

func TestGetAssetsDir_FallsBackWithoutEnvVar(t *testing.T) {
	t.Setenv("CMDX_ASSETS_DIR", "")
	os.Unsetenv("CMDX_ASSETS_DIR")

	got := getAssetsDir()
	if got == "" {
		t.Error("expected a non-empty fallback path when env var is unset")
	}
}
