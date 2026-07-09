<div align="center">

<pre>
 ██████╗███╗   ███╗██████╗ ██╗  ██╗
██╔════╝████╗ ████║██╔══██╗╚██╗██╔╝
██║     ██╔████╔██║██║  ██║ ╚███╔╝ 
██║     ██║╚██╔╝██║██║  ██║ ██╔██╗ 
╚██████╗██║ ╚═╝ ██║██████╔╝██╔╝ ██╗
 ╚═════╝╚═╝     ╚═╝╚═════╝ ╚═╝  ╚═╝
</pre>

**Break free from the boring terminal.**

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-blue?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Version](https://img.shields.io/badge/Version-1.0.0-purple?style=flat-square)

</div>

---

## What is cmdX?

**cmdX** (`cmd-customizer`) is a developer-focused terminal theming framework. Not just colors — full control over the visual primitives of your terminal: spinners, progress bars, prompts, banners, cursors, borders, PNG assets, reactive mascots, and composable status bars. One JSON file. Every render option overridable at runtime. Every shell. Every OS.

This isn't a "pick a preset and go" tool. Every asset command exposes runtime overrides (render mode, color depth, dimensions, dithering), every asset type supports a `--url` escape hatch for anything not in a curated catalog, and the theme/asset JSON schema is fully documented and IDE-validated — you're meant to bend this to your own setup, not just consume what ships in the box.

---

## Features

**Core theming**
- Full visual control via a single `.json` file — colors, prompt, loader, progress bar, cursor, borders, banner
- Live animated preview before you apply anything
- Instant validation with clear error messages
- Cross-platform shell injection: Windows (PowerShell/CMD), macOS (zsh), Linux (bash/zsh/fish)
- TUI theme editor with spring-physics animations (harmonica)
- Interactive theme creation wizard

**Graphics engine**
- LAB-space color gradients, not flat RGB interpolation
- 5 banner effects: glitch, rainbow, pulse, neon, typewriter
- 7 divider styles, pattern fills, complementary/analogous/triadic color generation

**PNG asset system**
- **Spinners, banners, dividers, icon sets** — your own PNG art rendered via chafa (braille, blocks, ASCII, sixel, or full-color)
- **Floaters** — decorative PNGs anchored to any terminal corner, static or animated
- **Mascots** — a reactive character with a full state machine: define your own states, wire them to exit codes, regex output matching, git status, idle time, or custom shell commands, with priority-based resolution when multiple triggers fire
- **Status bars** — composable segment-based bars (git, directory, time, exit code, battery, kubernetes context, virtualenv, and more) across 3 zones, 8 separator styles including Powerline, rendered as pure shell code — no binary call needed at prompt time
- Runtime override flags on every preview command (`--mode`, `--color`, `--width`, `--height`, `--dither`, `--stretch`, `--threshold`, `--position`, `--state`) — nothing is locked to the manifest
- `cmdx asset render` — test any PNG through the pipeline with zero manifest required
- `cmdx asset create` — interactive scaffolding wizard for any asset type
- **Sound themes** — audio feedback tied to shell events, the one asset type that doesn't use chafa at all: shells out to a platform audio player (afplay/paplay/aplay/ffplay/PowerShell SoundPlayer). Reuses the same trigger vocabulary as mascots (exit codes, output regex, git status, idle time, custom commands), with per-effect volume, cooldown, and async/blocking playback modes. Custom player command support for any format beyond WAV.

**Font management**
- Curated Nerd Fonts catalog (FiraCode, JetBrainsMono, Hack, CascadiaCode, Meslo, and more)
- Cross-platform install (Windows registry registration, macOS/Linux directory placement)
- `--url` install for any font not in the catalog — the catalog is a convenience, not a restriction

**Community & extensibility**
- Plugin API for contributing spinners, banners, prompts
- Community theme registry with search/fetch
- Developer tooling: [VS Code extension](vscode-extension/) (JSON schema validation, visual preview, sidebar browsing) and a [web-based theme builder](web-builder/) (full schema as live controls, no JSON hand-editing required)

---

## Requirements & Installation

cmdX has real dependencies — see **[REQUIREMENTS.md](REQUIREMENTS.md)** for the full breakdown (Go, git, chafa per platform, chafa version compatibility matrix, terminal sixel support table).

### Quick install

**Windows (PowerShell)**
```powershell
powershell -ExecutionPolicy Bypass -File scripts/install.ps1
```

**macOS / Linux**
```bash
curl -fsSL https://raw.githubusercontent.com/abhigyanwebber/cmdX/main/scripts/install.sh | bash
```

### Build from source

```bash
git clone https://github.com/abhigyanwebber/cmdX
cd cmdX
go mod tidy
go mod download
go build -o cmdx ./cmd/
```

Or with the included Makefile:
```bash
make deps    # go mod tidy && go mod download
make build   # compile the binary
make test    # run the full test suite
make lint    # golangci-lint (must be installed separately)
```

### Optional: chafa (for PNG assets)

Required for spinners, banners, dividers, icons, floaters, mascots — anything PNG-based. Without it, text-based theming still works fully.

```bash
winget install hpjansson.chafa     # Windows
brew install chafa                 # macOS
sudo apt install chafa             # Ubuntu/Debian
```

See [REQUIREMENTS.md](REQUIREMENTS.md) for the full per-platform matrix and minimum version requirements per render mode.

---

## Usage

```bash
# Themes
cmdx theme list                        # interactive picker
cmdx theme preview cyberpunk           # live animated preview
cmdx theme apply cyberpunk
cmdx theme info retro
cmdx theme validate ./my-theme.json
cmdx theme create                      # interactive wizard
cmdx theme inject cyberpunk            # persist across terminal restarts

# Registry
cmdx registry list
cmdx registry browse                   # interactive picker
cmdx registry search dark
cmdx registry fetch ocean

# Assets — spinners, banners, dividers, icons, floaters, mascots, status bars
cmdx asset create                      # scaffold a new asset (any type)
cmdx asset list
cmdx asset preview my-floater --mode sixel --width 20
cmdx asset render ./any-image.png --mode braille --width 16   # no manifest needed
cmdx asset use my-mascot --as mascot
cmdx asset mascot-info my-mascot       # full state machine breakdown
cmdx asset mascot-hooks my-mascot      # shell hook code for your profile
cmdx asset statusbar-preview my-bar --shell bash
cmdx asset use my-sounds --as sound
cmdx asset sound-info my-sounds        # every effect and its triggers
cmdx asset sound-hooks my-sounds       # shell hook code for your profile
cmdx asset status                      # what's currently active

# Fonts
cmdx font list
cmdx font search mono
cmdx font install firacode
cmdx font install --url <zip-url> --name my-font   # any font, not just the catalog

# Editor
cmdx edit cyberpunk                    # TUI editor with spring animations

# Config
cmdx config show
```

---

## Built-in Themes

| Theme | Description |
|-------|-------------|
| `default` | Clean and minimal — a solid starting point |
| `cyberpunk` | Neon soaked, magenta and cyan, high energy |
| `minimal` | Pure black and white, zero noise |
| `retro` | Old school green phosphor terminal |

Plus 13 community themes available via `cmdx registry list` — see [cmdX-themes](https://github.com/abhigyanwebber/cmdX-themes).

---

## Creating a Custom Theme

Every theme is a single `.json` file. Drop it in `themes/` and it's available instantly — no recompilation needed. Editing by hand? Install the [VS Code extension](vscode-extension/) for full autocomplete and validation, or use the [web theme builder](web-builder/) for a visual editor with zero JSON required.

```json
{
  "meta": {
    "name": "my-theme",
    "version": "1.0.0",
    "author": "your-name",
    "description": "My custom theme"
  },
  "colors": {
    "primary":    "#FF6B6B",
    "secondary":  "#4ECDC4",
    "background": "#1A1A2E",
    "foreground": "#EAEAEA",
    "accent":     "#FFE66D",
    "error":      "#FF4444",
    "success":    "#44FF88",
    "warning":    "#FFA500",
    "muted":      "#555555"
  },
  "prompt": {
    "symbol": "❯",
    "separator": "›",
    "style": "single",
    "segments": ["directory", "git"],
    "format": "{dir} {git} {symbol} "
  },
  "loader": {
    "frames": ["⠋","⠙","⠹","⠸","⠼","⠴","⠦","⠧","⠇","⠏"],
    "interval_ms": 80,
    "color": "primary"
  },
  "progress_bar": {
    "filled": "█",
    "empty":  "░",
    "width":  40,
    "color":  "primary"
  },
  "cursor": {
    "style": "bar",
    "blink": true,
    "color": "primary"
  },
  "borders": {
    "style": "rounded",
    "chars": {
      "top_left":     "╭",
      "top_right":    "╮",
      "bottom_left":  "╰",
      "bottom_right": "╯",
      "horizontal":   "─",
      "vertical":     "│"
    }
  },
  "banner": {
    "enabled": true,
    "text":    "Welcome, {user}",
    "style":   "plain",
    "color":   "accent"
  },
  "graphics": {
    "gradient": { "enabled": true, "from": "#FF6B6B", "to": "#4ECDC4", "direction": "horizontal" },
    "effects":  { "banner": "none", "prompt": "none" }
  }
}
```

Validate it before using:
```bash
cmdx theme validate ./my-theme.json
```

---

## Theme JSON Reference

### `meta`
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | yes | Theme identifier |
| `version` | string | yes | Semantic version |
| `author` | string | no | Creator name |
| `description` | string | no | Short description |

### `colors`
| Field | Type | Description |
|-------|------|-------------|
| `primary` / `secondary` / `accent` | hex | Core palette |
| `background` / `foreground` | hex | Base terminal colors |
| `error` / `success` / `warning` | hex | Status colors |
| `muted` | hex | Dimmed/secondary text |

### `prompt`
| Field | Type | Options | Description |
|-------|------|---------|-------------|
| `style` | string | `single`, `multiline` | Prompt layout |
| `symbol` | string | any | End character (❯ → $ etc) |
| `segments` | array | `user` `directory` `git` `time` | What shows in prompt |
| `format` | string | | Layout using `{tokens}` |

### `loader` / `progress_bar` / `cursor`
See the full field tables in [REQUIREMENTS.md](REQUIREMENTS.md), or install the VS Code extension for schema-driven autocomplete showing the complete, always-current reference as you type.

### `graphics`
| Field | Type | Options |
|-------|------|---------|
| `gradient.direction` | string | `horizontal`, `vertical` |
| `divider.style` | string | `line` `wave` `dots` `stars` `double` `arrow` `zigzag` |
| `effects.banner` | string | `glitch` `rainbow` `pulse` `neon` `typewriter` `none` |
| `effects.prompt` | string | `rainbow` `none` |

### Asset manifests (`asset.json`)

Each asset type (`spinner`, `banner`, `divider`, `icon`, `floater`, `mascot`, `status-bar`) has its own config block. The fastest way to see the full schema is:

```bash
cmdx asset create   # interactive wizard generates a fully-annotated manifest
```

or install the VS Code extension for live schema validation on any `asset.json`.

---

## Developer Tools

### [VS Code Extension](vscode-extension/)
Live JSON schema validation and autocomplete for theme/asset files, a visual webview preview, sidebar theme/asset browsing, and commands wrapping the real CLI (not a reimplementation). Build with `npm install && npm run compile`.

### [Web Theme Builder](web-builder/)
Browser-based visual editor — every theme field exposed as a real control, live animated preview (all 5 banner effects genuinely implemented, real loader timing), client-side validation mirroring the Go validator, shareable URLs. Build with `npm install && npm run build`.

---

## Project Structure

```
cmdX/
├── cmd/                    # CLI entry point and commands
│   ├── main.go, root.go, helpers.go
│   ├── theme.go, plugin.go, registry.go
│   ├── edit.go, create.go, config.go
│   ├── wallpaper.go, font.go
│   └── asset.go, asset_create.go, asset_mascot.go, asset_statusbar.go
├── internal/
│   ├── config/             # JSON loader, validator, centralized ResolveColor
│   ├── theme/               # Theme manager and renderer
│   ├── primitives/          # Spinner, progress bar, banner
│   ├── preview/              # Live preview engine (orchestrator + sections)
│   ├── graphics/             # Gradients, effects, dividers, patterns
│   ├── render/                # Glamour markdown rendering
│   ├── shells/                # Shell injection + sanitization
│   │   ├── powershell/, zsh/, bash/
│   ├── plugin/                 # Plugin API and manager
│   ├── registry/               # Community theme registry
│   ├── editor/                  # TUI theme editor (harmonica springs)
│   ├── tui/                      # Bubbles extended components
│   ├── wallpaper/                 # Terminal wallpaper support
│   ├── icons/                      # Icon sets (Nerd Fonts, Emoji, ASCII)
│   ├── fonts/                       # Font catalog, install, cross-platform registry
│   └── assets/                       # PNG asset engine
│       ├── types.go, loader.go, manager.go, chafa.go, override.go
│       ├── mascot.go, statusbar.go
├── vscode-extension/        # VS Code extension (TypeScript)
├── web-builder/             # Web-based theme builder (Vite + TypeScript)
├── assets/                  # User asset library
├── themes/                  # Built-in theme JSON files
├── plugins/                 # Plugin directory
├── scripts/                 # Install scripts
├── REQUIREMENTS.md          # Full dependency & platform documentation
├── Makefile                 # make deps/build/test/vet/lint/install
└── .llm/                    # Session continuity docs (todo.md, HANDOFF.md)
```

## Roadmap

- [x] Theme engine, live preview, validator, install scripts
- [x] Shell injection (PowerShell, zsh, bash) with sanitization
- [x] Graphics engine (gradients, effects, dividers, patterns)
- [x] Plugin API, community theme registry, TUI editor
- [x] Wallpaper support, icon sets
- [x] PNG asset system (spinners, banners, dividers, icons)
- [x] Full test suite across all core packages
- [x] Security audit (shell injection sanitization, path traversal fixes, download caps)
- [x] Exotic assets: floaters, mascots (full state machine), status bars
- [x] Font installer (curated catalog + custom URL support)
- [x] VS Code extension
- [x] Web-based theme builder
- [x] Sound themes

---

## Contributing

Contributions are welcome — themes, assets (floaters/mascots/status bars are especially fun to make), plugins, or code.

1. Fork the repo
2. Create your theme in `themes/` or asset in `assets/` — `cmdx asset create` scaffolds the manifest for you
3. Validate: `cmdx theme validate ./themes/your-theme.json` or `cmdx asset validate ./your-asset/`
4. Run the test suite: `make test` (or `go test ./...`)
5. Open a pull request

---

## License

MIT — do whatever you want with it.

---

<div align="center">
Built by <a href="https://github.com/abhigyanwebber">abhigyanwebber</a>
</div>
