package assets

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeManifest(t *testing.T, dir string, a Asset) {
	t.Helper()
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("could not marshal asset: %v", err)
	}
	path := filepath.Join(dir, "asset.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("could not write asset.json: %v", err)
	}
}

func validSpinnerAsset() Asset {
	return Asset{
		Name:    "pulse",
		Type:    AssetTypeSpinner,
		Version: "1.0.0",
		Render: RenderConfig{
			Mode:      RenderModeBraille,
			Width:     8,
			Height:    8,
			ColorMode: ColorModeTrueColor,
		},
		Spinner: &SpinnerConfig{
			Frames:     []string{"frame1.png", "frame2.png"},
			IntervalMs: 100,
			Loop:       true,
		},
	}
}

func TestLoadAsset_ValidManifest(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, validSpinnerAsset())

	a, err := LoadAsset(dir)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if a.Name != "pulse" {
		t.Errorf("expected name 'pulse', got %q", a.Name)
	}
	if a.Type != AssetTypeSpinner {
		t.Errorf("expected type spinner, got %q", a.Type)
	}
}

func TestLoadAsset_MissingManifest(t *testing.T) {
	dir := t.TempDir() // no asset.json written
	_, err := LoadAsset(dir)
	if err == nil {
		t.Fatal("expected error for missing asset.json, got nil")
	}
}

func TestLoadAsset_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "asset.json"), []byte("{broken json"), 0644)

	_, err := LoadAsset(dir)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestValidateAsset_MissingName(t *testing.T) {
	a := validSpinnerAsset()
	a.Name = ""
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for missing name, got nil")
	}
}

func TestValidateAsset_MissingType(t *testing.T) {
	a := validSpinnerAsset()
	a.Type = ""
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for missing type, got nil")
	}
}

func TestValidateAsset_MissingVersion(t *testing.T) {
	a := validSpinnerAsset()
	a.Version = ""
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
}

func TestValidateAsset_SpinnerMissingConfig(t *testing.T) {
	a := validSpinnerAsset()
	a.Spinner = nil
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for spinner asset with nil spinner config, got nil")
	}
}

func TestValidateAsset_SpinnerEmptyFrames(t *testing.T) {
	a := validSpinnerAsset()
	a.Spinner.Frames = []string{}
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for spinner with no frames, got nil")
	}
}

func TestValidateAsset_SpinnerZeroInterval(t *testing.T) {
	a := validSpinnerAsset()
	a.Spinner.IntervalMs = 0
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for zero interval_ms, got nil")
	}
}

func TestValidateAsset_SpinnerMissingFrameFile(t *testing.T) {
	dir := t.TempDir()
	a := validSpinnerAsset()
	// frame files don't actually exist on disk
	if err := ValidateAsset(&a, dir); err == nil {
		t.Fatal("expected error for missing frame file on disk, got nil")
	}
}

func TestValidateAsset_SpinnerValidWithRealFrameFiles(t *testing.T) {
	dir := t.TempDir()
	a := validSpinnerAsset()

	// create the actual frame files referenced
	for _, frame := range a.Spinner.Frames {
		os.WriteFile(filepath.Join(dir, frame), []byte("fake png data"), 0644)
	}

	if err := ValidateAsset(&a, dir); err != nil {
		t.Fatalf("expected no error with frame files present, got: %v", err)
	}
}

func TestValidateAsset_IconMissingConfig(t *testing.T) {
	a := Asset{Name: "icons", Type: AssetTypeIcon, Version: "1.0.0"}
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for icon asset with nil icon config, got nil")
	}
}

func TestValidateAsset_IconEmptyFiles(t *testing.T) {
	a := Asset{
		Name: "icons", Type: AssetTypeIcon, Version: "1.0.0",
		Icon: &IconConfig{Files: map[string]string{}},
	}
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for icon asset with no files defined, got nil")
	}
}

func TestValidateAsset_BannerMissingFile(t *testing.T) {
	dir := t.TempDir()
	a := Asset{
		Name: "banner", Type: AssetTypeBanner, Version: "1.0.0",
		Banner: &BannerConfig{File: "missing.png"},
	}
	if err := ValidateAsset(&a, dir); err == nil {
		t.Fatal("expected error for banner with nonexistent file, got nil")
	}
}

func TestValidateAsset_DividerMissingFile(t *testing.T) {
	dir := t.TempDir()
	a := Asset{
		Name: "div", Type: AssetTypeDivider, Version: "1.0.0",
		Divider: &DividerConfig{File: "missing.png"},
	}
	if err := ValidateAsset(&a, dir); err == nil {
		t.Fatal("expected error for divider with nonexistent file, got nil")
	}
}

func TestResolveFramePaths_ValidSpinner(t *testing.T) {
	a := validSpinnerAsset()
	dir := "/some/asset/dir"

	paths, err := ResolveFramePaths(&a, dir)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(paths))
	}
	expected := filepath.Join(dir, "frame1.png")
	if paths[0] != expected {
		t.Errorf("expected %q, got %q", expected, paths[0])
	}
}

