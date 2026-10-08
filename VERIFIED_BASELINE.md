# cmdX — Verified Baseline

Cumulative results across the full investigation: 4-phase black-box pass, mascot/sound
root-cause pass, a controlled repair session (BUG-001/002/004/006/007, GAP-001/002/003/004),
BUG-005's fix, and this session's close-out pass (config, wallpaper, sound, and — for the
first time — real bash execution). This file records what was **actually observed**, not
what documentation claims. See TEST_MATRIX.md for full per-command detail and evidence.

There is no remaining known defect from the original bug/gap list — BUG-001 through
BUG-007 and GAP-001 through GAP-004 are all fixed or correctly reclassified as intentional
design. What's below is the state of the genuinely-unverified surface after this session's
close-out pass against it.

## VERIFIED (observed working end-to-end)
- Theme lifecycle: `validate`, `info`, `apply`, `inject`, `remove` — full round-trip.
- Theme → asset auto-activation (spinner/banner/divider/icons/mascot/status-bar linking),
  including `--no-assets`/`--no-hooks` correctly suppressing only their respective behavior.
- Shell-hook marker isolation and idempotency — no cross-contamination, no duplicate blocks
  across repeated injection, clean removal of both blocks.
- Mascot `--state` force-override (BUG-002, fixed), `ResolveState`'s priority resolution
  (deterministic given distinct priorities), equal-priority tie tolerance (GAP-002).
- Floater — full lifecycle, no defects found (validate/info/list/preview/status/theme-link/
  `asset use --as floater --position X`), the one asset type with a clean result throughout.
- Status bar — CLI wiring (BUG-006) and bash/zsh generated-code compute/assembly
  consistency (BUG-007, critical fix) for the `left`/`right` zones.
- `asset use --as` — confirmed to already support all 8 asset types correctly; only its
  help text was stale (GAP-004, fixed).
- web-builder's linked-asset-name schema parity (GAP-001, fixed) and VS Code extension's
  schema — both current, both build/typecheck cleanly.
- Bubble Tea non-interactive hang (BUG-005) — fixed; `theme list`, `registry
  list/fetch/search` all confirmed returning promptly with correct output; `registry
  browse` fails clearly instead of hanging.
- **Config mutation — full lifecycle, no defects found** (this session): `config show`/
  `set`/`reset`, safely sandboxed via `$env:USERPROFILE` override (which *does* work for
  this — the file's prior claim that no safe isolation mechanism exists was wrong and is
  corrected here). `set` persists correctly across process invocations; `reset` correctly
  restores defaults; invalid key and invalid value both fail cleanly with exit 1 and a
  specific, correct error message.
- **Wallpaper lifecycle — full lifecycle, no defects found** (this session): `set`/`info`/
  `remove`, safely sandboxed by overriding `LOCALAPPDATA` directly (the actual env var
  `findWindowsTerminalSettings()` reads, confirmed from source — a different, also-safe
  mechanism from config's `USERPROFILE`-based sandbox). `set` correctly wrote all four
  fields (image path, opacity, alignment, stretch mode) into a sandboxed Windows Terminal
  `settings.json`; `remove` correctly cleared them; a missing image path failed cleanly
  with exit 1. The implementation also makes a `.cmdx.bak` backup of `settings.json` before
  ever modifying it — a real, confirmed safety feature.
- **Real bash execution — for the first time in this entire investigation** (this session,
  via WSL/Ubuntu, a genuine bash 5.3.9, not text inspection): cross-compiled a Linux cmdx
  binary, built a real mascot asset, and in a fully sandboxed `$HOME`: `theme inject`
  correctly wrote a working `.bash_profile` with a real, correctly-colored `PS1`; the
  generated `__cmdx_mascot_hook` correctly captured real exit codes (`exit_code=0` after
  `true`, `exit_code=1` after `false`) when manually invoked (simulating what
  `PROMPT_COMMAND` does before each real prompt); 3x repeated `theme inject` produced
  exactly one of each marker block (idempotent); `theme remove` cleanly emptied the profile
  (0 bytes) — all matching what had only ever been verified on PowerShell before.
- **Real audio playback — strong evidence, not just logic** (this session): using the
  sound system's documented custom-player escape hatch (`"player": "ffplay ... %f"`,
  confirmed safe — argv-based, never shell-interpreted) with a real WAV file and `ffplay`
  (confirmed installed on this machine), `sound-play` took 8.89s to return with exit 0 —
  far longer than the 0.5s clip itself, confirming a real blocking subprocess genuinely
  performed audio device I/O rather than silently no-opping. This is the strongest
  playback evidence obtainable without literally hearing it; true audible confirmation
  remains outside what this tooling can observe.

