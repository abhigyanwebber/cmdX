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

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-blue?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Version](https://img.shields.io/badge/Version-1.0.0-purple?style=flat-square)

</div>

---

## What is cmdx?

**cmd-customizer** (`cmdx`) is a developer-focused terminal theming framework.
Not just colors — full control over the visual primitives of your terminal:
spinners, progress bars, prompts, banners, cursors, borders.
One JSON file. Every shell. Every OS.

---

## Features

- **Theme engine** — full visual control via a single `.json` file
- **Live preview** — see any theme in action before applying it
- **Built-in themes** — cyberpunk, minimal, retro, default
- **Custom themes** — write your own and share them
- **Validation** — instant feedback if your theme JSON has errors
- **Cross-platform** — Windows (CMD + PowerShell), macOS (zsh), Linux (bash/zsh/fish)
- **Zero runtime** — single binary, no dependencies to install
- **PNG asset system** — upload your own graphics as spinners, icons, banners and dividers
- **Sixel rendering** — clean image-quality graphics inside your terminal via chafa

---

## Installation

### Windows (PowerShell)
```powershell
powershell -ExecutionPolicy Bypass -File scripts/install.ps1
```

### macOS / Linux
```bash
curl -fsSL https://raw.githubusercontent.com/abhigyanwebber/cmdX/main/scripts/install.sh | bash
```

### Manual (all platforms)
```bash
git clone https://github.com/abhigyanwebber/cmdX
cd cmdX
go build -o cmdx ./cmd/
```

---

## Usage

```bash
# List all available themes
cmdx theme list

# Preview a theme live before applying
cmdx theme preview cyberpunk

# Apply a theme
cmdx theme apply cyberpunk

# Get detailed info about a theme
cmdx theme info retro

# Validate your own custom theme file
cmdx theme validate ./my-theme.json
```

---

## Built-in Themes

| Theme | Description |
|-------|-------------|
| `default` | Clean and minimal — a solid starting point |
| `cyberpunk` | Neon soaked, magenta and cyan, high energy |
| `minimal` | Pure black and white, zero noise |
| `retro` | Old school green phosphor terminal |

---

## Creating a Custom Theme

Every theme is a single `.json` file. Drop it in the `themes/` folder and it
becomes available instantly — no recompilation needed.

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
| `primary` | hex | Main accent color |
| `secondary` | hex | Supporting color |
| `background` | hex | Terminal background |
| `foreground` | hex | Default text color |
| `accent` | hex | Highlight color |
| `error` | hex | Error messages |
| `success` | hex | Success messages |
| `warning` | hex | Warning messages |
| `muted` | hex | Dimmed/secondary text |

### `prompt`
| Field | Type | Options | Description |
|-------|------|---------|-------------|
| `style` | string | `single`, `multiline` | Prompt layout |
| `symbol` | string | any | End character (❯ → $ etc) |
| `segments` | array | `user` `directory` `git` `time` | What shows in prompt |
| `format` | string | | Layout using `{tokens}` |

### `loader`
| Field | Type | Description |
|-------|------|-------------|
| `frames` | array | Animation frames in order |
| `interval_ms` | int | Speed in milliseconds |
| `color` | string | Color key or hex value |

### `cursor`
| Field | Type | Options |
|-------|------|---------|
| `style` | string | `block` `bar` `underline` |
| `blink` | bool | `true` / `false` |
| `color` | string | color key or hex |

---

## Project Structure

```
cmdX/
├── cmd/                    # CLI entry point and commands
│   ├── main.go
│   ├── root.go
│   ├── theme.go
│   ├── plugin.go
│   ├── registry.go
│   ├── edit.go
│   ├── wallpaper.go
│   └── asset.go
├── internal/
│   ├── config/             # JSON loader and validator
│   ├── theme/              # Theme manager and renderer
│   ├── primitives/         # Spinner, progress bar, banner
│   ├── preview/            # Live preview engine
│   ├── graphics/           # Gradients, effects, dividers, patterns
│   ├── shells/             # Shell injection (PowerShell, zsh, bash)
│   │   ├── powershell/
│   │   ├── zsh/
│   │   └── bash/
│   ├── plugin/             # Plugin API and manager
│   ├── registry/           # Community theme registry
│   ├── editor/             # TUI theme editor
│   ├── wallpaper/          # Terminal wallpaper support
│   ├── icons/              # Icon sets (Nerd Fonts, Emoji, ASCII)
│   └── assets/             # PNG asset engine (chafa rendering)
│       ├── types.go
│       ├── chafa.go
│       ├── loader.go
│       └── manager.go
├── assets/                 # User asset library
│   ├── spinners/           # PNG spinner assets
│   │   └── pulse/          # Example spinner
│   ├── icons/              # PNG icon assets
│   ├── banners/            # PNG banner assets
│   └── dividers/           # PNG divider assets
├── themes/                 # Built-in theme JSON files
├── plugins/                # Plugin directory
│   └── example-plugin/
├── scripts/                # Install scripts
│   ├── install.sh
│   └── install.ps1
└── docs/                   # Documentation
```

## Roadmap

- [x] Theme engine with JSON config
- [x] Built-in themes (cyberpunk, minimal, retro, default)
- [x] Live preview command
- [x] Theme validator
- [x] Install scripts (Windows + Unix)
- [x] Shell injection (PowerShell, zsh, bash)
- [x] Graphics engine (gradients, effects, dividers, patterns)
- [x] Plugin API
- [x] Community theme registry
- [x] TUI theme editor
- [x] Wallpaper support (Windows Terminal, iTerm2, Kitty)
- [x] Icon sets (Nerd Fonts, Emoji, ASCII)
- [x] PNG asset system (spinners, icons, banners, dividers)
- [x] Sixel/braille/block rendering via chafa
- [ ] Sound themes
- [ ] Font installer
- [ ] Web-based theme builder
- [ ] VS Code extension

---

## Contributing

Contributions are welcome — especially new themes.

1. Fork the repo
2. Create your theme JSON in `themes/`
3. Run `cmdx theme validate ./themes/your-theme.json`
4. Open a pull request

---

## License

MIT — do whatever you want with it.

---

<div align="center">
Built by <a href="https://github.com/abhigyanwebber">abhigyanwebber</a>
</div>