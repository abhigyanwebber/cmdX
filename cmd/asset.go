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
			assets.AssetTypeFloater,
			assets.AssetTypeMascot,
			assets.AssetTypeStatusBar,
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
			assets.AssetTypeFloater,
			assets.AssetTypeMascot,
			assets.AssetTypeStatusBar,
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
			if a.Floater != nil {
				fmt.Printf("\n  Floater Config:\n")
				fmt.Printf("    File:       %s\n", a.Floater.File)
				fmt.Printf("    Position:   %s\n", a.Floater.Position)
				fmt.Printf("    Max Size:   %dx%d\n", a.Floater.MaxWidth, a.Floater.MaxHeight)
				fmt.Printf("    Margin:     x=%d y=%d\n", a.Floater.MarginX, a.Floater.MarginY)
				if len(a.Floater.AnimateFrames) > 0 {
					fmt.Printf("    Animated:   %d frames @ %dms\n", len(a.Floater.AnimateFrames), a.Floater.IntervalMs)
				} else {
					fmt.Printf("    Animated:   no (static)\n")
				}
			}
			if a.Mascot != nil {
				mc := a.Mascot
				fmt.Printf("\n  Mascot Config:\n")
				fmt.Printf("    States:     %d defined\n", len(mc.States))
				fmt.Printf("    Default:    %s\n", mc.DefaultState)
				fmt.Printf("    Position:   %s\n", mc.Position)
				fmt.Printf("    Max Size:   %dx%d\n", mc.MaxWidth, mc.MaxHeight)
				fmt.Printf("    Hook var:   %s\n", hookVarOrDefault(mc.HookVar))
				fmt.Printf("\n  Run 'cmdx asset mascot-info %s' for full state machine details.\n", a.Name)
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
	Short: "Preview an asset in the terminal",
	Long: `Preview an asset live in the terminal.

All render options can be overridden at the CLI level without editing asset.json:

  cmdx asset preview my-floater --mode sixel --width 32 --color truecolor
  cmdx asset preview my-floater --position bottom-left
  cmdx asset preview my-spinner --mode braille --width 16 --dither
  cmdx asset preview my-banner --stretch --threshold 0.3`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		duration, _ := cmd.Flags().GetInt("duration")
		iconKey, _ := cmd.Flags().GetString("icon")

		// render override flags
		modeStr, _ := cmd.Flags().GetString("mode")
		colorStr, _ := cmd.Flags().GetString("color")
		symbolsStr, _ := cmd.Flags().GetString("symbols")
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")
		dither, _ := cmd.Flags().GetBool("dither")
		stretch, _ := cmd.Flags().GetBool("stretch")
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		positionStr, _ := cmd.Flags().GetString("position")

		ditherSet := cmd.Flags().Changed("dither")
		stretchSet := cmd.Flags().Changed("stretch")

		overrides := assets.OverridesFromFlags(modeStr, colorStr, symbolsStr, width, height, dither, stretch, threshold, ditherSet, stretchSet)

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

		// search all types
		types := []assets.AssetType{
			assets.AssetTypeSpinner,
			assets.AssetTypeBanner,
			assets.AssetTypeDivider,
			assets.AssetTypeFloater,
			assets.AssetTypeMascot,
			assets.AssetTypeStatusBar,
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
				if err := m.PreviewSpinnerWithOverrides(name, time.Duration(duration)*time.Second, overrides); err != nil {
					fmt.Println("✗ Preview failed:", err)
					os.Exit(1)
				}

			case assets.AssetTypeBanner:
				if err := m.PreviewBannerWithOverrides(name, overrides); err != nil {
					fmt.Println("✗ Preview failed:", err)
					os.Exit(1)
				}

			case assets.AssetTypeDivider:
				if err := m.PreviewDividerWithOverrides(name, overrides); err != nil {
					fmt.Println("✗ Preview failed:", err)
					os.Exit(1)
				}

			case assets.AssetTypeFloater:
				pos := assets.FloaterPosition(positionStr)
				if err := m.PreviewFloaterWithOverrides(name, pos, overrides); err != nil {
					fmt.Println("✗ Preview failed:", err)
					os.Exit(1)
				}

			case assets.AssetTypeMascot:
				stateStr, _ := cmd.Flags().GetString("state")
				ctx := assets.MascotContext{
					LastExitCode: 0,
					Env:          map[string]string{},
				}
				if stateStr != "" {
					ctx.Env["CMDX_MASCOT_STATE"] = stateStr
				}
				if err := m.PreviewMascot(name, ctx, overrides); err != nil {
					fmt.Println("✗ Preview failed:", err)
					os.Exit(1)
				}

			case assets.AssetTypeIcon:
				if iconKey != "" {
					if err := m.PreviewIconWithOverrides(name, iconKey, overrides); err != nil {
						fmt.Println("✗ Preview failed:", err)
						os.Exit(1)
					}
				} else {
					fmt.Printf("  Icons in '%s':\n\n", name)
					for key := range a.Icon.Files {
						m.PreviewIconWithOverrides(name, key, overrides)
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

var assetUseCmd = &cobra.Command{
	Use:   "use [asset-name]",
	Short: "Use a standalone asset without applying a full theme",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		slot, _ := cmd.Flags().GetString("as")
		position, _ := cmd.Flags().GetString("position")

		if slot == "" {
			fmt.Println("✗ Specify what to use this asset as: --as spinner|banner|divider|icons|floater")
			os.Exit(1)
		}

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		validTypes := map[string]assets.AssetType{
			"spinner": assets.AssetTypeSpinner,
			"banner":  assets.AssetTypeBanner,
			"divider": assets.AssetTypeDivider,
			"icons":   assets.AssetTypeIcon,
			"floater": assets.AssetTypeFloater,
			"mascot":     assets.AssetTypeMascot,
			"status-bar": assets.AssetTypeStatusBar,
		}

		assetType, ok := validTypes[slot]
		if !ok {
			fmt.Printf("✗ Unknown slot '%s'. Use: spinner, banner, divider, icons, floater, mascot, status-bar\n", slot)
			os.Exit(1)
		}

		a, _, err := m.Get(name, assetType)
		if err != nil {
			fmt.Printf("✗ Asset '%s' not found as type '%s'\n", name, slot)
			os.Exit(1)
		}

		// Floaters are positioned independently per corner, so up to
		// four can be active at once. Every other slot is a single
		// global selection (e.g. one active spinner).
		stateDir := filepath.Join(getAssetsDir(), ".state")
		os.MkdirAll(stateDir, 0755)

		if slot == "floater" {
			resolvedPosition := position
			if resolvedPosition == "" {
				resolvedPosition = string(a.Floater.Position)
			}
			if !assets.IsValidFloaterPosition(assets.FloaterPosition(resolvedPosition)) {
				fmt.Printf("✗ Invalid position '%s'. Use: top-left, top-right, bottom-left, bottom-right\n", resolvedPosition)
				os.Exit(1)
			}

			statePath := filepath.Join(stateDir, "floater-"+resolvedPosition+".txt")
			if err := os.WriteFile(statePath, []byte(name), 0644); err != nil {
				fmt.Println("✗ Could not save floater state:", err)
				os.Exit(1)
			}

			fmt.Printf("✓ '%s' is now active as your %s floater\n", name, resolvedPosition)
			fmt.Printf("  Preview: cmdx asset preview %s\n", name)
			return
		}

		statePath := filepath.Join(stateDir, slot+".txt")
		if err := os.WriteFile(statePath, []byte(name), 0644); err != nil {
			fmt.Println("✗ Could not save asset state:", err)
			os.Exit(1)
		}

		fmt.Printf("✓ '%s' is now active as your %s\n", name, slot)
		fmt.Printf("  Preview: cmdx asset preview %s\n", name)
		if slot == "mascot" {
			fmt.Printf("  View states: cmdx asset mascot-info %s\n", name)
			fmt.Printf("  Get hooks:   cmdx asset mascot-hooks %s\n", name)
			fmt.Printf("  Shell hooks are automatically installed when you inject a theme.\n")
		}
		if slot == "status-bar" {
			fmt.Printf("  View layout: cmdx asset statusbar-info %s\n", name)
			fmt.Printf("  Get hooks:   cmdx asset statusbar-hooks %s\n", name)
			fmt.Printf("  Preview:     cmdx asset statusbar-preview %s --shell bash\n", name)
		}
	},
}

var assetRenderCmd = &cobra.Command{
	Use:   "render [image-path]",
	Short: "Render any PNG directly to the terminal (no manifest needed)",
	Long: `Render any PNG or image file through cmdX's chafa pipeline directly,
without needing an asset.json manifest. Useful for testing images and
exploring render options before committing to a manifest.

Examples:
  cmdx asset render ./my-image.png
  cmdx asset render ./sprite.png --mode braille --width 20 --height 10
  cmdx asset render ./banner.png --mode sixel --color truecolor
  cmdx asset render ./icon.png --mode ascii --color none --width 8`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]

		modeStr, _ := cmd.Flags().GetString("mode")
		colorStr, _ := cmd.Flags().GetString("color")
		symbolsStr, _ := cmd.Flags().GetString("symbols")
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")
		dither, _ := cmd.Flags().GetBool("dither")
		stretch, _ := cmd.Flags().GetBool("stretch")
		threshold, _ := cmd.Flags().GetFloat64("threshold")

		if !assets.ChafaAvailable() {
			fmt.Println("✗ chafa is not installed.")
			fmt.Println("  Install: winget install hpjansson.chafa")
			os.Exit(1)
		}

		// build defaults from the detected render mode
		mode := assets.RenderMode(modeStr)
		if mode == "" {
			mode = assets.BestRenderMode()
		}
		defaultCfg := assets.DefaultRenderConfig(mode)

		// apply any CLI overrides
		opts := assets.ChafaOptions{
			Mode:      defaultCfg.Mode,
			ColorMode: defaultCfg.ColorMode,
			SymbolSet: defaultCfg.SymbolSet,
			Width:     defaultCfg.Width,
			Height:    defaultCfg.Height,
			Dither:    defaultCfg.Dither,
			Stretch:   defaultCfg.Stretch,
			Threshold: defaultCfg.Threshold,
		}

		if colorStr != "" {
			opts.ColorMode = assets.ColorMode(colorStr)
		}
		if symbolsStr != "" {
			opts.SymbolSet = assets.SymbolSet(symbolsStr)
		}
		if width > 0 {
			opts.Width = width
		}
		if height > 0 {
			opts.Height = height
		}
		if cmd.Flags().Changed("dither") {
			opts.Dither = dither
		}
		if cmd.Flags().Changed("stretch") {
			opts.Stretch = stretch
		}
		if threshold > 0 {
			opts.Threshold = threshold
		}

		rendered, err := assets.RenderRaw(imagePath, opts)
		if err != nil {
			fmt.Println("✗ Render failed:", err)
			os.Exit(1)
		}

		fmt.Printf("\n  Rendered: %s  [%s / %s / %dx%d]\n\n",
			imagePath, opts.Mode, opts.ColorMode, opts.Width, opts.Height)
		fmt.Println(rendered)
	},
}

var assetStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show currently active assets",
	Run: func(cmd *cobra.Command, args []string) {
		stateDir := filepath.Join(getAssetsDir(), ".state")
		slots := []string{"spinner", "banner", "divider", "icons"}

		fmt.Println("\n  Active Assets:")
		for _, slot := range slots {
			statePath := filepath.Join(stateDir, slot+".txt")
			data, err := os.ReadFile(statePath)
			if err != nil {
				fmt.Printf("  %-10s  none\n", slot+":")
				continue
			}
			fmt.Printf("  %-10s  %s\n", slot+":", string(data))
		}

		// Floaters are tracked per corner rather than as a single slot,
		// so each of the four positions is checked individually.
		fmt.Println("\n  Active Floaters:")
		for _, pos := range assets.ValidFloaterPositions {
			statePath := filepath.Join(stateDir, "floater-"+string(pos)+".txt")
			data, err := os.ReadFile(statePath)
			label := string(pos) + ":"
			if err != nil {
				fmt.Printf("  %-14s  none\n", label)
				continue
			}
			fmt.Printf("  %-14s  %s\n", label, string(data))
		}
		fmt.Println()
	},
}