## PARTIAL / documented-as-not-implemented (not bugs)
- Mascot trigger coverage via generated hooks: `exit_code`/`command` confirmed reachable
  (now live-verified in real bash too, not just PowerShell); `output_regex`/`git_status`/
  `idle_time` are correctly implemented in the resolver but not populated by the shipped
  hook generator (GAP-003, evidence-supported intentional boundary). `env_var` already
  works ambiently — not actually gapped.
- Status bar: center-zone rendering (all shells) and right-zone rendering (PowerShell) are
  confirmed *not implemented*, not merely untested (BUG-007's documented remainder).
- `theme info`: works but its asset-slot display doesn't show mascot/floater/status_bar/sound
  even when linked (a known, minor, unfixed display gap).

## NEW FINDINGS, this session (minor, not fixed)
- `theme inject`'s startup banner injects unconditionally regardless of the theme's
  `banner.enabled: false` setting — observed live in the real-bash test (a theme with
  `"enabled": false` still printed its banner text on every new shell). Small, real,
  worth a future fix; not chased further this session given scope.
- In a genuinely non-interactive bash invocation (not a real interactive shell), the
  mascot hook's `command=$(history 1 | ...)` capture is empty, since bash doesn't track
  history non-interactively by default. This is very likely a testing-methodology
  artifact rather than a real product gap — a real interactive terminal session does
  track history automatically — but it was not independently confirmed in a genuinely
  interactive bash session (not possible in this environment), so it's noted rather than
  dismissed outright.
- Confirmed intentional, not a bug: `mascot-state`/`sound-play`-family commands silently
  exit 0 with zero output when `chafa` isn't installed (`cmd/asset_mascot.go` explicitly
  checks `ChafaAvailable()` and exits early) — sensible given this command fires on every
  single shell prompt via the hook; a noisy per-prompt error would be worse UX than silence.

## BLOCKED (environment limitation, not a proven product defect)
- TUI editor (`cmdx edit`) — never launched, to avoid orphaning an uninterruptible process;
  no real TTY available in this environment throughout the entire investigation.
- Interactive `theme create`/`asset create` wizards — not attempted, same reason.
- zsh — confirmed genuinely unavailable in this environment (not installed; only bash is,
  via WSL). Generated zsh code is regression-tested for internal consistency (BUG-007) but
  was never executed live, unlike bash this session.
- True audible confirmation of sound playback (see above — subprocess-level evidence is as
  far as this tooling can go).

## NOT TESTED (no meaningful attempt made)
- Status bar: narrow-terminal behavior, conditional visibility, segment types beyond
  `directory`/`git`/`exit_code`/`time`/`text`.
- Mascot `output_regex`/`git_status`/`idle_time` triggers in a live shell session (only
  exercised via direct `mascot-state` flag injection).
- Font installation. Registry's `fetch`/`browse` against a real, valid (not just
  validation-rejected) theme name — `list`'s real network fetch was exercised; a full
  successful download-and-write via `fetch` was not separately exercised.
- Floater animation (`animate_frames`) — structurally present, covered by `floater_test.go`,
  not exercised live (only static preview was).
- The "everything" maximal integration theme (deferred pending status-bar center/right-zone
  completion).

Do not upgrade any BLOCKED or NOT TESTED item to VERIFIED without an actual observed test.
