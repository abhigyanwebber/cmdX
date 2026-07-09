import type { Theme } from './theme';

// Mirrors internal/config/validator.go's ValidateTheme. Kept in sync
// deliberately — the goal is that a theme which passes validation here
// also passes `cmdx theme validate`, so builder users don't get a
// surprise failure after exporting.

const hexColorPattern = /^#([A-Fa-f0-9]{3}|[A-Fa-f0-9]{6})$/;

export interface ValidationError {
  field: string;
  message: string;
}

export function validateTheme(t: Theme): ValidationError[] {
  const errors: ValidationError[] = [];

  if (!t.meta.name.trim()) errors.push({ field: 'meta.name', message: 'name is required' });
  if (!t.meta.version.trim()) errors.push({ field: 'meta.version', message: 'version is required' });

  const colorFields: [keyof Theme['colors'], string][] = [
    ['primary', 'primary color'],
    ['secondary', 'secondary color'],
    ['background', 'background color'],
    ['foreground', 'foreground color'],
    ['accent', 'accent color'],
    ['error', 'error color'],
    ['success', 'success color'],
    ['warning', 'warning color'],
    ['muted', 'muted color'],
  ];
  for (const [key, label] of colorFields) {
    const value = t.colors[key];
    if (!value) {
      errors.push({ field: `colors.${key}`, message: `${label} is required` });
    } else if (!hexColorPattern.test(value)) {
      errors.push({ field: `colors.${key}`, message: `${label} must be a valid hex color (e.g. #FF00FF)` });
    }
  }

  if (!t.prompt.symbol.trim()) errors.push({ field: 'prompt.symbol', message: 'symbol is required' });
  if (!t.prompt.format.trim()) errors.push({ field: 'prompt.format', message: 'format is required' });
  if (t.prompt.style !== 'single' && t.prompt.style !== 'multiline') {
    errors.push({ field: 'prompt.style', message: 'style must be "single" or "multiline"' });
  }

  if (t.loader.frames.length === 0) errors.push({ field: 'loader.frames', message: 'at least one frame is required' });
  if (t.loader.interval_ms <= 0) errors.push({ field: 'loader.interval_ms', message: 'interval_ms must be greater than 0' });

  if (!t.progress_bar.filled) errors.push({ field: 'progress_bar.filled', message: 'filled character is required' });
  if (!t.progress_bar.empty) errors.push({ field: 'progress_bar.empty', message: 'empty character is required' });
  if (t.progress_bar.width <= 0) errors.push({ field: 'progress_bar.width', message: 'width must be greater than 0' });

  if (!['block', 'bar', 'underline'].includes(t.cursor.style)) {
    errors.push({ field: 'cursor.style', message: 'style must be block, bar, or underline' });
  }

  return errors;
}

export function isHexColor(value: string): boolean {
  return hexColorPattern.test(value);
}
