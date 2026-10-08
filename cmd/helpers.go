package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abhigyanwebber/cmd-customizer/internal/assets"
	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/shells"
	"github.com/abhigyanwebber/cmd-customizer/internal/shells/bash"
	"github.com/abhigyanwebber/cmd-customizer/internal/shells/powershell"
	"github.com/abhigyanwebber/cmd-customizer/internal/shells/zsh"
	"github.com/abhigyanwebber/cmd-customizer/internal/theme"
)

// detectShell auto-detects the running shell and returns the right
// implementation along with its identifier name. Used by every command
// that injects or removes shell configuration (theme inject/remove,
// edit, create).
func detectShell() (shells.Shell, string) {
	ps := powershell.New()
	if ps.Detect() {
		return ps, "powershell"
	}

	z := zsh.New()
	if z.Detect() {
		return z, "zsh"
	}

	b := bash.New()
	if b.Detect() {
		return b, "bash"
	}

	return nil, ""
}

// getThemesDir resolves the themes directory relative to either the
// current working directory (dev mode) or the binary's own location
// (installed mode).
// getThemesDir resolves the themes directory. Resolution order:
//  1. CMDX_THEMES_DIR environment variable, if set — lets external
//     tooling (the VS Code extension, custom scripts) point cmdx at a
//     specific themes directory without relying on cwd or binary
//     location.
//  2. ./themes relative to the current working directory (dev mode).
//  3. themes/ next to the binary's own location (installed mode).
func getThemesDir() string {
	if dir := os.Getenv("CMDX_THEMES_DIR"); dir != "" {
		return dir
	}
	wd, err := os.Getwd()
	if err == nil {
		local := filepath.Join(wd, "themes")
		if _, err := os.Stat(local); err == nil {
			return local
		}
	}
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "themes")
}

// loadThemeOrExit opens a theme manager rooted at the themes directory and
// loads the named theme, printing a formatted error and exiting the
// process on any failure. This collapses the repeated
// "new manager -> load -> check err -> os.Exit" pattern that previously
// appeared in nearly every theme subcommand.
func loadThemeOrExit(name string) (*theme.Manager, *config.Theme) {
	m, err := theme.NewManager(getThemesDir())
	if err != nil {
		fmt.Println("✗ Error:", err)
		os.Exit(1)
	}

	t, err := m.Load(name)
	if err != nil {
		fmt.Println("✗ Error:", err)
		os.Exit(1)
	}

	return m, t
}

// showActiveFloaters renders every currently-active floater (one per
// corner, tracked independently via "cmdx asset use --as floater") so
// commands like "theme preview" can display them alongside the rest of
// a theme's visual elements.
//
// Floaters can be linked into a theme's `assets.floater` field like any
// other slot (see activateThemeAssets) and get their name auto-activated
// on `theme apply`/`inject` — but a floater's *corner position* always
// comes from its own manifest (or an explicit `--position` override via
// `asset use`), never from the theme JSON. This function specifically
// reads live per-corner state directly from the asset state directory
// (not from the loaded *config.Theme) because it needs to reflect
// whichever floaters are actually active right now, including ones set
// via a direct `asset use --as floater --position X` outside of any
// theme. Missing or unconfigured corners are silently skipped; this is
// a best-effort visual extra, not a required step, so errors here are
// reported but never fatal to the calling command.
//
// Uses getAssetsDir, defined in cmd/asset.go.
func showActiveFloaters() {
	stateDir := filepath.Join(getAssetsDir(), ".state")

	m, err := assets.NewManager(getAssetsDir())
	if err != nil {
		return
	}

	var shown bool
	for _, pos := range assets.ValidFloaterPositions {
		statePath := filepath.Join(stateDir, "floater-"+string(pos)+".txt")
		data, err := os.ReadFile(statePath)
		if err != nil {
			continue
		}

		name := string(data)
		if !shown {
			fmt.Println("[ Active Floaters ]")
			shown = true
		}

		if err := m.PreviewFloater(name); err != nil {
			fmt.Printf("  ✗ Could not render floater '%s' (%s): %v\n", name, pos, err)
		}
	}

	if shown {
		fmt.Println()
	}
}

// showActiveStatusBar renders the currently active status bar (if any)
// after a theme preview, mirroring showActiveFloaters above. Reads the
// name written to .state/status-bar.txt by activateAsset and previews
// it for the current shell; errors here are reported but never fatal
// to the calling command.
//
// Uses getAssetsDir, defined in cmd/asset.go.
func showActiveStatusBar() {
	statePath := filepath.Join(getAssetsDir(), ".state", "status-bar.txt")
	data, err := os.ReadFile(statePath)
	if err != nil {
		return
	}
	name := string(data)

	m, err := assets.NewManager(getAssetsDir())
	if err != nil {
		return
	}

	_, shell := detectShell()
	colors := loadThemeColors()

	fmt.Println("[ Active Status Bar ]")
	if err := m.PreviewStatusBar(name, shell, colors); err != nil {
		fmt.Printf("  ✗ Could not render status bar '%s': %v\n", name, err)
	}
	fmt.Println()
}

