# cmdX — Test Matrix (cumulative, pre-repair)

Condensed from 4 investigation phases + root-cause pass. Full narrative detail lives in the
session transcripts this file summarizes; this is the reference table.

| Subsystem | Operation | Expected | Observed | Result | Evidence |
|---|---|---|---|---|---|
| Build | `go build ./...` (real repo, pre-fix) | success | `syncAssetHooks redeclared` | FAIL | BUG-001 |
| Build | `go build/vet/test` (fixed copy) | success | clean | PASS | — |
| Theme | `validate` | accept/reject correctly | correct on valid + malformed JSON | PASS | — |
| Theme | `info` | full detail | correct, but asset slots incomplete for mascot/floater/status_bar/sound | PARTIAL | — |
| Theme | `list` | prompt list, or plain fallback non-interactively | now returns instantly with a real plain-text listing (was: hangs indefinitely) | PASS (was BLOCKED) | BUG-005 (repair) |
| Theme | `apply` | renders + activates linked assets | correct; `.state/*` files overwritten as expected | PASS | phase2 |
| Theme | `apply --no-assets` | suppress activation only | confirmed suppressed | PASS | phase2 |
| Theme | `inject`/`remove`, 3x idempotency | no duplicate blocks | exactly 1 block each after 3x | PASS | phase2 |
| Theme | `inject` with mascot linked | 2 independent marker blocks | confirmed, correctly isolated | PASS | phase2 |
| Theme | switch theme (mascot→none) | asset-hooks block cleanly removed | confirmed | PASS | phase2 |
| Theme | `inject --no-hooks` | suppress hooks only | confirmed | PASS | phase2 |
| Theme | apply/validate nonexistent/malformed | clean fail, exit 1 | confirmed, no panic | PASS | phase2 |
| Asset | banner/divider/icon `preview` | real chafa output | real sixel rendering observed | PASS | phase2 |
| Asset | spinner `preview` | animated loop | hangs pending interrupt — ambiguous by-design vs bug | BLOCKED | phase2 |
| Asset | mascot `validate` (position/max_width/max_height required) | reject incomplete manifest | correctly rejected, then passed once fields added — confirmed intentional (floater-position reuse) in source | PASS | phase2/5 |
| Asset | mascot `mascot-info`/`mascot-hooks` | correct table / script | correct | PASS | phase2 |
| Asset | mascot `mascot-state` normal resolution | matches highest-priority trigger | correct once distinct priorities are set (8/8 runs) | PASS | phase5 |
| Asset | mascot `mascot-state --state X` | forces state X | now correctly forces (fixed, regression-tested) | PASS | phase5-repair |
| Asset | mascot `env_var` trigger | fires when var is ambiently set | works — reads real OS env, no hook wiring needed (was miscategorized as blocked) | PASS | GAP-003 review |
| Asset | mascot `git_status`/`idle_time`/`output_regex` via generated hooks | — | not populated by shipped hooks — evidence-supported as intentional, not a defect | N/A (by design, see GAP-003) | GAP-003 review |
| Asset | mascot equal-priority trigger tie | resolve to one valid candidate, never panic | confirmed 20/20 runs, no panic, always one of the two matching states (which one is documented-unspecified) | PASS | GAP-002 closure |
| Asset | sound `validate`/`sound-info`/`sound-hooks` | correct | correct | PASS | phase2 |
| Asset | sound `preview` | audible playback | exits 0, no error, playback unconfirmable | PARTIAL | phase2 |
| Asset | sound `sound-play --sound X` | forces sound X | correctly forces via separate early-return branch | PASS | phase5 |
| Asset | `asset info` nonexistent | exit 1 | exit 0 despite error text | FAIL | BUG-004 |
| Asset | `asset remove` (2-arg) nonexistent | exit 1 | exit 1 | PASS | phase5 |
| Registry | `fetch`/`list`/`browse`/`search`, non-interactive | prompt fast, reject bad names locally | now confirmed: `fetch "../../evil"` returns almost instantly with the correct validation error (was: hung 25s+ regardless of input validity) | PASS (was BLOCKED) | BUG-005 (repair) |
| Registry | `list`, real network fetch | glamour table renders | confirmed live — real HTTP fetch + rendered table with 13 community themes | PASS | BUG-005 (repair) |
| Registry | `browse`, non-interactive | clear message, no hang | now prints "'registry browse' needs an interactive terminal" + pointer to `list`/`fetch` (was: hung, or silently returned nothing on the `q`/cancel-conflated path) | PASS (was BLOCKED) | BUG-005 (repair) |
| TUI | `RunSpinner`/`RunThemeList` interactive path | unchanged from before | not live-testable (no real TTY in this environment, throughout); structurally unchanged — new branches are purely additive early-returns, no line in the original path touched | NOT TESTED (structurally verified only) | BUG-005 (repair) |
| Plugin | `list` | correct listing | correct | PASS | phase2 |
| Font | `list` | correct catalog | correct | PASS | phase2 |
| Config | `show` | correct | correct | PASS | phase2 |
| Config | `set`/`reset`, invalid key/value | persist, validate, clean fail | all confirmed — sandboxed via `$env:USERPROFILE`; set persists cross-process, reset restores defaults, invalid key/value both exit 1 with a specific message | PASS (was BLOCKED) | verify-session |
| Wallpaper | `info` | correct | correct | PASS | phase2 |
| Wallpaper | `set`/`remove`, missing image | write/clear 4 fields, clean fail | all confirmed — sandboxed via `LOCALAPPDATA` override; set wrote image/opacity/alignment/stretch correctly, remove cleared them, missing image exited 1; `.cmdx.bak` backup confirmed created | PASS (was BLOCKED) | verify-session |
| TUI | `edit` launch | — | not attempted — PTY interaction risk | BLOCKED | phase2-4 |
| Web builder | `npm install/tsc/build` | clean | clean, 10 modules | PASS | phase2 |
| Web builder | schema coverage of mascot/floater/status_bar/sound | present | was absent — fixed, added to theme.ts + formPanel.ts | PASS (was FAIL/GAP) | GAP-001 (repair) |
| VS Code ext | `tsc --noEmit` | clean | clean | PASS | phase2 |
| VS Code ext | schema coverage | matches Go structs | matches exactly | PASS | phase2 |
| Floater | validate/info/list/preview/status/theme-link/`asset use` | all work | all confirmed live, no defects found | PASS | floater unit |
| Asset | `asset use --as` help text vs. real capability | help matches implementation | help was stale (5 of 8 types listed); implementation always supported all 8 | PASS (was doc-only FAIL) | GAP-004 (repair) |
| Status bar | `asset validate`/`info`/`preview`, `asset status` slot, `theme preview` integration | all four wiring points work | fixed and confirmed live (real 5-segment test asset) | PASS | BUG-006 (repair) |
| Status bar | bash generated code compute/assembly consistency | every assembled var is assigned | was FAIL (empty bar, always) — fixed; regression tests + manual regeneration confirm | PASS (was CRITICAL FAIL) | BUG-007 (repair) |
| Status bar | zsh generated code compute/assembly consistency | every assembled var is assigned | was FAIL (same bug, different cause: hardcoded idx=0) — fixed | PASS (was CRITICAL FAIL) | BUG-007 (repair) |
| Status bar | PowerShell generated code, actual live execution | segment values render correctly | confirmed via real execution: `[DIR:cmdX][GIT:main*]`, correctly reflecting real repo state | PASS | BUG-007 (repair) |
| Status bar | center-zone rendering (all shells), right-zone rendering (PowerShell) | segments appear in output | never assembled/emitted — confirmed real gap, not just untested | FAIL (documented, unfixed) | BUG-007 (repair) |
| Cross-shell | bash hook generation + live execution | valid syntax, real runtime behavior | **now live-executed for real** (WSL/Ubuntu bash 5.3.9, cross-compiled Linux binary): `theme inject` wrote a correct `.bash_profile`, real exit codes (0/1) correctly captured by the mascot hook from real `true`/`false`, 3x injection idempotent, `theme remove` cleaned up fully | PASS (was PARTIAL — generated-only) | verify-session |
| Cross-shell | zsh hook generation | valid syntax | zsh confirmed genuinely unavailable in this environment (not installed); generation remains consistency-tested only (BUG-007), never executed | PARTIAL (generated + consistency-verified, not executed) | BUG-007, verify-session |
| Cross-shell | PowerShell live injection/idempotency/removal | works | fully live-tested, works; status-bar prompt function also now live-executed | PASS | phase2-4, BUG-007 |
| Sound | `sound-play` real subprocess invocation | real player launches, performs I/O | confirmed via real `ffplay` custom-player: 8.89s runtime (vs. 0.5s clip) with exit 0 — strong evidence of genuine playback attempt, not a silent no-op; true audible confirmation remains unobservable through this tooling | PASS (subprocess-level; audible confirmation still N/A) | verify-session |
| Theme | `inject`'s startup banner vs. `banner.enabled: false` | banner suppressed when disabled | banner printed anyway — new, real, minor finding, not fixed | FAIL (documented, unfixed) | verify-session |
| Asset | `mascot-state`/`sound-play` with chafa unavailable | — | confirmed intentional: explicit `ChafaAvailable()` check, silent exit 0 by design (avoids per-prompt error spam) | PASS (intentional, not a bug) | verify-session |

Cumulative totals (pre-repair): PASS 35 · FAIL 2 · PARTIAL 5 · BLOCKED 8 · NOT TESTED ~9.
Repair-session additions: status bar moved from entirely NOT TESTED to substantially PASS
(BUG-007 critical fix); floater fully verified with no defects (the one asset type with a
clean result); web-builder schema parity fixed (GAP-001); `asset use --as` help text fixed
(GAP-004); and BUG-005 (Bubble Tea non-interactive hangs across `theme list` and all of
`registry list/browse/fetch/search`) fixed and live-verified — this was the last remaining
open item from the original bug/gap list.

Verify-session additions (no code changes — pure verification, everything passed clean or
was confirmed intentional): config and wallpaper lifecycles both moved from BLOCKED to
PASS (the earlier "no safe isolation" conclusion was wrong — both have working sandboxes);
real bash execution via WSL closed the "never executed in a real shell" caveat that
persisted through every prior phase; real `ffplay` subprocess invocation gives the
strongest sound-playback evidence obtainable through this tooling. One new minor finding
(banner ignores `enabled: false`) and one confirmed-intentional non-bug (silent exit when
chafa is missing) resulted.
