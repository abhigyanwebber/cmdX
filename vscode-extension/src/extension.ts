import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs/promises';
import { runCmdx, listAssetsRaw, CmdxNotFoundError } from './cmdxClient';
import { ThemePreviewPanel } from './previewPanel';
import { ThemeTreeProvider } from './themeTreeProvider';
import { AssetTreeProvider } from './assetTreeProvider';

let outputChannel: vscode.OutputChannel;

export function activate(context: vscode.ExtensionContext) {
  outputChannel = vscode.window.createOutputChannel('cmdX');

  const themeTreeProvider = new ThemeTreeProvider(resolveThemesDir);
  const assetTreeProvider = new AssetTreeProvider(resolveAssetsDir);

  vscode.window.registerTreeDataProvider('cmdxThemes', themeTreeProvider);
  vscode.window.registerTreeDataProvider('cmdxAssets', assetTreeProvider);

  context.subscriptions.push(
    vscode.commands.registerCommand('cmdx.refreshThemes', () => themeTreeProvider.refresh()),
    vscode.commands.registerCommand('cmdx.refreshAssets', () => assetTreeProvider.refresh()),

    vscode.commands.registerCommand('cmdx.openThemeFile', async (filePath: string) => {
      await vscode.window.showTextDocument(vscode.Uri.file(filePath));
    }),

    vscode.commands.registerCommand('cmdx.previewTheme', async (filePathArg?: string) => {
      const filePath = filePathArg ?? vscode.window.activeTextEditor?.document.uri.fsPath;
      if (!filePath) {
        vscode.window.showWarningMessage('cmdX: Open a theme JSON file first, or select one from the Themes view.');
        return;
      }
      try {
        const raw = await fs.readFile(filePath, 'utf-8');
        const theme = JSON.parse(raw);
        ThemePreviewPanel.createOrShow(context.extensionUri, theme, filePath);
      } catch (err) {
        vscode.window.showErrorMessage(`cmdX: Could not parse theme file — ${(err as Error).message}`);
      }
    }),

    vscode.commands.registerCommand('cmdx.applyTheme', async (filePathArg?: string) => {
      const filePath = filePathArg ?? vscode.window.activeTextEditor?.document.uri.fsPath;
      if (!filePath) {
        vscode.window.showWarningMessage('cmdX: Open a theme JSON file first.');
        return;
      }
      const themeName = await themeNameFromFile(filePath);
      if (!themeName) return;

      await withCmdxErrorHandling(async () => {
        const result = await runCmdx(['theme', 'apply', themeName]);
        outputChannel.appendLine(result.stdout);
        if (result.exitCode !== 0) {
          outputChannel.appendLine(result.stderr);
          outputChannel.show();
          vscode.window.showErrorMessage(`cmdX: Failed to apply theme '${themeName}'. See cmdX output for details.`);
        } else {
          vscode.window.showInformationMessage(`cmdX: Applied theme '${themeName}'. Restart your terminal to see changes.`);
        }
      });
    }),

    vscode.commands.registerCommand('cmdx.validateTheme', async (filePathArg?: string) => {
      const filePath = filePathArg ?? vscode.window.activeTextEditor?.document.uri.fsPath;
      if (!filePath) {
        vscode.window.showWarningMessage('cmdX: Open a theme JSON file first.');
        return;
      }

      await withCmdxErrorHandling(async () => {
        const result = await runCmdx(['theme', 'validate', filePath]);
        outputChannel.appendLine(result.stdout);
        if (result.exitCode !== 0) {
          outputChannel.appendLine(result.stderr);
          outputChannel.show();
          vscode.window.showErrorMessage(`cmdX: Theme validation failed. See cmdX output for details.`);
        } else {
          vscode.window.showInformationMessage(`cmdX: Theme is valid.`);
        }
      });
    }),

    vscode.commands.registerCommand('cmdx.validateAsset', async (filePathArg?: string) => {
      const filePath = filePathArg ?? vscode.window.activeTextEditor?.document.uri.fsPath;
      if (!filePath) {
        vscode.window.showWarningMessage('cmdX: Open an asset.json file first.');
        return;
      }
      const assetDir = path.dirname(filePath);

      await withCmdxErrorHandling(async () => {
        const result = await runCmdx(['asset', 'validate', assetDir]);
        outputChannel.appendLine(result.stdout);
        if (result.exitCode !== 0) {
          outputChannel.appendLine(result.stderr);
          outputChannel.show();
          vscode.window.showErrorMessage(`cmdX: Asset validation failed. See cmdX output for details.`);
        } else {
          vscode.window.showInformationMessage(`cmdX: Asset is valid.`);
        }
      });
    }),

    // "cmdx theme create" is an interactive terminal wizard (huh forms) —
    // running it in a real integrated terminal, rather than reimplementing
    // it as a webview form, keeps the extension honest about what the CLI
    // actually does and avoids duplicating logic that could drift out of
    // sync with the Go implementation.
    vscode.commands.registerCommand('cmdx.newTheme', () => {
      const terminal = getOrCreateCmdxTerminal();
      terminal.show();
      terminal.sendText(`${quoteIfNeeded(getBinaryPathSetting())} theme create`);
    }),

    // Asset preview renders real chafa graphics (braille/sixel/truecolor
    // blocks) — a webview cannot faithfully reproduce actual terminal
    // rendering, so this opens a real terminal rather than faking the
    // output. The command still supports the full override flag set via
    // a multi-step quick-pick, matching the CLI's own flexibility.
    vscode.commands.registerCommand('cmdx.previewAsset', async () => {
      const assetsDir = resolveAssetsDir();
      const assets = await listAssetsRaw(assetsDir);
      if (assets.length === 0) {
        vscode.window.showWarningMessage('cmdX: No assets found. Set cmdx.assetsDirectory in settings, or run "cmdx asset create" first.');
        return;
      }

      const picked = await vscode.window.showQuickPick(
        assets.map((a) => ({ label: a.name, description: a.type, detail: a.description })),
        { placeHolder: 'Select an asset to preview' }
      );
      if (!picked) return;

      const mode = await vscode.window.showQuickPick(
        ['(manifest default)', 'braille', 'blocks', 'ascii', 'sixel', 'color'],
        { placeHolder: 'Render mode override (optional)' }
      );
      if (mode === undefined) return;

      const args = ['asset', 'preview', picked.label];
      if (mode && mode !== '(manifest default)') {
        args.push('--mode', mode);
      }

      const terminal = getOrCreateCmdxTerminal();
      terminal.show();
      terminal.sendText(`${quoteIfNeeded(getBinaryPathSetting())} ${args.map(quoteIfNeeded).join(' ')}`);
    }),

    vscode.commands.registerCommand('cmdx.installFont', () => {
      const terminal = getOrCreateCmdxTerminal();
      terminal.show();
      terminal.sendText(`${quoteIfNeeded(getBinaryPathSetting())} font list`);
    }),

    vscode.commands.registerCommand('cmdx.setBinaryPath', async () => {
      const current = getBinaryPathSetting();
      const value = await vscode.window.showInputBox({
        prompt: 'Path to the cmdx executable',
        value: current,
        placeHolder: 'cmdx (or an absolute path if not on PATH)',
      });
      if (value === undefined) return;
      await vscode.workspace.getConfiguration('cmdx').update('executablePath', value, vscode.ConfigurationTarget.Global);
      vscode.window.showInformationMessage(`cmdX: Executable path set to "${value}".`);
    })
  );

  // Optional live preview refresh on save.
  context.subscriptions.push(
    vscode.workspace.onDidSaveTextDocument(async (doc) => {
      const autoPreview = vscode.workspace.getConfiguration('cmdx').get<boolean>('autoPreviewOnSave', false);
      if (!autoPreview) return;
      if (!doc.fileName.endsWith('.json')) return;
      if (!ThemePreviewPanel.currentPanel) return;

      try {
        const theme = JSON.parse(doc.getText());
        ThemePreviewPanel.createOrShow(context.extensionUri, theme, doc.fileName);
      } catch {
        // invalid JSON mid-edit — ignore, user will see validation
        // errors from VS Code's own JSON schema support already
      }
    })
  );
}

