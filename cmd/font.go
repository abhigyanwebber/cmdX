// Font commands: list, search, info, install, remove, installed, path.
//
// cmdX ships a curated catalog of popular Nerd Fonts for convenience,
// but every command supports installing any font from a direct URL via
// --url, so the catalog is a starting point, not a restriction.
package main

import (
	"fmt"
	"os"

	"github.com/abhigyanwebber/cmd-customizer/internal/fonts"
	"github.com/spf13/cobra"
)

var fontCmd = &cobra.Command{
	Use:   "font",
	Short: "Discover and install terminal fonts (Nerd Fonts + custom)",
	Long: `Discover, install, and manage terminal fonts.

cmdX ships a curated catalog of popular Nerd Fonts (fonts patched with
icon glyphs, needed for the icons and status-bar segment types). You
are not limited to the catalog — install any font from any direct
archive URL with --url.

Examples:
  cmdx font list
  cmdx font install firacode
  cmdx font install jetbrainsmono --variant Bold
  cmdx font install --url https://example.com/my-font.zip --name my-font
  cmdx font installed
  cmdx font remove firacode`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var fontListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the curated font catalog",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n  Curated Font Catalog (Nerd Fonts %s)\n\n", fonts.NerdFontsVersion)
		for _, f := range fonts.Catalog {
			marker := "  "
			if fonts.IsInstalled(f.Name) {
				marker = "✓ "
			}
			fmt.Printf("  %s%-16s  %-28s  %s\n", marker, f.Name, f.DisplayName, f.Description)
		}
		fmt.Println()
		fmt.Println("  Not seeing what you want? Install any font directly:")
		fmt.Println("    cmdx font install --url <zip-url> --name <your-font-name>")
		fmt.Println()
	},
}

var fontSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the curated font catalog",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		results := fonts.Search(args[0])
		if len(results) == 0 {
			fmt.Printf("\n  No catalog fonts matched %q.\n", args[0])
			fmt.Println("  Try 'cmdx font list' to see all options, or install any font with --url.")
			return
		}
		fmt.Printf("\n  Found %d match(es):\n\n", len(results))
		for _, f := range results {
			fmt.Printf("  %-16s  %-28s  %s\n", f.Name, f.DisplayName, f.Description)
		}
		fmt.Println()
	},
}

var fontInfoCmd = &cobra.Command{
	Use:   "info [font-name]",
	Short: "Show details about a catalog font",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f := fonts.Find(args[0])
		if f == nil {
			fmt.Printf("✗ Font %q not found in catalog. Run 'cmdx font list' to see options.\n", args[0])
			os.Exit(1)
		}

		fmt.Printf("\n  %s\n", f.DisplayName)
		fmt.Printf("  %s\n\n", f.Description)
		fmt.Printf("  Catalog name:  %s\n", f.Name)
		fmt.Printf("  License:       %s\n", f.License)
		fmt.Printf("  Monospace:     %v\n", f.Monospace)
		fmt.Printf("  Download URL:  %s\n", f.DownloadURL())
		fmt.Printf("  Installed:     %v\n\n", fonts.IsInstalled(f.Name))
	},
}

var fontInstallCmd = &cobra.Command{
	Use:   "install [font-name]",
	Short: "Install a font from the catalog or a custom URL",
	Long: `Install a font.

From the curated catalog:
  cmdx font install firacode
  cmdx font install jetbrainsmono --variant Bold
  cmdx font install hack --all-variants

From any direct URL (developer freedom — not limited to the catalog):
  cmdx font install --url https://example.com/my-font.zip --name my-font
  cmdx font install --url https://github.com/user/repo/releases/download/v1/Font.zip --name custom --max-size 200

Fonts install to your per-user font directory:
  Windows: %LOCALAPPDATA%\Microsoft\Windows\Fonts
  macOS:   ~/Library/Fonts
  Linux:   ~/.local/share/fonts`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		name, _ := cmd.Flags().GetString("name")
		variant, _ := cmd.Flags().GetString("variant")
		allVariants, _ := cmd.Flags().GetBool("all-variants")
		force, _ := cmd.Flags().GetBool("force")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		maxSizeMB, _ := cmd.Flags().GetInt("max-size")

		var fontName string
		if len(args) > 0 {
			fontName = args[0]
		} else if url != "" {
			fontName = name
		}

		if fontName == "" {
			fmt.Println("✗ Provide a catalog font name, or use --url with --name for a custom font.")
			os.Exit(1)
		}

		opts := fonts.InstallOptions{
			URL:         url,
			Variant:     variant,
			AllVariants: allVariants,
			Force:       force,
			DryRun:      dryRun,
		}
		if maxSizeMB > 0 {
			opts.MaxDownloadBytes = int64(maxSizeMB) << 20
		}

		fmt.Printf("\n  Installing %s...\n", fontName)

		installed, err := fonts.InstallFont(fontName, opts)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		if dryRun {
			return
		}

		for _, rec := range installed {
			fmt.Printf("\n✓ Installed %d file(s):\n", len(rec.Files))
			for _, f := range rec.Files {
				fmt.Printf("    %s\n", f)
			}
		}

		fmt.Printf("\n  Set this as your terminal font in your terminal emulator's settings.\n")
		if !fonts.RegistrySupported() {
			fmt.Printf("  You may need to restart your terminal for the font to appear in font pickers.\n")
		}
		fmt.Println()
	},
}

var fontRemoveCmd = &cobra.Command{
	Use:   "remove [font-name]",
	Short: "Remove a cmdX-installed font",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := fonts.RemoveFont(name); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Removed font '%s'\n", name)
	},
}

var fontInstalledCmd = &cobra.Command{
	Use:   "installed",
	Short: "List fonts cmdX has installed",
	Run: func(cmd *cobra.Command, args []string) {
		installed, err := fonts.ListInstalledFonts()
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		if len(installed) == 0 {
			fmt.Println("\n  No fonts installed via cmdX yet. Run 'cmdx font list' to browse the catalog.")
			return
		}

		fmt.Printf("\n  Installed Fonts (%d):\n\n", len(installed))
		for _, f := range installed {
			fmt.Printf("  %-16s  %d file(s)  installed %s\n", f.Name, len(f.Files), f.InstalledAt.Format("2006-01-02 15:04"))
		}
		fmt.Println()
	},
}

var fontPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show where fonts install on this platform",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := fonts.FontsDir()
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}
		fmt.Println(dir)
	},
}

func init() {
	fontInstallCmd.Flags().String("url", "", "Install from a custom archive URL instead of the catalog")
	fontInstallCmd.Flags().String("name", "", "Name to track this font under (required with --url)")
	fontInstallCmd.Flags().String("variant", "", "Only install files matching this variant (e.g. 'Bold', 'Italic', 'Mono')")
	fontInstallCmd.Flags().Bool("all-variants", false, "Install every variant in the archive (default if --variant is unset)")
	fontInstallCmd.Flags().Bool("force", false, "Reinstall even if already tracked as installed")
	fontInstallCmd.Flags().Bool("dry-run", false, "Show what would be installed without writing any files")
	fontInstallCmd.Flags().Int("max-size", 0, "Override the download size limit in MB (default: 150MB)")

	fontCmd.AddCommand(fontListCmd)
	fontCmd.AddCommand(fontSearchCmd)
	fontCmd.AddCommand(fontInfoCmd)
	fontCmd.AddCommand(fontInstallCmd)
	fontCmd.AddCommand(fontRemoveCmd)
	fontCmd.AddCommand(fontInstalledCmd)
	fontCmd.AddCommand(fontPathCmd)
	rootCmd.AddCommand(fontCmd)
}
