# cmdX — Session Handoff Document
> Save this as .llm/HANDOFF.md in the project root

---

## What is cmdX
A Go-based terminal theming framework that gives developers full control over terminal visual primitives — spinners, prompts, progress bars, banners, cursors, borders, PNG assets — via a single JSON theme file. Cross-platform: Windows (PowerShell/CMD), macOS (zsh), Linux (bash/zsh/fish).

- **Repo:** github.com/abhigyanwebber/cmdX
- **Themes repo:** github.com/abhigyanwebber/cmdX-themes
- **Author:** abhigyanwebber
- **Language:** Go (92.4%)
- **Machine:** C:\Users\ASUS\Desktop\cmdX

---

## What's Been Built (Complete)

### Core Engine
- JSON theme system with full schema (colors, prompt, loader, progress bar, cursor, borders, banner, graphics, wallpaper, assets)
- Theme manager (load, apply, list, validate, preview)
- 4 built-in themes: default, cyberpunk, minimal, retro
- Live preview engine with gradients, effects, dividers, patterns

### CLI Commands
- `cmdx theme list/apply/info/validate/preview/inject/remove/create`
- `cmdx registry list/fetch/search`
- `cmdx plugin list/info/validate/spinners`
- `cmdx asset list/import/info/preview/validate/remove/use/status/chafa`
- `cmdx wallpaper set/remove/info`
- `cmdx edit [theme]` — TUI theme editor
- `cmdx config show/set/reset` — global user config

### Shell Injection
- PowerShell, zsh, bash
- Auto-detects running shell
- Uses InjectStart/InjectEnd markers
- Persists theme across terminal restarts

### Graphics Engine
- **go-colorful LAB gradients** (upgraded from basic RGB)
- Effects: glitch, rainbow, neon, typewriter, pulse
- Dividers: line, wave, dots, zigzag, stars, double
- Patterns: dots, grid, diagonal, bricks, circuit, stars

### PNG Asset System (via chafa)
- Asset types: spinner, icon, banner, divider
- Render modes: sixel, braille, blocks, ascii
- **termenv** capability detection — auto-selects best render mode
- `cmdx asset use --as [slot]` — composable, standalone usage
- State persisted in assets/.state/

### TUI Theme Editor
- Bubble Tea Model-View-Update
- 24 editable fields
- Live preview panel (colors, prompt, loader, progress, banner)
- Save to JSON with S key

### Interactive Theme Creator
- **huh forms** wizard — 6 steps
- Meta → Colors → Prompt → Loader → Cursor/Effects → Gradient
- Validates hex colors at each step
- Saves directly to themes/ folder

### Plugin API
- Manifest-based (plugin.json)
- Supports: spinners, banners, prompts
- `cmdx plugin list/info/validate/spinners`

### Community Registry
- Pulls from github.com/abhigyanwebber/cmdX-themes
- `cmdx registry list/fetch/search`
- 13 themes available

### Community Themes (cmdX-themes repo)
synthwave, ocean, forest, ocean-depths, sunset-boulevard, forest-canopy, modern-minimalist, golden-hour, arctic-frost, desert-rose, tech-innovation, botanical-garden, midnight-galaxy

### Wallpaper Support
- Windows Terminal (settings.json)
- iTerm2 (AppleScript)
- Kitty (kitty.conf)

### Global Config (viper)
- ~/.cmdx/config.yaml
- Keys: default_theme, render_mode, chafa_path, assets_dir, themes_dir, auto_inject, show_banner
- Env var overrides: CMDX_DEFAULT_THEME, CMDX_RENDER_MODE, CMDX_CHAFA_PATH

### Foundation Files
- CLAUDE.md, AGENTS.md, .claudeignore, .golangci.yml
- .llm/todo.md, .llm/HANDOFF_TEMPLATE.md
- Nested .CLAUDE.md in: internal/shells/, themes/, internal/assets/, internal/graphics/, internal/editor/, internal/config/

---

## Current Dependencies
```
github.com/spf13/cobra
github.com/charmbracelet/bubbletea
github.com/charmbracelet/lipgloss
github.com/charmbracelet/bubbles
github.com/charmbracelet/huh
github.com/charmbracelet/glamour      ← installed, not integrated yet
github.com/charmbracelet/harmonica    ← installed, not integrated yet
github.com/lucasb-eyer/go-colorful   ← integrated in graphics/gradient.go
github.com/muesli/termenv            ← integrated in assets/capabilities.go
github.com/spf13/viper               ← integrated in config/global.go
github.com/rivo/tview                ← installed, evaluation pending
External: chafa (winget install hpjansson.chafa)
```

