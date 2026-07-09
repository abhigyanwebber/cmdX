import * as vscode from 'vscode';

interface ThemeColors {
  primary?: string;
  secondary?: string;
  background?: string;
  foreground?: string;
  accent?: string;
  error?: string;
  success?: string;
  warning?: string;
  muted?: string;
}

interface ThemeJson {
  meta?: { name?: string; author?: string; description?: string; version?: string };
  colors?: ThemeColors;
  prompt?: { symbol?: string; format?: string; style?: string };
  loader?: { frames?: string[]; interval_ms?: number };
  progress_bar?: { filled?: string; empty?: string; width?: number };
  banner?: { enabled?: boolean; text?: string };
  graphics?: {
    gradient?: { enabled?: boolean; from?: string; to?: string };
    effects?: { banner?: string; prompt?: string };
  };
}

/**
 * ThemePreviewPanel owns the single webview panel used to render a live
 * visual preview of a theme JSON document. Reuses one panel across
 * calls (VS Code convention for "preview" style panels) rather than
 * spawning a new tab per preview.
 */
export class ThemePreviewPanel {
  public static currentPanel: ThemePreviewPanel | undefined;
  private readonly panel: vscode.WebviewPanel;
  private disposables: vscode.Disposable[] = [];

  public static createOrShow(extensionUri: vscode.Uri, theme: ThemeJson, sourcePath?: string) {
    const column = vscode.window.activeTextEditor?.viewColumn;

    if (ThemePreviewPanel.currentPanel) {
      ThemePreviewPanel.currentPanel.panel.reveal(column);
      ThemePreviewPanel.currentPanel.update(theme, sourcePath);
      return;
    }

    const panel = vscode.window.createWebviewPanel(
      'cmdxThemePreview',
      'cmdX Theme Preview',
      column ?? vscode.ViewColumn.Beside,
      {
        enableScripts: true,
        localResourceRoots: [vscode.Uri.joinPath(extensionUri, 'media')],
        retainContextWhenHidden: true,
      }
    );

    ThemePreviewPanel.currentPanel = new ThemePreviewPanel(panel, theme, sourcePath);
  }

  private constructor(panel: vscode.WebviewPanel, theme: ThemeJson, sourcePath?: string) {
    this.panel = panel;
    this.update(theme, sourcePath);

    this.panel.onDidDispose(() => this.dispose(), null, this.disposables);
  }

  /** Re-renders the panel with a new theme document — called when the
   * source file changes (if cmdx.autoPreviewOnSave is enabled) or when
   * the user re-runs the preview command. */
  public update(theme: ThemeJson, sourcePath?: string) {
    this.panel.title = theme.meta?.name ? `Preview: ${theme.meta.name}` : 'cmdX Theme Preview';
    this.panel.webview.html = this.render(theme, sourcePath);
  }

  private dispose() {
    ThemePreviewPanel.currentPanel = undefined;
    this.panel.dispose();
    while (this.disposables.length) {
      const d = this.disposables.pop();
      d?.dispose();
    }
  }

