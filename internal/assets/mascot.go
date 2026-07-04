package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ValidateMascot checks a MascotConfig for required fields and
// verifies that every referenced PNG file exists under assetDir.
func ValidateMascot(m *MascotConfig, assetDir string) error {
	if m == nil {
		return fmt.Errorf("mascot config is nil")
	}
	if len(m.States) == 0 {
		return fmt.Errorf("mascot must define at least one state")
	}
	if !IsValidFloaterPosition(m.Position) {
		return fmt.Errorf("invalid mascot position %q: must be one of top-left, top-right, bottom-left, bottom-right", m.Position)
	}
	if m.MaxWidth <= 0 {
		return fmt.Errorf("mascot max_width must be greater than 0")
	}
	if m.MaxHeight <= 0 {
		return fmt.Errorf("mascot max_height must be greater than 0")
	}

	// default state must be defined
	defaultState := m.DefaultState
	if defaultState == "" {
		defaultState = MascotStateIdle
	}
	if _, ok := m.States[defaultState]; !ok {
		return fmt.Errorf("default_state %q is not defined in states", defaultState)
	}

	// validate each state
	for stateName, state := range m.States {
		if len(state.Frames) == 0 {
			return fmt.Errorf("state %q must have at least one frame", stateName)
		}
		for _, frame := range state.Frames {
			path := filepath.Join(assetDir, frame)
			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("state %q frame %q not found at %s", stateName, frame, path)
			}
		}
		for _, frame := range state.TransitionFrames {
			path := filepath.Join(assetDir, frame)
			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("state %q transition frame %q not found at %s", stateName, frame, path)
			}
		}
		// validate render override tint if set
		if state.RenderOverride != nil && state.RenderOverride.Tint != "" {
			if !strings.HasPrefix(state.RenderOverride.Tint, "#") || (len(state.RenderOverride.Tint) != 7 && len(state.RenderOverride.Tint) != 4) {
				return fmt.Errorf("state %q render_override tint %q is not a valid hex color", stateName, state.RenderOverride.Tint)
			}
		}
		// validate trigger values
		for _, trigger := range state.Triggers {
			if err := validateTrigger(trigger, stateName); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateTrigger checks that a trigger's value is syntactically valid
// for its type.
func validateTrigger(t MascotTrigger, stateName MascotState) error {
	switch t.Type {
	case TriggerExitCode:
		if err := validateExitCodeValue(t.Value); err != nil {
			return fmt.Errorf("state %q exit_code trigger: %w", stateName, err)
		}
	case TriggerOutputRegex:
		if _, err := regexp.Compile(t.Value); err != nil {
			return fmt.Errorf("state %q output_regex trigger: invalid regex %q: %w", stateName, t.Value, err)
		}
	case TriggerIdleTime:
		if _, err := strconv.Atoi(t.Value); err != nil {
			return fmt.Errorf("state %q idle_time trigger: value must be an integer number of seconds, got %q", stateName, t.Value)
		}
	case TriggerGitStatus:
		valid := map[string]bool{"dirty": true, "clean": true, "untracked": true, "ahead": true, "behind": true}
		if !valid[t.Value] {
			return fmt.Errorf("state %q git_status trigger: invalid value %q (must be dirty, clean, untracked, ahead, or behind)", stateName, t.Value)
		}
	case TriggerEnvVar:
		if t.Value == "" {
			return fmt.Errorf("state %q env_var trigger: value must be VAR_NAME or VAR_NAME=expected_value", stateName)
		}
	case TriggerCommand:
		// glob patterns are hard to validate statically — just check non-empty
		if t.Value == "" {
			return fmt.Errorf("state %q command trigger: value must be a glob pattern like 'go build*'", stateName)
		}
	case TriggerAlways:
		// no value needed
	default:
		return fmt.Errorf("state %q: unknown trigger type %q", stateName, t.Type)
	}
	return nil
}

// validateExitCodeValue checks that an exit_code trigger value is one of:
// "0", "1", "1-127", "128+" style patterns.
func validateExitCodeValue(v string) error {
	if v == "" {
		return fmt.Errorf("exit_code trigger value must not be empty")
	}
	// "128+" style
	if strings.HasSuffix(v, "+") {
		_, err := strconv.Atoi(strings.TrimSuffix(v, "+"))
		if err != nil {
			return fmt.Errorf("invalid exit code threshold %q (expected integer+)", v)
		}
		return nil
	}
	// "1-127" style range
	if strings.Contains(v, "-") {
		parts := strings.SplitN(v, "-", 2)
		if _, err := strconv.Atoi(parts[0]); err != nil {
			return fmt.Errorf("invalid exit code range start %q", parts[0])
		}
		if _, err := strconv.Atoi(parts[1]); err != nil {
			return fmt.Errorf("invalid exit code range end %q", parts[1])
		}
		return nil
	}
	// single value
	_, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("invalid exit code %q (must be integer, range like 1-127, or threshold like 128+)", v)
	}
	return nil
}

