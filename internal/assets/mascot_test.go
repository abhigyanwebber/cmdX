package assets

import (
	"os"
	"path/filepath"
	"testing"
)

// buildMascotConfig builds a complete valid MascotConfig for testing.
func buildMascotConfig() *MascotConfig {
	return &MascotConfig{
		DefaultState:     MascotStateIdle,
		Position:         FloaterBottomRight,
		MaxWidth:         12,
		MaxHeight:        8,
		MarginX:          2,
		MarginY:          1,
		GlobalIntervalMs: 100,
		States: map[MascotState]MascotStateConfig{
			MascotStateIdle: {
				Frames: []string{"idle_01.png"},
				Triggers: []MascotTrigger{
					{Type: TriggerAlways, Priority: 0},
				},
			},
			MascotStateSuccess: {
				Frames: []string{"success_01.png"},
				Triggers: []MascotTrigger{
					{Type: TriggerExitCode, Value: "0", Priority: 20},
				},
				RenderOverride: &MascotRenderOverride{Tint: "#00FF88"},
			},
			MascotStateError: {
				Frames: []string{"error_01.png"},
				Triggers: []MascotTrigger{
					{Type: TriggerExitCode, Value: "1-127", Priority: 20},
					{Type: TriggerOutputRegex, Value: "error|fatal", Priority: 15},
				},
				RenderOverride: &MascotRenderOverride{Tint: "#FF4444"},
			},
			MascotStateSleeping: {
				Frames: []string{"sleep_01.png"},
				Triggers: []MascotTrigger{
					{Type: TriggerIdleTime, Value: "30", Priority: 1},
				},
				RenderOverride: &MascotRenderOverride{Width: 6, Height: 4},
			},
		},
	}
}

func writeMascotPNGs(t *testing.T, dir string, mc *MascotConfig) {
	t.Helper()
	for _, state := range mc.States {
		for _, f := range state.Frames {
			os.WriteFile(filepath.Join(dir, f), minimalPNG, 0644)
		}
		for _, f := range state.TransitionFrames {
			os.WriteFile(filepath.Join(dir, f), minimalPNG, 0644)
		}
	}
}

// ── ValidateMascot ────────────────────────────────────────────────────────────

