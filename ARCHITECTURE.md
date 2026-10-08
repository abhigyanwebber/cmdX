# cmdX — Architecture (Verified)

## Theme → Asset → Shell Hooks → Runtime

```
theme.json (assets: {spinner, banner, divider, icons, mascot, floater, status_bar, sound})
    ↓  cmdx theme apply/inject
activateThemeAssets()  (cmd/helpers.go)
    ↓
activateAsset() per linked slot → writes internal/assets/.state/<type>.txt
    +
syncAssetHooks() → installs/removes a SEPARATE shell-profile block for mascot/sound hooks
    ↓
shell profile ends up with TWO independent marker-delimited blocks:
    # cmdx theme start ... # cmdx theme end          (InjectStart/InjectEnd)
    # cmdx asset hooks start ... # cmdx asset hooks end   (AssetHooksStart/AssetHooksEnd)
```
Verified: these two blocks are written/removed independently. Switching to a theme with no
mascot/sound linked cleanly removes the asset-hooks block while leaving the theme block
intact. `--no-assets` suppresses `activateAsset()`; `--no-hooks` suppresses `syncAssetHooks()`
only. Repeated `apply`/`inject` is idempotent (no duplicate blocks, verified 3x).

## Mascot Resolution — verified correct core, one confirmed CLI defect

```
MascotContext{ LastExitCode, LastCommand, LastOutput, IdleSeconds, GitStatus, Env }
    ↓
ResolveState(mascotConfig, ctx)   [internal/assets/mascot.go]
    — iterates all states' triggers, matchesTrigger() each, keeps highest-Priority match
    — falls back to any `always` trigger, then DefaultState, then "idle"
    — CORRECT AND DETERMINISTIC given triggers with distinct priorities (verified 8/8 runs)
    ↓
RenderMascotState() → frames, transition, interval
```

Two entry points both ultimately call `ResolveState` on the same `ctx`:
- `ResolveMascotState(name, ctx)` — thin wrapper, used by `mascot-state`'s override-skip
  branch and available for future callers that want a bare answer without rendering.
- `PreviewMascot(name, ctx, overrides)` — renders AND independently re-calls
  `ResolveState(a.Mascot, ctx)` itself; it does **not** accept a pre-resolved state.

**Defect (BUG-002, fixed this session — see KNOWN_ISSUES.md):** `assetMascotStateCmd`
computed a `state` value (honoring `--state`) and then discarded it (`_ = state`), calling
`PreviewMascot(name, ctx, overrides)` which re-resolves from `ctx` alone — `ctx` has no field
to carry a forced state, so `--state` never reached the display path. This did **not** affect
real shell hooks (they never pass `--state`), only the manual/debugging override affordance.

**Non-bug (was misdiagnosed mid-investigation, corrected):** apparent non-determinism in
priority resolution was caused by a test manifest placing `"priority"` at the wrong JSON
nesting level (state level instead of inside each trigger object), leaving all priorities at
the zero-value default and creating genuine ties broken by Go's randomized map iteration
order. The resolver itself has no bug here — but tied/default priorities are a real
authoring hazard worth documenting (see GAP-002).

## Sound — separate, simpler, no equivalent defect
`sound-play --sound X` takes an early-return branch: direct `a.Sound.Sounds[X]` lookup →
`assets.PlaySound()` → `return`. It never falls through to the ctx-based `PreviewSound()`
resolution, so it does not share mascot's discard bug.

## Shell Hook Data Flow (as actually generated)
```
PROMPT_COMMAND / prompt function
    → captures $?/exit code and last command text
    → cmdx asset mascot-state "<name>" --exit-code "$code" --command "$cmd"
```
Confirmed: generated hooks (bash inspected directly, PowerShell live-tested) populate only
`exit_code` and `command`. `output_regex`, `git_status`, `idle_time`, `env_var` triggers are
correctly implemented in the resolver but are **not wired up by the generated hooks** — see
BUG-003. `env_var` is partially live via `buildEnvSnapshot()` reading real OS env vars
(ambient, not hook-supplied), so it does work for real ambient environment variables, just
not for anything a hook would need to compute per-command.