---

## Project Structure
```
cmdX/
├── .llm/                   # Session continuity files
├── cmd/                    # CLI commands
│   ├── main.go, root.go, theme.go, plugin.go
│   ├── registry.go, edit.go, wallpaper.go
│   ├── asset.go, config.go, create.go
├── internal/
│   ├── config/             # types.go, loader.go, validator.go, global.go
│   ├── theme/              # manager.go, renderer.go
│   ├── primitives/         # spinner.go, progress.go, banner.go, utils.go
│   ├── preview/            # preview.go
│   ├── graphics/           # gradient.go, effects.go, divider.go, patterns.go
│   ├── shells/             # shell.go + powershell/, zsh/, bash/
│   ├── plugin/             # types.go, loader.go, manager.go
│   ├── editor/             # model.go, marshal.go
│   ├── assets/             # types.go, chafa.go, loader.go, manager.go, capabilities.go
│   ├── registry/           # registry.go
│   ├── wallpaper/          # wallpaper.go, windows_terminal.go, unix.go
│   └── icons/              # icons.go
├── themes/                 # Built-in theme JSON files
├── assets/                 # User PNG asset library
│   ├── spinners/pulse/     # Example spinner with 8 frames
│   ├── icons/neon-icons/
│   ├── banners/cyberpunk-banner/
│   └── dividers/neon-divider/
├── plugins/example-plugin/
├── scripts/                # install.sh, install.ps1
└── docs/
```

---

## What's Left — Priority Order

### 🟠 NEXT — Remaining Library Integration
- [ ] **glamour** — upgrade cmdx theme info, registry list to rich markdown output
- [ ] **harmonica** — spring animations in TUI editor cursor/panel transitions
- [ ] **bubbles extended** — list, table, progress, spinner, viewport components
- [ ] **tview** — evaluate for editor v2

### 🟡 THEN — Tests (zero coverage currently — critical gap)
- [ ] internal/config/validator_test.go
- [ ] internal/config/loader_test.go
- [ ] internal/theme/manager_test.go
- [ ] internal/graphics/gradient_test.go
- [ ] internal/graphics/effects_test.go
- [ ] internal/assets/loader_test.go
- [ ] internal/shells/powershell/powershell_test.go
- [ ] internal/registry/registry_test.go
- [ ] internal/plugin/manager_test.go

### 🟡 THEN — Code Cleanup
- [ ] Centralize resolveColor() — duplicated in theme/renderer.go, primitives/utils.go, preview/preview.go
- [ ] Refactor preview.go (300+ lines)
- [ ] Reorganize theme.go
- [ ] Add Go doc comments to all exported functions
- [ ] Run golangci-lint and fix all warnings

### 🟢 THEN — Security Audit
- [ ] Shell injection markers validation
- [ ] Registry download validation (checksum/size limits)
- [ ] Chafa path validation (prevent path injection)
- [ ] settings.json write safety (backup before write)

### 🔵 THEN — Exotic Assets
- [ ] **Floaters** — corner decorative PNG graphics
- [ ] **Mascots** — reactive character (idle/working/success/error states)
- [ ] **Status bar** — bottom bar with git status, time, directory

### 🔵 FUTURE — Major Features
- [ ] Web theme builder
- [ ] VS Code extension
- [ ] Sound themes
- [ ] Font installer

---

## Session Rules (always follow)
1. `git push origin HEAD:main` — always use this exact command
2. Stop at 80% context, write handoff before ending
3. Karpathy 4 rules: Think Before Coding, Simplicity First, Surgical Changes, Goal-Driven Execution
4. Paste .llm/todo.md at start of every session
5. Never duplicate resolveColor() across packages
6. Every new package needs a _test.go file
7. Run `go build -o cmdx.exe ./cmd/` to verify clean build

---

## Known Issues / Watch Out For
- Branch naming: always `git push origin HEAD:main` (master/main conflict)
- Two machines: user machine (C:\Users\user\Desktop\cmd-customizer) and ASUS machine (C:\Users\ASUS\Desktop\cmdX)
- Themes with spaces in names need quotes: `./cmdx.exe theme preview "blue sky"`
- Theme creator now validates no spaces in name
- PNG frame files should never be committed to git (.claudeignore handles this)
- go.sum should never be manually edited

---

*Last updated: After huh interactive theme creator integration*
*Next session starts with: glamour integration*