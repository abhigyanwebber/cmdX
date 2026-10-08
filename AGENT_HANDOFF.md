# cmdX — Agent Handoff

Read this before editing anything. It supersedes conversational memory from any prior
session. If something here conflicts with `.llm/HANDOFF.md` or `.llm/todo.md`, **this file
and its siblings (PROJECT_BIBLE.md, ARCHITECTURE.md, VERIFIED_BASELINE.md,
KNOWN_ISSUES.md, TEST_MATRIX.md) are authoritative** — they were produced by an independent
forensic audit plus real command execution, not self-reported completion claims.

## Current Repository State (at the start of this repair session)
- Branch `main`, HEAD `b483170`, **in sync with `origin/main`**.
- `origin/master` is a separate, stale branch, 30 commits behind / 1 commit ahead of `main`
  — do not attempt to merge/reconcile it without explicit instruction; it appears to be an
  old default-branch artifact, not a parallel line of real work in progress.
- 5 files were uncommitted at session start implementing theme→asset auto-activation
  (`activateAsset`, `activateThemeAssets`, `syncAssetHooks`, `AssetHooksStart/End` markers).
  This repair session fixes the build blocker in that work and commits it — see the repair
  report for exact diffs and the proposed commit breakdown.

## Architectural Invariants — do not casually change these
1. Theme injection (`InjectStart/InjectEnd`) and asset hooks (`AssetHooksStart/AssetHooksEnd`)
   use **independent** marker blocks. Never merge them into one block or key removal of one
   off the presence of the other.
2. Repeated `theme inject`/`apply` must be idempotent — no duplicate marker blocks ever.
3. Theme→asset activation must not duplicate `.state/*` files or hook blocks; switching
   themes must cleanly replace, not append.
4. `--no-assets` suppresses `activateAsset()` only. `--no-hooks` suppresses `syncAssetHooks()`
   only. They are independent flags — do not conflate them.
5. `ResolveState`'s priority resolution (highest `Priority` wins, ties broken by map
   iteration order — currently undocumented/unspecified) must not be changed without
   updating GAP-002 and adding a tie-break test, since real themes may depend on either the
   current behavior or its absence of guarantee.
6. Mascot's `position`/`max_width`/`max_height` requirement is **intentional** (reuses the
   floater four-corner model, per an explicit source comment in `internal/assets/types.go`)
   — do not "fix" this as if it were a bug.
7. Shell integration profile paths are hardcoded via `os.UserHomeDir()` with no env-var
   override — any testing must sandbox `$HOME`/`$USERPROFILE`, never run against a real
   profile directly.

## Verified Behavior
See VERIFIED_BASELINE.md's "VERIFIED" section — treat it as ground truth, not as an
exhaustive feature list.

## Known Bugs
See KNOWN_ISSUES.md. BUG-001, BUG-002, BUG-004 are fixed as of this session, with BUG-002
now backed by regression tests in `internal/assets/manager_mascot_test.go` (verify against
the repair report's build/test results before trusting that fix status if time has passed).
Former BUG-003 (mascot hooks only populate exit_code/command) has been reclassified as
GAP-003 — a documented, evidence-supported capability boundary rather than a bug; do not
implement `git_status`/`idle_time`/`output_regex` hook-population without re-reading
KNOWN_ISSUES.md's GAP-003 entry first, since that decision was made on inference from
`.llm/HANDOFF.md`'s phrasing and `MascotContext`'s graceful-degradation design comment, not
a single explicit repository statement — if you find stronger evidence either way, update
GAP-003 rather than silently acting against it. `env_var` triggers already work (ambient OS
env, no hook wiring needed) — this was a correction to the original bug report, not new work.
GAP-002 (priority tie-break semantics) is now also closed — documented in
`internal/assets/.CLAUDE.md`'s "Mascot Triggers" section (which also fixed that file's own
stale floater/mascot/status-bar/sound TODOs), with a regression test
(`TestResolveState_EqualPriorityTieIsToleratedNotFatal`) confirming ties resolve to a valid
candidate without asserting which one, since which one is correctly left unspecified.
BUG-005, GAP-001 remain open and documented, not fixed.