func init() {
	// ── asset render flags ─────────────────────────────────────────────────
	assetRenderCmd.Flags().String("mode", "", "Render mode: braille, blocks, ascii, sixel, color (default: auto-detect)")
	assetRenderCmd.Flags().String("color", "", "Color depth: truecolor, 256, ansi, none (default: auto-detect)")
	assetRenderCmd.Flags().String("symbols", "", "chafa symbol set: braille, block, border, ascii, all")
	assetRenderCmd.Flags().Int("width", 0, "Output width in terminal columns (default: mode-appropriate)")
	assetRenderCmd.Flags().Int("height", 0, "Output height in terminal rows (default: mode-appropriate)")
	assetRenderCmd.Flags().Bool("dither", false, "Enable dithering")
	assetRenderCmd.Flags().Bool("stretch", false, "Stretch image to fill dimensions")
	assetRenderCmd.Flags().Float64("threshold", 0, "Alpha/contrast threshold (0.0–1.0, default: 0.5)")

	// ── asset preview flags ────────────────────────────────────────────────
	// All render options are exposed as CLI flags so developers can test
	// different settings without editing asset.json.
	assetPreviewCmd.Flags().IntP("duration", "d", 5, "Preview duration in seconds (spinners only)")
	assetPreviewCmd.Flags().StringP("icon", "i", "", "Icon key to preview (icons only)")
	assetPreviewCmd.Flags().String("mode", "", "Override render mode: braille, blocks, ascii, sixel, color")
	assetPreviewCmd.Flags().String("color", "", "Override color depth: truecolor, 256, ansi, none")
	assetPreviewCmd.Flags().String("symbols", "", "Override chafa symbol set: braille, block, border, ascii, all")
	assetPreviewCmd.Flags().Int("width", 0, "Override output width in terminal columns")
	assetPreviewCmd.Flags().Int("height", 0, "Override output height in terminal rows")
	assetPreviewCmd.Flags().Bool("dither", false, "Override: enable dithering")
	assetPreviewCmd.Flags().Bool("stretch", false, "Override: stretch image to fill dimensions")
	assetPreviewCmd.Flags().Float64("threshold", 0, "Override alpha/contrast threshold (0.0–1.0)")
	assetPreviewCmd.Flags().String("position", "", "Override floater position: top-left, top-right, bottom-left, bottom-right")
	assetPreviewCmd.Flags().String("state", "", "Override mascot state: idle, working, success, error, warning, sleeping (or custom)")

	// ── asset use flags ────────────────────────────────────────────────────
	assetUseCmd.Flags().StringP("as", "a", "", "Asset slot: spinner, banner, divider, icons, floater")
	assetUseCmd.Flags().String("position", "", "Floater corner position: top-left, top-right, bottom-left, bottom-right (defaults to the asset's configured position)")

	assetCmd.AddCommand(assetListCmd)
	assetCmd.AddCommand(assetImportCmd)
	assetCmd.AddCommand(assetInfoCmd)
	assetCmd.AddCommand(assetPreviewCmd)
	assetCmd.AddCommand(assetValidateCmd)
	assetCmd.AddCommand(assetRemoveCmd)
	assetCmd.AddCommand(assetChafaCmd)
	assetCmd.AddCommand(assetUseCmd)
	assetCmd.AddCommand(assetStatusCmd)
	assetCmd.AddCommand(assetRenderCmd)
	rootCmd.AddCommand(assetCmd)
}
