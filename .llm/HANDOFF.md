# cmdX — Session Handoff
> Last updated: After VS Code extension, web theme builder, and main README overhaul

---

## What is cmdX
A Go-based terminal theming framework — spinners, prompts, progress bars, banners, cursors, borders, PNG assets, reactive mascots, composable status bars — via a single JSON theme file. Cross-platform: Windows (PowerShell/CMD), macOS (zsh), Linux (bash/zsh/fish).

- **Repo:** github.com/abhigyanwebber/cmdX
- **Themes repo:** github.com/abhigyanwebber/cmdX-themes
- **Go module:** github.com/abhigyanwebber/cmd-customizer
- **Local path:** C:\Users\ASUS\Desktop\cmdX

---

## What's Complete

### Library Integrations
go-colorful, termenv, viper, huh, glamour, harmonica, bubbles extended (list/spinner/progress/viewport). tview evaluated and skipped.

### CLI Commands
- `cmdx theme list/apply/info/validate/preview/inject/remove/create`
- `cmdx registry list/browse/fetch/search`
- `cmdx plugin list/info/validate/spinners`
- `cmdx asset list/import/info/preview/validate/remove/use/status/chafa/render/create`
- `cmdx asset mascot-state/mascot-info/mascot-hooks`
- `cmdx asset statusbar-preview/statusbar-hooks/statusbar-info`
- `cmdx font list/search/info/install/remove/installed/path`
- `cmdx wallpaper set/remove/info`
- `cmdx edit [theme]`, `cmdx config show/set/reset`

### Test Suite
Full coverage across config, theme, graphics, assets (incl. floater/mascot/statusbar/chafa), shells (incl. sanitize), registry, plugin, fonts. `cmd` package has helpers_test.go for the CMDX_THEMES_DIR/CMDX_ASSETS_DIR env var overrides (added during VS Code extension work — **not yet re-verified with a full `go test ./...` run since being added, verify this on next session start**).

### Code Cleanup & Security Audit
resolveColor centralized, preview.go split, cmd/helpers.go extracted, doc comments everywhere, lint pass against .golangci.yml. Security: shell injection sanitization (SanitizeForShell), registry path traversal fix (ValidThemeName + size caps), chafa argument injection fix (validateImagePath + `--` separator), install script hardening (set -euo pipefail), profile path detection audited safe.

### Exotic Assets (all 3 complete)
- **Floaters** — fixed 4 integration gaps (preview switch, status display, info block, theme-preview display)
- **Mascots** — full state machine, 7 trigger types (exit_code/output_regex/env_var/idle_time/git_status/command/always), priority resolution, per-state render overrides with tinting, shell hooks for bash/zsh/powershell
- **Status bar** — 16 segment types, 3 zones, 8 separator styles (incl. Powerline), pure shell code generation, no binary call at prompt time

### Font Installer (complete)
10-font curated Nerd Fonts catalog, safe zip extraction (zip-slip protection matching cross-platform absolute-path edge cases, size/count caps), cross-platform install dirs, Windows registry registration (build-tagged), install state tracking, `--url` escape hatch for any font. Two real Windows bugs found and fixed during testing: (1) `filepath.IsAbs` doesn't catch POSIX-style absolute paths on Windows, (2) Windows Defender briefly locks newly-written .ttf files during real-time scan — fixed with temp-file + atomic-rename + exponential-backoff retry (up to ~13s worst case).

