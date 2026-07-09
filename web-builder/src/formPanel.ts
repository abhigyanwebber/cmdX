import type { Theme } from './theme';
import { isHexColor } from './validate';

type OnChange = () => void;

/**
 * FormPanel renders every field of the Theme schema as an editable
 * control and mutates the shared Theme object in place on input,
 * calling onChange() so the caller can re-render the preview and
 * re-run validation. Deliberately exposes every field the CLI/schema
 * supports — no simplified subset.
 */
export class FormPanel {
  constructor(
    private container: HTMLElement,
    private theme: Theme,
    private onChange: OnChange
  ) {}

  render(): void {
    this.container.innerHTML = '';
    this.container.appendChild(this.section('Meta', this.metaFields()));
    this.container.appendChild(this.section('Colors', this.colorFields()));
    this.container.appendChild(this.section('Prompt', this.promptFields()));
    this.container.appendChild(this.section('Loader', this.loaderFields()));
    this.container.appendChild(this.section('Progress Bar', this.progressFields()));
    this.container.appendChild(this.section('Cursor', this.cursorFields()));
    this.container.appendChild(this.section('Borders', this.borderFields()));
    this.container.appendChild(this.section('Banner', this.bannerFields()));
    this.container.appendChild(this.section('Gradient', this.gradientFields()));
    this.container.appendChild(this.section('Effects', this.effectsFields()));
    this.container.appendChild(this.section('Divider', this.dividerFields()));
    this.container.appendChild(this.section('Wallpaper', this.wallpaperFields()));
    this.container.appendChild(this.section('Icons (Nerd Font glyphs)', this.iconsFields()));
    this.container.appendChild(this.section('Linked Assets', this.assetsFields()));
  }

  // ── section scaffolding ──────────────────────────────────────────

  private section(title: string, body: HTMLElement): HTMLElement {
    const details = document.createElement('details');
    details.className = 'section';
    details.open = ['Meta', 'Colors', 'Prompt'].includes(title);
    const summary = document.createElement('summary');
    summary.textContent = title;
    details.appendChild(summary);
    details.appendChild(body);
    return details;
  }

  private fieldGroup(): HTMLElement {
    const div = document.createElement('div');
    div.className = 'field-group';
    return div;
  }

  private textField(label: string, value: string, onInput: (v: string) => void): HTMLElement {
    const wrap = document.createElement('label');
    wrap.className = 'field';
    wrap.innerHTML = `<span>${label}</span>`;
    const input = document.createElement('input');
    input.type = 'text';
    input.value = value;
    input.addEventListener('input', () => {
      onInput(input.value);
      this.onChange();
    });
    wrap.appendChild(input);
    return wrap;
  }

  private colorField(label: string, value: string, onInput: (v: string) => void): HTMLElement {
    const wrap = document.createElement('label');
    wrap.className = 'field field-color';
    wrap.innerHTML = `<span>${label}</span>`;

    const row = document.createElement('div');
    row.className = 'color-row';

    const swatch = document.createElement('input');
    swatch.type = 'color';
    swatch.value = isHexColor(value) && value.length === 7 ? value : '#888888';

    const text = document.createElement('input');
    text.type = 'text';
    text.value = value;
    text.className = isHexColor(value) ? '' : 'invalid';

    swatch.addEventListener('input', () => {
      text.value = swatch.value;
      text.classList.remove('invalid');
      onInput(swatch.value);
      this.onChange();
    });
    text.addEventListener('input', () => {
      text.classList.toggle('invalid', !isHexColor(text.value));
      onInput(text.value);
      if (isHexColor(text.value) && text.value.length === 7) swatch.value = text.value;
      this.onChange();
    });

    row.appendChild(swatch);
    row.appendChild(text);
    wrap.appendChild(row);
    return wrap;
  }

  private numberField(label: string, value: number, onInput: (v: number) => void, min = 0): HTMLElement {
    const wrap = document.createElement('label');
    wrap.className = 'field';
    wrap.innerHTML = `<span>${label}</span>`;
    const input = document.createElement('input');
    input.type = 'number';
    input.min = String(min);
    input.value = String(value);
    input.addEventListener('input', () => {
      onInput(Number(input.value));
      this.onChange();
    });
    wrap.appendChild(input);
    return wrap;
  }

  private boolField(label: string, value: boolean, onInput: (v: boolean) => void): HTMLElement {
    const wrap = document.createElement('label');
    wrap.className = 'field field-checkbox';
    const input = document.createElement('input');
    input.type = 'checkbox';
    input.checked = value;
    input.addEventListener('change', () => {
      onInput(input.checked);
      this.onChange();
    });
    wrap.appendChild(input);
    const span = document.createElement('span');
    span.textContent = label;
    wrap.appendChild(span);
    return wrap;
  }

