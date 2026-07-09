import type { Theme } from './theme';

export function themeToJson(t: Theme): string {
  return JSON.stringify(t, null, 2);
}

/** Triggers a browser download of the theme as "<name>.json". */
export function downloadTheme(t: Theme): void {
  const filename = `${sanitizeFilename(t.meta.name || 'theme')}.json`;
  const blob = new Blob([themeToJson(t)], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

export async function copyThemeToClipboard(t: Theme): Promise<void> {
  await navigator.clipboard.writeText(themeToJson(t));
}

/** Encodes the theme into a URL hash fragment so it can be shared as a
 * link — the recipient's browser loads the theme straight from the URL,
 * no server or database required. Uses base64 of the UTF-8 JSON bytes. */
export function themeToShareUrl(t: Theme): string {
  const json = themeToJson(t);
  const encoded = base64EncodeUnicode(json);
  const url = new URL(window.location.href);
  url.hash = `theme=${encoded}`;
  return url.toString();
}

/** Reads a theme from the current URL hash, if present. Returns null
 * if there's no theme in the URL or it fails to parse. */
export function themeFromCurrentUrl(): Theme | null {
  const hash = window.location.hash.replace(/^#/, '');
  const params = new URLSearchParams(hash);
  const encoded = params.get('theme');
  if (!encoded) return null;

  try {
    const json = base64DecodeUnicode(encoded);
    return JSON.parse(json) as Theme;
  } catch {
    return null;
  }
}

/** Reads a theme from a user-selected .json file. */
export function loadThemeFromFile(file: File): Promise<Theme> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      try {
        const parsed = JSON.parse(reader.result as string);
        resolve(parsed as Theme);
      } catch (err) {
        reject(new Error(`Could not parse "${file.name}" as JSON: ${(err as Error).message}`));
      }
    };
    reader.onerror = () => reject(new Error(`Could not read file "${file.name}"`));
    reader.readAsText(file);
  });
}

function sanitizeFilename(name: string): string {
  return name.trim().toLowerCase().replace(/[^a-z0-9-]+/g, '-').replace(/^-+|-+$/g, '') || 'theme';
}

function base64EncodeUnicode(str: string): string {
  const bytes = new TextEncoder().encode(str);
  let binary = '';
  for (const byte of bytes) {
    binary += String.fromCharCode(byte);
  }
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

function base64DecodeUnicode(encoded: string): string {
  const restored = encoded.replace(/-/g, '+').replace(/_/g, '/');
  const padded = restored + '='.repeat((4 - (restored.length % 4)) % 4);
  const binary = atob(padded);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return new TextDecoder().decode(bytes);
}
