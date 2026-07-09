import type { Theme } from './theme';

/**
 * Preview owns a single DOM container and re-renders a live terminal
 * mockup from a Theme object: color palette, gradient bar, animated
 * banner (approximating each graphics.effects.banner style), animated
 * loader cycling through the real configured frames at the real
 * interval_ms, progress bar, cursor blink, and a border box sample.
 *
 * Animation timers are tracked and cleared on every re-render so
 * repeated theme edits never leak intervals.
 */
export class Preview {
  private container: HTMLElement;
  private timers: number[] = [];
  private typewriterIndex = 0;

  constructor(container: HTMLElement) {
    this.container = container;
  }

  render(theme: Theme): void {
    this.clearTimers();
    this.typewriterIndex = 0;

    const c = theme.colors;
    this.container.style.setProperty('--primary', c.primary);
    this.container.style.setProperty('--secondary', c.secondary);
    this.container.style.setProperty('--accent', c.accent);
    this.container.style.setProperty('--background', c.background);
    this.container.style.setProperty('--foreground', c.foreground);
    this.container.style.setProperty('--error', c.error);
    this.container.style.setProperty('--success', c.success);
    this.container.style.setProperty('--warning', c.warning);
    this.container.style.setProperty('--muted', c.muted);

    this.container.innerHTML = `
      <div class="preview-swatches">${this.renderSwatches(theme)}</div>
      ${theme.graphics.gradient.enabled ? this.renderGradient(theme) : ''}
      <div class="preview-terminal">
        <div class="preview-banner" id="pv-banner"></div>
        <div class="preview-prompt">
          <span class="pv-dir">~/projects/my-app</span>
          <span class="pv-git">(main)</span>
          <span class="pv-symbol" style="color:${c.primary}">${escapeHtml(theme.prompt.symbol)}</span>
          <span class="pv-cursor pv-cursor-${theme.cursor.style}${theme.cursor.blink ? ' pv-blink' : ''}" style="background:${theme.cursor.color || c.primary}"></span>
        </div>
        <div class="preview-loader" id="pv-loader"></div>
        <div class="preview-progress" id="pv-progress"></div>
        <div class="preview-status">
          <span style="color:${c.success}">✓ Build successful</span>
          <span style="color:${c.error}">✗ Connection refused</span>
          <span style="color:${c.warning}">⚠ Deprecated API</span>
        </div>
        <div class="preview-border">${this.renderBorder(theme)}</div>
      </div>
    `;

    this.animateBanner(theme);
    this.animateLoader(theme);
    this.renderProgress(theme);
  }

  private renderSwatches(theme: Theme): string {
    const entries: [string, string][] = [
      ['primary', theme.colors.primary],
      ['secondary', theme.colors.secondary],
      ['accent', theme.colors.accent],
      ['success', theme.colors.success],
      ['error', theme.colors.error],
      ['warning', theme.colors.warning],
      ['muted', theme.colors.muted],
    ];
    return entries
      .map(
        ([label, hex]) => `
      <div class="swatch">
        <div class="swatch-color" style="background:${hex}"></div>
        <div class="swatch-label">${label}<br><span class="swatch-hex">${hex}</span></div>
      </div>`
      )
      .join('');
  }

  private renderGradient(theme: Theme): string {
    const dir = theme.graphics.gradient.direction === 'vertical' ? '180deg' : '90deg';
    return `<div class="preview-gradient" style="background:linear-gradient(${dir}, ${theme.graphics.gradient.from}, ${theme.graphics.gradient.to})"></div>`;
  }

  private renderBorder(theme: Theme): string {
    const b = theme.borders.chars;
    const width = 30;
    const top = `${b.top_left}${b.horizontal.repeat(width)}${b.top_right}`;
    const mid = `${b.vertical}${' '.repeat(width)}${b.vertical}`;
    const bottom = `${b.bottom_left}${b.horizontal.repeat(width)}${b.bottom_right}`;
    return `<pre class="pv-border-box">${escapeHtml(top)}\n${escapeHtml(mid)}\n${escapeHtml(bottom)}</pre>`;
  }

  private renderProgress(theme: Theme): void {
    const el = this.container.querySelector('#pv-progress');
    if (!el) return;
    const filled = Math.round(theme.progress_bar.width * 0.65);
    const empty = theme.progress_bar.width - filled;
    const color = theme.progress_bar.color || theme.colors.primary;
    el.innerHTML = `[<span style="color:${color}">${theme.progress_bar.filled.repeat(Math.max(filled, 0))}</span><span style="color:${theme.colors.muted}">${theme.progress_bar.empty.repeat(Math.max(empty, 0))}</span>] 65%`;
  }

  private animateLoader(theme: Theme): void {
    const el = this.container.querySelector('#pv-loader');
    if (!el || theme.loader.frames.length === 0) return;

    let i = 0;
    const color = theme.loader.color || theme.colors.primary;
    const tick = () => {
      el.innerHTML = `<span style="color:${color}">${escapeHtml(theme.loader.frames[i % theme.loader.frames.length])}</span> <span class="pv-loader-label">Installing dependencies...</span>`;
      i++;
    };
    tick();
    const interval = window.setInterval(tick, Math.max(theme.loader.interval_ms, 16));
    this.timers.push(interval);
  }

  private animateBanner(theme: Theme): void {
    const el = this.container.querySelector('#pv-banner') as HTMLElement | null;
    if (!el) return;

    if (!theme.banner.enabled) {
      el.innerHTML = `<span class="pv-banner-disabled">banner disabled</span>`;
      return;
    }

    const text = theme.banner.text.replace('{user}', 'developer') || '(no banner text)';
    const color = theme.banner.color || theme.colors.primary;

    switch (theme.graphics.effects.banner) {
      case 'rainbow':
        this.animateRainbow(el, text);
        break;
      case 'glitch':
        el.textContent = text;
        el.style.color = color;
        el.className = 'preview-banner pv-glitch';
        break;
      case 'pulse':
        el.textContent = text;
        el.style.color = color;
        el.className = 'preview-banner pv-pulse';
        break;
      case 'neon':
        el.textContent = text;
        el.style.color = color;
        el.style.textShadow = `0 0 6px ${color}, 0 0 12px ${color}`;
        el.className = 'preview-banner';
        break;
      case 'typewriter':
        this.animateTypewriter(el, text, color);
        break;
      default:
        el.textContent = text;
        el.style.color = color;
        el.className = 'preview-banner';
    }
  }

  private animateRainbow(el: HTMLElement, text: string): void {
    const hues = [0, 30, 60, 120, 200, 260, 300];
    let offset = 0;
    const tick = () => {
      el.innerHTML = [...text]
        .map((ch, i) => {
          if (ch === ' ') return ' ';
          const hue = hues[(i + offset) % hues.length];
          return `<span style="color:hsl(${hue},90%,60%)">${escapeHtml(ch)}</span>`;
        })
        .join('');
      offset = (offset + 1) % hues.length;
    };
    tick();
    this.timers.push(window.setInterval(tick, 150));
  }

  private animateTypewriter(el: HTMLElement, text: string, color: string): void {
    el.style.color = color;
    el.className = 'preview-banner';
    const tick = () => {
      this.typewriterIndex = (this.typewriterIndex + 1) % (text.length + 8);
      el.textContent = text.slice(0, this.typewriterIndex);
    };
    this.timers.push(window.setInterval(tick, 60));
  }

  private clearTimers(): void {
    for (const t of this.timers) window.clearInterval(t);
    this.timers = [];
  }

  dispose(): void {
    this.clearTimers();
  }
}

function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}
