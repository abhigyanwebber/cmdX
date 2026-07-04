# cmdX — Session Handoff
> Last updated: After floater integration, developer freedom expansion, security audit

---

## What is cmdX
A Go-based terminal theming framework — spinners, prompts, progress bars, banners, cursors, borders, PNG assets — via a single JSON theme file. Cross-platform: Windows (PowerShell/CMD), macOS (zsh), Linux (bash/zsh/fish).

- **Repo:** github.com/abhigyanwebber/cmdX
- **Themes repo:** github.com/abhigyanwebber/cmdX-themes
- **Go module:** github.com/abhigyanwebber/cmd-customizer
- **Local path:** C:\Users\ASUS\Desktop\cmdX

---

## What's Complete

### Library Integrations
- go-colorful (LAB gradients), termenv (capability detection), viper (global config)
- huh (interactive theme creator), glamour (markdown output), harmonica (spring animations)
- bubbles extended: list, spinner, progress, viewport — all wired into CLI
- tview: evaluated, skipped (Bubble Tea editor is solid, no net gain)

### CLI Commands
- `cmdx theme list/apply/info/validate/preview/inject/remove/create`
- `cmdx registry list/browse/fetch/search`
- `cmdx plugin list/info/validate/spinners`
- `cmdx asset list/import/info/preview/validate/remove/use/status/chafa/render/create`
- `cmdx wallpaper set/remove/info`
- `cmdx edit [theme]` — TUI editor with harmonica spring animations (60fps)
- `cmdx config show/set/reset`

### Test Suite (full coverage)
- internal/config: validator_test.go, loader_test.go, colors_test.go
- internal/theme: manager_test.go
- internal/graphics: gradient_test.go, effects_test.go
- internal/assets: loader_test.go, chafa_test.go, floater_test.go
- internal/shells: sanitize_test.go, powershell/powershell_test.go
- internal/registry: registry_test.go
- internal/plugin: manager_test.go
- Real bugs found in testing: registry.Search() case-sensitive (fixed)

### Code Cleanup
- resolveColor() centralized in internal/config/colors.go (was 4 duplicates)
- preview.go split: preview.go (orchestrator) + sections.go + colors.go
- cmd/helpers.go extracted: detectShell(), getThemesDir(), loadThemeOrExit()
- Doc comments on every exported type/function across entire codebase
- Manual lint pass against .golangci.yml: fixed noctx (HTTP context), errcheck (Sscanf)

### Security Audit (all items closed)
- Shell injection: SanitizeForShell() applied to all theme text fields in bash/zsh/powershell injection
- Registry path traversal: ValidThemeName() allowlist + 1MiB download cap + JSON validity gate
- Chafa argument injection: validateImagePath() + "--" separator
- settings.json write: already safe (backup-before-write), confirmed
- Install scripts: set -euo pipefail added, confirmed no injection vectors
- Profile path detection: confirmed safe by design (hardcoded filenames + os.UserHomeDir)

### Floaters (complete)
- Type system, validation, CLI plumbing already scaffolded
- Fixed 4 integration gaps: asset preview switch, asset status corners, asset info detail block, theme preview showing active floaters
- floater_test.go: 8 tests, chafa tests use real PNG from repo assets

### Developer Freedom Expansion
- internal/assets/override.go: RenderOverrides struct + ApplyOverrides() + OverridesFromFlags()
- Every Manager method has a WithOverrides variant
- cmd/asset_create.go: interactive scaffolding wizard via huh (cmdx asset create)
- cmdx asset render <path>: raw PNG→terminal, no manifest needed, all render flags exposed
- cmdx asset preview: 9 override flags (--mode, --color, --symbols, --width, --height, --dither, --stretch, --threshold, --position)
- REQUIREMENTS.md: Go/git/chafa install commands per platform, render mode matrix, sixel terminal support table
- Makefile: make deps/build/test/vet/lint/clean/install

### Community Themes (cmdX-themes repo)
13 themes: synthwave, ocean, forest, ocean-depths, sunset-boulevard, forest-canopy, modern-minimalist, golden-hour, arctic-frost, desert-rose, tech-innovation, botanical-garden, midnight-galaxy

---

## What's Next (Priority Order)

### 🔵 NOW — Exotic Assets (remaining)
- [ ] Mascots — reactive character (idle/working/success/error/warning/sleeping)
  - Different PNG per state in same asset folder
  - asset.json type: mascot
  - Shell injection triggers state changes based on exit codes
  - cmdx asset use my-mascot --as mascot
- [ ] Status bar — bottom bar with git status, time, directory
  - asset.json type: status-bar
  - Real-time updates via shell hooks
  - Configurable segments

### 🔵 FUTURE — Major Features
- [ ] Web theme builder
- [ ] VS Code extension
- [ ] Sound themes
- [ ] Font installer

---

## Session Rules
1. Always: git push origin HEAD:main
2. Always: go mod tidy && go mod download after extracting zip
3. Stop at 80% context, write handoff
4. Karpathy 4 rules: Think Before Coding, Simplicity First, Surgical Changes, Goal-Driven Execution
5. Never duplicate resolveColor() across packages
6. Every new package needs a _test.go file
7. Verify build: go build -o cmdx.exe ./cmd/
8. go test uses go vet automatically — fix vet warnings or tests appear to fail

## Known Patterns
- go.sum gets stale when I add go.mod entries (can't run go mod tidy in sandbox)
  → always run `go mod tidy && go mod download` on your machine after unzipping
- Chafa tests need real PNGs from repo assets, not synthetic bytes
  → findRealPNG() helper in floater_test.go walks to ../../assets/
- Branch naming: always git push origin HEAD:main (master/main conflict)
