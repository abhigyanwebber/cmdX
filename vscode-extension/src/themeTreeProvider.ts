import * as vscode from 'vscode';
import { listThemesRaw, ThemeSummary } from './cmdxClient';

export class ThemeTreeItem extends vscode.TreeItem {
  constructor(public readonly theme: ThemeSummary) {
    super(theme.name, vscode.TreeItemCollapsibleState.None);
    this.description = theme.author ? `by ${theme.author}` : undefined;
    this.tooltip = theme.description ?? theme.name;
    this.contextValue = 'theme';
    this.iconPath = new vscode.ThemeIcon('symbol-color');
    this.command = {
      command: 'cmdx.openThemeFile',
      title: 'Open Theme File',
      arguments: [theme.filePath],
    };
  }
}

/**
 * Provides the "cmdX: Themes" sidebar view — lists every theme JSON
 * file found in the resolved themes directory, with inline preview
 * and apply actions (contributed via package.json's view/item/context).
 */
export class ThemeTreeProvider implements vscode.TreeDataProvider<ThemeTreeItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<void>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  constructor(private themesDirResolver: () => string) {}

  refresh(): void {
    this._onDidChangeTreeData.fire();
  }

  getTreeItem(element: ThemeTreeItem): vscode.TreeItem {
    return element;
  }

  async getChildren(): Promise<ThemeTreeItem[]> {
    const themesDir = this.themesDirResolver();
    if (!themesDir) return [];

    const themes = await listThemesRaw(themesDir);
    return themes
      .sort((a, b) => a.name.localeCompare(b.name))
      .map((t) => new ThemeTreeItem(t));
  }
}
