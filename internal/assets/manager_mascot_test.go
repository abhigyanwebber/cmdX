package assets

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupMascotManager writes a real mascot asset (manifest + PNG frames)
// under a temp assets directory and returns a Manager rooted there,
// plus the asset name. Two competing states are configured so tests can
// tell forced-state override apart from natural trigger resolution:
// "idle" (always, low priority) and "error" (exit_code 1-127, high
// priority).
func setupMascotManager(t *testing.T) (*Manager, string) {
	t.Helper()
	assetsDir := t.TempDir()
	mascotDir := filepath.Join(assetsDir, "mascots", "regtest")
	if err := os.MkdirAll(mascotDir, 0755); err != nil {
		t.Fatalf("could not create mascot dir: %v", err)
	}

	a := &Asset{
		Name:        "regtest",
		Type:        AssetTypeMascot,
		Version:     "1.0.0",
		Author:      "test",
		Description: "regression test mascot",
		Render: RenderConfig{
			Mode:      "sixel",
			Width:     8,
			Height:    4,
			ColorMode: "truecolor",
			SymbolSet: "block",
		},
		Mascot: &MascotConfig{
			DefaultState: MascotStateIdle,
			Position:     FloaterBottomRight,
			MaxWidth:     8,
			MaxHeight:    4,
			States: map[MascotState]MascotStateConfig{
				MascotStateIdle: {
					Frames:   []string{"idle.png"},
					Triggers: []MascotTrigger{{Type: TriggerAlways, Priority: 0}},
				},
				MascotStateError: {
					Frames:   []string{"error.png"},
					Triggers: []MascotTrigger{{Type: TriggerExitCode, Value: "1-127", Priority: 10}},
				},
			},
		},
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("could not marshal test manifest: %v", err)
	}
	if err := os.WriteFile(filepath.Join(mascotDir, "asset.json"), data, 0644); err != nil {
		t.Fatalf("could not write asset.json: %v", err)
	}
	// minimalPNG (a synthetic 1x1 stub used elsewhere in this package
	// for manifest/path validation tests, which only os.Stat the file)
	// is rejected by chafa when actually asked to render it. This test
	// exercises the real render pipeline end-to-end, so it needs a
	// genuine decodable PNG — reuse one already bundled with the repo's
	// sample assets rather than embedding a second fixture.
	realPNG := findRealFramePNG(t)
	writeFileOrFatal(t, filepath.Join(mascotDir, "idle.png"), realPNG)
	writeFileOrFatal(t, filepath.Join(mascotDir, "error.png"), realPNG)

	m, err := NewManager(assetsDir)
	if err != nil {
		t.Fatalf("could not create manager: %v", err)
	}
	return m, "regtest"
}

// findRealFramePNG locates a genuine, chafa-decodable PNG bundled with
// the repo's sample assets (walking up from the package directory to
// the repo root's assets/ folder) so render-pipeline tests don't depend
// on a synthetic fixture that chafa itself can't decode.
func findRealFramePNG(t *testing.T) []byte {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}
	dir := wd
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(dir, "assets", "spinners", "pulse", "frame1.png")
		if data, err := os.ReadFile(candidate); err == nil {
			return data
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Skip("could not locate a bundled sample PNG (assets/spinners/pulse/frame1.png) to render against")
	return nil
}

func writeFileOrFatal(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("could not write %s: %v", path, err)
	}
}

// captureStdout redirects os.Stdout for the duration of fn and returns
// everything written to it. The reader runs concurrently in a goroutine
// so a write larger than the OS pipe buffer (real sixel/chafa output
// easily exceeds the typical 64KB pipe capacity) cannot deadlock fn().
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("could not create pipe: %v", err)
	}
	os.Stdout = w

	outCh := make(chan string, 1)
	go func() {
		data, _ := io.ReadAll(r)
		outCh <- string(data)
	}()

	fn()

	w.Close()
	os.Stdout = old
	return <-outCh
}

// TestPreviewMascotState_ForcedStateOverridesResolution is the regression
// test for BUG-002: a forced state must be displayed as-is, even when
// natural ctx-based trigger resolution would select a different state.
// Before the fix, the forced state was silently discarded and the
// independently-resolved ctx-based state was shown instead.
func TestPreviewMascotState_ForcedStateOverridesResolution(t *testing.T) {
	if !ChafaAvailable() {
		t.Skip("chafa not installed in this environment")
	}
	m, name := setupMascotManager(t)

	// ctx alone (exit code 0) would naturally resolve to "idle", not
	// "error" — forcing "error" must still display "error".
	ctx := MascotContext{LastExitCode: 0}
	out := captureStdout(t, func() {
		if err := m.PreviewMascotState(name, ctx, RenderOverrides{}, MascotStateError); err != nil {
			t.Fatalf("PreviewMascotState returned error: %v", err)
		}
	})

	if !strings.Contains(out, "state: error") {
		t.Errorf("expected forced state 'error' to be displayed, got output: %q", out)
	}
}

// TestPreviewMascotState_EmptyForcedStateFallsBackToResolution ensures
// the normal (non-override) path still resolves from ctx exactly as
// before — an empty forcedState must not force anything.
func TestPreviewMascotState_EmptyForcedStateFallsBackToResolution(t *testing.T) {
	if !ChafaAvailable() {
		t.Skip("chafa not installed in this environment")
	}
	m, name := setupMascotManager(t)

	ctx := MascotContext{LastExitCode: 1} // should naturally resolve to "error"
	out := captureStdout(t, func() {
		if err := m.PreviewMascotState(name, ctx, RenderOverrides{}, ""); err != nil {
			t.Fatalf("PreviewMascotState returned error: %v", err)
		}
	})

	if !strings.Contains(out, "state: error") {
		t.Errorf("expected natural resolution to select 'error' for exit code 1, got output: %q", out)
	}
}

// TestPreviewMascot_CompatibilityWrapperStillResolvesNormally checks that
// the pre-existing PreviewMascot function (kept as a thin wrapper around
// PreviewMascotState with no forced state) still behaves exactly as
// before the fix — this guards against the refactor accidentally
// changing its public contract.
func TestPreviewMascot_CompatibilityWrapperStillResolvesNormally(t *testing.T) {
	if !ChafaAvailable() {
		t.Skip("chafa not installed in this environment")
	}
	m, name := setupMascotManager(t)

	ctx := MascotContext{LastExitCode: 0} // should naturally resolve to "idle"
	out := captureStdout(t, func() {
		if err := m.PreviewMascot(name, ctx, RenderOverrides{}); err != nil {
			t.Fatalf("PreviewMascot returned error: %v", err)
		}
	})

	if !strings.Contains(out, "state: idle") {
		t.Errorf("expected PreviewMascot to still resolve naturally to 'idle', got output: %q", out)
	}
}
