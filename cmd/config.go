package main

import (
	"fmt"
	"os"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cmdX global configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current global configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.InitGlobalConfig(); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		cfg, err := config.GetGlobalConfig()
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		configDir, _ := config.ConfigDir()

		fmt.Println("\n  cmdX Global Configuration")
		fmt.Printf("  Config dir:    %s\n", configDir)
		fmt.Printf("  Default theme: %s\n", cfg.DefaultTheme)
		fmt.Printf("  Render mode:   %s\n", cfg.RenderMode)
		fmt.Printf("  Chafa path:    %s\n", cfg.ChafaPath)
		fmt.Printf("  Auto inject:   %v\n", cfg.AutoInject)
		fmt.Printf("  Show banner:   %v\n", cfg.ShowBanner)
		if cfg.AssetsDir != "" {
			fmt.Printf("  Assets dir:    %s\n", cfg.AssetsDir)
		}
		if cfg.ThemesDir != "" {
			fmt.Printf("  Themes dir:    %s\n", cfg.ThemesDir)
		}
		fmt.Println()
		fmt.Println("  Override with environment variables:")
		fmt.Println("  CMDX_DEFAULT_THEME, CMDX_RENDER_MODE, CMDX_CHAFA_PATH")
		fmt.Println()
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a global configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if err := config.InitGlobalConfig(); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		cfg, err := config.GetGlobalConfig()
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		switch key {
		case "default_theme":
			cfg.DefaultTheme = value
		case "render_mode":
			valid := map[string]bool{"auto": true, "sixel": true, "braille": true, "blocks": true, "ascii": true}
			if !valid[value] {
				fmt.Println("✗ render_mode must be: auto, sixel, braille, blocks, ascii")
				os.Exit(1)
			}
			cfg.RenderMode = value
		case "chafa_path":
			cfg.ChafaPath = value
		case "assets_dir":
			cfg.AssetsDir = value
		case "themes_dir":
			cfg.ThemesDir = value
		case "auto_inject":
			cfg.AutoInject = value == "true"
		case "show_banner":
			cfg.ShowBanner = value == "true"
		default:
			fmt.Printf("✗ Unknown config key '%s'\n", key)
			fmt.Println("  Valid keys: default_theme, render_mode, chafa_path, assets_dir, themes_dir, auto_inject, show_banner")
			os.Exit(1)
		}

		if err := config.SaveGlobalConfig(cfg); err != nil {
			fmt.Println("✗ Could not save config:", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Set %s = %s\n", key, value)
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.DefaultGlobalConfig()
		if err := config.SaveGlobalConfig(cfg); err != nil {
			fmt.Println("✗ Could not reset config:", err)
			os.Exit(1)
		}
		fmt.Println("✓ Configuration reset to defaults")
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	rootCmd.AddCommand(configCmd)
}