export function deactivate() {
  outputChannel?.dispose();
}

let cmdxTerminal: vscode.Terminal | undefined;

function getOrCreateCmdxTerminal(): vscode.Terminal {
  if (cmdxTerminal && cmdxTerminal.exitStatus === undefined) {
    return cmdxTerminal;
  }
  cmdxTerminal = vscode.window.createTerminal('cmdX');
  return cmdxTerminal;
}

function getBinaryPathSetting(): string {
  return vscode.workspace.getConfiguration('cmdx').get<string>('executablePath', 'cmdx');
}

/** Quotes an argument for the integrated terminal shell if it contains
 * spaces or other characters that would otherwise be split incorrectly. */
function quoteIfNeeded(arg: string): string {
  if (/^[A-Za-z0-9_\-./:\\]+$/.test(arg)) {
    return arg;
  }
  return `"${arg.replace(/"/g, '\\"')}"`;
}

function resolveThemesDir(): string {
  const configured = vscode.workspace.getConfiguration('cmdx').get<string>('themesDirectory', '');
  if (configured) return configured;

  const ws = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
  if (ws) return path.join(ws, 'themes');

  return '';
}

function resolveAssetsDir(): string {
  const configured = vscode.workspace.getConfiguration('cmdx').get<string>('assetsDirectory', '');
  if (configured) return configured;

  const ws = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
  if (ws) return path.join(ws, 'assets');

  return '';
}

/** Determines the theme name to pass to `cmdx theme apply/validate`
 * from a theme JSON file: prefers meta.name from the file content,
 * falls back to the filename stem if the file can't be parsed. */
async function themeNameFromFile(filePath: string): Promise<string | undefined> {
  try {
    const raw = await fs.readFile(filePath, 'utf-8');
    const parsed = JSON.parse(raw);
    if (parsed?.meta?.name) return parsed.meta.name;
  } catch {
    // fall through to filename-based guess
  }
  const base = path.basename(filePath, '.json');
  return base || undefined;
}

async function withCmdxErrorHandling(fn: () => Promise<void>): Promise<void> {
  try {
    await fn();
  } catch (err) {
    if (err instanceof CmdxNotFoundError) {
      const choice = await vscode.window.showErrorMessage(err.message, 'Set Executable Path');
      if (choice === 'Set Executable Path') {
        vscode.commands.executeCommand('cmdx.setBinaryPath');
      }
    } else {
      vscode.window.showErrorMessage(`cmdX: ${(err as Error).message}`);
    }
  }
}
