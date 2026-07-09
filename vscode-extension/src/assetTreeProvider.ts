import * as vscode from 'vscode';
import { listAssetsRaw, AssetSummary } from './cmdxClient';

const typeIcons: Record<string, string> = {
  spinner: 'loading',
  banner: 'symbol-string',
  divider: 'horizontal-rule',
  icon: 'symbol-misc',
  floater: 'layout-panel',
  mascot: 'smiley',
  'status-bar': 'layout-statusbar',
};

export class AssetTreeItem extends vscode.TreeItem {
  constructor(public readonly asset: AssetSummary) {
    super(asset.name, vscode.TreeItemCollapsibleState.None);
    this.description = asset.type;
    this.tooltip = asset.description ?? asset.name;
    this.contextValue = 'asset';
    this.iconPath = new vscode.ThemeIcon(typeIcons[asset.type] ?? 'file-media');
    this.command = {
      command: 'vscode.open',
      title: 'Open Asset Manifest',
      arguments: [vscode.Uri.file(asset.manifestPath)],
    };
  }
}

/**
 * Provides the "cmdX: Assets" sidebar view — lists every asset found
 * across all type subfolders in the resolved assets directory, with an
 * inline preview action (contributed via package.json's
 * view/item/context, invoking cmdx.previewAsset).
 */
export class AssetTreeProvider implements vscode.TreeDataProvider<AssetTreeItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<void>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  constructor(private assetsDirResolver: () => string) {}

  refresh(): void {
    this._onDidChangeTreeData.fire();
  }

  getTreeItem(element: AssetTreeItem): vscode.TreeItem {
    return element;
  }

  async getChildren(): Promise<AssetTreeItem[]> {
    const assetsDir = this.assetsDirResolver();
    if (!assetsDir) return [];

    const assets = await listAssetsRaw(assetsDir);
    return assets
      .sort((a, b) => a.type.localeCompare(b.type) || a.name.localeCompare(b.name))
      .map((a) => new AssetTreeItem(a));
  }
}
