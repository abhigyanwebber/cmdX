package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abhigyanwebber/cmd-customizer/internal/plugin"
	"github.com/spf13/cobra"
)

func getPluginsDir() string {
	wd, err := os.Getwd()
	if err == nil {
		local := filepath.Join(wd, "plugins")
		if _, err := os.Stat(local); err == nil {
			return local
		}
	}
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "plugins")
}

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage cmdX plugins",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed plugins",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := plugin.NewManager(getPluginsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		if err := m.Discover(); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		names := m.List()
		if len(names) == 0 {
			fmt.Println("\n  No plugins installed.")
			fmt.Printf("  Drop a plugin folder into: %s\n\n", getPluginsDir())
			return
		}

		fmt.Println("\n  Installed Plugins:")
		for _, name := range names {
			p, _ := m.Get(name)
			fmt.Printf("  %-20s  %s\n", p.Meta.Name, p.Meta.Description)
			fmt.Printf("  %-20s  v%s by %s\n\n", "", p.Meta.Version, p.Meta.Author)
		}
	},
}

var pluginInfoCmd = &cobra.Command{
	Use:   "info [plugin-name]",
	Short: "Show detailed info about a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		m, err := plugin.NewManager(getPluginsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		if err := m.Discover(); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		p, err := m.Get(name)
		if err != nil {
			fmt.Printf("✗ Plugin '%s' not found. Run 'cmdx plugin list' to see installed plugins.\n", name)
			os.Exit(1)
		}

		fmt.Printf("\n  Plugin: %s\n", p.Meta.Name)
		fmt.Printf("  Version:     %s\n", p.Meta.Version)
		fmt.Printf("  Author:      %s\n", p.Meta.Author)
		fmt.Printf("  Description: %s\n", p.Meta.Description)
		if p.Meta.Homepage != "" {
			fmt.Printf("  Homepage:    %s\n", p.Meta.Homepage)
		}

		if len(p.Spinners) > 0 {
			fmt.Printf("\n  Spinners (%d):\n", len(p.Spinners))
			for _, s := range p.Spinners {
				fmt.Printf("    %-15s  %v  %dms\n", s.Name, s.Frames, s.IntervalMs)
			}
		}

		if len(p.Banners) > 0 {
			fmt.Printf("\n  Banners (%d):\n", len(p.Banners))
			for _, b := range p.Banners {
				fmt.Printf("    %-15s  %s\n", b.Name, b.Text)
			}
		}

		if len(p.Prompts) > 0 {
			fmt.Printf("\n  Prompts (%d):\n", len(p.Prompts))
			for _, pr := range p.Prompts {
				fmt.Printf("    %-15s  %s\n", pr.Name, pr.Format)
			}
		}
		fmt.Println()
	},
}

var pluginValidateCmd = &cobra.Command{
	Use:   "validate [plugin-name]",
	Short: "Validate an installed plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		pluginDir := filepath.Join(getPluginsDir(), name)
		p, err := plugin.LoadPlugin(pluginDir)
		if err != nil {
			fmt.Println("✗ Load error:", err)
			os.Exit(1)
		}

		m, _ := plugin.NewManager(getPluginsDir())
		if err := m.ValidatePlugin(p); err != nil {
			fmt.Println("✗ Validation failed:", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Plugin '%s' is valid\n", p.Meta.Name)
	},
}

var pluginSpinnersCmd = &cobra.Command{
	Use:   "spinners [plugin-name]",
	Short: "List all spinners provided by a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		m, err := plugin.NewManager(getPluginsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		if err := m.Discover(); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		p, err := m.Get(name)
		if err != nil {
			fmt.Printf("✗ Plugin '%s' not found.\n", name)
			os.Exit(1)
		}

		if len(p.Spinners) == 0 {
			fmt.Printf("  Plugin '%s' has no spinners.\n", name)
			return
		}

		fmt.Printf("\n  Spinners from '%s':\n\n", name)
		for _, s := range p.Spinners {
			frames := ""
			for _, f := range s.Frames {
				frames += f + " "
			}
			fmt.Printf("  %-15s  frames: %s  speed: %dms\n", s.Name, frames, s.IntervalMs)
		}
		fmt.Println()
	},
}

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginInfoCmd)
	pluginCmd.AddCommand(pluginValidateCmd)
	pluginCmd.AddCommand(pluginSpinnersCmd)
	rootCmd.AddCommand(pluginCmd)
}
