package bash

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/shells"
)

type Bash struct{}

func New() *Bash {
	return &Bash{}
}

func (b *Bash) Name() string {
	return "bash"
}

func (b *Bash) Detect() bool {
	shell := os.Getenv("SHELL")
	return strings.Contains(shell, "bash")
}

func (b *Bash) ProfilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}

	// prefer .bashrc, fallback to .bash_profile
	rc := filepath.Join(home, ".bashrc")
	if _, err := os.Stat(rc); err == nil {
		return rc, nil
	}
	return filepath.Join(home, ".bash_profile"), nil
}

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

func buildBashBlock(t *config.Theme) string {
	primary := hexToRgbAnsi(t.Colors.Primary)
	accent := hexToRgbAnsi(t.Colors.Accent)
	muted := hexToRgbAnsi(t.Colors.Muted)
	symbol := t.Prompt.Symbol

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
		t.Meta.Name, t.Meta.Author,
		promptLine,
		hexToRgbAnsi(t.Colors.Primary), t.Banner.Text,
		hexToRgbAnsi(t.Colors.Muted), t.Meta.Description,
		t.Colors.Primary,
		t.Colors.Accent,
		t.Colors.Success,
		t.Colors.Muted,
		shells.InjectEnd,
	)
}

func hexToRgbAnsi(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "255;255;255"
	}
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("%d;%d;%d", r, g, b)
}

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
