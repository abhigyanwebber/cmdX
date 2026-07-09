# cmdX — Requirements & Installation

## Required Dependencies

### Go (≥ 1.22)
All core functionality. Required to build cmdX from source.

| Platform | Install |
|----------|---------|
| Windows  | `winget install GoLang.Go` |
| macOS    | `brew install go` |
| Linux    | `sudo apt install golang-go` or [go.dev/dl](https://go.dev/dl) |

Verify: `go version`

### Git
Required for `cmdx registry fetch` and the install script.

| Platform | Install |
|----------|---------|
| Windows  | `winget install Git.Git` |
| macOS    | `brew install git` (or Xcode CLT) |
| Linux    | `sudo apt install git` |

Verify: `git version`

---

## Optional Dependencies

### chafa (PNG-to-terminal rendering)
Required for **all PNG asset commands** (`cmdx asset preview`, `cmdx asset render`, floaters, spinners, banners, dividers, icons). Without chafa, text-based theme features still work — only PNG asset rendering is unavailable.

| Platform | Install |
|----------|---------|
| Windows  | `winget install hpjansson.chafa` |
| macOS    | `brew install chafa` |
| Ubuntu   | `sudo apt install chafa` |
| Fedora   | `sudo dnf install chafa` |
| Arch     | `sudo pacman -S chafa` |

Check install: `cmdx asset chafa`

Minimum version: **chafa 1.10** (for braille and sixel support).  
Recommended: **chafa 1.14+** (best sixel quality and symbol coverage).

**Render modes by chafa version:**
| Mode | Minimum chafa | Notes |
|------|---------------|-------|
| blocks | 1.0+ | Always works |
| ascii | 1.0+ | Always works |
| braille | 1.6+ | Fine-detail rendering |
| sixel | 1.10+ | Pixel-accurate (terminal must support sixel) |
| color | 1.0+ | 24-bit block rendering |

**Terminal sixel support:**
- Windows Terminal (1.18+) — yes
- iTerm2 (macOS) — yes
- Kitty — yes (use `--format=kitty` instead)
- standard xterm/gnome-terminal — no

---

### Audio players (for sound themes)

Sound theme assets shell out to a platform audio player — unlike every other asset type, no chafa involved. **WAV** is the only format guaranteed to work with zero extra installs; other formats need a custom player configured via a sound asset's `player` field (see below).

| Platform | Default player | Extra install needed? |
|----------|----------------|------------------------|
| Windows  | PowerShell's built-in `System.Media.SoundPlayer` | No — always available, WAV only |
| macOS    | `afplay` | No — built into macOS |
| Linux    | `paplay` → `aplay` → `ffplay` (tried in that order) | Usually no — `paplay`/`aplay` ship with most desktop distros (PulseAudio/ALSA). Install `ffmpeg` for the `ffplay` fallback if neither is present, or for broader format support |

For MP3/OGG/etc. on any platform, or a completely custom playback setup, set a `player` template in the sound asset's manifest:
```json
{
  "sound": {
    "player": "ffplay -nodisp -autoexit -loglevel quiet %f"
  }
}
```
`%f` is replaced with the sound file's path. Requires `ffmpeg` (for `ffplay`) or whatever player binary you point it at.

```bash
# Linux — install ffmpeg for the ffplay fallback / broader format support
sudo apt install ffmpeg      # Debian/Ubuntu
sudo dnf install ffmpeg      # Fedora
```

---

## Building from source

```bash
git clone https://github.com/abhigyanwebber/cmdX
cd cmdX
go mod tidy
go build -o cmdx ./cmd/

# Linux/macOS — install to PATH
cp cmdx ~/.local/bin/

# Windows — add to PATH manually or use install.ps1
.\scripts\install.ps1
```

---

## For Asset Developers

If you're building custom assets (floaters, spinners, banners, dividers, icons):

1. **chafa** (required — see above)
2. **Any PNG editor** — Aseprite, GIMP, Photoshop, etc.
3. For animated assets: a tool that can export individual PNG frames (Aseprite does this natively)

Quick start:
```bash
# Scaffold a new asset interactively
cmdx asset create

# Test any PNG through the render pipeline before creating a manifest
cmdx asset render ./my-image.png --mode braille --width 16

# Try different render modes to find what looks best
cmdx asset render ./my-image.png --mode sixel --color truecolor
cmdx asset render ./my-image.png --mode blocks --width 20 --dither
cmdx asset render ./my-image.png --mode ascii --color none

# Once happy, import and use
cmdx asset import ./my-asset/
cmdx asset use my-asset --as floater --position top-right
```

---

## Verifying your setup

```bash
go version          # should print go1.22+
git version         # any recent version
cmdx asset chafa    # should print chafa version if installed
cmdx theme list     # should show built-in themes
```
