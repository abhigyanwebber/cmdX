package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func buildValidSoundTheme() *SoundThemeConfig {
	return &SoundThemeConfig{
		Enabled:      true,
		GlobalVolume: 0.8,
		Sounds: map[string]SoundEffect{
			"success": {
				File:   "success.wav",
				Volume: 1.0,
				Async:  true,
				Triggers: []SoundTrigger{
					{Type: TriggerExitCode, Value: "0", Priority: 20},
				},
			},
			"error": {
				File:       "error.wav",
				Volume:     1.0,
				Async:      true,
				CooldownMs: 2000,
				Triggers: []SoundTrigger{
					{Type: TriggerExitCode, Value: "1-127", Priority: 20},
					{Type: TriggerOutputRegex, Value: "error|fatal", Priority: 15},
				},
			},
		},
	}
}

func writeSoundWavs(t *testing.T, dir string, s *SoundThemeConfig) {
	t.Helper()
	for _, effect := range s.Sounds {
		os.WriteFile(filepath.Join(dir, effect.File), []byte("fake wav data"), 0644)
	}
}

// ── ValidateSoundTheme ────────────────────────────────────────────────────────

func TestValidateSoundTheme_Valid(t *testing.T) {
	dir := t.TempDir()
	s := buildValidSoundTheme()
	writeSoundWavs(t, dir, s)

	if err := ValidateSoundTheme(s, dir); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateSoundTheme_Nil(t *testing.T) {
	if err := ValidateSoundTheme(nil, t.TempDir()); err == nil {
		t.Fatal("expected error for nil config, got nil")
	}
}

func TestValidateSoundTheme_NoSounds(t *testing.T) {
	s := buildValidSoundTheme()
	s.Sounds = map[string]SoundEffect{}
	if err := ValidateSoundTheme(s, t.TempDir()); err == nil {
		t.Fatal("expected error for zero sounds, got nil")
	}
}

func TestValidateSoundTheme_InvalidGlobalVolume(t *testing.T) {
	dir := t.TempDir()
	s := buildValidSoundTheme()
	writeSoundWavs(t, dir, s)
	s.GlobalVolume = 1.5
	if err := ValidateSoundTheme(s, dir); err == nil {
		t.Fatal("expected error for global_volume > 1.0, got nil")
	}
}

func TestValidateSoundTheme_MissingFile(t *testing.T) {
	dir := t.TempDir()
	s := buildValidSoundTheme()
	// don't write the wav files
	if err := ValidateSoundTheme(s, dir); err == nil {
		t.Fatal("expected error for missing sound files, got nil")
	}
}

func TestValidateSoundTheme_EmptyFileField(t *testing.T) {
	dir := t.TempDir()
	s := buildValidSoundTheme()
	effect := s.Sounds["success"]
	effect.File = ""
	s.Sounds["success"] = effect
	if err := ValidateSoundTheme(s, dir); err == nil {
		t.Fatal("expected error for empty file field, got nil")
	}
}

func TestValidateSoundTheme_InvalidEffectVolume(t *testing.T) {
	dir := t.TempDir()
	s := buildValidSoundTheme()
	writeSoundWavs(t, dir, s)
	effect := s.Sounds["success"]
	effect.Volume = 2.0
	s.Sounds["success"] = effect
	if err := ValidateSoundTheme(s, dir); err == nil {
		t.Fatal("expected error for volume > 1.0, got nil")
	}
}

func TestValidateSoundTheme_NegativeCooldown(t *testing.T) {
	dir := t.TempDir()
	s := buildValidSoundTheme()
	writeSoundWavs(t, dir, s)
	effect := s.Sounds["error"]
	effect.CooldownMs = -100
	s.Sounds["error"] = effect
	if err := ValidateSoundTheme(s, dir); err == nil {
		t.Fatal("expected error for negative cooldown_ms, got nil")
	}
}

// ── validateSoundTrigger ──────────────────────────────────────────────────────

func TestValidateSoundTrigger_ExitCodeReusesSharedValidator(t *testing.T) {
	// exercises the bridge to validateExitCodeValue (mascot.go)
	valid := SoundTrigger{Type: TriggerExitCode, Value: "1-127"}
	if err := validateSoundTrigger(valid, "test"); err != nil {
		t.Errorf("expected valid exit_code range to pass, got: %v", err)
	}

	invalid := SoundTrigger{Type: TriggerExitCode, Value: "not-a-number"}
	if err := validateSoundTrigger(invalid, "test"); err == nil {
		t.Error("expected invalid exit_code value to fail, got nil")
	}
}

func TestValidateSoundTrigger_IdleTimeStrictParsing(t *testing.T) {
	if err := validateSoundTrigger(SoundTrigger{Type: TriggerIdleTime, Value: "30"}, "test"); err != nil {
		t.Errorf("expected '30' to be valid idle_time, got: %v", err)
	}
	// trailing garbage must be rejected — strconv.Atoi is strict,
	// unlike fmt.Sscanf which would silently accept "30abc"
	if err := validateSoundTrigger(SoundTrigger{Type: TriggerIdleTime, Value: "30abc"}, "test"); err == nil {
		t.Error("expected '30abc' to be rejected as invalid idle_time, got nil")
	}
}

func TestValidateSoundTrigger_InvalidRegex(t *testing.T) {
	if err := validateSoundTrigger(SoundTrigger{Type: TriggerOutputRegex, Value: "["}, "test"); err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestValidateSoundTrigger_UnknownType(t *testing.T) {
	if err := validateSoundTrigger(SoundTrigger{Type: "on-hover"}, "test"); err == nil {
		t.Fatal("expected error for unknown trigger type, got nil")
	}
}

func TestValidateSoundTrigger_ValidGitStatusValues(t *testing.T) {
	for _, v := range []string{"dirty", "clean", "untracked", "ahead", "behind"} {
		if err := validateSoundTrigger(SoundTrigger{Type: TriggerGitStatus, Value: v}, "test"); err != nil {
			t.Errorf("expected git_status %q to be valid, got: %v", v, err)
		}
	}
}

// ── ResolveSound ──────────────────────────────────────────────────────────────

func TestResolveSound_ExitCodeZeroResolvesToSuccess(t *testing.T) {
	s := buildValidSoundTheme()
	name, _, ok := ResolveSound(s, MascotContext{LastExitCode: 0})
	if !ok || name != "success" {
		t.Errorf("expected 'success', got name=%q ok=%v", name, ok)
	}
}

func TestResolveSound_NonZeroExitResolvesToError(t *testing.T) {
	s := buildValidSoundTheme()
	name, _, ok := ResolveSound(s, MascotContext{LastExitCode: 1})
	if !ok || name != "error" {
		t.Errorf("expected 'error', got name=%q ok=%v", name, ok)
	}
}

func TestResolveSound_NoMatchReturnsFalse(t *testing.T) {
	s := &SoundThemeConfig{
		Sounds: map[string]SoundEffect{
			"only-on-command": {
				File:     "x.wav",
				Triggers: []SoundTrigger{{Type: TriggerCommand, Value: "docker*", Priority: 5}},
			},
		},
	}
	_, _, ok := ResolveSound(s, MascotContext{LastExitCode: 0, LastCommand: "ls -la"})
	if ok {
		t.Error("expected no match for unrelated command, got a match")
	}
}

func TestResolveSound_PriorityDeterminesWinner(t *testing.T) {
	s := &SoundThemeConfig{
		Sounds: map[string]SoundEffect{
			"low":  {File: "low.wav", Triggers: []SoundTrigger{{Type: TriggerAlways, Priority: 1}}},
			"high": {File: "high.wav", Triggers: []SoundTrigger{{Type: TriggerAlways, Priority: 100}}},
		},
	}
	name, _, ok := ResolveSound(s, MascotContext{})
	if !ok || name != "high" {
		t.Errorf("expected higher-priority 'high' to win, got name=%q ok=%v", name, ok)
	}
}

// ── SoundShellHooks ───────────────────────────────────────────────────────────

func TestSoundShellHooks_Bash(t *testing.T) {
	hooks := SoundShellHooks("my-sounds", "bash")
	if hooks == "" {
		t.Fatal("expected non-empty bash hooks")
	}
	if !containsStr(hooks, "PROMPT_COMMAND") {
		t.Error("expected bash hooks to set PROMPT_COMMAND")
	}
	if !containsStr(hooks, "sound-play") {
		t.Error("expected bash hooks to call 'asset sound-play'")
	}
}

func TestSoundShellHooks_Zsh(t *testing.T) {
	hooks := SoundShellHooks("my-sounds", "zsh")
	if !containsStr(hooks, "add-zsh-hook") {
		t.Error("expected zsh hooks to use add-zsh-hook")
	}
}

func TestSoundShellHooks_PowerShell(t *testing.T) {
	hooks := SoundShellHooks("my-sounds", "powershell")
	if !containsStr(hooks, "LASTEXITCODE") {
		t.Error("expected powershell hooks to check $LASTEXITCODE")
	}
}

func TestSoundShellHooks_UnknownShell(t *testing.T) {
	if hooks := SoundShellHooks("x", "fish"); hooks != "" {
		t.Error("expected empty string for unsupported shell")
	}
}
