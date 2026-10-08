# cmdX — Project Bible

Authoritative description of what cmdX actually is, derived from source inspection and
live functional testing (not from README claims alone). Last verified: this repair session.

## What cmdX Is

A cross-platform (Windows/macOS/Linux) Go CLI + TUI that turns a single JSON "theme" file
into a full terminal aesthetic — colors, prompt, spinners, progress bars, banners, cursors,
borders, PNG-rendered assets (via `chafa`), reactive "mascots" driven by shell-event
triggers, composable status bars, and sound effects.

Primary user: developers who want a fast, ricer-style terminal setup without hand-editing
shell profiles. Core abstraction: the `Asset` (typed, named, PNG/behavior bundle) and the
`Theme` (JSON document that can reference assets by name).

Module path: `github.com/abhigyanwebber/cmd-customizer` (historical name; repo/product is
branded `cmdX` — never renamed).

## Major Subsystems

- **Theme system** (`internal/config`, `internal/theme`, `cmd/theme.go`): schema, validation,
  centralized `resolveColor()`, render pipeline.
- **Asset system** (`internal/assets`): the `Asset` interface; spinner, icon, banner,
  divider, floater, mascot, status_bar, sound types; chafa-based PNG rendering
  (ascii/blocks/braille/sixel/truecolor modes).
- **Reactive mascot system**: per-state PNG frames selected by a trigger-priority resolver
  (`ResolveState` in `internal/assets/mascot.go`). Trigger types: `always`, `exit_code`,
  `command`, `output_regex`, `env_var`, `git_status`, `idle_time`. See ARCHITECTURE.md for
  exact resolution semantics and a documented defect in the CLI's forced-state override.
- **Sound system**: theme-linked WAV effects, triggered by the same trigger-type vocabulary
  as mascots (shared trigger schema, separate implementation — see ARCHITECTURE.md).
- **Shell integration** (`internal/shells`): per-shell (bash/zsh/powershell) injection via
  marker-delimited blocks. Two independent marker pairs: `InjectStart/InjectEnd` (theme) and
  `AssetHooksStart/AssetHooksEnd` (mascot/sound reactive hooks) — kept separate so switching
  themes never orphans or clobbers hooks and vice versa. Verified idempotent.
- **Theme → asset auto-activation** (currently uncommitted — see KNOWN_ISSUES.md BUG-001):
  linking assets in a theme's `assets` block activates them on `theme apply`/`inject`,
  with `--no-assets`/`--no-hooks` flags to suppress. Verified working once BUG-001 is fixed.
- **Registry** (`internal/registry`): fetches community themes from GitHub. Name validation
  (path-traversal-safe allowlist) is implemented and correct in source, but is unreachable
  through the CLI in non-fully-interactive contexts — see BUG-005.
- **Plugins** (`internal/plugin`): local plugin discovery/listing; example plugin bundled.
- **Fonts** (`internal/fonts`): curated font catalog, list/search/install.
- **Wallpaper** (`internal/wallpaper`): Windows Terminal / iTerm2 / Kitty background image
  support.
- **TUI** (`internal/tui`): Bubble Tea theme editor and theme-list browser. Never
  successfully driven end-to-end in this investigation — see UNVERIFIED list.
- **Web builder** (`web-builder/`, TypeScript/Vite): visual theme builder. Builds/typechecks
  clean but its schema is stale relative to the Go core — see GAP-001.
- **VS Code extension** (`vscode-extension/`, TypeScript): schema + editor support. Schema
  is current and matches the Go structs exactly, including mascot/floater/status_bar/sound.

## What Differentiates cmdX
Assets as a first-class composable abstraction (usable standalone or bound into a theme);
marker-delimited shell injection that persists across restarts without manually edited
profiles; a community registry; outward-facing tooling (VS Code extension, web builder)
that wraps the same JSON schema rather than reimplementing it.