### VS Code Extension (vscode-extension/, complete)
JSON schemas for theme.json and asset.json matching the Go structs exactly. Webview theme preview panel. Sidebar Themes + Assets tree views. Commands wrap the real CLI (apply/validate via child_process + output channel; theme-create and asset-preview open a real integrated terminal since chafa/sixel can't be faithfully reproduced in a webview). Added CMDX_THEMES_DIR/CMDX_ASSETS_DIR env var support to the Go CLI so the extension's directory settings actually work (didn't exist before). Verified: `tsc` 0 errors, `eslint` 0 warnings (Node 22 in sandbox).

### Web Theme Builder (web-builder/, complete)
Vite + vanilla TypeScript static site. Full theme.json schema as editable controls (not a simplified subset). Live terminal mockup with real working animations for all 5 banner effects (glitch/rainbow/pulse/neon/typewriter) and real loader frame timing. Client-side validation mirroring internal/config/validator.go. 4 presets (cyberpunk/ocean/minimal/forest). Import/export/clipboard/shareable-URL (theme base64-encoded in URL hash, no backend). Verified: `npm install`, `tsc --noEmit`, `npm run build` all clean, dev server confirmed serving correctly (Node 22 in sandbox).

### Main README.md
Fully rewritten — was stale (missing all of the above, roadmap checkboxes wrong). Now covers exotic assets, font installer, both dev tools, updated project structure tree, corrected roadmap.

---

## What's Next

### ✅ ALL MAJOR FEATURES COMPLETE
- [x] Sound themes (internal/assets/sound.go, player.go) — the one asset
      type with no chafa involved at all: shells out to a platform audio
      player (afplay/paplay/aplay/ffplay/PowerShell SoundPlayer, with
      paplay→aplay→ffplay fallback chain on Linux). Reuses the mascot
      trigger vocabulary via a type-bridge to matchesTrigger (no
      duplicated matching logic). Per-effect volume/cooldown/async,
      cooldown persisted via .state/ files since each shell hook
      invocation is a fresh process. Custom "player" command template
      is the developer-freedom escape hatch for non-WAV formats. CLI:
      asset sound-play/sound-info/sound-hooks. VS Code schema + 
      REQUIREMENTS.md updated.

### 🔵 FUTURE (not part of the original 4 major features, optional)
- Nothing currently planned — all roadmap items are complete. Natural
  next directions if picked back up: community sound theme packs in
  cmdX-themes, a "sound theme create" preset library similar to the
  community theme registry, or extending the web builder to cover
  asset.json (currently theme.json only).

### Before push
- [ ] Run full `go test ./...` — this covers both the CMDX_THEMES_DIR/CMDX_ASSETS_DIR changes from earlier in the session (cmd/helpers.go, cmd/asset.go, cmd/helpers_test.go) and the brand-new sound theme code (internal/assets/sound.go, player.go, sound_test.go, player_test.go). None of this has been run through the real Go toolchain yet — only manual brace-balance and duplicate-symbol checks in the sandbox (no `go` binary available here).
- [ ] Sound test note: player_test.go's tests only verify command *construction* (args, paths), never actually invoke a real audio player — that's deliberate (CI/sandbox environments won't reliably have paplay/aplay/ffplay installed), so a passing test suite doesn't guarantee audio actually plays on your machine. Worth a manual `cmdx asset preview <sound-asset> --state success` smoke test after import.

---

## Session Rules
1. Always: `git push origin HEAD:main`
2. Always: `go mod tidy && go mod download` after extracting zip (go.sum goes stale every time go.mod changes in the sandbox, since there's no Go toolchain here to run tidy directly)
3. Stop at 80% context, write handoff
4. Karpathy 4 rules: Think Before Coding, Simplicity First, Surgical Changes, Goal-Driven Execution
5. Never duplicate resolveColor() across packages — centralized in internal/config/colors.go
6. Every new Go package needs a _test.go file
7. `go test` runs `go vet` automatically — vet warnings fail the whole package's tests, not just show as warnings
8. TypeScript projects (vscode-extension/, web-builder/) — sandbox HAS Node 22 + npm, so verify with real `tsc`/`eslint`/`vite build` rather than guessing

## Known Patterns
- Chafa tests need real PNGs from repo assets, not synthetic bytes — see findRealPNG() in internal/assets/floater_test.go
- Windows file writes to fresh files can hit AV-scan locks — write-to-temp + atomic-rename + retry pattern, see internal/fonts/zip.go's extractOneFile
- filepath.IsAbs behaves differently per-OS — always test path-validation logic against POSIX-style paths explicitly, don't rely on IsAbs alone (see internal/fonts/zip.go's sanitizeExtractPath)
- Branch naming: always `git push origin HEAD:main` (master/main conflict across the two dev machines)