Two more were found and fixed this session while closing status bar's integration gaps
(the same class of gap floaters already had): BUG-006 (four CLI wiring gaps — generic
`asset preview`/`info`/`status`, and `theme preview` never showed status bar) and, far more
importantly, **BUG-007**: bash and zsh status-bar generation computed segment values into
one variable-naming scheme and assembled the visible bar by checking a different,
disconnected one — so any status bar with more than one segment per zone rendered
completely empty, for every user, since it shipped. Fixed, regression-tested (new tests in
`internal/assets/statusbar_test.go` parse real generated code and assert compute/assembly
consistency), and verified via actual live PowerShell execution (PowerShell wasn't affected
— different, sound design — but was the closest available proxy for confirming the shared
segment-computation logic actually works end-to-end). Read BUG-007's full KNOWN_ISSUES.md
entry before touching `internal/assets/statusbar.go` again: it documents a deliberately
out-of-scope remainder (center-zone never assembled in any shell; PowerShell only ever
processes the `left` zone) that is a real, separate, larger feature-completeness gap, not
covered by this fix.

Floater's entire lifecycle was verified end-to-end this session with **no defects found**
(validate/info/list/preview/status/theme-linking/`asset use --as floater --position`, all
live-tested) — a genuinely different, positive outcome from status bar, consistent with
floater already having received the integration-gap pass mascot/status-bar hadn't. One
stale doc comment was fixed along the way (`showActiveFloaters` in `cmd/helpers.go`
wrongly claimed floaters aren't theme-linkable; `internal/config/types.go`'s comment, and
the actual working code, say otherwise — corrected, no functional change).

GAP-001 (web-builder schema staleness) is now **closed and fixed**: evidence
(`web-builder/README.md`'s own "full schema coverage... including linked asset names"
claim, plus the fact web-builder was built in `b483170`, *after* mascot/floater/status_bar
already existed) settled it as unintentional staleness, not narrower-by-design scope. Added
the 4 missing linked-asset-name fields to `web-builder/src/theme.ts` and
`formPanel.ts`, mirroring the existing 4 fields exactly. This does **not** mean the web
builder authors full mascot/status-bar/sound *configs* (state machines, segments, etc.) —
only that it can reference them by name in a theme, matching its own stated scope.

GAP-004 (new, also closed): `asset use --as`'s help text and one error message understated
what the command already did — the actual implementation already handled all 8 asset
types correctly (confirmed by reading the full function, not just `--help`); only two
strings in `cmd/asset.go` were stale. Fixed and live-reverified.

