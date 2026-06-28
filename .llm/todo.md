# cmdX — Active Todo List
> Paste this at the start of every session for instant context

## Current Status
- Last completed: PNG asset system (spinners, icons, banners, dividers with sixel rendering)
- Build status: Clean
- Branch: main
- Machines: user machine (cmd-customizer folder) + ASUS machine (cmdX folder)

## Priority Queue (in order)

### 🔴 NOW — Foundation (current session)
- [x] CLAUDE.md
- [x] AGENTS.md  
- [x] .claudeignore
- [x] .llm/todo.md
- [ ] .llm/HANDOFF_TEMPLATE.md
- [ ] .golangci.yml
- [ ] internal/shells/.CLAUDE.md
- [ ] themes/.CLAUDE.md
- [ ] internal/assets/.CLAUDE.md
- [ ] internal/graphics/.CLAUDE.md
- [ ] internal/editor/.CLAUDE.md
- [ ] internal/config/.CLAUDE.md

### 🟠 NEXT — Library Upgrades
- [ ] go get github.com/lucasb-eyer/go-colorful
- [ ] go get github.com/muesli/termenv
- [ ] go get github.com/charmbracelet/huh
- [ ] go get github.com/charmbracelet/glamour
- [ ] go get github.com/charmbracelet/harmonica
- [ ] go get github.com/spf13/viper
- [ ] Integrate go-colorful into internal/graphics/gradient.go
- [ ] Integrate termenv into internal/assets/chafa.go
- [ ] Build cmdx theme create with huh forms
- [ ] Upgrade cmdx theme info with glamour
- [ ] Add harmonica to TUI editor transitions
- [ ] Add viper global config (~/.cmdx/config.yaml)
- [ ] Evaluate tview for editor v2

### 🟡 THEN — Tests
- [ ] internal/config/validator_test.go
- [ ] internal/config/loader_test.go
- [ ] internal/theme/manager_test.go
- [ ] internal/graphics/gradient_test.go
- [ ] internal/graphics/effects_test.go
- [ ] internal/assets/loader_test.go
- [ ] internal/shells/powershell/powershell_test.go
- [ ] internal/registry/registry_test.go
- [ ] internal/plugin/manager_test.go

### 🟡 THEN — Community Themes (cmdX-themes repo)
- [ ] ocean-depths
- [ ] sunset-boulevard
- [ ] forest-canopy
- [ ] modern-minimalist
- [ ] golden-hour
- [ ] arctic-frost
- [ ] desert-rose
- [ ] tech-innovation
- [ ] botanical-garden
- [ ] midnight-galaxy

### 🟡 THEN — Code Cleanup
- [ ] Centralize resolveColor() — remove duplicates from renderer.go, utils.go, preview.go
- [ ] Refactor preview.go (300+ lines)
- [ ] Reorganize theme.go
- [ ] Add Go doc comments to all exported functions
- [ ] Run golangci-lint and fix warnings

### 🟢 THEN — Security Audit
- [ ] Shell injection markers validation
- [ ] Registry download validation
- [ ] Chafa path validation
- [ ] settings.json write safety

### 🔵 THEN — Exotic Assets
- [ ] Floaters (corner decorative graphics)
- [ ] Mascots (reactive character with states)
- [ ] Status bar (bottom info bar with PNG graphics)

### 🔵 FUTURE — Major Features
- [ ] Web theme builder
- [ ] VS Code extension
- [ ] Sound themes
- [ ] Font installer

## Session Rules
- Always: git push origin HEAD:main
- Always: stop at 80% context
- Always: write handoff doc at end of session
- Always: apply Karpathy 4 rules