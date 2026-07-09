import * as vscode from 'vscode';
import { spawn } from 'child_process';

/**
 * CmdxNotFoundError is thrown when the configured cmdx executable
 * cannot be found or fails to run at all (as opposed to running and
 * returning a non-zero exit code for a legitimate reason, e.g. theme
 * validation failure).
 */
export class CmdxNotFoundError extends Error {
  constructor(binaryPath: string) {
    super(
      `Could not find or run the cmdx executable at "${binaryPath}". ` +
        `Set "cmdx.executablePath" in settings, or install cmdx and ensure it's on your PATH. ` +
        `See https://github.com/abhigyanwebber/cmdX/blob/main/REQUIREMENTS.md`
    );
    this.name = 'CmdxNotFoundError';
  }
}

export interface CmdxResult {
  stdout: string;
  stderr: string;
  exitCode: number;
}

/** Returns the configured cmdx executable path (default: "cmdx" on PATH). */
export function getBinaryPath(): string {
  return vscode.workspace.getConfiguration('cmdx').get<string>('executablePath', 'cmdx');
}

/**
 * Runs the cmdx CLI with the given arguments and returns stdout/stderr
 * and exit code. Does not throw on a non-zero exit code — callers
 * decide what a given exit code means for their command. Throws
 * CmdxNotFoundError only if the executable itself couldn't be spawned.
 *
 * If the user has set cmdx.themesDirectory / cmdx.assetsDirectory in
 * settings, those are passed through as CMDX_THEMES_DIR / CMDX_ASSETS_DIR
 * environment variables so the invoked binary resolves the same
 * directories the extension itself is reading from.
 */
export function runCmdx(args: string[], cwd?: string): Promise<CmdxResult> {
  const binaryPath = getBinaryPath();
  const extraEnv = envOverrides();

  return new Promise((resolve, reject) => {
    const proc = spawn(binaryPath, args, {
      cwd: cwd ?? vscode.workspace.workspaceFolders?.[0]?.uri.fsPath,
      shell: false,
      env: { ...process.env, ...extraEnv },
    });

    let stdout = '';
    let stderr = '';

    proc.stdout.on('data', (chunk) => {
      stdout += chunk.toString();
    });
    proc.stderr.on('data', (chunk) => {
      stderr += chunk.toString();
    });

    proc.on('error', () => {
      reject(new CmdxNotFoundError(binaryPath));
    });

    proc.on('close', (code) => {
      resolve({ stdout, stderr, exitCode: code ?? -1 });
    });
  });
}

/** Extra environment variables derived from user settings, passed to
 * every cmdx invocation so the binary resolves the same themes/assets
 * directories the extension itself reads from directly. */
function envOverrides(): NodeJS.ProcessEnv {
  const cfg = vscode.workspace.getConfiguration('cmdx');
  const env: NodeJS.ProcessEnv = {};
  const themesDir = cfg.get<string>('themesDirectory', '');
  const assetsDir = cfg.get<string>('assetsDirectory', '');
  if (themesDir) env.CMDX_THEMES_DIR = themesDir;
  if (assetsDir) env.CMDX_ASSETS_DIR = assetsDir;
  return env;
}

export interface ThemeSummary {
  name: string;
  description?: string;
  author?: string;
  filePath: string;
}

/** Lists installed themes by reading the themes directory directly
 * (avoids depending on cmdx's interactive TUI list command, which
 * isn't suitable for non-interactive extension use). */
export async function listThemesRaw(themesDir: string): Promise<ThemeSummary[]> {
  const fs = await import('fs/promises');
  const path = await import('path');
  try {
    const entries = await fs.readdir(themesDir);
    const themes: ThemeSummary[] = [];
    for (const entry of entries) {
      if (!entry.endsWith('.json')) continue;
      const filePath = path.join(themesDir, entry);
      try {
        const raw = await fs.readFile(filePath, 'utf-8');
        const parsed = JSON.parse(raw);
        themes.push({
          name: parsed?.meta?.name ?? entry.replace(/\.json$/, ''),
          description: parsed?.meta?.description,
          author: parsed?.meta?.author,
          filePath,
        });
      } catch {
        // skip unparsable theme files — surfaced via validate command instead
      }
    }
    return themes;
  } catch {
    return [];
  }
}

export interface AssetSummary {
  name: string;
  type: string;
  description?: string;
  manifestPath: string;
}

/** Lists installed assets by scanning the assets directory's
 * type subfolders (spinners/, banners/, dividers/, icons/, floaters/,
 * mascots/, status-bars/), reading each asset.json manifest. */
export async function listAssetsRaw(assetsDir: string): Promise<AssetSummary[]> {
  const fs = await import('fs/promises');
  const path = await import('path');

  const typeFolders = ['spinners', 'banners', 'dividers', 'icons', 'floaters', 'mascots', 'status-bars'];
  const assets: AssetSummary[] = [];

  for (const folder of typeFolders) {
    const folderPath = path.join(assetsDir, folder);
    let entries: string[];
    try {
      entries = await fs.readdir(folderPath);
    } catch {
      continue;
    }
    for (const entry of entries) {
      const manifestPath = path.join(folderPath, entry, 'asset.json');
      try {
        const raw = await fs.readFile(manifestPath, 'utf-8');
        const parsed = JSON.parse(raw);
        assets.push({
          name: parsed?.name ?? entry,
          type: parsed?.type ?? folder,
          description: parsed?.description,
          manifestPath,
        });
      } catch {
        // skip unparsable/missing manifests
      }
    }
  }

  return assets;
}
