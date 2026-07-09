# cmdX Theme Tools for VS Code

Edit, preview, and manage [cmdX](https://github.com/abhigyanwebber/cmdX) terminal themes and assets directly in VS Code.

This extension doesn't reimplement cmdX's logic — it wraps the real `cmdx` CLI and gives you VS Code-native tooling on top of it: live JSON schema validation while you type, a visual preview panel, and sidebar browsing for your themes and assets. Every command that involves actual terminal rendering (asset preview, the interactive theme creator) opens a real integrated terminal running the real `cmdx` binary, rather than faking the output in a webview.

## Requirements

- **The `cmdx` CLI** must be installed and either on your `PATH` or pointed to via the `cmdx.executablePath` setting. See [REQUIREMENTS.md](https://github.com/abhigyanwebber/cmdX/blob/main/REQUIREMENTS.md) in the main repo for install instructions (Go, chafa, etc.).
- VS Code 1.85 or later.

This extension does not bundle or reimplement any part of the cmdX Go codebase — it is a thin client around the CLI.

## Features

### Live JSON schema validation
Open any theme file under a `themes/` folder, or any `asset.json`, and get real autocomplete, hover docs, and validation errors matching cmdX's exact schema — no more guessing field names or trailing-comma hunting.

### Visual theme preview
Run **cmdX: Preview Theme** on an open theme JSON file to see a live rendering of its color palette, gradient, prompt mockup, banner, and progress bar in a side panel. Enable `cmdx.autoPreviewOnSave` to have it refresh automatically as you edit.

### Sidebar browsing
The cmdX activity bar icon opens two views:
- **Themes** — every theme in your resolved themes directory, click to open, inline preview/apply actions
- **Assets** — every spinner, banner, divider, icon, floater, mascot, and status bar asset, click to open its manifest

### Full CLI depth, not a simplified subset
- **cmdX: Preview Asset** opens a real terminal with `cmdx asset preview`, including a render-mode override prompt — the same flags the CLI exposes (`--mode`, `--color`, etc.), not a locked-down GUI
- **cmdX: Create New Theme** runs the real interactive `cmdx theme create` wizard in a terminal
- **cmdX: Install Font** opens the font catalog via `cmdx font list`

### Commands
| Command | What it does |
|---|---|
| `cmdX: Preview Theme` | Visual preview panel for the active theme file |
| `cmdX: Apply Theme` | Runs `cmdx theme apply` |
| `cmdX: Validate Theme` | Runs `cmdx theme validate` |
| `cmdX: Validate Asset` | Runs `cmdx asset validate` on the containing folder |
| `cmdX: Create New Theme...` | Opens the interactive theme wizard in a terminal |
| `cmdX: Preview Asset (with overrides)...` | Quick-pick asset + render mode, runs in a terminal |
| `cmdX: Install Font...` | Opens the font catalog |
| `cmdX: Set cmdx Executable Path...` | Configure the binary location |

## Settings

| Setting | Default | Description |
|---|---|---|
| `cmdx.executablePath` | `"cmdx"` | Path to the cmdx binary. Set an absolute path if it's not on PATH. |
| `cmdx.themesDirectory` | `""` | Explicit themes directory. Empty resolves to `<workspace>/themes`. |
| `cmdx.assetsDirectory` | `""` | Explicit assets directory. Empty resolves to `<workspace>/assets`. |
| `cmdx.autoPreviewOnSave` | `false` | Refresh the preview panel automatically on save. |
| `cmdx.defaultRenderMode` | `""` | Default render mode for asset preview overrides. |

`cmdx.themesDirectory` / `cmdx.assetsDirectory` are passed to the CLI as `CMDX_THEMES_DIR` / `CMDX_ASSETS_DIR` environment variables, so the extension and the binary always resolve the same directories.

## Building from source

```bash
git clone https://github.com/abhigyanwebber/cmdX
cd cmdX/vscode-extension   # or wherever this extension lives in the repo
npm install
npm run compile
```

To package a `.vsix` for local install:

```bash
npm run package
code --install-extension cmdx-vscode-*.vsix
```

To develop with live reload, open this folder in VS Code and press F5 to launch an Extension Development Host.

## Development requirements

- Node.js 18+ and npm
- TypeScript (installed via `npm install`, no global install needed)
- [`@vscode/vsce`](https://github.com/microsoft/vscode-vsce) for packaging (installed via `npm install`)

## License

MIT — see [LICENSE](./LICENSE).
