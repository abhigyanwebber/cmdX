package assets

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const cooldownStateSubdir = ".state"

// PlaySound plays a resolved sound effect. Handles the platform-specific
// player selection, volume application (where supported), per-effect
// cooldown enforcement (persisted across process invocations, since
// each shell hook call is a fresh `cmdx` process), and blocking vs
// async playback.
//
// assetDir is the sound asset's own directory (for resolving effect.File
// and locating the cooldown state file). assetName + soundName together
// form the cooldown state file's identity.
func PlaySound(assetDir, assetName, soundName string, effect SoundEffect, cfg *SoundThemeConfig) error {
	if effect.CooldownMs > 0 {
		onCooldown, err := isOnCooldown(assetDir, assetName, soundName, effect.CooldownMs)
		if err != nil {
			// cooldown tracking failure shouldn't block playback entirely —
			// treat as "not on cooldown" and proceed, since silently
			// refusing to ever play a sound due to a state-file glitch
			// would be a worse user experience than an occasional
			// cooldown miss.
			onCooldown = false
		}
		if onCooldown {
			return nil
		}
	}

	filePath := filepath.Join(assetDir, effect.File)
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("sound file not found: %w", err)
	}

	volume := effect.Volume
	if volume <= 0 {
		volume = 1.0
	}
	if cfg.GlobalVolume > 0 {
		volume *= cfg.GlobalVolume
	}

	cmd, err := buildPlayerCommand(filePath, volume, cfg.Player)
	if err != nil {
		return err
	}

	if effect.Async {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("could not start playback: %w", err)
		}
		// intentionally not waiting — fire and forget
	} else {
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("playback failed: %w", err)
		}
	}

	if effect.CooldownMs > 0 {
		_ = recordPlayed(assetDir, assetName, soundName)
	}

	return nil
}

// buildPlayerCommand resolves which audio player to invoke and builds
// the exec.Cmd for it. If customPlayer is set, it's used verbatim (with
// "%f" substituted for filePath) rather than the platform default —
// this is the developer-freedom escape hatch documented on
// SoundThemeConfig.Player. The custom template is split on whitespace
// and executed directly via argv, never through a shell, so it isn't
// vulnerable to shell injection — the tradeoff is that template
// arguments containing spaces aren't supported.
func buildPlayerCommand(filePath string, volume float64, customPlayer string) (*exec.Cmd, error) {
	if customPlayer != "" {
		tokens := strings.Fields(customPlayer)
		if len(tokens) == 0 {
			return nil, fmt.Errorf("custom player command is empty")
		}
		args := make([]string, len(tokens)-1)
		foundPlaceholder := false
		for i, tok := range tokens[1:] {
			if tok == "%f" {
				args[i] = filePath
				foundPlaceholder = true
			} else {
				args[i] = tok
			}
		}
		if !foundPlaceholder {
			// no %f placeholder anywhere — append the file path as the
			// final argument, the most common convention for CLI
			// media players when no template is given
			args = append(args, filePath)
		}
		return exec.Command(tokens[0], args...), nil
	}

	switch runtime.GOOS {
	case "darwin":
		// afplay accepts a fractional volume via -v
		return exec.Command("afplay", "-v", formatVolume(volume), filePath), nil

	case "windows":
		return buildWindowsPlayerCommand(filePath)

	default: // linux and other unix-likes
		return buildLinuxPlayerCommand(filePath, volume)
	}
}

// buildLinuxPlayerCommand tries paplay, then aplay, then ffplay, in
// that order — covers PulseAudio (most modern distros), plain ALSA
// (minimal/server distros), and a broader-format fallback if ffmpeg is
// installed. This "try several, use whichever exists" approach is
// itself a developer-freedom concession: we don't assume one audio
// stack, we adapt to whatever the user's system actually has.
func buildLinuxPlayerCommand(filePath string, volume float64) (*exec.Cmd, error) {
	if _, err := exec.LookPath("paplay"); err == nil {
		// paplay's --volume takes 0-65536, where 65536 = 100%
		vol := strconv.Itoa(int(volume * 65536))
		return exec.Command("paplay", "--volume="+vol, filePath), nil
	}
	if _, err := exec.LookPath("aplay"); err == nil {
		// aplay has no volume flag — plays at system mixer level
		return exec.Command("aplay", "-q", filePath), nil
	}
	if _, err := exec.LookPath("ffplay"); err == nil {
		volPercent := strconv.Itoa(int(volume * 100))
		return exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", "-volume", volPercent, filePath), nil
	}
	return nil, fmt.Errorf("no audio player found (tried paplay, aplay, ffplay) — install one, or set sound.player to a custom command")
}

// buildWindowsPlayerCommand uses PowerShell's built-in
// System.Media.SoundPlayer, which requires no extra install but only
// supports WAV files. The file path is embedded into a PowerShell
// single-quoted string literal, so any single quotes in the path are
// escaped (doubled) first — the same class of care taken for shell
// injection elsewhere in cmdX (see internal/shells.SanitizeForShell),
// applied here because a sound asset's file path could in principle
// come from a downloaded/shared theme.
func buildWindowsPlayerCommand(filePath string) (*exec.Cmd, error) {
	escaped := strings.ReplaceAll(filePath, "'", "''")
	script := fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", escaped)
	return exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script), nil
}

func formatVolume(v float64) string {
	return strconv.FormatFloat(v, 'f', 3, 64)
}

// ── cooldown state tracking ──────────────────────────────────────────

// cooldownStatePath returns the path to the timestamp file tracking
// when a given sound effect was last played. Sits alongside the other
// per-asset state files (see cmd/asset.go's floater/spinner/banner
// state tracking) at <assetsDir>/.state/ — assetDir here is the sound
// asset's own directory (<assetsDir>/sounds/<name>), so we go up two
// levels to reach the assets root.
func cooldownStatePath(assetDir, assetName, soundName string) string {
	assetsRoot := filepath.Dir(filepath.Dir(assetDir))
	stateDir := filepath.Join(assetsRoot, cooldownStateSubdir)
	return filepath.Join(stateDir, fmt.Sprintf("sound-%s-%s-lastplayed.txt", assetName, soundName))
}

// isOnCooldown reports whether soundName was played more recently than
// cooldownMs milliseconds ago.
func isOnCooldown(assetDir, assetName, soundName string, cooldownMs int) (bool, error) {
	path := cooldownStatePath(assetDir, assetName, soundName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // never played before — not on cooldown
		}
		return false, err
	}

	lastMs, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return false, err
	}

	elapsed := time.Now().UnixMilli() - lastMs
	return elapsed < int64(cooldownMs), nil
}

// recordPlayed writes the current timestamp to the cooldown state file
// for soundName, creating the state directory if needed.
func recordPlayed(assetDir, assetName, soundName string) error {
	path := cooldownStatePath(assetDir, assetName, soundName)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	now := strconv.FormatInt(time.Now().UnixMilli(), 10)
	return os.WriteFile(path, []byte(now), 0644)
}
