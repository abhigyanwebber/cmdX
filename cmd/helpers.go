package main

import (
	"fmt"
	"os"
	"path/filepath"

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
// Floaters are deliberately not part of the theme JSON schema — they're
// a per-machine asset selection, not theme-portable config — so this
// reads directly from the asset state directory rather than from the
// loaded *config.Theme. Missing or unconfigured corners are silently
// skipped; this is a best-effort visual extra, not a required step, so
// errors here are reported but never fatal to the calling command.
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
