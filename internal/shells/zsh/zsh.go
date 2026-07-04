// Package zsh implements the shells.Shell interface for zsh, injecting
// cmdX theme configuration into ~/.zshrc.
package zsh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/shells"
)

// Zsh implements shells.Shell for the zsh shell.
type Zsh struct{}

// New creates a new Zsh shell handler.
func New() *Zsh {
	return &Zsh{}
}

// Name returns the shell identifier "zsh".
func (z *Zsh) Name() string {
	return "zsh"
}

// Detect reports whether zsh is the user's currently configured shell,
// based on the $SHELL environment variable.
func (z *Zsh) Detect() bool {
	shell := os.Getenv("SHELL")
	return strings.Contains(shell, "zsh")
}

// ProfilePath returns the path to the user's ~/.zshrc file.
func (z *Zsh) ProfilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(home, ".zshrc"), nil
}

// Inject writes the theme's prompt, banner, and color exports into
// ~/.zshrc, replacing any previous cmdx injection first.
func (z *Zsh) Inject(t *config.Theme) error {
	path, err := z.ProfilePath()
	if err != nil {
		return err
	}

	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	}

	existing = removeInjection(existing)
	block := buildZshBlock(t)
	final := existing + "\n" + block + "\n"

	return os.WriteFile(path, []byte(final), 0644)
}

// Remove strips any cmdx-injected block from ~/.zshrc.
func (z *Zsh) Remove() error {
	path, err := z.ProfilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read .zshrc: %w", err)
	}

	cleaned := removeInjection(string(data))
	return os.WriteFile(path, []byte(cleaned), 0644)
}

// IsInjected reports whether a cmdx theme block is currently present in
// ~/.zshrc. A missing or unreadable profile is treated as "not injected"
// rather than an error.
func (z *Zsh) IsInjected() (bool, error) {
	path, err := z.ProfilePath()
	if err != nil {
		return false, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}

	return strings.Contains(string(data), shells.InjectStart), nil
}

// buildZshBlock renders the full injected zsh config block for a theme:
// prompt function, git branch helper, startup banner, and color exports.
//
// All free-text theme fields (name, author, description, banner text,
// prompt symbol) pass through shells.SanitizeForShell first, since theme
// JSON may come from an untrusted source (the community registry or a
// shared file) and is otherwise embedded directly into this script.
func buildZshBlock(t *config.Theme) string {
	primary := hexToAnsi(t.Colors.Primary)
	accent := hexToAnsi(t.Colors.Accent)
	muted := hexToAnsi(t.Colors.Muted)
	symbol := shells.SanitizeForShell(t.Prompt.Symbol)
	name := shells.SanitizeForShell(t.Meta.Name)
	author := shells.SanitizeForShell(t.Meta.Author)
	bannerText := shells.SanitizeForShell(t.Banner.Text)
	description := shells.SanitizeForShell(t.Meta.Description)

	promptStyle := "single"
	if t.Prompt.Style == "multiline" {
		promptStyle = "multiline"
	}

	var promptLine string
	if promptStyle == "multiline" {
		promptLine = fmt.Sprintf(
			`PROMPT='%%F{%s}%%~%%f %%F{%s}$(git_branch)%%f
%%F{%s}%s%%f '`,
			primary, muted, accent, symbol,
		)
	} else {
		promptLine = fmt.Sprintf(
			`PROMPT='%%F{%s}%%~%%f %%F{%s}$(git_branch)%%f %%F{%s}%s%%f '`,
			primary, muted, accent, symbol,
		)
	}

	return fmt.Sprintf(`%s
# Theme: %s by %s

# Git branch helper
git_branch() {
    git branch 2>/dev/null | grep '^*' | sed 's/* /(/;s/$/)/'
}

# Prompt
autoload -Uz colors && colors
%s

# Startup Banner
echo ""
echo "  \033[38;2;%sm%s $USER\033[0m"
echo "  \033[38;2;%sm%s\033[0m"
echo ""

# Color exports
export CMDX_PRIMARY="%s"
export CMDX_ACCENT="%s"
export CMDX_SUCCESS="%s"
export CMDX_MUTED="%s"
%s`,
		shells.InjectStart,
		name, author,
		promptLine,
		hexToRgbAnsi(t.Colors.Primary), bannerText,
		hexToRgbAnsi(t.Colors.Muted), description,
		t.Colors.Primary,
		t.Colors.Accent,
		t.Colors.Success,
		t.Colors.Muted,
		shells.InjectEnd,
	)
}

// hexToAnsi converts a hex color to a zsh-compatible 256-color index
// (returned as a string), or "white" if the hex string is malformed.
func hexToAnsi(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "white"
	}

	var r, g, b int
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return "white"
	}

	// convert to 256-color xterm index
	if r == g && g == b {
		if r < 8 {
			return "0"
		}
		if r > 248 {
			return "15"
		}
		idx := ((r-8)/247)*24 + 232
		return fmt.Sprintf("%d", idx)
	}

	r256 := (r * 5) / 255
	g256 := (g * 5) / 255
	b256 := (b * 5) / 255
	return fmt.Sprintf("%d", 16+36*r256+6*g256+b256)
}

// hexToRgbAnsi returns "R;G;B" for use in \033[38;2;R;G;Bm true-color
// escape codes, or "255;255;255" (white) if the hex string is malformed.
func hexToRgbAnsi(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "255;255;255"
	}
	var r, g, b int
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return "255;255;255"
	}
	return fmt.Sprintf("%d;%d;%d", r, g, b)
}

// removeInjection strips every cmdx-marked block from content, including
// repeated/stale injections from prior runs.
func removeInjection(content string) string {
	for {
		start := strings.Index(content, shells.InjectStart)
		end := strings.Index(content, shells.InjectEnd)
		if start == -1 || end == -1 {
			break
		}
		content = content[:start] + content[end+len(shells.InjectEnd):]
	}
	return strings.TrimSpace(content)
}