// ResolveState determines which mascot state should be active given the
// current environment context. It evaluates all states' triggers in
// priority order (highest Priority wins) and returns the first match.
// Falls back to the configured DefaultState (or "idle") if nothing matches.
func ResolveState(m *MascotConfig, ctx MascotContext) MascotState {
	type candidate struct {
		state    MascotState
		priority int
	}

	var best *candidate

	for stateName, stateCfg := range m.States {
		for _, trigger := range stateCfg.Triggers {
			if matchesTrigger(trigger, ctx) {
				if best == nil || trigger.Priority > best.priority {
					s := stateName
					p := trigger.Priority
					best = &candidate{state: s, priority: p}
				}
			}
		}
	}

	if best != nil {
		return best.state
	}

	// check for TriggerAlways state (lowest priority fallback)
	for stateName, stateCfg := range m.States {
		for _, trigger := range stateCfg.Triggers {
			if trigger.Type == TriggerAlways {
				return stateName
			}
		}
	}

	defaultState := m.DefaultState
	if defaultState == "" {
		defaultState = MascotStateIdle
	}
	return defaultState
}

// MascotContext holds the runtime values that trigger resolution uses
// to decide which state is active. All fields are optional — zero values
// are treated as "not available" and triggers that depend on them won't
// fire.
type MascotContext struct {
	// LastExitCode is the exit code of the most recent shell command.
	LastExitCode int

	// LastCommand is the text of the most recently executed command.
	LastCommand string

	// LastOutput is the combined stdout+stderr of the last command,
	// truncated to a reasonable size for regex matching.
	LastOutput string

	// IdleSeconds is how long the shell has been idle (no commands run).
	IdleSeconds int

	// GitStatus is the current git status of the working directory:
	// "dirty", "clean", "untracked", "ahead", "behind", or "" if not in a repo.
	GitStatus string

	// Env is a snapshot of relevant environment variables, used for
	// TriggerEnvVar matching without repeated os.Getenv calls.
	Env map[string]string
}

// matchesTrigger reports whether a trigger condition is satisfied by ctx.
func matchesTrigger(t MascotTrigger, ctx MascotContext) bool {
	switch t.Type {
	case TriggerAlways:
		return true

	case TriggerExitCode:
		return matchesExitCode(t.Value, ctx.LastExitCode)

	case TriggerOutputRegex:
		if ctx.LastOutput == "" {
			return false
		}
		re, err := regexp.Compile(t.Value)
		if err != nil {
			return false
		}
		return re.MatchString(ctx.LastOutput)

	case TriggerCommand:
		if ctx.LastCommand == "" {
			return false
		}
		matched, err := filepath.Match(t.Value, ctx.LastCommand)
		return err == nil && matched

	case TriggerIdleTime:
		threshold, err := strconv.Atoi(t.Value)
		if err != nil {
			return false
		}
		return ctx.IdleSeconds >= threshold

	case TriggerGitStatus:
		return ctx.GitStatus == t.Value

	case TriggerEnvVar:
		if ctx.Env == nil {
			return false
		}
		if strings.Contains(t.Value, "=") {
			parts := strings.SplitN(t.Value, "=", 2)
			return ctx.Env[parts[0]] == parts[1]
		}
		_, set := ctx.Env[t.Value]
		return set
	}
	return false
}

// matchesExitCode checks whether exitCode satisfies the trigger value pattern.
func matchesExitCode(pattern string, exitCode int) bool {
	if strings.HasSuffix(pattern, "+") {
		threshold, err := strconv.Atoi(strings.TrimSuffix(pattern, "+"))
		if err != nil {
			return false
		}
		return exitCode >= threshold
	}
	if strings.Contains(pattern, "-") {
		parts := strings.SplitN(pattern, "-", 2)
		lo, err1 := strconv.Atoi(parts[0])
		hi, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return false
		}
		return exitCode >= lo && exitCode <= hi
	}
	expected, err := strconv.Atoi(pattern)
	if err != nil {
		return false
	}
	return exitCode == expected
}

