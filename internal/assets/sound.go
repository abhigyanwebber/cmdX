package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// ValidateSoundTheme checks a SoundThemeConfig for required fields and
// verifies every referenced audio file exists under assetDir.
func ValidateSoundTheme(s *SoundThemeConfig, assetDir string) error {
	if s == nil {
		return fmt.Errorf("sound config is nil")
	}
	if len(s.Sounds) == 0 {
		return fmt.Errorf("sound theme must define at least one sound effect")
	}
	if s.GlobalVolume < 0 || s.GlobalVolume > 1 {
		if s.GlobalVolume != 0 { // 0 is the "unset, use default" zero value
			return fmt.Errorf("global_volume must be between 0.0 and 1.0")
		}
	}

	for name, effect := range s.Sounds {
		if effect.File == "" {
			return fmt.Errorf("sound %q must specify a file", name)
		}
		path := filepath.Join(assetDir, effect.File)
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("sound %q file %q not found at %s", name, effect.File, path)
		}
		if effect.Volume < 0 || effect.Volume > 1 {
			return fmt.Errorf("sound %q volume must be between 0.0 and 1.0", name)
		}
		if effect.CooldownMs < 0 {
			return fmt.Errorf("sound %q cooldown_ms must not be negative", name)
		}
		for _, trigger := range effect.Triggers {
			if err := validateSoundTrigger(trigger, name); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateSoundTrigger checks a trigger's value is syntactically valid
// for its type. Deliberately reuses validateExitCodeValue from mascot.go
// (same package) rather than duplicating exit-code range/threshold
// parsing — sound and mascot triggers share the exact same vocabulary.
func validateSoundTrigger(t SoundTrigger, soundName string) error {
	switch t.Type {
	case TriggerExitCode:
		if err := validateExitCodeValue(t.Value); err != nil {
			return fmt.Errorf("sound %q exit_code trigger: %w", soundName, err)
		}
	case TriggerOutputRegex:
		if _, err := regexp.Compile(t.Value); err != nil {
			return fmt.Errorf("sound %q output_regex trigger: invalid regex %q: %w", soundName, t.Value, err)
		}
	case TriggerIdleTime:
		if _, err := parseIdleSeconds(t.Value); err != nil {
			return fmt.Errorf("sound %q idle_time trigger: %w", soundName, err)
		}
	case TriggerGitStatus:
		valid := map[string]bool{"dirty": true, "clean": true, "untracked": true, "ahead": true, "behind": true}
		if !valid[t.Value] {
			return fmt.Errorf("sound %q git_status trigger: invalid value %q", soundName, t.Value)
		}
	case TriggerEnvVar:
		if t.Value == "" {
			return fmt.Errorf("sound %q env_var trigger: value must be VAR_NAME or VAR_NAME=expected_value", soundName)
		}
	case TriggerCommand:
		if t.Value == "" {
			return fmt.Errorf("sound %q command trigger: value must be a glob pattern", soundName)
		}
	case TriggerAlways:
		// no value needed
	default:
		return fmt.Errorf("sound %q: unknown trigger type %q", soundName, t.Type)
	}
	return nil
}

// parseIdleSeconds validates that v is a plain integer number of
// seconds, matching the same strictness as mascot.go's TriggerIdleTime
// validation (strconv.Atoi rejects trailing garbage like "30abc",
// unlike fmt.Sscanf which would silently accept it).
func parseIdleSeconds(v string) (int, error) {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("value must be an integer number of seconds, got %q", v)
	}
	return n, nil
}

// ResolveSound determines which sound effect (if any) should play
// given the current shell context, using the same priority-based
// trigger resolution as mascot state resolution: the highest-priority
// matching trigger across all sounds wins. Returns ok=false if nothing
// matches (most invocations won't have a sound to play — that's normal,
// not an error).
func ResolveSound(s *SoundThemeConfig, ctx MascotContext) (name string, effect SoundEffect, ok bool) {
	bestPriority := -1

	for soundName, soundEffect := range s.Sounds {
		for _, trigger := range soundEffect.Triggers {
			// SoundTrigger and MascotTrigger share an identical field
			// shape (Type, Value, Priority) by design, so we bridge to
			// mascot.go's matchesTrigger directly rather than
			// duplicating trigger-matching logic for every trigger type.
			bridged := MascotTrigger{Type: trigger.Type, Value: trigger.Value, Priority: trigger.Priority}
			if matchesTrigger(bridged, ctx) {
				if trigger.Priority > bestPriority {
					name = soundName
					effect = soundEffect
					ok = true
					bestPriority = trigger.Priority
				}
			}
		}
	}

	return name, effect, ok
}

// SoundShellHooks returns the shell code that should be injected into
// a shell profile to enable sound theme playback. Mirrors
// MascotShellHooks' structure exactly (same hook points: bash's
// PROMPT_COMMAND, zsh's precmd/preexec, PowerShell's prompt function
// override) but calls "asset sound-play" instead of "asset
// mascot-state". Kept as a separate function rather than a shared one
// because the two commands take different flags and the two hook
// blocks are independently toggle-able (a theme can use a mascot
// without sounds, or vice versa).
func SoundShellHooks(assetName string, shell string) string {
	switch shell {
	case "bash":
		return fmt.Sprintf(`
# cmdX sound theme hook
__cmdx_sound_hook() {
    local exit_code=$?
    local last_cmd=$(history 1 | sed 's/^[ ]*[0-9]*[ ]*//')
    cmdx asset sound-play "%s" --exit-code "$exit_code" --command "$last_cmd" >/dev/null 2>&1 &
}
PROMPT_COMMAND="__cmdx_sound_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
`, assetName)

	case "zsh":
		return fmt.Sprintf(`
# cmdX sound theme hook
__cmdx_sound_exit_code=0
__cmdx_sound_last_cmd=""

__cmdx_sound_preexec() { __cmdx_sound_last_cmd="$1" }
__cmdx_sound_precmd() {
    __cmdx_sound_exit_code=$?
    cmdx asset sound-play "%s" --exit-code "${__cmdx_sound_exit_code}" --command "${__cmdx_sound_last_cmd}" >/dev/null 2>&1 &
}

autoload -Uz add-zsh-hook
add-zsh-hook preexec __cmdx_sound_preexec
add-zsh-hook precmd __cmdx_sound_precmd
`, assetName)

	case "powershell":
		return fmt.Sprintf(`
# cmdX sound theme hook
function __cmdx_sound_prompt {
    $exitCode = $LASTEXITCODE
    Start-Job -ScriptBlock { param($code) cmdx asset sound-play "%s" --exit-code $code } -ArgumentList $exitCode | Out-Null
}
$__cmdx_sound_original_prompt = $function:prompt
function prompt {
    __cmdx_sound_prompt
    & $__cmdx_sound_original_prompt
}
`, assetName)

	default:
		return ""
	}
}