  private render(theme: ThemeJson, sourcePath?: string): string {
    const c = theme.colors ?? {};
    const primary = c.primary ?? '#888888';
    const secondary = c.secondary ?? '#888888';
    const accent = c.accent ?? '#888888';
    const background = c.background ?? '#0D0D0D';
    const foreground = c.foreground ?? '#FFFFFF';
    const error = c.error ?? '#FF4444';
    const success = c.success ?? '#00FF88';
    const warning = c.warning ?? '#FFA500';
    const muted = c.muted ?? '#666666';

    const promptSymbol = escapeHtml(theme.prompt?.symbol ?? '▶');
    const bannerText = escapeHtml(theme.banner?.text?.replace('{user}', 'developer') ?? '');
    const bannerEnabled = theme.banner?.enabled ?? false;
    const gradientEnabled = theme.graphics?.gradient?.enabled ?? false;
    const gradientFrom = theme.graphics?.gradient?.from ?? primary;
    const gradientTo = theme.graphics?.gradient?.to ?? secondary;

    const progressFilled = theme.progress_bar?.width ? Math.floor(theme.progress_bar.width * 0.65) : 13;
    const progressEmpty = (theme.progress_bar?.width ?? 20) - progressFilled;

    const swatches = [
      ['primary', primary],
      ['secondary', secondary],
      ['accent', accent],
      ['success', success],
      ['error', error],
      ['warning', warning],
      ['muted', muted],
    ]
      .map(
        ([label, hex]) => `
        <div class="swatch">
          <div class="swatch-color" style="background:${hex}"></div>
          <div class="swatch-label">${label}<br><span class="swatch-hex">${hex}</span></div>
        </div>`
      )
      .join('');

    return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<style>
  :root {
    --primary: ${primary};
    --secondary: ${secondary};
    --accent: ${accent};
    --background: ${background};
    --foreground: ${foreground};
    --error: ${error};
    --success: ${success};
    --warning: ${warning};
    --muted: ${muted};
  }
  body {
    font-family: 'Cascadia Code', 'Consolas', monospace;
    background: var(--vscode-editor-background);
    color: var(--vscode-editor-foreground);
    padding: 20px;
    line-height: 1.6;
  }
  h1 { font-size: 1.1em; margin-bottom: 4px; }
  .meta { color: var(--muted); font-size: 0.85em; margin-bottom: 24px; }
  .section-title {
    font-size: 0.75em;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--accent);
    margin: 24px 0 10px 0;
    border-bottom: 1px solid var(--muted);
    padding-bottom: 4px;
  }
  .terminal-mockup {
    background: var(--background);
    color: var(--foreground);
    border-radius: 6px;
    padding: 16px;
    font-size: 0.95em;
  }
  .prompt-line { display: flex; align-items: center; gap: 8px; }
  .prompt-dir { color: var(--accent); font-weight: bold; }
  .prompt-git { color: var(--muted); }
  .prompt-symbol { color: var(--primary); font-weight: bold; }
  .banner-text {
    color: var(--primary);
    font-weight: bold;
    padding: 8px 0;
  }
  .banner-disabled { color: var(--muted); font-style: italic; font-size: 0.85em; }
  .swatch-row { display: flex; flex-wrap: wrap; gap: 12px; }
  .swatch { text-align: center; }
  .swatch-color {
    width: 48px; height: 48px; border-radius: 8px;
    border: 1px solid rgba(255,255,255,0.1);
    margin-bottom: 4px;
  }
  .swatch-label { font-size: 0.7em; color: var(--muted); }
  .swatch-hex { font-size: 0.65em; }
  .progress-bar { font-family: monospace; }
  .progress-filled { color: var(--primary); }
  .progress-empty { color: var(--muted); }
  .gradient-demo {
    height: 24px; border-radius: 4px;
    background: linear-gradient(90deg, ${gradientFrom}, ${gradientTo});
  }
  .status-row { display: flex; gap: 16px; margin-top: 8px; font-size: 0.85em; }
  .status-success { color: var(--success); }
  .status-error { color: var(--error); }
  .status-warning { color: var(--warning); }
  .source-path { color: var(--muted); font-size: 0.7em; margin-top: 32px; word-break: break-all; }
</style>
</head>
<body>
  <h1>${escapeHtml(theme.meta?.name ?? 'Untitled Theme')}</h1>
  <div class="meta">
    ${theme.meta?.author ? `by ${escapeHtml(theme.meta.author)} · ` : ''}${theme.meta?.version ?? ''}<br>
    ${escapeHtml(theme.meta?.description ?? '')}
  </div>

  <div class="section-title">Color Palette</div>
  <div class="swatch-row">${swatches}</div>

  ${
    gradientEnabled
      ? `<div class="section-title">Gradient</div><div class="gradient-demo"></div>`
      : ''
  }

  <div class="section-title">Terminal Mockup</div>
  <div class="terminal-mockup">
    ${
      bannerEnabled
        ? `<div class="banner-text">${bannerText || '(no banner text set)'}</div>`
        : `<div class="banner-disabled">banner disabled</div>`
    }
    <div class="prompt-line">
      <span class="prompt-dir">~/projects/my-app</span>
      <span class="prompt-git">(main)</span>
      <span class="prompt-symbol">${promptSymbol}</span>
    </div>
    <div class="progress-bar">
      [<span class="progress-filled">${'█'.repeat(Math.max(progressFilled, 0))}</span><span class="progress-empty">${'░'.repeat(
      Math.max(progressEmpty, 0)
    )}</span>] 65%
    </div>
    <div class="status-row">
      <span class="status-success">✓ Build successful</span>
      <span class="status-error">✗ Connection refused</span>
      <span class="status-warning">⚠ Deprecated API</span>
    </div>
  </div>

  ${sourcePath ? `<div class="source-path">${escapeHtml(sourcePath)}</div>` : ''}
</body>
</html>`;
  }
}

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
