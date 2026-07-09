# Changelog

## 0.1.0 — Initial release

- JSON schemas for `theme.json` and `asset.json` with full autocomplete/validation matching cmdX's Go structs exactly (all 12 theme sections, all 7 asset types including mascots' state-machine triggers and status bar segments)
- Visual theme preview webview panel: color palette, gradient, prompt mockup, banner, progress bar
- Sidebar Themes and Assets tree views with inline preview/apply actions
- Commands: Preview Theme, Apply Theme, Validate Theme, Validate Asset, Create New Theme, Preview Asset (with render overrides), Install Font, Set Executable Path
- Configurable `cmdx` binary path, themes/assets directory overrides (passed through as `CMDX_THEMES_DIR`/`CMDX_ASSETS_DIR`)
- Optional auto-refresh preview on save
