package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/abhigyanwebber/cmd-customizer/internal/assets"
	"github.com/spf13/cobra"
)

func getAssetsDir() string {
	wd, err := os.Getwd()
	if err == nil {
		local := filepath.Join(wd, "assets")
		if _, err := os.Stat(local); err == nil {
			return local
		}
	}
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "assets")
}

var assetCmd = &cobra.Command{
	Use:   "asset",
	Short: "Manage custom graphic assets (spinners, icons, banners, dividers)",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var assetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed assets",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		types := []assets.AssetType{
			assets.AssetTypeSpinner,
			assets.AssetTypeIcon,
			assets.AssetTypeBanner,
			assets.AssetTypeDivider,
		}

		found := false
		for _, t := range types {
			list, err := m.List(t)
			if err != nil || len(list) == 0 {
				continue
			}
			found = true
			fmt.Printf("\n  %ss:\n", string(t))
			for _, a := range list {
				fmt.Printf("    %-20s  %s  v%s by %s\n",
					a.Name, a.Description, a.Version, a.Author)
			}
		}

		if !found {
			fmt.Println("\n  No assets installed.")
			fmt.Printf("  Drop asset folders into: %s\n", getAssetsDir())
			fmt.Println("  Then run: cmdx asset import <path>")
		}
		fmt.Println()
	},
}

var assetImportCmd = &cobra.Command{
	Use:   "import [path-to-asset-folder]",
	Short: "Import an asset folder into cmdX",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := args[0]

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		fmt.Printf("  Importing asset from: %s\n", sourcePath)

		a, err := m.Import(sourcePath)
		if err != nil {
			fmt.Println("✗ Import failed:", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Asset '%s' imported successfully\n", a.Name)
		fmt.Printf("  Type:    %s\n", a.Type)
		fmt.Printf("  Version: %s\n", a.Version)
		fmt.Printf("  Author:  %s\n", a.Author)
		fmt.Printf("\n  Run 'cmdx asset preview %s' to preview it\n", a.Name)
	},
}

var assetInfoCmd = &cobra.Command{
	Use:   "info [asset-name]",
	Short: "Show detailed info about an asset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		types := []assets.AssetType{
			assets.AssetTypeSpinner,
			assets.AssetTypeIcon,
			assets.AssetTypeBanner,
			assets.AssetTypeDivider,
		}

		for _, t := range types {
			a, assetDir, err := m.Get(name, t)
			if err != nil {
				continue
			}

			fmt.Printf("\n  Asset: %s\n", a.Name)
			fmt.Printf("  Type:        %s\n", a.Type)
			fmt.Printf("  Version:     %s\n", a.Version)
			fmt.Printf("  Author:      %s\n", a.Author)
			fmt.Printf("  Description: %s\n", a.Description)
			if a.Homepage != "" {
				fmt.Printf("  Homepage:    %s\n", a.Homepage)
			}
			fmt.Printf("  Location:    %s\n", assetDir)
			fmt.Printf("\n  Render Config:\n")
			fmt.Printf("    Mode:       %s\n", a.Render.Mode)
			fmt.Printf("    Size:       %dx%d\n", a.Render.Width, a.Render.Height)
			fmt.Printf("    Color:      %s\n", a.Render.ColorMode)
			fmt.Printf("    Symbols:    %s\n", a.Render.SymbolSet)

			if a.Spinner != nil {
				fmt.Printf("\n  Spinner Config:\n")
				fmt.Printf("    Frames:     %d\n", len(a.Spinner.Frames))
				fmt.Printf("    Speed:      %dms\n", a.Spinner.IntervalMs)
				fmt.Printf("    Bounce:     %v\n", a.Spinner.Bounce)
				fmt.Printf("    Reverse:    %v\n", a.Spinner.Reverse)
				fmt.Printf("    On Complete: %s\n", a.Spinner.OnComplete)
			}
			fmt.Println()
			return
		}

		fmt.Printf("✗ Asset '%s' not found.\n", name)
		fmt.Println("  Run 'cmdx asset list' to see installed assets.")
	},
}

