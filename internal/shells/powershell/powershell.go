package powershell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/shells"
)

type PowerShell struct{}

func New() *PowerShell {
	return &PowerShell{}
}

func (p *PowerShell) Name() string {
	return "powershell"
}

func (p *PowerShell) Detect() bool {
	_, exists := os.LookupEnv("PSModulePath")
	return exists
}

func (p *PowerShell) ProfilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil
}

func (p *PowerShell) Inject(t *config.Theme) error {
	path, err := p.ProfilePath()
	if err != nil {
		return err
	}

	// create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("could not create profile directory: %w", err)
	}

	// read existing profile if it exists
	existing := ""
	if data, err := os.ReadFile(path); err == nil {
		existing = string(data)
	}

	// remove any previous cmdx injection
	existing = removeInjection(existing)

	// build the injection block
	block := buildPowerShellBlock(t)

	// append to profile
	final := existing + "\n" + block + "\n"
	if err := os.WriteFile(path, []byte(final), 0644); err != nil {
		return fmt.Errorf("could not write profile: %w", err)
	}

	return nil
}

func (p *PowerShell) Remove() error {
	path, err := p.ProfilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read profile: %w", err)
	}

	cleaned := removeInjection(string(data))
	return os.WriteFile(path, []byte(cleaned), 0644)
}

func (p *PowerShell) IsInjected() (bool, error) {
	path, err := p.ProfilePath()
	if err != nil {
		return false, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}

	return strings.Contains(string(data), shells.InjectStart), nil
}

func buildPowerShellBlock(t *config.Theme) string {
	primary := t.Colors.Primary
	accent := t.Colors.Accent
	success := t.Colors.Success
	errorCol := t.Colors.Error
	muted := t.Colors.Muted
	symbol := t.Prompt.Symbol

	return fmt.Sprintf(`%s
# Theme: %s by %s

# Colors
$cmdx_primary  = "%s"
$cmdx_accent   = "%s"
$cmdx_success  = "%s"
$cmdx_error    = "%s"
$cmdx_muted    = "%s"

# Custom Prompt
function prompt {
    $dir = Split-Path -Leaf (Get-Location)
    $branch = ""

    if (Test-Path .git) {
        $branch = " (" + (git branch --show-current 2>$null) + ")"
    }

    Write-Host $dir -ForegroundColor Cyan -NoNewline
    Write-Host $branch -ForegroundColor DarkGray -NoNewline
    Write-Host " %s " -ForegroundColor Magenta -NoNewline
    return " "
}

# Startup Banner
Write-Host ""
Write-Host "  %s $env:USERNAME" -ForegroundColor Cyan
Write-Host "  %s" -ForegroundColor DarkGray
Write-Host ""
%s`,
		shells.InjectStart,
		t.Meta.Name, t.Meta.Author,
		primary, accent, success, errorCol, muted,
		symbol,
		"SYSTEM ONLINE //",
		t.Meta.Description,
		shells.InjectEnd,
	)
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