// RenderMascotState renders the frames for a specific mascot state,
// applying any per-state render overrides on top of the base config.
// Returns the rendered frames, frame interval in ms, and any transition
// frames (played once on state entry, before the main loop).
func RenderMascotState(a *Asset, assetDir string, state MascotState, baseOverrides RenderOverrides) (frames []string, transitionFrames []string, intervalMs int, err error) {
	if a.Mascot == nil {
		return nil, nil, 0, fmt.Errorf("asset has no mascot config")
	}

	stateCfg, ok := a.Mascot.States[state]
	if !ok {
		return nil, nil, 0, fmt.Errorf("state %q not defined in mascot config", state)
	}

	// merge render overrides: base → per-state override → caller override
	opts := ApplyOverrides(optionsFromConfig(a.Render), baseOverrides)
	if stateCfg.RenderOverride != nil {
		stateOverride := RenderOverrides{
			Width:     stateCfg.RenderOverride.Width,
			Height:    stateCfg.RenderOverride.Height,
			ColorMode: stateCfg.RenderOverride.ColorMode,
		}
		opts = ApplyOverrides(opts, stateOverride)
	}

	// render main frames
	framePaths := make([]string, len(stateCfg.Frames))
	for i, f := range stateCfg.Frames {
		framePaths[i] = filepath.Join(assetDir, f)
	}
	rendered, err := RenderFrames(framePaths, opts)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("could not render state %q frames: %w", state, err)
	}

	// render transition frames if any
	var renderedTransition []string
	if len(stateCfg.TransitionFrames) > 0 {
		tPaths := make([]string, len(stateCfg.TransitionFrames))
		for i, f := range stateCfg.TransitionFrames {
			tPaths[i] = filepath.Join(assetDir, f)
		}
		renderedTransition, err = RenderFrames(tPaths, opts)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not render state %q transition frames: %w", state, err)
		}
	}

	// determine interval: state-specific → global → default
	interval := stateCfg.IntervalMs
	if interval == 0 {
		interval = a.Mascot.GlobalIntervalMs
	}
	if interval == 0 {
		interval = 100
	}

	return rendered, renderedTransition, interval, nil
}

// MascotShellHooks returns the shell code that should be injected into
// a shell profile to enable mascot state tracking. The hook captures
// the exit code and last command after each command runs, then calls
// cmdx to update the mascot state.
//
// This is separate from the normal theme injection — mascot hooks can
// be installed independently via "cmdx asset use <name> --as mascot".
func MascotShellHooks(assetName string, hookVar string, shell string) string {
	if hookVar == "" {
		hookVar = "CMDX_MASCOT_TRIGGER"
	}

	switch shell {
	case "bash":
		return fmt.Sprintf(`
# cmdX mascot hook — tracks exit code and command for state resolution
__cmdx_mascot_hook() {
    local exit_code=$?
    local last_cmd=$(history 1 | sed 's/^[ ]*[0-9]*[ ]*//')
    export %s="exit_code=$exit_code,command=$last_cmd"
    cmdx asset mascot-state "%s" --exit-code "$exit_code" --command "$last_cmd" 2>/dev/null || true
}
PROMPT_COMMAND="__cmdx_mascot_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
`, hookVar, assetName)

	case "zsh":
		return fmt.Sprintf(`
# cmdX mascot hook — tracks exit code and command for state resolution
__cmdx_mascot_exit_code=0
__cmdx_mascot_last_cmd=""

__cmdx_mascot_preexec() { __cmdx_mascot_last_cmd="$1" }
__cmdx_mascot_precmd() {
    __cmdx_mascot_exit_code=$?
    export %s="exit_code=${__cmdx_mascot_exit_code},command=${__cmdx_mascot_last_cmd}"
    cmdx asset mascot-state "%s" --exit-code "${__cmdx_mascot_exit_code}" --command "${__cmdx_mascot_last_cmd}" 2>/dev/null || true
}

autoload -Uz add-zsh-hook
add-zsh-hook preexec __cmdx_mascot_preexec
add-zsh-hook precmd __cmdx_mascot_precmd
`, hookVar, assetName)

	case "powershell":
		return fmt.Sprintf(`
# cmdX mascot hook
function __cmdx_mascot_prompt {
    $exitCode = $LASTEXITCODE
    $env:%s = "exit_code=$exitCode"
    cmdx asset mascot-state "%s" --exit-code $exitCode 2>$null
}
$__cmdx_original_prompt = $function:prompt
function prompt {
    __cmdx_mascot_prompt
    & $__cmdx_original_prompt
}
`, hookVar, assetName)

	default:
		return ""
	}
}
