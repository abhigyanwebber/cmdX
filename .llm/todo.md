# cmdX — Active Todo List
> Paste this at the start of every session for instant context

## Current Status
- Last completed: bubbles extended components integration
- Build status: needs `go mod tidy` after adding tui package
- Branch: main
- Machines: user machine (cmd-customizer folder) + ASUS machine (cmdX folder)

## Priority Queue (in order)

### 🔴 COMPLETED — Foundation
- [x] CLAUDE.md, AGENTS.md, .claudeignore, .golangci.yml
- [x] .llm/todo.md, .llm/HANDOFF_TEMPLATE.md
- [x] Nested .CLAUDE.md files in all internal packages

### 🔴 COMPLETED — Library Integrations
- [x] go-colorful — LAB gradients in graphics/gradient.go
- [x] termenv — capability detection in assets/capabilities.go
- [x] viper — global config in internal/config/global.go
- [x] huh — interactive theme creator (cmdx theme create)
- [x] glamour — rich markdown output (theme info, registry list/search/fetch)
- [x] harmonica — spring animations in TUI editor (cursor + panel slide-in)
- [x] bubbles extended:
  - [x] bubbles/list → internal/tui/list.go (theme picker + registry browser)
  - [x] bubbles/spinner → internal/tui/spinner.go (loading states)
  - [x] bubbles/progress → internal/tui/progress.go (download progress)
  - [x] bubbles/viewport → internal/tui/viewport.go (scrollable panels)
- [x] tview — EVALUATED: skip. Bubble Tea editor is solid; tview rewrite offers
      no net gain. Verdict: stay with current stack.

### 🔴 COMPLETED — Tests (zero coverage — critical gap)
- [x] internal/config/validator_test.go
- [x] internal/config/loader_test.go
- [x] internal/config/colors_test.go (added during cleanup, centralized resolveColor)
- [x] internal/theme/manager_test.go
- [x] internal/graphics/gradient_test.go
- [x] internal/graphics/effects_test.go
- [x] internal/assets/loader_test.go
- [x] internal/assets/chafa_test.go (added during security audit)
- [x] internal/shells/powershell/powershell_test.go
- [x] internal/shells/sanitize_test.go (added during security audit)
- [x] internal/registry/registry_test.go
- [x] internal/plugin/manager_test.go
- Real bug caught: registry.Search() was case-sensitive — fixed,
  hand-rolled contains() replaced with stdlib strings.Contains

### 🔴 COMPLETED — Code Cleanup
- [x] Centralized resolveColor() into internal/config/colors.go —
      theme/renderer.go, theme/manager.go, primitives/utils.go, and
      preview/preview.go all delegate to it now
- [x] Refactored preview.go (323 lines) into preview.go (129, orchestrator
      only) + sections.go (219, section renderers) + colors.go (9)
- [x] Reorganized theme.go (291 → 231 lines) — pulled detectShell(),
      getThemesDir(), and new loadThemeOrExit() into cmd/helpers.go
      since they're shared across create.go, edit.go, registry.go
- [x] Doc comments added to every exported type/function across the
      codebase, including the full theme JSON schema in config/types.go
- [x] Manual lint-aware pass against .golangci.yml's enabled linters —
      fixed noctx (registry HTTP calls now use context) and errcheck
      (unchecked fmt.Sscanf in zsh/bash hex parsing)

### 🟢 NOW — Security Audit
- [x] Shell injection markers validation — added shells.SanitizeForShell(),
      applied to every theme text field (name, author, description,
      banner text, prompt symbol) embedded into bash/zsh/PowerShell
      profile injection blocks. Strips backticks, $, quotes, semicolons,
      pipes, and the injection markers themselves.
- [x] Registry download validation — FetchTheme now validates theme
      names against a strict allowlist (alphanumeric/hyphen/underscore
      only) before using them in a URL or filesystem path, closing a
      path traversal vulnerability (e.g. "../../.bashrc" could
      previously overwrite arbitrary files). Downloads are capped at
      1 MiB and validated as JSON before being written to disk.
- [x] Chafa path validation — Render() now rejects paths starting with
      "-" (flag injection) and verifies the path is a real, accessible
      regular file before shelling out. Added "--" separator in
      buildChafaArgs as defense in depth.
- [x] settings.json write safety — already correct (backup-before-write
      at internal/wallpaper/windows_terminal.go), no changes needed.
- [x] Install scripts — audited scripts/install.sh and install.ps1.
      Neither processes untrusted/network-derived input (REPO and all
      paths are hardcoded), so no injection vector exists. Hardened
      install.sh with `set -euo pipefail` (was just `set -e`) so an
      unset variable fails loudly instead of silently expanding empty.
- [x] Profile path detection — audited ProfilePath() in bash.go, zsh.go,
      powershell.go. All three build the path from os.UserHomeDir()
      (OS-provided, not attacker-influenced) joined with a hardcoded
      literal filename (".bashrc", ".zshrc", profile.ps1). No user or
      theme-derived input reaches these paths, so no traversal vector
      exists — confirmed safe, no code change needed.

### 🔵 NOW — Exotic Assets
- [x] Floaters (corner decorative graphics) — type system, validation,
      and CLI plumbing (asset.json type:floater, asset use --as floater
      --position) were already scaffolded from an earlier session, but
      had real integration gaps, now closed:
      - cmd/asset.go assetPreviewCmd had no case for AssetTypeFloater in
        its render switch, so "cmdx asset preview <floater>" silently
        did nothing after printing the header — fixed
      - assetStatusCmd's slots list omitted floaters entirely (they're
        tracked per-corner, not as a single slot) — fixed, now shows
        all four corner states under "Active Floaters"
      - assetInfoCmd had no floater-specific detail block (file,
        position, max size, margin, animated/static) — added
      - "cmdx theme preview" never displayed active floaters alongside
        the rest of the theme — added showActiveFloaters() in
        cmd/helpers.go, called after preview.Run()
      - Added internal/assets/floater_test.go (position validation,
        PreviewFloater, RenderFloaterFrames) and floater coverage in
        loader_test.go (manifest validation: missing config, missing
        file, invalid position, animated frame validation)
- [ ] Mascots (reactive character with states)
- [ ] Status bar (bottom info bar)

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
- Never: duplicate resolveColor() across packages
- Every new package needs a _test.go file
- Verify build: go build -o cmdx.exe ./cmd/
