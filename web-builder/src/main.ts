import { defaultTheme, Theme } from './theme';
import { FormPanel } from './formPanel';
import { Preview } from './preview';
import { validateTheme } from './validate';
import { presets } from './presets';
import { downloadTheme, copyThemeToClipboard, themeToShareUrl, themeFromCurrentUrl, loadThemeFromFile } from './exportImport';

const app = document.querySelector<HTMLDivElement>('#app')!;

app.innerHTML = `
  <header class="topbar">
    <div class="brand">cmdX <span class="brand-accent">Theme Builder</span></div>
    <div class="toolbar">
      <select id="preset-select" class="toolbar-select">
        <option value="">Load preset…</option>
        ${Object.keys(presets)
          .map((name) => `<option value="${name}">${name}</option>`)
          .join('')}
      </select>
      <label class="toolbar-button" for="import-file">Import JSON</label>
      <input type="file" id="import-file" accept="application/json" hidden />
      <button id="btn-copy" class="toolbar-button">Copy JSON</button>
      <button id="btn-share" class="toolbar-button">Share Link</button>
      <button id="btn-download" class="toolbar-button toolbar-button-primary">Download theme.json</button>
      <button id="btn-reset" class="toolbar-button toolbar-button-danger">Reset</button>
    </div>
  </header>
  <main class="layout">
    <div class="panel-form" id="form-panel"></div>
    <div class="panel-preview">
      <div class="preview-container" id="preview-container"></div>
      <div class="validation-panel" id="validation-panel"></div>
    </div>
  </main>
`;

let theme: Theme = themeFromCurrentUrl() ?? defaultTheme();

const formContainer = document.querySelector<HTMLDivElement>('#form-panel')!;
const previewContainer = document.querySelector<HTMLDivElement>('#preview-container')!;
const validationContainer = document.querySelector<HTMLDivElement>('#validation-panel')!;

const preview = new Preview(previewContainer);
let form: FormPanel;

function rerender(): void {
  preview.render(theme);
  renderValidation();
}

function renderValidation(): void {
  const errors = validateTheme(theme);
  if (errors.length === 0) {
    validationContainer.innerHTML = `<div class="validation-ok">✓ Valid — matches what "cmdx theme validate" would accept</div>`;
    return;
  }
  validationContainer.innerHTML = `
    <div class="validation-errors">
      <strong>${errors.length} issue${errors.length > 1 ? 's' : ''}:</strong>
      <ul>${errors.map((e) => `<li><code>${e.field}</code> — ${e.message}</li>`).join('')}</ul>
    </div>`;
}

function rebuildForm(): void {
  form = new FormPanel(formContainer, theme, rerender);
  form.render();
}

function loadTheme(next: Theme): void {
  theme = next;
  rebuildForm();
  rerender();
}

rebuildForm();
rerender();

// ── toolbar wiring ─────────────────────────────────────────────────

document.querySelector<HTMLSelectElement>('#preset-select')!.addEventListener('change', (e) => {
  const name = (e.target as HTMLSelectElement).value;
  if (!name) return;
  loadTheme(structuredClone(presets[name]));
  (e.target as HTMLSelectElement).value = '';
});

document.querySelector<HTMLInputElement>('#import-file')!.addEventListener('change', async (e) => {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (!file) return;
  try {
    const imported = await loadThemeFromFile(file);
    loadTheme({ ...defaultTheme(), ...imported });
  } catch (err) {
    alert((err as Error).message);
  }
  (e.target as HTMLInputElement).value = '';
});

document.querySelector<HTMLButtonElement>('#btn-copy')!.addEventListener('click', async () => {
  await copyThemeToClipboard(theme);
  flashButton('#btn-copy', 'Copied!');
});

document.querySelector<HTMLButtonElement>('#btn-share')!.addEventListener('click', async () => {
  const url = themeToShareUrl(theme);
  await navigator.clipboard.writeText(url);
  window.history.replaceState(null, '', url);
  flashButton('#btn-share', 'Link copied!');
});

document.querySelector<HTMLButtonElement>('#btn-download')!.addEventListener('click', () => {
  downloadTheme(theme);
});

document.querySelector<HTMLButtonElement>('#btn-reset')!.addEventListener('click', () => {
  if (confirm('Reset to the default theme? Unsaved changes will be lost.')) {
    loadTheme(defaultTheme());
    window.history.replaceState(null, '', window.location.pathname);
  }
});

function flashButton(selector: string, message: string): void {
  const btn = document.querySelector<HTMLButtonElement>(selector)!;
  const original = btn.textContent;
  btn.textContent = message;
  setTimeout(() => {
    btn.textContent = original;
  }, 1200);
}