func TestResolveFramePaths_NotASpinner(t *testing.T) {
	a := Asset{Name: "icon-set", Type: AssetTypeIcon, Version: "1.0.0"}
	_, err := ResolveFramePaths(&a, "/some/dir")
	if err == nil {
		t.Fatal("expected error for non-spinner asset, got nil")
	}
}

func TestResolveFramePaths_NilSpinnerConfig(t *testing.T) {
	a := Asset{Name: "x", Type: AssetTypeSpinner, Version: "1.0.0", Spinner: nil}
	_, err := ResolveFramePaths(&a, "/some/dir")
	if err == nil {
		t.Fatal("expected error for nil spinner config, got nil")
	}
}

func validFloaterManifest() Asset {
	return Asset{
		Name: "corner-mascot", Type: AssetTypeFloater, Version: "1.0.0",
		Floater: &FloaterConfig{
			File:     "mascot.png",
			Position: FloaterBottomRight,
		},
	}
}

func TestValidateAsset_FloaterMissingConfig(t *testing.T) {
	a := Asset{Name: "x", Type: AssetTypeFloater, Version: "1.0.0", Floater: nil}
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for floater asset with nil Floater config, got nil")
	}
}

func TestValidateAsset_FloaterMissingFile(t *testing.T) {
	a := validFloaterManifest()
	a.Floater.File = ""
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for floater with no file specified, got nil")
	}
}

func TestValidateAsset_FloaterFileNotFound(t *testing.T) {
	a := validFloaterManifest()
	// mascot.png is never written to disk
	if err := ValidateAsset(&a, t.TempDir()); err == nil {
		t.Fatal("expected error for floater file missing on disk, got nil")
	}
}

func TestValidateAsset_FloaterInvalidPosition(t *testing.T) {
	dir := t.TempDir()
	a := validFloaterManifest()
	a.Floater.Position = "center"
	os.WriteFile(filepath.Join(dir, a.Floater.File), []byte("fake png"), 0644)

	if err := ValidateAsset(&a, dir); err == nil {
		t.Fatal("expected error for invalid floater position, got nil")
	}
}

func TestValidateAsset_FloaterValidAllPositions(t *testing.T) {
	for _, pos := range ValidFloaterPositions {
		dir := t.TempDir()
		a := validFloaterManifest()
		a.Floater.Position = pos
		os.WriteFile(filepath.Join(dir, a.Floater.File), []byte("fake png"), 0644)

		if err := ValidateAsset(&a, dir); err != nil {
			t.Errorf("expected position %q to be valid, got error: %v", pos, err)
		}
	}
}

func TestValidateAsset_FloaterAnimatedZeroInterval(t *testing.T) {
	dir := t.TempDir()
	a := validFloaterManifest()
	a.Floater.AnimateFrames = []string{"f1.png", "f2.png"}
	a.Floater.IntervalMs = 0
	os.WriteFile(filepath.Join(dir, a.Floater.File), []byte("fake png"), 0644)
	for _, f := range a.Floater.AnimateFrames {
		os.WriteFile(filepath.Join(dir, f), []byte("fake png"), 0644)
	}

	if err := ValidateAsset(&a, dir); err == nil {
		t.Fatal("expected error for animated floater with zero interval_ms, got nil")
	}
}

func TestValidateAsset_FloaterAnimatedMissingFrameFile(t *testing.T) {
	dir := t.TempDir()
	a := validFloaterManifest()
	a.Floater.AnimateFrames = []string{"f1.png", "f2.png"}
	a.Floater.IntervalMs = 100
	os.WriteFile(filepath.Join(dir, a.Floater.File), []byte("fake png"), 0644)
	// only write f1.png, not f2.png
	os.WriteFile(filepath.Join(dir, "f1.png"), []byte("fake png"), 0644)

	if err := ValidateAsset(&a, dir); err == nil {
		t.Fatal("expected error for missing animation frame file, got nil")
	}
}

func TestValidateAsset_FloaterAnimatedValid(t *testing.T) {
	dir := t.TempDir()
	a := validFloaterManifest()
	a.Floater.AnimateFrames = []string{"f1.png", "f2.png"}
	a.Floater.IntervalMs = 100
	os.WriteFile(filepath.Join(dir, a.Floater.File), []byte("fake png"), 0644)
	for _, f := range a.Floater.AnimateFrames {
		os.WriteFile(filepath.Join(dir, f), []byte("fake png"), 0644)
	}

	if err := ValidateAsset(&a, dir); err != nil {
		t.Fatalf("expected no error for valid animated floater, got: %v", err)
	}
}

func TestValidateAsset_FloaterStaticValid(t *testing.T) {
	dir := t.TempDir()
	a := validFloaterManifest()
	os.WriteFile(filepath.Join(dir, a.Floater.File), []byte("fake png"), 0644)

	if err := ValidateAsset(&a, dir); err != nil {
		t.Fatalf("expected no error for valid static floater, got: %v", err)
	}
}