var assetPreviewCmd = &cobra.Command{
	Use:   "preview [asset-name]",
	Short: "Preview an asset live in the terminal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		duration, _ := cmd.Flags().GetInt("duration")
		iconKey, _ := cmd.Flags().GetString("icon")

		if !assets.ChafaAvailable() {
			fmt.Println("✗ chafa is not installed.")
			fmt.Println("  Install it with: winget install hpjansson.chafa")
			os.Exit(1)
		}

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		// try each asset type
		types := []assets.AssetType{
			assets.AssetTypeSpinner,
			assets.AssetTypeBanner,
			assets.AssetTypeDivider,
			assets.AssetTypeIcon,
		}

		for _, t := range types {
			a, _, err := m.Get(name, t)
			if err != nil {
				continue
			}

			fmt.Printf("\n  Preview: %s (%s)\n\n", a.Name, a.Type)

			switch a.Type {
			case assets.AssetTypeSpinner:
				if err := m.PreviewSpinner(name, time.Duration(duration)*time.Second); err != nil {
					fmt.Println("✗ Preview failed:", err)
					os.Exit(1)
				}

			case assets.AssetTypeBanner:
				if err := m.PreviewBanner(name); err != nil {
					fmt.Println("✗ Preview failed:", err)
					os.Exit(1)
				}

			case assets.AssetTypeDivider:
				if err := m.PreviewDivider(name); err != nil {
					fmt.Println("✗ Preview failed:", err)
					os.Exit(1)
				}

			case assets.AssetTypeIcon:
				if iconKey != "" {
					if err := m.PreviewIcon(name, iconKey); err != nil {
						fmt.Println("✗ Preview failed:", err)
						os.Exit(1)
					}
				} else {
					// preview all icons
					fmt.Printf("  Icons in '%s':\n\n", name)
					for key := range a.Icon.Files {
						m.PreviewIcon(name, key)
					}
				}
			}
			return
		}

		fmt.Printf("✗ Asset '%s' not found.\n", name)
	},
}

var assetValidateCmd = &cobra.Command{
	Use:   "validate [path-to-asset-folder]",
	Short: "Validate an asset folder before importing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := args[0]

		a, err := assets.LoadAsset(sourcePath)
		if err != nil {
			fmt.Println("✗ Load failed:", err)
			os.Exit(1)
		}

		if err := assets.ValidateAsset(a, sourcePath); err != nil {
			fmt.Println("✗ Validation failed:", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Asset '%s' is valid\n", a.Name)
		fmt.Printf("  Type: %s  Version: %s  Author: %s\n",
			a.Type, a.Version, a.Author)
	},
}

var assetRemoveCmd = &cobra.Command{
	Use:   "remove [asset-name] [type]",
	Short: "Remove an installed asset",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		assetType := assets.AssetType(args[1])

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		if err := m.Remove(name, assetType); err != nil {
			fmt.Println("✗ Remove failed:", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Asset '%s' removed\n", name)
	},
}

var assetChafaCmd = &cobra.Command{
	Use:   "chafa",
	Short: "Check chafa installation status",
	Run: func(cmd *cobra.Command, args []string) {
		if !assets.ChafaAvailable() {
			fmt.Println("✗ chafa is not installed")
			fmt.Println("  Install: winget install hpjansson.chafa")
			os.Exit(1)
		}

		version, _ := assets.ChafaVersion()
		fmt.Println("✓ chafa is installed:", version)
	},
}

func init() {
	assetPreviewCmd.Flags().IntP("duration", "d", 5, "Preview duration in seconds (spinners only)")
	assetPreviewCmd.Flags().StringP("icon", "i", "", "Icon key to preview (icons only, e.g. directory, error)")

	assetCmd.AddCommand(assetListCmd)
	assetCmd.AddCommand(assetImportCmd)
	assetCmd.AddCommand(assetInfoCmd)
	assetCmd.AddCommand(assetPreviewCmd)
	assetCmd.AddCommand(assetValidateCmd)
	assetCmd.AddCommand(assetRemoveCmd)
	assetCmd.AddCommand(assetChafaCmd)
	rootCmd.AddCommand(assetCmd)
}
