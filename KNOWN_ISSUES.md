# cmdX — Known Issues

Confirmed issues only. Status updated as repairs land in this session.

### BUG-001 — Build-breaking duplicate function
`cmd/helpers.go` contained `syncAssetHooks` declared twice, verbatim, in the uncommitted
working tree. `go build ./...` failed. **STATUS: FIXED this session** (see the repair
report for the exact diff).

### BUG-002 — Mascot `--state` override discarded before rendering
`cmd/asset_mascot.go`'s `mascot-state` command computed a resolved/forced state into a local
`state` variable, then discarded it (`_ = state`) and called `PreviewMascot(name, ctx, ...)`,
which independently re-resolves from `ctx` alone (no field exists on `MascotContext` to carry
a forced state). Confirmed: `--state error` always displayed whatever plain ctx-resolution
produced, never `error`. Scope: **CLI-only** — real shell hooks never pass `--state`, so this
did not affect the reactive-mascot experience for normal users, only manual/debug use of the
documented override flag. **STATUS: FIXED, with regression coverage** — `internal/assets/manager.go`
now has `PreviewMascotState(name, ctx, overrides, forcedState)`; `PreviewMascot` is a thin
compatibility wrapper. Both real call sites (`cmd/asset_mascot.go`'s `mascot-state`, and
`cmd/asset.go`'s generic `asset preview` — which had an independent, also-broken
`ctx.Env["CMDX_MASCOT_STATE"]` workaround) now thread the override through correctly.
Regression tests: `internal/assets/manager_mascot_test.go`
(`TestPreviewMascotState_ForcedStateOverridesResolution`,
`TestPreviewMascotState_EmptyForcedStateFallsBackToResolution`,
`TestPreviewMascot_CompatibilityWrapperStillResolvesNormally`) — all pass, exercised against
a real chafa render (skipped automatically if chafa isn't installed).

### GAP-003 (was BUG-003) — Mascot hook context is intentionally minimal, not incomplete
Reclassified from "bug" to "documented capability boundary" after evidence review (see
AGENT_HANDOFF.md for the full reasoning). Shipped `MascotShellHooks` (bash/zsh/powershell)
capture and pass only `--exit-code` and `--command`. The resolver additionally supports
`output_regex`, `env_var`, `git_status`, and `idle_time`.

Findings:
- **`env_var` was miscategorized — it isn't blocked at all.** `buildEnvSnapshot()` reads
  real ambient OS environment variables directly; it needs no hook wiring and already works
  today for anyone who exports a variable before running a command.
- **`git_status`, `idle_time`, `output_regex` are, on the best-available evidence,
  intentionally hook-light rather than unfinished:** `.llm/HANDOFF.md`'s own mascot summary
  lists "7 trigger types" and "shell hooks for bash/zsh/powershell" as two separate,
  coordinate claims — never asserting hooks populate all 7 contexts, in a document that is
  otherwise very specific about integration completeness elsewhere. `MascotContext`'s doc
  comment in `internal/assets/mascot.go` explicitly designs for partial context ("zero
  values are treated as 'not available' and triggers that depend on them won't fire").
  These three triggers are also categorically more expensive/invasive to capture per-prompt
  than exit_code/command (a subprocess call every prompt for git status; new
  cross-invocation persistent state for idle time; broader output-capture machinery that
  risks interfering with interactive programs for output regex) — consistent with a
  deliberate cost-conscious design, mirrored by the status-bar feature's explicit
  "no binary call at prompt time" design goal from the same commit.
- This is an inference from converging evidence, not a single explicit repository
  statement — if the project owner intended fuller hook coverage, this should be corrected
  rather than treated as settled.

**STATUS: Downgraded from bug to documented gap.** Not implemented into hooks. The
resolver's graceful-degradation contract (absent context → trigger doesn't fire, no error)
already has test coverage in `internal/assets/mascot_test.go`. A theme/mascot author who
wants `git_status`/`idle_time`/`output_regex` triggers to actually fire needs a custom hook
that populates those `mascot-state` flags itself — this is achievable today (the flags
exist and work, confirmed in the BUG-002 regression work) but is not what the generated
hooks do automatically.

### BUG-004 — `asset info` on a missing asset exits 0
Prints `✗ Asset '...' not found.` but returns exit code 0. Confirmed isolated to this one
command (`theme validate`, `asset remove`, `font info` all correctly exit 1 on equivalent
failures). **STATUS: FIXED this session.**

### BUG-005 — Bubble Tea spinner/list wrappers hang without a fallback (FIXED)
`registry list/browse/fetch/search` and `theme list` are wrapped in `tui.RunSpinner`/
`tui.RunThemeList`. In non-fully-interactive contexts, they hung indefinitely (confirmed
25s+, both for instantly-rejectable and plausible inputs — the hang was in the TUI wrapper
itself, before the wrapped logic's result mattered).

**Root cause, confirmed from source, not inferred from symptoms:** `tea.NewProgram(m).Run()`
blocks on its own internal terminal/input initialization when stdin isn't a real TTY, and
does this by hanging rather than returning an error — so `p.Run()`'s error return can never
be used to detect the case. `RunSpinner`'s own model (`internal/tui/spinner.go`) has zero
keypress handling and calls `tea.Quit` the instant its wrapped task finishes, confirming the
hang is entirely inside Bubble Tea's own setup, upstream of any application logic.

**Intended contract, also confirmed from source:** `registry list`'s own code comment reads
*"show as glamour table — interactive browse is via `registry browse`"* — explicitly
establishing `list`/`fetch`/`search` as the intended non-interactive-safe commands, with
`browse` as the deliberately separate interactive alternative. `theme list` already had a
`printPlainThemeList` fallback, proving non-interactive usability was intended, but it only
triggered on a `p.Run()` *error*, which a hang never produces. `registry browse` and
`theme list`'s picker stage are genuinely, permanently interactive by design (no
auto-completion condition — they require a keypress); a non-interactive hang wasn't
violating that intent, but hanging forever instead of failing fast was still the wrong
failure mode.

**STATUS: FIXED.** Added `internal/tui/interactive.go`'s `IsInteractive()` (using
`github.com/mattn/go-isatty`, already an indirect dependency, checked *before*
`tea.NewProgram` is ever called — checking `p.Run()`'s error afterward can't work, per the
root cause above). `RunSpinner` now skips Bubble Tea entirely when non-interactive: prints
the message plainly and runs the task directly, with an identical return contract (zero
call-site changes anywhere). `RunThemeList` now returns an error immediately when
non-interactive, which routes automatically into `theme list`'s *existing*
`printPlainThemeList` fallback with zero changes needed to `cmd/theme.go`. `registry
browse`'s silent no-op on error was replaced with a clear message pointing to `registry
list`/`fetch` (and correctly separated from the `q`/cancel case, which had been wrongly
conflated with genuine errors before).

Both interactive and non-interactive paths were live-verified with the real binary:
`registry fetch "../../evil"` (the exact command that hung 25s+ before) now returns almost
instantly with the plain-text spinner message and the correct, previously-unreachable
path-traversal validation error; `theme list` now returns instantly with a real plain-text
listing of all 7 bundled themes; `registry list` returns correctly including a real network
fetch and glamour-rendered table; `registry browse` now prints the new clear message instead
of hanging or silently returning nothing. The interactive TUI path itself could not be
live-tested (no real TTY available in this environment, true throughout this entire
investigation) but is structurally unchanged — the `!IsInteractive()` branches are purely
additive early-returns; no line inside the original interactive path was touched.
Regression coverage: `internal/tui/interactive_test.go` (new — `internal/tui` had zero
prior test files) exercises the real non-interactive branch directly, since `go test`
itself runs with non-TTY stdin/stdout: confirms `RunSpinner` runs its task and returns
within 3s (not hanging) with the task's error correctly propagated, and confirms
`RunThemeList` returns an error within 3s rather than hanging.

### BUG-006 — Status bar CLI wiring gaps (same pattern floaters already had)
Status bar never received the integration-gap-closing pass floaters got (per `.llm/todo.md`,
whose own checklist left Mascots and Status Bar unchecked under "Exotic Assets" while
Floaters' 4 fixed integration gaps are itemized in detail). Four identical-pattern gaps
found and fixed: (1) `assetPreviewCmd`'s generic switch had no `AssetTypeStatusBar` case —
silent no-op; (2) `assetStatusCmd`'s "Active Assets" slots list omitted `mascot`,
`status-bar`, and `sound` entirely; (3) `assetInfoCmd` had no `StatusBar` detail block;
(4) `theme preview` never showed the active status bar (only floaters, via
`showActiveFloaters`). **STATUS: FIXED, live-verified** — added `showActiveStatusBar()`
mirroring `showActiveFloaters()`, wired into `theme preview`; all four confirmed working via
a real status-bar asset built during testing.

Separate, minor finding surfaced while testing: `asset validate` doesn't cross-check a
manifest's `"type"` field against the sub-config object actually present, so a manifest
with a wrong/mismatched `type` string can pass validation while still being correctly
routed by the type inferred from its directory — worth tightening later, not fixed here
(no user-facing impact found; it only masked a mistake in a hand-written test manifest).

### BUG-007 — Generated bash/zsh status-bar code never actually renders (CRITICAL, FIXED)
Discovered live-testing BUG-006. The bash and zsh generators compute each segment's value
into one variable-naming scheme, then assemble the visible bar by checking a **different**,
disconnected variable-naming scheme — so the assembled bar always checked variables the
compute step never assigned. **Any status bar with more than one segment in a zone (i.e.
almost any real one) would render completely empty for every user, in both bash and zsh,
since this shipped in commit `d469f5f`.**

- **Bash** (`bashSegmentCode`, `internal/assets/statusbar.go`): computed into hardcoded,
  per-*type* names (`__seg_dir`, `__seg_exit`, ...) with no per-instance uniqueness at all,
  while assembly checked `bashSegVarName(seg, idx)`-style indexed names
  (`__seg_directory_0`, `__seg_exitcode_100`, ...) that compute never produced.
- **Zsh** (`zshSegmentCode`): did call the same indexed-naming helper, but with `idx`
  hardcoded to `0` for every segment regardless of position, while assembly used the real
  zone-relative loop index — so only the first segment in each zone ever had a chance of
  matching, and even that only by coincidence (index 0 in both).
- **PowerShell was not affected** — its generator uses a fundamentally different, sound
  design (compute one segment into a single reused `$__seg` variable, immediately
  consume/print it, then move to the next segment) that never needed unique per-segment
  names. Confirmed correct via genuine live execution in this session (not just generated-
  text inspection): a hand-built copy of the real generated function, run in an actual
  PowerShell process in this repo, printed `[DIR:cmdX][GIT:main*]` — correctly reflecting
  this repo's real directory name and real dirty-git-status at the time of the test.

**Root cause:** `bashSegVarName(seg, idx)` (the indexed, collision-free naming helper) was
clearly written to be the single source of truth for segment variable names, but the
value-computation functions were never connected to it (bash) or were connected with a
hardcoded index (zsh) — two independently-written code paths that were never integration-
tested together. No prior test caught this: existing tests called the generators and
checked superficial properties (non-empty, contains `PROMPT_COMMAND`) but never traced
whether assembly's checks matched compute's assignments.

**STATUS: FIXED, with regression coverage, verified three ways:** (1) full test suite
passes; (2) two new tests parse real generated bash/zsh code and assert every variable
assembly checks was actually assigned by compute (`TestStatusBarShellCode_
BashAssemblyVarsAreAllAssigned`, `...ZshAssemblyVarsAreAllAssigned`, using the existing
2-segments-per-zone fixture, which is exactly the shape that triggered the bug) plus
`TestBashSegVarName_UniquePerZoneRelativeIndex`, all in `internal/assets/statusbar_test.go`;
(3) manually regenerated bash code for a real 5-segment test status bar and confirmed by
inspection that every `[ -n "$X" ]` assembly check now has a matching `X=` compute line.

**Known, deliberately out-of-scope remainder (not fixed, documented only):** center-zone
segments are still never assembled into the visible bar in bash or zsh (compute now
produces a harmless, correctly-named but unused line for them); PowerShell's generator only
ever processes the `left` zone — center and right segments are silently dropped
(`segmentsByZone`'s second and third return values are discarded in
`statusBarPowerShell`). This is a materially larger feature-completeness gap (implementing
full 3-zone rendering across three shell dialects) than the naming-mismatch bug above, and
was intentionally left as a separate, clearly-flagged next unit rather than silently
expanded into during this fix.

### GAP-001 — Web builder schema stale relative to Go core / VS Code extension (CLOSED — FIXED)
`web-builder/src` had zero references to `mascot`, `floater`, `status_bar`, or `sound`.
Evidence resolved this cleanly as **Option A (parity intended, not narrower by design)**:
`web-builder/README.md` explicitly claims *"Full schema coverage — every field from
`internal/config/types.go`'s `Theme` struct: ... and linked asset names"* as a stated
feature. The web builder was also built in commit `b483170` — the same commit that added
sound, i.e. *after* mascot/floater/status_bar already existed in the Go schema (added in
the earlier `d469f5f`) — so this cannot be explained by "predates the feature." The scope
is precisely "linked asset **names**" (simple string fields referencing an asset by name,
matching a theme's `assets` block), not full asset-type authoring (e.g. building a mascot's
state machine or a status bar's segments in the browser) — that broader claim was never made
anywhere and was correctly left out of this fix.

**STATUS: FIXED.** Added `mascot?`, `floater?`, `status_bar?`, `sound?` to the TypeScript
`ThemeAssets` interface (`web-builder/src/theme.ts`) and four matching text-field inputs to
`assetsFields()` (`web-builder/src/formPanel.ts`), mirroring the existing
spinner/banner/divider/icons fields exactly. `validate.ts` and `exportImport.ts` needed no
changes — neither has any per-asset-slot-specific logic; both are generic over whatever
fields `Theme`/`ThemeAssets` declares. Verified: `npx tsc --noEmit` clean, `npm run build`
clean (10 modules, same as before).

### GAP-004 — `asset use --as` help text understated the real capability (FIXED)
The `--as` flag's help string and one error message both said "spinner, banner, divider,
icons, floater" — but the actual `assetUseCmd` implementation's `validTypes` map, switch
logic, and even a *different* error message inside the same function already correctly
handled all 8 types, including mascot and status-bar (with helpful extra hints like
`View states: cmdx asset mascot-info <name>`). This was pure help-text/error-message drift,
not a functional gap — confirmed by reading the full function body rather than just
`--help` output, then live-verified: `asset use verifymascot --as mascot` worked correctly
end-to-end, printing the expected mascot-specific hints. **STATUS: FIXED** — both stale
strings in `cmd/asset.go` now list all 8 types; live-reverified via `asset use --help` and
a real `asset use --as mascot` invocation.

### GAP-002 — Mascot trigger priority tie-break semantics (CLOSED — documented)
`ResolveState`'s tie-break (`trigger.Priority > best.priority`, strict greater-than) means
that two triggers left at equal priority (including the common case of both defaulting to 0)
resolve via Go's randomized map iteration order — non-deterministic across process
invocations. This is not a bug in the resolver's contract (undocumented behavior is
technically "correct" by omission) but is a real authoring hazard for theme/mascot authors
who don't realize they need to set explicit distinct priorities.

**STATUS: CLOSED — documented, with regression coverage.** No runtime behavior was changed
(the resolver was already correct; only its unspecified-tie-break case needed writing down).
Documented in `internal/assets/.CLAUDE.md`'s new "Mascot Triggers" section (which also
corrects that file's own stale TODOs for floater/mascot/status-bar/sound — those shipped in
commit `d469f5f`, this nested doc had never been updated since). Regression test added:
`internal/assets/mascot_test.go`'s `TestResolveState_EqualPriorityTieIsToleratedNotFatal` —
asserts a tie always resolves to one of the genuinely-matching candidates (never panics,
never returns an unrelated/empty state) without asserting which one, since which one is
correctly unspecified.

## UNVERIFIED (not bugs — genuinely untested, do not assume working or broken)
- Status bar — core generation verified working (bash/zsh compute+assembly consistency,
  PowerShell *and real bash* live-executed) for the `left`+`right` zones; center-zone
  rendering across all three shells and right-zone rendering in PowerShell remain
  unimplemented, not merely untested (see BUG-007's deliberately-out-of-scope note).
  Narrow-terminal behavior, conditional visibility, and segment types beyond
  `directory`/`git`/`exit_code`/`time`/`text` were not individually exercised.
- TUI editor — never launched (no real TTY available in this environment, throughout).
- True audible confirmation of sound playback — subprocess-level evidence (a real `ffplay`
  invocation taking 8.89s, far longer than the 0.5s clip, with exit 0) is the strongest
  evidence obtainable through this tooling; see the dedicated session below.
- zsh live execution — confirmed genuinely unavailable in this environment (not installed;
  bash is, via WSL, and was live-executed this session). Generated zsh code remains
  regression-tested for internal consistency only (BUG-007).
- Registry's `fetch`/`browse` against a real, valid (not validation-rejected) theme name —
  `list`'s real network fetch was exercised and confirmed working; a full successful
  download via `fetch` was not separately exercised.

## Config / wallpaper / bash / sound verification session
Closes out the items the prior session marked "no safe isolation mechanism found" — that
conclusion was wrong for config and wallpaper; both have safe, working sandboxes, just not
the `CMDX_THEMES_DIR`-style ones already in use elsewhere. No source code changes resulted
from this session; everything tested either passed cleanly or was confirmed intentional.

- **Config**: sandboxed via `$env:USERPROFILE`/`HOME` override (same technique already used
  for shell-profile testing — `os.UserHomeDir()` resolves `~/.cmdx` the same way). `config
  show`/`set`/`reset` all verified: `set` persists across process invocations, `reset`
  restores defaults, invalid key and invalid value both fail cleanly with exit 1 and a
  specific error naming the valid options. No defects found.
- **Wallpaper**: sandboxed by overriding `LOCALAPPDATA` directly — confirmed from source
  (`internal/wallpaper/windows_terminal.go`'s `findWindowsTerminalSettings()` reads
  `os.Getenv("LOCALAPPDATA")`/`"APPDATA"` directly, not `os.UserHomeDir()`, so it needed a
  different override than config). `set`/`info`/`remove` all verified against a sandboxed
  Windows Terminal `settings.json`: `set` correctly wrote image path, opacity, alignment,
  and stretch mode; `remove` correctly cleared them; a missing image path failed cleanly
  with exit 1. The implementation backs up `settings.json` to `.cmdx.bak` before ever
  modifying it. No defects found.
- **Real bash execution**: this environment turns out to have genuine bash available via
  WSL (Ubuntu, bash 5.3.9) — previously missed because the investigation had only checked
  for a native Windows bash/zsh. Cross-compiled a Linux `cmdx` binary (`GOOS=linux go
  build`), built a real mascot asset, and in a fully sandboxed `$HOME`: `theme inject`
  wrote a correct, correctly-colored `.bash_profile`; the generated mascot hook correctly
  captured real exit codes from real `true`/`false` commands; 3x repeated injection was
  idempotent; `theme remove` cleanly emptied the profile. This closes the "never executed
  in a real bash" caveat that had persisted through every prior phase. One caveat: the
  hook's `command=$(history 1 ...)` capture was empty in this non-interactive-script
  invocation (bash doesn't track history non-interactively by default) — very likely a
  testing-methodology limitation rather than a real gap (a genuine interactive terminal
  session tracks history automatically), but not independently confirmed in a real
  interactive bash session, since none was available.
- **Real audio playback**: using the sound system's documented custom-player escape hatch
  (`"player": "ffplay -nodisp -autoexit -loglevel error %f"` — confirmed safe, argv-based,
  never shell-interpreted) with a real WAV file and `ffplay` (confirmed installed on this
  machine), `sound-play` took 8.89 seconds to return with exit 0 — far longer than the
  0.5-second clip itself, confirming a real blocking subprocess genuinely performed audio
  device I/O. This is the strongest playback evidence obtainable through this tooling;
  true audible confirmation remains outside what it can observe.

### New findings (minor, not fixed this session)
- `theme inject`'s startup banner prints unconditionally regardless of the theme's
  `banner.enabled: false` setting — observed live during the real-bash test (a theme with
  `enabled: false` still printed its banner text on every new shell). Real, small, worth a
  future fix.
- Confirmed **intentional, not a bug**: `mascot-state`/`sound-play`-family commands exit 0
  with zero output when `chafa` isn't installed (`cmd/asset_mascot.go` explicitly checks
  `ChafaAvailable()` first) — sensible since these fire on every single shell prompt via
  the hook; a noisy per-prompt error would be worse UX than silence.
