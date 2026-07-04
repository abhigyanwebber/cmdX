package main

import (
	"fmt"
	"os"

	"github.com/abhigyanwebber/cmd-customizer/internal/assets"
	"github.com/abhigyanwebber/cmd-customizer/internal/theme"
	"github.com/spf13/cobra"
)

var assetStatusBarPreviewCmd = &cobra.Command{
	Use:   "statusbar-preview [asset-name]",
	Short: "Preview a status bar asset — shows segment mockup and generated shell code",
	Long: `Preview a status bar asset. Shows a visual mockup of the bar layout
and the generated shell hook code for a specific shell.

Examples:
  cmdx asset statusbar-preview my-bar
  cmdx asset statusbar-preview my-bar --shell bash
  cmdx asset statusbar-preview my-bar --shell zsh
  cmdx asset statusbar-preview my-bar --shell powershell`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		shell, _ := cmd.Flags().GetString("shell")

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		// load active theme colors for segment color resolution
		colors := loadThemeColors()

		if err := m.PreviewStatusBar(name, shell, colors); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}
	},
}

var assetStatusBarHooksCmd = &cobra.Command{
	Use:   "statusbar-hooks [asset-name]",
	Short: "Print the shell hook code for a status bar asset",
	Long: `Print the generated shell hook code for a status bar asset.
Pipe to a file or inspect before injecting into your shell profile.

Examples:
  cmdx asset statusbar-hooks my-bar --shell bash
  cmdx asset statusbar-hooks my-bar --shell zsh >> ~/.zshrc
  cmdx asset statusbar-hooks my-bar --shell powershell`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		shell, _ := cmd.Flags().GetString("shell")

		if shell == "" {
			_, shellName := detectShell()
			if shellName == "" {
				fmt.Println("✗ Could not detect shell. Use --shell to specify one.")
				fmt.Println("  Options: bash, zsh, powershell")
				os.Exit(1)
			}
			shell = shellName
		}

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		colors := loadThemeColors()

		hooks, err := m.StatusBarHooks(name, shell, colors)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		fmt.Printf("# Status bar hooks for '%s' (%s)\n", name, shell)
		fmt.Printf("# Add to your profile, or inject with: cmdx asset use %s --as status-bar\n\n", name)
		fmt.Println(hooks)
	},
}

var assetStatusBarInfoCmd = &cobra.Command{
	Use:   "statusbar-info [asset-name]",
	Short: "Show detailed segment breakdown of a status bar asset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		a, _, err := m.Get(name, assets.AssetTypeStatusBar)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}
		if a.StatusBar == nil {
			fmt.Println("✗ Asset has no status_bar config")
			os.Exit(1)
		}

		sb := a.StatusBar
		fmt.Printf("\n  Status Bar: %s\n", name)
		fmt.Printf("  Position:        %s\n", sb.Position)
		fmt.Printf("  Separator:       %s\n", sb.SeparatorStyle)
		if sb.CustomSeparatorLeft != "" {
			fmt.Printf("  Custom sep (L):  %s\n", sb.CustomSeparatorLeft)
		}
		if sb.CustomSeparatorRight != "" {
			fmt.Printf("  Custom sep (R):  %s\n", sb.CustomSeparatorRight)
		}
		fmt.Printf("  Height:          %d line(s)\n", sb.Height)
		fmt.Printf("  Fill background: %v\n", sb.FillBackground)
		if sb.HideOnNarrow > 0 {
			fmt.Printf("  Hide on narrow:  <%d columns\n", sb.HideOnNarrow)
		}
		fmt.Printf("\n  Segments (%d):\n\n", len(sb.Segments))

		zones := []assets.StatusBarZone{assets.ZoneLeft, assets.ZoneCenter, assets.ZoneRight}
		for _, zone := range zones {
			var segs []assets.SegmentConfig
			for _, s := range sb.Segments {
				if s.Zone == zone {
					segs = append(segs, s)
				}
			}
			if len(segs) == 0 {
				continue
			}
			fmt.Printf("  ── %s ──\n", zone)
			for _, seg := range segs {
				fmt.Printf("  %-14s  color: %-10s", seg.Type, seg.Color)
				if seg.Label != "" {
					fmt.Printf("  label: %q", seg.Label)
				}
				if seg.MaxLength > 0 {
					fmt.Printf("  max: %d", seg.MaxLength)
				}
				if seg.Bold {
					fmt.Printf("  bold")
				}
				if len(seg.Conditions) > 0 {
					fmt.Printf("  [conditional]")
				}
				if seg.EnvVar != "" {
					fmt.Printf("  env: %s", seg.EnvVar)
				}
				fmt.Println()
			}
			fmt.Println()
		}
	},
}

// loadThemeColors returns the active theme's color palette for segment
// color resolution. Returns empty map if no theme is active.
func loadThemeColors() map[string]string {
	m, err := theme.NewManager(getThemesDir())
	if err != nil {
		return map[string]string{}
	}
	t, err := m.GetActive()
	if err != nil {
		return map[string]string{}
	}
	return map[string]string{
		"primary":    t.Colors.Primary,
		"secondary":  t.Colors.Secondary,
		"accent":     t.Colors.Accent,
		"background": t.Colors.Background,
		"foreground": t.Colors.Foreground,
		"error":      t.Colors.Error,
		"success":    t.Colors.Success,
		"warning":    t.Colors.Warning,
		"muted":      t.Colors.Muted,
	}
}

func init() {
	assetStatusBarPreviewCmd.Flags().String("shell", "", "Shell type: bash, zsh, powershell")
	assetStatusBarHooksCmd.Flags().String("shell", "", "Shell type: bash, zsh, powershell (default: auto-detect)")

	assetCmd.AddCommand(assetStatusBarPreviewCmd)
	assetCmd.AddCommand(assetStatusBarHooksCmd)
	assetCmd.AddCommand(assetStatusBarInfoCmd)
}