// activateAsset writes the ".state/<slot>.txt" file that marks name as
// the active asset for the given slot, exactly matching what
// "cmdx asset use --as <slot>" does. Extracted here so both the manual
// "asset use" command and the automatic "theme apply" activation (see
// activateThemeAssets) share one code path rather than two copies that
// could silently drift apart.
//
// slot must be one of: spinner, banner, divider, icons, mascot,
// status-bar, sound, floater. For floater, position may be empty to
// use the floater's own configured position.
func activateAsset(assetsDir, name, slot, position string) error {
	validTypes := map[string]assets.AssetType{
		"spinner":    assets.AssetTypeSpinner,
		"banner":     assets.AssetTypeBanner,
		"divider":    assets.AssetTypeDivider,
		"icons":      assets.AssetTypeIcon,
		"floater":    assets.AssetTypeFloater,
		"mascot":     assets.AssetTypeMascot,
		"status-bar": assets.AssetTypeStatusBar,
		"sound":      assets.AssetTypeSound,
	}

	assetType, ok := validTypes[slot]
	if !ok {
		return fmt.Errorf("unknown slot %q", slot)
	}

	m, err := assets.NewManager(assetsDir)
	if err != nil {
		return err
	}

	a, _, err := m.Get(name, assetType)
	if err != nil {
		return fmt.Errorf("asset %q not found as type %q: %w", name, slot, err)
	}

	stateDir := filepath.Join(assetsDir, ".state")
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return err
	}

	if slot == "floater" {
		resolvedPosition := position
		if resolvedPosition == "" {
			resolvedPosition = string(a.Floater.Position)
		}
		if !assets.IsValidFloaterPosition(assets.FloaterPosition(resolvedPosition)) {
			return fmt.Errorf("invalid floater position %q", resolvedPosition)
		}
		statePath := filepath.Join(stateDir, "floater-"+resolvedPosition+".txt")
		return os.WriteFile(statePath, []byte(name), 0644)
	}

	statePath := filepath.Join(stateDir, slot+".txt")
	return os.WriteFile(statePath, []byte(name), 0644)
}

// activateThemeAssets activates every asset slot a theme declares in
// its "assets" block (spinner, banner, divider, icons, mascot, floater,
// status_bar, sound), continuing past individual failures rather than
// aborting the whole theme apply — a broken asset reference shouldn't
// prevent the theme's colors/prompt/etc. from applying. Returns the
// list of slots that were successfully activated and any per-slot
// errors encountered, so the caller can report both.
func activateThemeAssets(assetsDir string, ta config.ThemeAssets) (activated []string, errs []error) {
	type slotRef struct {
		slot string
		name string
	}
	slots := []slotRef{
		{"spinner", ta.Spinner},
		{"banner", ta.Banner},
		{"divider", ta.Divider},
		{"icons", ta.Icons},
		{"mascot", ta.Mascot},
		{"floater", ta.Floater},
		{"status-bar", ta.StatusBar},
		{"sound", ta.Sound},
	}

	for _, s := range slots {
		if s.name == "" {
			continue
		}
		if err := activateAsset(assetsDir, s.name, s.slot, ""); err != nil {
			errs = append(errs, fmt.Errorf("%s %q: %w", s.slot, s.name, err))
			continue
		}
		activated = append(activated, fmt.Sprintf("%s (%s)", s.name, s.slot))
	}

	return activated, errs
}

// syncAssetHooks installs or removes the mascot/sound theme shell hook
// block in the given shell's profile, based on what the theme's assets
// block currently links. Writes to a separate marked block
// (shells.AssetHooksStart/End) from the theme's own injection block,
// so removing/reinjecting a theme (which only touches
// shells.InjectStart/End) never orphans or clobbers the asset hooks,
// and vice versa. If the theme links neither a mascot nor a sound
// asset, any previously-injected asset hooks block is cleanly removed
// rather than left stale from an earlier theme.
func syncAssetHooks(profilePath string, shellName string, ta config.ThemeAssets) error {
	existing := ""
	if data, err := os.ReadFile(profilePath); err == nil {
		existing = string(data)
	}
	existing = stripAssetHooksBlock(existing)

	m, err := assets.NewManager(getAssetsDir())
	if err != nil {
		return err
	}

	var hookBlocks []string
	if ta.Mascot != "" {
		if hooks, err := m.MascotHooks(ta.Mascot, shellName); err == nil {
			hookBlocks = append(hookBlocks, hooks)
		} else {
			fmt.Printf("  ! Could not generate mascot hooks for '%s': %v\n", ta.Mascot, err)
		}
	}
	if ta.Sound != "" {
		if hooks, err := m.SoundHooks(ta.Sound, shellName); err == nil {
			hookBlocks = append(hookBlocks, hooks)
		} else {
			fmt.Printf("  ! Could not generate sound hooks for '%s': %v\n", ta.Sound, err)
		}
	}

	if len(hookBlocks) == 0 {
		return os.WriteFile(profilePath, []byte(existing), 0644)
	}

	block := shells.AssetHooksStart + "\n" + joinHookBlocks(hookBlocks) + "\n" + shells.AssetHooksEnd
	final := existing + "\n" + block + "\n"
	return os.WriteFile(profilePath, []byte(final), 0644)
}

// stripAssetHooksBlock removes any existing asset-hooks block(s) from
// content, mirroring the same start/end-marker stripping pattern each
// shell package uses internally for its own theme injection block.
func stripAssetHooksBlock(content string) string {
	for {
		start := strings.Index(content, shells.AssetHooksStart)
		end := strings.Index(content, shells.AssetHooksEnd)
		if start == -1 || end == -1 {
			break
		}
		content = content[:start] + content[end+len(shells.AssetHooksEnd):]
	}
	return strings.TrimSpace(content)
}

func joinHookBlocks(blocks []string) string {
	return strings.Join(blocks, "\n")
}
