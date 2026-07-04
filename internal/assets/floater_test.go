package assets

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// minimalPNG is a valid 1x1 red-pixel PNG (69 bytes). Used in tests that
// only need a file to exist on disk (ValidateAsset uses os.Stat, not chafa).
var minimalPNG = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
	0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
	0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
	0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
	0x00, 0x00, 0x02, 0x00, 0x01, 0xE2, 0x21, 0xBC,
	0x33, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
	0x44, 0xAE, 0x42, 0x60, 0x82,
}

// findRealPNG locates a real PNG from the repo's assets directory so
// chafa-dependent tests have a valid image to render. Returns "" if none found.
func findRealPNG(t *testing.T) string {
	t.Helper()
	// walk up from the test file to find the repo root assets dir
	candidates := []string{
		"../../assets/spinners/pulse/frame1.png",
		"../../assets/dividers/neon-divider/divider.png",
		"../../assets/spinners/pulse/frame4.png",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	return ""
}

func validFloaterAsset() Asset {
	return Asset{
		Name:    "corner-cat",
		Type:    AssetTypeFloater,
		Version: "1.0.0",
		Render: RenderConfig{
			Mode:      RenderModeBraille,
			Width:     8,
			Height:    8,
			ColorMode: ColorModeTrueColor,
		},
		Floater: &FloaterConfig{
			File:      "cat.png",
			Position:  FloaterTopRight,
			MaxWidth:  10,
			MaxHeight: 6,
			MarginX:   1,
			MarginY:   1,
		},
	}
}

func setupFloaterAsset(t *testing.T, assetsDir string, a Asset) string {
	t.Helper()
	assetDir := filepath.Join(assetsDir, "floaters", a.Name)
	if err := os.MkdirAll(assetDir, 0755); err != nil {
		t.Fatalf("could not create asset dir: %v", err)
	}
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("could not marshal asset: %v", err)
	}
	if err := os.WriteFile(filepath.Join(assetDir, "asset.json"), data, 0644); err != nil {
		t.Fatalf("could not write asset.json: %v", err)
	}
	if a.Floater != nil {
		os.WriteFile(filepath.Join(assetDir, a.Floater.File), minimalPNG, 0644)
		for _, frame := range a.Floater.AnimateFrames {
			os.WriteFile(filepath.Join(assetDir, frame), minimalPNG, 0644)
		}
	}
	return assetDir
}

func TestManager_Get_FindsFloaterAsset(t *testing.T) {
	assetsDir := t.TempDir()
	setupFloaterAsset(t, assetsDir, validFloaterAsset())

	m, err := NewManager(assetsDir)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, dir, err := m.Get("corner-cat", AssetTypeFloater)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if got.Name != "corner-cat" {
		t.Errorf("expected name 'corner-cat', got %q", got.Name)
	}
	if dir == "" {
		t.Error("expected non-empty asset directory")
	}
}

func TestManager_PreviewFloater_NoFloaterConfig(t *testing.T) {
	assetsDir := t.TempDir()
	a := validFloaterAsset()
	a.Floater = nil
	setupFloaterAsset(t, assetsDir, a)

	m, _ := NewManager(assetsDir)
	if err := m.PreviewFloater("corner-cat"); err == nil {
		t.Fatal("expected error for floater asset with nil Floater config, got nil")
	}
}

func TestManager_PreviewFloater_AssetNotFound(t *testing.T) {
	m, _ := NewManager(t.TempDir())
	if err := m.PreviewFloater("nonexistent-floater"); err == nil {
		t.Fatal("expected error for missing floater asset, got nil")
	}
}

func TestManager_PreviewFloater_RendersWithChafa(t *testing.T) {
	if !ChafaAvailable() {
		t.Skip("chafa not installed")
	}
	realPNG := findRealPNG(t)
	if realPNG == "" {
		t.Skip("no real PNG found in repo assets")
	}
	pngData, err := os.ReadFile(realPNG)
	if err != nil {
		t.Fatalf("could not read real PNG: %v", err)
	}

	assetsDir := t.TempDir()
	a := validFloaterAsset()
	assetDir := filepath.Join(assetsDir, "floaters", a.Name)
	os.MkdirAll(assetDir, 0755)
	data, _ := json.Marshal(a)
	os.WriteFile(filepath.Join(assetDir, "asset.json"), data, 0644)
	os.WriteFile(filepath.Join(assetDir, a.Floater.File), pngData, 0644)

	m, _ := NewManager(assetsDir)
	if err := m.PreviewFloater("corner-cat"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestManager_RenderFloaterFrames_NoAnimateFrames(t *testing.T) {
	assetsDir := t.TempDir()
	setupFloaterAsset(t, assetsDir, validFloaterAsset())

	m, _ := NewManager(assetsDir)
	_, _, _, err := m.RenderFloaterFrames("corner-cat")
	if err == nil {
		t.Fatal("expected error for floater with no animate_frames, got nil")
	}
}

func TestManager_RenderFloaterFrames_ReturnsPositionAndInterval(t *testing.T) {
	if !ChafaAvailable() {
		t.Skip("chafa not installed")
	}
	realPNG := findRealPNG(t)
	if realPNG == "" {
		t.Skip("no real PNG found in repo assets")
	}
	pngData, err := os.ReadFile(realPNG)
	if err != nil {
		t.Fatalf("could not read real PNG: %v", err)
	}

	assetsDir := t.TempDir()
	a := validFloaterAsset()
	a.Floater.AnimateFrames = []string{"frame1.png", "frame2.png"}
	a.Floater.IntervalMs = 150

	assetDir := filepath.Join(assetsDir, "floaters", a.Name)
	os.MkdirAll(assetDir, 0755)
	data, _ := json.Marshal(a)
	os.WriteFile(filepath.Join(assetDir, "asset.json"), data, 0644)
	os.WriteFile(filepath.Join(assetDir, a.Floater.File), pngData, 0644)
	os.WriteFile(filepath.Join(assetDir, "frame1.png"), pngData, 0644)
	os.WriteFile(filepath.Join(assetDir, "frame2.png"), pngData, 0644)

	m, _ := NewManager(assetsDir)
	frames, interval, pos, err := m.RenderFloaterFrames("corner-cat")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(frames) != 2 {
		t.Errorf("expected 2 frames, got %d", len(frames))
	}
	if interval != 150 {
		t.Errorf("expected interval 150, got %d", interval)
	}
	if pos != FloaterTopRight {
		t.Errorf("expected position %q, got %q", FloaterTopRight, pos)
	}
}

func TestIsValidFloaterPosition_AllFourCorners(t *testing.T) {
	for _, pos := range ValidFloaterPositions {
		if !IsValidFloaterPosition(pos) {
			t.Errorf("expected %q to be a valid floater position", pos)
		}
	}
}

func TestIsValidFloaterPosition_RejectsInvalid(t *testing.T) {
	cases := []FloaterPosition{"center", "middle", "", "top-center"}
	for _, pos := range cases {
		if IsValidFloaterPosition(pos) {
			t.Errorf("expected %q to be an invalid floater position", pos)
		}
	}
}