  private selectField(label: string, value: string, options: string[], onInput: (v: string) => void): HTMLElement {
    const wrap = document.createElement('label');
    wrap.className = 'field';
    wrap.innerHTML = `<span>${label}</span>`;
    const select = document.createElement('select');
    for (const opt of options) {
      const o = document.createElement('option');
      o.value = opt;
      o.textContent = opt;
      if (opt === value) o.selected = true;
      select.appendChild(o);
    }
    select.addEventListener('change', () => {
      onInput(select.value);
      this.onChange();
    });
    wrap.appendChild(select);
    return wrap;
  }

  private framesField(label: string, frames: string[], onInput: (v: string[]) => void): HTMLElement {
    const wrap = document.createElement('label');
    wrap.className = 'field';
    wrap.innerHTML = `<span>${label} (comma-separated)</span>`;
    const input = document.createElement('input');
    input.type = 'text';
    input.value = frames.join(',');
    input.addEventListener('input', () => {
      onInput(input.value.split(',').map((s) => s.trim()).filter(Boolean));
      this.onChange();
    });
    wrap.appendChild(input);
    return wrap;
  }

  // ── field groups per theme section ──────────────────────────────

  private metaFields(): HTMLElement {
    const g = this.fieldGroup();
    const m = this.theme.meta;
    g.appendChild(this.textField('Name', m.name, (v) => (m.name = v)));
    g.appendChild(this.textField('Version', m.version, (v) => (m.version = v)));
    g.appendChild(this.textField('Author', m.author, (v) => (m.author = v)));
    g.appendChild(this.textField('Description', m.description, (v) => (m.description = v)));
    return g;
  }

  private colorFields(): HTMLElement {
    const g = this.fieldGroup();
    const c = this.theme.colors;
    (Object.keys(c) as (keyof typeof c)[]).forEach((key) => {
      g.appendChild(this.colorField(key, c[key], (v) => (c[key] = v)));
    });
    return g;
  }

  private promptFields(): HTMLElement {
    const g = this.fieldGroup();
    const p = this.theme.prompt;
    g.appendChild(this.textField('Symbol', p.symbol, (v) => (p.symbol = v)));
    g.appendChild(this.textField('Separator', p.separator, (v) => (p.separator = v)));
    g.appendChild(this.selectField('Style', p.style, ['single', 'multiline'], (v) => (p.style = v as typeof p.style)));
    g.appendChild(
      this.textField('Format (use {user} {dir} {git} {time} {symbol})', p.format, (v) => (p.format = v))
    );
    return g;
  }

  private loaderFields(): HTMLElement {
    const g = this.fieldGroup();
    const l = this.theme.loader;
    g.appendChild(this.framesField('Frames', l.frames, (v) => (l.frames = v)));
    g.appendChild(this.numberField('Interval (ms)', l.interval_ms, (v) => (l.interval_ms = v), 1));
    g.appendChild(this.colorField('Color (blank = primary)', l.color, (v) => (l.color = v)));
    return g;
  }

  private progressFields(): HTMLElement {
    const g = this.fieldGroup();
    const p = this.theme.progress_bar;
    g.appendChild(this.textField('Filled char', p.filled, (v) => (p.filled = v)));
    g.appendChild(this.textField('Empty char', p.empty, (v) => (p.empty = v)));
    g.appendChild(this.numberField('Width', p.width, (v) => (p.width = v), 1));
    g.appendChild(this.colorField('Color (blank = primary)', p.color, (v) => (p.color = v)));
    return g;
  }

  private cursorFields(): HTMLElement {
    const g = this.fieldGroup();
    const c = this.theme.cursor;
    g.appendChild(this.selectField('Style', c.style, ['block', 'bar', 'underline'], (v) => (c.style = v as typeof c.style)));
    g.appendChild(this.boolField('Blink', c.blink, (v) => (c.blink = v)));
    g.appendChild(this.colorField('Color (blank = primary)', c.color, (v) => (c.color = v)));
    return g;
  }

  private borderFields(): HTMLElement {
    const g = this.fieldGroup();
    const b = this.theme.borders;
    g.appendChild(this.textField('Style name', b.style, (v) => (b.style = v)));
    g.appendChild(this.textField('Top-left', b.chars.top_left, (v) => (b.chars.top_left = v)));
    g.appendChild(this.textField('Top-right', b.chars.top_right, (v) => (b.chars.top_right = v)));
    g.appendChild(this.textField('Bottom-left', b.chars.bottom_left, (v) => (b.chars.bottom_left = v)));
    g.appendChild(this.textField('Bottom-right', b.chars.bottom_right, (v) => (b.chars.bottom_right = v)));
    g.appendChild(this.textField('Horizontal', b.chars.horizontal, (v) => (b.chars.horizontal = v)));
    g.appendChild(this.textField('Vertical', b.chars.vertical, (v) => (b.chars.vertical = v)));
    return g;
  }