BUG-005 is now **fixed**: added `internal/tui/interactive.go`'s `IsInteractive()`
(`github.com/mattn/go-isatty`, checked *before* `tea.NewProgram` is ever called, since
`p.Run()`'s own error return can never detect this case — it hangs rather than erroring).
`RunSpinner` skips Bubble Tea entirely when non-interactive (prints the message, runs the
task directly, identical return contract); `RunThemeList` returns an error immediately,
routing into `theme list`'s pre-existing `printPlainThemeList` fallback with zero changes
to `cmd/theme.go`; `registry browse`'s silent no-op was replaced with a clear message.
Live-verified with the real binary: `registry fetch "../../evil"` (the exact command that
hung 25s+ before) now returns almost instantly; `theme list` and `registry list` (real
network fetch) both return correctly; `registry browse` prints the new message. The
interactive TUI path itself was **not** live-tested (no real TTY available in this
environment, true throughout) but is structurally unchanged — the new branches are purely
additive early-returns. Regression tests added in `internal/tui/interactive_test.go` (this
package had zero prior tests) exercise the real non-interactive branch directly, since
`go test` itself runs with non-TTY stdin/stdout. Read BUG-005's full KNOWN_ISSUES.md entry
if you need to touch `internal/tui/spinner.go` or `list.go` again.

There is no remaining open item from the original bug/gap list — BUG-001 through BUG-007
and GAP-001 through GAP-004 are all either fixed or correctly reclassified/closed as
intentional design. What's left is the genuinely-unverified surface below, not known defects.

A later close-out session verified config, wallpaper, real bash execution (via WSL — this
environment has genuine bash 5.3.9, previously missed), and sound more thoroughly. No code
changed — everything tested passed cleanly or was confirmed intentional. Key corrections to
earlier conclusions: **config and wallpaper do have safe sandboxes** (config via
`$env:USERPROFILE`/`HOME`, same as shell-profile testing; wallpaper via `LOCALAPPDATA`
directly, since `internal/wallpaper/windows_terminal.go` reads that env var rather than
`os.UserHomeDir()`) — the earlier "no safe isolation mechanism" conclusion in this file was
wrong and should not be repeated. Real bash execution (cross-compiled Linux binary run in
WSL) confirmed theme injection, real exit-code capture by the mascot hook, idempotency, and
clean removal all work exactly as the PowerShell-only testing had suggested. Real `ffplay`
subprocess invocation (8.89s runtime vs. a 0.5s clip, exit 0) is the strongest sound-
playback evidence obtainable through this tooling. One new minor finding, not fixed:
`theme inject`'s startup banner ignores `banner.enabled: false`. One confirmed-intentional
non-bug: `mascot-state`/`sound-play` silently exit 0 when chafa is missing (explicit
`ChafaAvailable()` check — correct UX for a per-prompt hook, not an oversight).

## Unverified Areas — do not assume working OR broken
TUI editor (no real TTY available in this environment, throughout). True audible
confirmation of sound playback (subprocess-level evidence exists — see above — but genuine
audible confirmation is outside what this tooling can observe). zsh literal execution
(confirmed genuinely unavailable in this environment — not installed; bash is, via WSL, and
was live-executed). Registry's `fetch`/`browse` against a real, valid theme name (`list`'s
real network fetch was exercised and confirmed working; a full successful download via
`fetch` wasn't separately exercised). Status bar's center-zone (all shells) and right-zone
(PowerShell) rendering are confirmed *not implemented*, not merely untested — don't build
on top of them without adding that first.
Do not write documentation or make claims about these without first testing them.

## Development Rules
1. Read this document (and its siblings) before editing.
2. Inspect current source before implementing — do not assume prior session summaries
   (including this one) are still accurate; things may have changed.
3. Compile (`go build ./...`) and test (`go test ./...`, `go vet ./...`) after any
   meaningful change. Do not declare something fixed or complete without running it.
4. Update VERIFIED_BASELINE.md / KNOWN_ISSUES.md when behavior changes — mark items FIXED
   with evidence, don't just delete the historical record of the defect.
5. Preserve unrelated uncommitted work — never `git reset`/`checkout`/`stash --drop`
   without explicit instruction; branch divergence is not yours to resolve unilaterally.
6. Avoid scope expansion — every fix in this project so far was deliberately scoped to
   one confirmed defect at a time, with evidence for intended behavior established first
   (see GAP-003 and GAP-001 for examples where the evidence-gathering itself, not the
   code change, was the hard part). Don't silently implement documented-but-unfixed items
   (status-bar center/right zones, the `banner.enabled` finding) in an unrelated session;
   treat each as its own deliberately-scoped unit.
7. When testing shell integration, sandbox `$HOME`/`$USERPROFILE` to a temp directory —
   never let a test write to the real user's shell profile or global `~/.cmdx` config.
8. Bubble Tea commands (`theme list`, `registry list/browse/fetch/search`) no longer hang
   non-interactively as of BUG-005's fix — but if you're testing a *new* Bubble Tea-wrapped
   command, still use a hard timeout (e.g. PowerShell `Start-Job`/`Wait-Job -Timeout`)
   rather than a bare blocking call until you've confirmed it also checks
   `tui.IsInteractive()` before calling `tea.NewProgram`.
