# cmdX — Terminal Theming Framework

A Go-based CLI tool giving developers full control over terminal visual primitives
(spinners, prompts, progress bars, banners, cursors, borders, PNG assets) via JSON themes.
Cross-platform: Windows (PowerShell/CMD), macOS (zsh), Linux (bash/zsh/fish).

## Build & Run

```bash
go build -o cmdx.exe ./cmd/
go test ./...
go mod tidy
git push origin HEAD:main
```

## Project Structure

cmdX/
├── cmd/
│   ├── main.go
│   ├── root.go
│   ├── theme.go
│   ├── plugin.go
│   ├── registry.go
│   ├── edit.go
│   ├── wallpaper.go
│   └── asset.go
├── internal/
│   ├── config/
│   ├── theme/
│   ├── primitives/
│   ├── preview/
│   ├── graphics/
│   ├── shells/
│   │   ├── powershell/
│   │   ├── zsh/
│   │   └── bash/
│   ├── plugin/
│   ├── editor/
│   ├── assets/
│   ├── registry/
│   ├── wallpaper/
│   └── icons/
├── themes/
├── assets/
│   ├── spinners/
│   ├── icons/
│   ├── banners/
│   └── dividers/
├── plugins/
├── scripts/
└── docs/

## Conventions

- Standard Go project layout (cmd/, internal/)
- Table-driven tests for all packages
- Every exported function must have a Go doc comment
- Wrap errors with fmt.Errorf("context: %w", err)
- Color keys: primary, secondary, background, foreground, accent, error, success, warning, muted
- Asset types: spinner, icon, banner, divider, floater, mascot, status-bar
- IMPORTANT: Run go fmt and go vet before committing
- IMPORTANT: Never commit cmdx.exe or PNG frame files
- IMPORTANT: Always git push origin HEAD:main
- IMPORTANT: Never duplicate resolveColor() — centralize it
- IMPORTANT: Every new package needs a _test.go file

## Key Dependencies

- github.com/spf13/cobra — CLI framework
- github.com/charmbracelet/bubbletea — TUI framework
- github.com/charmbracelet/lipgloss — Terminal styling
- github.com/charmbracelet/bubbles — UI components
- github.com/charmbracelet/huh — Interactive forms
- github.com/charmbracelet/glamour — Rich markdown output
- github.com/charmbracelet/harmonica — Spring animations
- github.com/lucasb-eyer/go-colorful — LAB color space gradients
- github.com/muesli/termenv — Terminal capability detection
- github.com/spf13/viper — Global user config
- External: chafa (winget install hpjansson.chafa)

## Active Priorities

1. Foundation files — CLAUDE.md, AGENTS.md, .claudeignore, .llm/todo.md
2. Library upgrades — go-colorful, termenv, huh, glamour, harmonica, viper
3. Tests — zero coverage currently
4. Community themes — 10 new palettes for cmdX-themes repo
5. Code cleanup — centralize color resolution, refactor preview.go
6. Security audit — shell injection, registry downloads, chafa path
7. Exotic assets — floaters, mascots, status-bar
8. Major features — web builder, VS Code extension