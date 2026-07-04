package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateImagePath_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "frame.png")
	os.WriteFile(path, []byte("fake png data"), 0644)

	if err := validateImagePath(path); err != nil {
		t.Fatalf("expected no error for valid file, got: %v", err)
	}
}

func TestValidateImagePath_EmptyPath(t *testing.T) {
	if err := validateImagePath(""); err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestValidateImagePath_RejectsLeadingDash(t *testing.T) {
	cases := []string{"-o/etc/passwd", "--exec=evil", "-rf"}
	for _, p := range cases {
		if err := validateImagePath(p); err == nil {
			t.Errorf("expected error for flag-like path %q, got nil", p)
		}
	}
}

func TestValidateImagePath_NonexistentFile(t *testing.T) {
	if err := validateImagePath("/nonexistent/path/frame.png"); err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestValidateImagePath_RejectsDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := validateImagePath(dir); err == nil {
		t.Fatal("expected error when path is a directory, got nil")
	}
}

func TestRender_RejectsFlagInjectionPath(t *testing.T) {
	if !ChafaAvailable() {
		t.Skip("chafa not installed in this environment")
	}
	_, err := Render("-x", ChafaOptions{})
	if err == nil {
		t.Fatal("expected error for flag-injection-shaped path, got nil")
	}
}

func TestRender_RejectsNonexistentFile(t *testing.T) {
	if !ChafaAvailable() {
		t.Skip("chafa not installed in this environment")
	}
	_, err := Render("/nonexistent/asset/frame.png", ChafaOptions{})
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}
