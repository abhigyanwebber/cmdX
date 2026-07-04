// Package bash implements the shells.Shell interface for bash, injecting
// cmdX theme configuration into ~/.bashrc or ~/.bash_profile.
package bash

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/shells"
)

// Bash implements shells.Shell for the bash shell.
type Bash struct{}

// New creates a new Bash shell handler.
func New() *Bash {
	return &Bash{}
}

// Name returns the shell identifier "bash".
func (b *Bash) Name() string {
	return "bash"
}

// Detect reports whether bash is the user's currently configured shell,
// based on the $SHELL environment variable.
func (b *Bash) Detect() bool {
	shell := os.Getenv("SHELL")
	return strings.Contains(shell, "bash")
}

// ProfilePath returns the path to the user's bash profile, preferring
// ~/.bashrc and falling back to ~/.bash_profile if .bashrc doesn't exist.
func (b *Bash) ProfilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}

	rc := filepath.Join(home, ".bashrc")
	if _, err := os.Stat(rc); err == nil {
		return rc, nil
	}
	return filepath.Join(home, ".bash_profile"), nil
}

// Inject writes the theme's prompt, banner, and color exports into the
// bash profile, replacing any previous cmdx injection first.
func (b *Bash) Inject(t *config.Theme) error {
	path, err := b.ProfilePath()
	if err != nil {
		return err
	}

	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	}

	existing = removeInjection(existing)
	block := buildBashBlock(t)
	final := existing + "\n" + block + "\n"

	return os.WriteFile(path, []byte(final), 0644)
}

// Remove strips any cmdx-injected block from the bash profile.
func (b *Bash) Remove() error {
	path, err := b.ProfilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read bash profile: %w", err)
	}

	cleaned := removeInjection(string(data))
	return os.WriteFile(path, []byte(cleaned), 0644)
}

// IsInjected reports whether a cmdx theme block is currently present in
// the bash profile. A missing or unreadable profile is treated as
// "not injected" rather than an error.
func (b *Bash) IsInjected() (bool, error) {
	path, err := b.ProfilePath()
	if err != nil {
		return false, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}

	return strings.Contains(string(data), shells.InjectStart), nil
}

// buildBashBlock renders the full injected bash config block for a theme:
// PS1 prompt, git branch helper, startup banner, and color exports.
//
// All free-text theme fields (name, author, description, banner text,
// prompt symbol) pass through shells.SanitizeForShell first, since theme
// JSON may come from an untrusted source (the community registry or a
// shared file) and is otherwise embedded directly into this script.
func buildBashBlock(t *config.Theme) string {
	primary := hexToRgbAnsi(t.Colors.Primary)
	accent := hexToRgbAnsi(t.Colors.Accent)
	muted := hexToRgbAnsi(t.Colors.Muted)
	symbol := shells.SanitizeForShell(t.Prompt.Symbol)
	name := shells.SanitizeForShell(t.Meta.Name)
	author := shells.SanitizeForShell(t.Meta.Author)
	bannerText := shells.SanitizeForShell(t.Banner.Text)
	description := shells.SanitizeForShell(t.Meta.Description)

	var promptLine string
	if t.Prompt.Style == "multiline" {
		promptLine = fmt.Sprintf(
			`PS1='\[\033[38;2;%sm\]\w\[\033[0m\] \[\033[38;2;%sm\]$(git_branch)\[\033[0m\]\n\[\033[38;2;%sm\]%s\[\033[0m\] '`,
			primary, muted, accent, symbol,
		)
	} else {
		promptLine = fmt.Sprintf(
			`PS1='\[\033[38;2;%sm\]\w\[\033[0m\] \[\033[38;2;%sm\]$(git_branch)\[\033[0m\] \[\033[38;2;%sm\]%s\[\033[0m\] '`,
			primary, muted, accent, symbol,
		)
	}

	return fmt.Sprintf(`%s
# Theme: %s by %s

# Git branch helper
git_branch() {
    git branch 2>/dev/null | grep '^\*' | sed 's/\* /(/;s/$/)/'
}

# Prompt
%s

# Startup Banner
echo ""
echo -e "  \033[38;2;%sm%s $USER\033[0m"
echo -e "  \033[38;2;%sm%s\033[0m"
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
