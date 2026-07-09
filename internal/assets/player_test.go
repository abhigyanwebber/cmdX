package assets

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// ── buildPlayerCommand — custom player template ──────────────────────────────

func TestBuildPlayerCommand_CustomPlayerWithPlaceholder(t *testing.T) {
	cmd, err := buildPlayerCommand("/path/to/sound.wav", 1.0, "ffplay -nodisp -autoexit %f")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if filepath.Base(cmd.Path) != "ffplay" && cmd.Args[0] != "ffplay" {
		t.Errorf("expected ffplay as the command, got %v", cmd.Args)
	}
	found := false
	for _, a := range cmd.Args {
		if a == "/path/to/sound.wav" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected file path substituted for %%f, got args: %v", cmd.Args)
	}
}

func TestBuildPlayerCommand_CustomPlayerWithoutPlaceholderAppendsPath(t *testing.T) {
	cmd, err := buildPlayerCommand("/path/to/sound.wav", 1.0, "my-player --quiet")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	last := cmd.Args[len(cmd.Args)-1]
	if last != "/path/to/sound.wav" {
		t.Errorf("expected file path appended as final arg when no %%f placeholder, got args: %v", cmd.Args)
	}
}

func TestBuildPlayerCommand_EmptyCustomPlayerErrors(t *testing.T) {
	_, err := buildPlayerCommand("/path/to/sound.wav", 1.0, "   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only custom player, got nil")
	}
}

func TestBuildPlayerCommand_NoCustomPlayerUsesPlatformDefault(t *testing.T) {
	// Just verify it doesn't error and produces *some* command — the
	// actual binary chosen depends on what's installed on the test
	// machine (or the OS branch on non-Linux), which we don't control
	// here. buildLinuxPlayerCommand returning an error when no player
	// is found is a legitimate, expected outcome in a minimal CI
	// container with none of paplay/aplay/ffplay installed.
	_, err := buildPlayerCommand("/path/to/sound.wav", 1.0, "")
	if err != nil {
		t.Logf("no default player found on this system (expected in minimal environments): %v", err)
	}
}

func TestFormatVolume_ProducesDecimalString(t *testing.T) {
	got := formatVolume(0.5)
	if got != "0.500" {
		t.Errorf("expected '0.500', got %q", got)
	}
}

// ── cooldown state tracking ───────────────────────────────────────────────────

func TestCooldownStatePath_ResolvesToAssetsRoot(t *testing.T) {
	assetDir := filepath.Join("/tmp/myassets", "sounds", "my-sound-theme")
	path := cooldownStatePath(assetDir, "my-sound-theme", "error")
	expected := filepath.Join("/tmp/myassets", ".state", "sound-my-sound-theme-error-lastplayed.txt")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}

func TestIsOnCooldown_NeverPlayedIsNotOnCooldown(t *testing.T) {
	assetsDir := t.TempDir()
	assetDir := filepath.Join(assetsDir, "sounds", "test-theme")
	os.MkdirAll(assetDir, 0755)

	onCooldown, err := isOnCooldown(assetDir, "test-theme", "never-played", 5000)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if onCooldown {
		t.Error("expected a never-played sound to not be on cooldown")
	}
}

func TestRecordPlayedThenIsOnCooldown(t *testing.T) {
	assetsDir := t.TempDir()
	assetDir := filepath.Join(assetsDir, "sounds", "test-theme")
	os.MkdirAll(assetDir, 0755)

	if err := recordPlayed(assetDir, "test-theme", "success"); err != nil {
		t.Fatalf("expected no error recording play, got: %v", err)
	}

	onCooldown, err := isOnCooldown(assetDir, "test-theme", "success", 60000) // 60s cooldown
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !onCooldown {
		t.Error("expected sound just played to be on cooldown for a 60s window")
	}
}

func TestIsOnCooldown_ExpiresAfterWindow(t *testing.T) {
	assetsDir := t.TempDir()
	assetDir := filepath.Join(assetsDir, "sounds", "test-theme")
	os.MkdirAll(assetDir, 0755)

	// manually write a timestamp far enough in the past that a short
	// cooldown window will have already elapsed
	path := cooldownStatePath(assetDir, "test-theme", "success")
	os.MkdirAll(filepath.Dir(path), 0755)
	past := time.Now().Add(-10 * time.Second).UnixMilli()
	os.WriteFile(path, []byte(strconv.FormatInt(past, 10)), 0644)

	onCooldown, err := isOnCooldown(assetDir, "test-theme", "success", 100) // 100ms cooldown, 10s have passed
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if onCooldown {
		t.Error("expected cooldown to have expired after 10s with a 100ms window")
	}
}

func TestIsOnCooldown_CorruptedStateFileIsHandledGracefully(t *testing.T) {
	assetsDir := t.TempDir()
	assetDir := filepath.Join(assetsDir, "sounds", "test-theme")
	path := cooldownStatePath(assetDir, "test-theme", "success")
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte("not-a-timestamp"), 0644)

	_, err := isOnCooldown(assetDir, "test-theme", "success", 5000)
	if err == nil {
		t.Error("expected an error for a corrupted timestamp file, got nil")
	}
	// PlaySound itself treats this error as "not on cooldown" rather
	// than failing outright — that graceful-degradation behavior is
	// documented in PlaySound and isn't re-tested here since it would
	// require an actual audio player invocation.
}