  private bannerFields(): HTMLElement {
    const g = this.fieldGroup();
    const b = this.theme.banner;
    g.appendChild(this.boolField('Enabled', b.enabled, (v) => (b.enabled = v)));
    g.appendChild(this.textField('Text (use {user})', b.text, (v) => (b.text = v)));
    g.appendChild(this.textField('Style label', b.style, (v) => (b.style = v)));
    g.appendChild(this.colorField('Color (blank = primary)', b.color, (v) => (b.color = v)));
    return g;
  }

  private gradientFields(): HTMLElement {
    const g = this.fieldGroup();
    const gr = this.theme.graphics.gradient;
    g.appendChild(this.boolField('Enabled', gr.enabled, (v) => (gr.enabled = v)));
    g.appendChild(this.colorField('From', gr.from, (v) => (gr.from = v)));
    g.appendChild(this.colorField('To', gr.to, (v) => (gr.to = v)));
    g.appendChild(
      this.selectField('Direction', gr.direction, ['horizontal', 'vertical'], (v) => (gr.direction = v as typeof gr.direction))
    );
    return g;
  }

  private effectsFields(): HTMLElement {
    const g = this.fieldGroup();
    const e = this.theme.graphics.effects;
    g.appendChild(
      this.selectField(
        'Banner effect',
        e.banner,
        ['none', 'glitch', 'rainbow', 'pulse', 'neon', 'typewriter'],
        (v) => (e.banner = v as typeof e.banner)
      )
    );
    g.appendChild(this.selectField('Prompt effect', e.prompt, ['none', 'rainbow'], (v) => (e.prompt = v as typeof e.prompt)));
    return g;
  }

  private dividerFields(): HTMLElement {
    const g = this.fieldGroup();
    const d = this.theme.graphics.divider;
    g.appendChild(
      this.selectField(
        'Style',
        d.style,
        ['line', 'wave', 'dots', 'stars', 'double', 'arrow', 'zigzag'],
        (v) => (d.style = v as typeof d.style)
      )
    );
    g.appendChild(this.colorField('Color (blank = primary)', d.color, (v) => (d.color = v)));
    return g;
  }

  private wallpaperFields(): HTMLElement {
    const g = this.fieldGroup();
    const w = this.theme.wallpaper;
    g.appendChild(this.boolField('Enabled', w.enabled, (v) => (w.enabled = v)));
    g.appendChild(this.textField('Image path', w.path, (v) => (w.path = v)));
    g.appendChild(this.numberField('Opacity (0–1)', w.opacity, (v) => (w.opacity = v), 0));
    g.appendChild(
      this.selectField('Stretch', w.stretch, ['fill', 'uniform', 'uniformToFill', 'none'], (v) => (w.stretch = v as typeof w.stretch))
    );
    g.appendChild(
      this.selectField(
        'Alignment',
        w.alignment,
        ['center', 'topLeft', 'topRight', 'bottomLeft', 'bottomRight'],
        (v) => (w.alignment = v as typeof w.alignment)
      )
    );
    return g;
  }

  private iconsFields(): HTMLElement {
    const g = this.fieldGroup();
    const i = this.theme.icons;
    g.appendChild(this.boolField('Enabled (requires a Nerd Font — see cmdx font install)', i.enabled, (v) => (i.enabled = v)));
    g.appendChild(this.selectField('Font', i.font, ['nerd-fonts', 'emoji', 'ascii'], (v) => (i.font = v as typeof i.font)));
    const glyphKeys: (keyof typeof i)[] = [
      'directory', 'file', 'git_branch', 'git_dirty', 'git_clean',
      'error', 'success', 'warning', 'time', 'package', 'docker',
      'python', 'node', 'rust', 'go',
    ];
    for (const key of glyphKeys) {
      g.appendChild(this.textField(key, String(i[key] ?? ''), (v) => ((i as any)[key] = v)));
    }
    return g;
  }

  private assetsFields(): HTMLElement {
    const g = this.fieldGroup();
    const a = this.theme.assets;
    g.appendChild(this.textField('Spinner asset name', a.spinner ?? '', (v) => (a.spinner = v || undefined)));
    g.appendChild(this.textField('Banner asset name', a.banner ?? '', (v) => (a.banner = v || undefined)));
    g.appendChild(this.textField('Divider asset name', a.divider ?? '', (v) => (a.divider = v || undefined)));
    g.appendChild(this.textField('Icons asset name', a.icons ?? '', (v) => (a.icons = v || undefined)));
    return g;
  }
}
