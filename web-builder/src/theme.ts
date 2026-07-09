// Mirrors internal/config/types.go exactly. Every field name here must
// match the Go struct's `json:"..."` tag — this is the source of truth
// the builder edits and exports.

export interface Meta {
  name: string;
  version: string;
  author: string;
  description: string;
}

export interface Colors {
  primary: string;
  secondary: string;
  background: string;
  foreground: string;
  accent: string;
  error: string;
  success: string;
  warning: string;
  muted: string;
}

export interface Prompt {
  symbol: string;
  separator: string;
  style: 'single' | 'multiline';
  segments: string[];
  format: string;
}

export interface Loader {
  frames: string[];
  interval_ms: number;
  color: string;
}

export interface ProgressBar {
  filled: string;
  empty: string;
  width: number;
  color: string;
}

export interface Cursor {
  style: 'block' | 'bar' | 'underline';
  blink: boolean;
  color: string;
}

export interface BorderChars {
  top_left: string;
  top_right: string;
  bottom_left: string;
  bottom_right: string;
  horizontal: string;
  vertical: string;
}

export interface Borders {
  style: string;
  chars: BorderChars;
}

export interface Banner {
  enabled: boolean;
  text: string;
  style: string;
  color: string;
}

export interface GradientConfig {
  enabled: boolean;
  from: string;
  to: string;
  direction: 'horizontal' | 'vertical';
}

export interface DividerConfig {
  style: 'line' | 'wave' | 'dots' | 'stars' | 'double' | 'arrow' | 'zigzag';
  color: string;
}

export interface EffectsConfig {
  banner: 'glitch' | 'rainbow' | 'pulse' | 'neon' | 'typewriter' | 'none';
  prompt: 'rainbow' | 'none';
}

export interface IconsConfig {
  enabled: boolean;
  directory: string;
  git_branch: string;
  error: string;
  success: string;
  time: string;
}

export interface Graphics {
  gradient: GradientConfig;
  divider: DividerConfig;
  effects: EffectsConfig;
  icons: IconsConfig;
}

export interface Wallpaper {
  enabled: boolean;
  path: string;
  opacity: number;
  stretch: 'fill' | 'uniform' | 'uniformToFill' | 'none';
  alignment: 'center' | 'topLeft' | 'topRight' | 'bottomLeft' | 'bottomRight';
}

export interface IconSet {
  enabled: boolean;
  font: 'nerd-fonts' | 'emoji' | 'ascii';
  directory: string;
  file: string;
  git_branch: string;
  git_dirty: string;
  git_clean: string;
  error: string;
  success: string;
  warning: string;
  time: string;
  package: string;
  docker: string;
  python: string;
  node: string;
  rust: string;
  go: string;
}

export interface ThemeAssets {
  spinner?: string;
  banner?: string;
  divider?: string;
  icons?: string;
}

export interface Theme {
  meta: Meta;
  colors: Colors;
  prompt: Prompt;
  loader: Loader;
  progress_bar: ProgressBar;
  cursor: Cursor;
  borders: Borders;
  banner: Banner;
  graphics: Graphics;
  wallpaper: Wallpaper;
  icons: IconSet;
  assets: ThemeAssets;
}

/** A complete, valid starting point — every field populated so the
 * builder never has to guess at a partial theme. */
export function defaultTheme(): Theme {
  return {
    meta: { name: 'my-theme', version: '1.0.0', author: '', description: 'A custom cmdX theme' },
    colors: {
      primary: '#FF00FF',
      secondary: '#00FFFF',
      background: '#0D0D0D',
      foreground: '#FFFFFF',
      accent: '#FFD700',
      error: '#FF4444',
      success: '#00FF88',
      warning: '#FFA500',
      muted: '#444444',
    },
    prompt: { symbol: '▶', separator: '›', style: 'single', segments: [], format: '{user}@{dir} {symbol}' },
    loader: { frames: ['◐', '◓', '◑', '◒'], interval_ms: 100, color: '' },
    progress_bar: { filled: '█', empty: '░', width: 20, color: '' },
    cursor: { style: 'block', blink: true, color: '' },
    borders: {
      style: 'rounded',
      chars: { top_left: '╭', top_right: '╮', bottom_left: '╰', bottom_right: '╯', horizontal: '─', vertical: '│' },
    },
    banner: { enabled: true, text: 'SYSTEM ONLINE // {user}', style: 'bold', color: '' },
    graphics: {
      gradient: { enabled: true, from: '#FF00FF', to: '#00FFFF', direction: 'horizontal' },
      divider: { style: 'line', color: '' },
      effects: { banner: 'none', prompt: 'none' },
      icons: { enabled: false, directory: '', git_branch: '', error: '', success: '', time: '' },
    },
    wallpaper: { enabled: false, path: '', opacity: 0.3, stretch: 'fill', alignment: 'center' },
    icons: {
      enabled: false,
      font: 'nerd-fonts',
      directory: '',
      file: '',
      git_branch: '',
      git_dirty: '',
      git_clean: '',
      error: '',
      success: '',
      warning: '',
      time: '',
      package: '',
      docker: '',
      python: '',
      node: '',
      rust: '',
      go: '',
    },
    assets: {},
  };
}
