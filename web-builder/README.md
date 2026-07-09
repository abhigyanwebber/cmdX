# cmdX Web Theme Builder

A browser-based visual editor for cmdX `theme.json` files — every field the schema supports, exposed as a real control (color pickers, dropdowns, checkboxes), with a live terminal mockup that updates as you edit. No hand-written JSON required, but nothing is hidden or simplified either: this exposes the full schema, not a "starter" subset.

## Requirements

- **Node.js 18+** and **npm** — same as the [VS Code extension](../vscode-extension), this is a real dependency-managed project, not a single dropped-in HTML file.

## Getting started

```bash
cd web-builder
npm install
npm run dev
```

Opens a dev server (Vite) with hot reload. Edit `src/*.ts` and the browser updates instantly.

## Building for deployment

```bash
npm run build
```

Outputs a fully static site to `dist/` — no backend, no server-side code. Deploy it anywhere that serves static files: GitHub Pages, Netlify, Vercel, S3, or just open `dist/index.html` via any static file server.

```bash
npm run preview   # serve the production build locally to sanity-check it
```

## Features

- **Full schema coverage** — every field from `internal/config/types.go`'s `Theme` struct: meta, colors, prompt, loader, progress bar, cursor, borders, banner, gradient, effects, divider, wallpaper, the full Nerd Font icon glyph set, and linked asset names.
- **Live preview** — color swatches, gradient bar, an actual terminal mockup with your configured prompt/banner/loader/progress bar, including working animations for all five banner effects (glitch, rainbow, pulse, neon, typewriter) and the loader cycling through your real configured frames at your real configured speed.
- **Client-side validation** — mirrors `internal/config/validator.go`'s rules (hex color format, required fields, positive numbers) so what passes here also passes `cmdx theme validate`.
- **Built-in presets** — cyberpunk, ocean, minimal, forest — real starting points, not blank forms.
- **Import** — load an existing `theme.json` to continue editing it.
- **Export** — download as `<name>.json`, copy to clipboard, or generate a shareable URL (the theme is base64-encoded into the URL hash — no server, no database, the link *is* the theme).

## Project structure

```
src/
  theme.ts        — TypeScript types matching the Go schema exactly, + defaultTheme()
  validate.ts      — client-side mirror of Go's ValidateTheme
  presets.ts       — built-in starting themes
  formPanel.ts     — renders every field as an editable control
  preview.ts        — live terminal mockup renderer + animations
  exportImport.ts  — download/clipboard/share-URL/file-import logic
  main.ts          — wires everything together
  style.css        — dark UI styling
```

## Why no framework?

Deliberately vanilla TypeScript + Vite, no React/Vue/etc. Keeps the dependency footprint small and the code approachable for anyone who wants to fork and extend it — you don't need to learn a framework's conventions to add a field or change the preview, just read `formPanel.ts` and `preview.ts` directly.