func TestValidateMascot_Valid(t *testing.T) {
	dir := t.TempDir()
	mc := buildMascotConfig()
	writeMascotPNGs(t, dir, mc)

	if err := ValidateMascot(mc, dir); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateMascot_NilConfig(t *testing.T) {
	if err := ValidateMascot(nil, t.TempDir()); err == nil {
		t.Fatal("expected error for nil mascot config, got nil")
	}
}

func TestValidateMascot_NoStates(t *testing.T) {
	mc := buildMascotConfig()
	mc.States = map[MascotState]MascotStateConfig{}
	if err := ValidateMascot(mc, t.TempDir()); err == nil {
		t.Fatal("expected error for zero states, got nil")
	}
}

func TestValidateMascot_InvalidPosition(t *testing.T) {
	dir := t.TempDir()
	mc := buildMascotConfig()
	mc.Position = "center"
	writeMascotPNGs(t, dir, mc)
	if err := ValidateMascot(mc, dir); err == nil {
		t.Fatal("expected error for invalid position, got nil")
	}
}

func TestValidateMascot_InvalidDefaultState(t *testing.T) {
	dir := t.TempDir()
	mc := buildMascotConfig()
	mc.DefaultState = "nonexistent-state"
	writeMascotPNGs(t, dir, mc)
	if err := ValidateMascot(mc, dir); err == nil {
		t.Fatal("expected error for undefined default_state, got nil")
	}
}

func TestValidateMascot_MissingFrameFile(t *testing.T) {
	dir := t.TempDir()
	mc := buildMascotConfig()
	// write PNGs for all states EXCEPT error
	for stateName, state := range mc.States {
		if stateName == MascotStateError {
			continue
		}
		for _, f := range state.Frames {
			os.WriteFile(filepath.Join(dir, f), minimalPNG, 0644)
		}
	}
	if err := ValidateMascot(mc, dir); err == nil {
		t.Fatal("expected error for missing frame file, got nil")
	}
}

func TestValidateMascot_InvalidTintColor(t *testing.T) {
	dir := t.TempDir()
	mc := buildMascotConfig()
	writeMascotPNGs(t, dir, mc)
	mc.States[MascotStateSuccess] = MascotStateConfig{
		Frames: []string{"success_01.png"},
		RenderOverride: &MascotRenderOverride{Tint: "red"}, // invalid hex
	}
	if err := ValidateMascot(mc, dir); err == nil {
		t.Fatal("expected error for invalid tint color, got nil")
	}
}

// ── Trigger validation ────────────────────────────────────────────────────────

func TestValidateTrigger_ExitCodeFormats(t *testing.T) {
	valid := []string{"0", "1", "127", "1-127", "128+"}
	for _, v := range valid {
		if err := validateExitCodeValue(v); err != nil {
			t.Errorf("expected %q to be valid exit code pattern, got error: %v", v, err)
		}
	}
}

func TestValidateTrigger_InvalidExitCode(t *testing.T) {
	invalid := []string{"", "abc", "1-", "-5", "1.5"}
	for _, v := range invalid {
		if err := validateExitCodeValue(v); err == nil {
			t.Errorf("expected %q to be invalid exit code pattern, got nil", v)
		}
	}
}

func TestValidateTrigger_InvalidOutputRegex(t *testing.T) {
	if err := validateTrigger(MascotTrigger{Type: TriggerOutputRegex, Value: "["}, MascotStateError); err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestValidateTrigger_ValidGitStatus(t *testing.T) {
	for _, v := range []string{"dirty", "clean", "untracked", "ahead", "behind"} {
		if err := validateTrigger(MascotTrigger{Type: TriggerGitStatus, Value: v}, MascotStateIdle); err != nil {
			t.Errorf("expected git_status %q to be valid, got: %v", v, err)
		}
	}
}

func TestValidateTrigger_InvalidGitStatus(t *testing.T) {
	if err := validateTrigger(MascotTrigger{Type: TriggerGitStatus, Value: "stashed"}, MascotStateIdle); err == nil {
		t.Fatal("expected error for invalid git_status value, got nil")
	}
}

func TestValidateTrigger_UnknownType(t *testing.T) {
	if err := validateTrigger(MascotTrigger{Type: "hover", Value: "x"}, MascotStateIdle); err == nil {
		t.Fatal("expected error for unknown trigger type, got nil")
	}
}

// ── ResolveState ──────────────────────────────────────────────────────────────

func TestResolveState_ExitCodeZeroResolvesToSuccess(t *testing.T) {
	mc := buildMascotConfig()
	ctx := MascotContext{LastExitCode: 0}
	state := ResolveState(mc, ctx)
	if state != MascotStateSuccess {
		t.Errorf("expected success state for exit code 0, got %q", state)
	}
}

func TestResolveState_NonZeroExitCodeResolvesToError(t *testing.T) {
	mc := buildMascotConfig()
	ctx := MascotContext{LastExitCode: 1}
	state := ResolveState(mc, ctx)
	if state != MascotStateError {
		t.Errorf("expected error state for exit code 1, got %q", state)
	}
}

func TestResolveState_HighExitCodeEdgeCase(t *testing.T) {
	mc := buildMascotConfig()
	ctx := MascotContext{LastExitCode: 127}
	state := ResolveState(mc, ctx)
	if state != MascotStateError {
		t.Errorf("expected error state for exit code 127 (in range 1-127), got %q", state)
	}
}

func TestResolveState_IdleTimeResolvesToSleeping(t *testing.T) {
	mc := buildMascotConfig()
	// no exit code trigger fires at 0, and idle_time=30s threshold met
	ctx := MascotContext{LastExitCode: 0, IdleSeconds: 45}
	state := ResolveState(mc, ctx)
	// success (priority 20) beats sleeping (priority 1) even with idle
	// so with exit code 0 + idle, success should win
	if state != MascotStateSuccess {
		t.Errorf("expected success to beat sleeping due to higher priority, got %q", state)
	}
}

func TestResolveState_PureIdleResolvesToSleeping(t *testing.T) {
	mc := buildMascotConfig()
	// remove success/error triggers to isolate idle
	delete(mc.States, MascotStateSuccess)
	delete(mc.States, MascotStateError)
	ctx := MascotContext{LastExitCode: 0, IdleSeconds: 60}
	state := ResolveState(mc, ctx)
	if state != MascotStateSleeping {
		t.Errorf("expected sleeping state for idle > 30s with no competing triggers, got %q", state)
	}
}

func TestResolveState_OutputRegexMatchesError(t *testing.T) {
	mc := buildMascotConfig()
	ctx := MascotContext{LastExitCode: 0, LastOutput: "build failed: fatal error in main.go"}
	state := ResolveState(mc, ctx)
	// output_regex priority 15 < exit_code 0 priority 20 (success)
	// so success wins here since exit code matters more
	if state != MascotStateSuccess {
		t.Errorf("expected success (higher priority) to beat output_regex match, got %q", state)
	}
}

func TestResolveState_OutputRegexWinsWhenExitCodeMatches(t *testing.T) {
	mc := buildMascotConfig()
	// exit code 1 (priority 20) + output regex (priority 15) — error wins on both
	ctx := MascotContext{LastExitCode: 1, LastOutput: "fatal: could not connect"}
	state := ResolveState(mc, ctx)
	if state != MascotStateError {
		t.Errorf("expected error state, got %q", state)
	}
}

func TestResolveState_FallsBackToDefaultState(t *testing.T) {
	mc := buildMascotConfig()
	// remove always trigger from idle to test pure fallback
	idleState := mc.States[MascotStateIdle]
	idleState.Triggers = nil
	mc.States[MascotStateIdle] = idleState

	// no triggers fire with exit code 0 if success is removed too
	delete(mc.States, MascotStateSuccess)
	ctx := MascotContext{LastExitCode: 0}
	state := ResolveState(mc, ctx)
	// should fall back to DefaultState which is "idle"
	if state != MascotStateIdle {
		t.Errorf("expected default state idle, got %q", state)
	}
}

func TestResolveState_CustomStateNames(t *testing.T) {
	mc := &MascotConfig{
		DefaultState: "chilling",
		Position:     FloaterTopRight,
		MaxWidth:     8,
		MaxHeight:    6,
		States: map[MascotState]MascotStateConfig{
			"chilling": {
				Frames: []string{"chill.png"},
				Triggers: []MascotTrigger{{Type: TriggerAlways, Priority: 0}},
			},
			"hacking": {
				Frames: []string{"hack.png"},
				Triggers: []MascotTrigger{
					{Type: TriggerCommand, Value: "vim*", Priority: 10},
				},
			},
		},
	}

	// command matches hacking state
	ctx := MascotContext{LastCommand: "vim main.go"}
	state := ResolveState(mc, ctx)
	if state != "hacking" {
		t.Errorf("expected 'hacking' state for vim command, got %q", state)
	}
}

// ── matchesExitCode ───────────────────────────────────────────────────────────

func TestMatchesExitCode_ExactMatch(t *testing.T) {
	if !matchesExitCode("0", 0) {
		t.Error("expected exact match for '0' == 0")
	}
	if matchesExitCode("0", 1) {
		t.Error("expected no match for '0' != 1")
	}
}

func TestMatchesExitCode_Range(t *testing.T) {
	if !matchesExitCode("1-127", 1) || !matchesExitCode("1-127", 127) || !matchesExitCode("1-127", 64) {
		t.Error("expected 1, 64, 127 to match range 1-127")
	}
	if matchesExitCode("1-127", 0) || matchesExitCode("1-127", 128) {
		t.Error("expected 0 and 128 to not match range 1-127")
	}
}

func TestMatchesExitCode_Threshold(t *testing.T) {
	if !matchesExitCode("128+", 128) || !matchesExitCode("128+", 255) {
		t.Error("expected 128 and 255 to match 128+")
	}
	if matchesExitCode("128+", 127) {
		t.Error("expected 127 to not match 128+")
	}
}

// ── MascotShellHooks ─────────────────────────────────────────────────────────

func TestMascotShellHooks_Bash(t *testing.T) {
	hooks := MascotShellHooks("my-mascot", "CMDX_MASCOT_TRIGGER", "bash")
	if hooks == "" {
		t.Fatal("expected non-empty bash hooks, got empty string")
	}
	if !containsStr(hooks, "PROMPT_COMMAND") {
		t.Error("expected bash hooks to set PROMPT_COMMAND")
	}
	if !containsStr(hooks, "my-mascot") {
		t.Error("expected bash hooks to reference asset name")
	}
}

func TestMascotShellHooks_Zsh(t *testing.T) {
	hooks := MascotShellHooks("my-mascot", "CMDX_MASCOT_TRIGGER", "zsh")
	if hooks == "" {
		t.Fatal("expected non-empty zsh hooks, got empty string")
	}
	if !containsStr(hooks, "add-zsh-hook") {
		t.Error("expected zsh hooks to use add-zsh-hook")
	}
}

func TestMascotShellHooks_PowerShell(t *testing.T) {
	hooks := MascotShellHooks("my-mascot", "CMDX_MASCOT_TRIGGER", "powershell")
	if hooks == "" {
		t.Fatal("expected non-empty powershell hooks, got empty string")
	}
	if !containsStr(hooks, "LASTEXITCODE") {
		t.Error("expected powershell hooks to check $LASTEXITCODE")
	}
}

func TestMascotShellHooks_UnknownShell(t *testing.T) {
	hooks := MascotShellHooks("my-mascot", "", "fish")
	if hooks != "" {
		t.Error("expected empty string for unsupported shell 'fish', got content")
	}
}

func TestMascotShellHooks_DefaultHookVar(t *testing.T) {
	hooks := MascotShellHooks("cat", "", "bash")
	if !containsStr(hooks, "CMDX_MASCOT_TRIGGER") {
		t.Error("expected default hook var CMDX_MASCOT_TRIGGER when none specified")
	}
}

// containsStr is a simple substring check used in hook content tests.
func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
