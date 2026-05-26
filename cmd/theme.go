package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/preview"
	"github.com/abhigyanwebber/cmd-customizer/internal/theme"
	"github.com/spf13/cobra"
)

// getThemesDir resolves themes directory relative to binary or dev working dir
func getThemesDir() string {
	// in dev, use working directory
	wd, err := os.Getwd()
	if err == nil {
		local := filepath.Join(wd, "themes")
		if _, err := os.Stat(local); err == nil {
			return local
		}
	}
	// fallback to binary location
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "themes")
}

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage and apply themes",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available themes",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := theme.NewManager(getThemesDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		names, err := m.List()
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		fmt.Println("\n Available Themes:\n")
		for _, name := range names {
			t, err := m.Load(name)
			if err != nil {
				fmt.Printf("  %-15s  (invalid)\n", name)
				continue
			}
			fmt.Printf("  %-15s  %s\n", name, t.Meta.Description)
		}
		fmt.Println()
	},
}

var themeApplyCmd = &cobra.Command{
	Use:   "apply [theme-name]",
	Short: "Apply a theme to your terminal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		m, err := theme.NewManager(getThemesDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		if !m.Exists(name) {
			fmt.Printf("✗ Theme '%s' not found. Run 'cmdx theme list' to see available themes.\n", name)
			os.Exit(1)
		}

		if err := m.Apply(name); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		t, _ := m.GetActive()
		r := theme.NewRenderer(t)
		r.RenderThemeInfo()
	},
}

var themeInfoCmd = &cobra.Command{
	Use:   "info [theme-name]",
	Short: "Show detailed info about a theme",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		m, err := theme.NewManager(getThemesDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		t, err := m.Load(name)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		r := theme.NewRenderer(t)
		r.RenderThemeInfo()
	},
}

var themeValidateCmd = &cobra.Command{
	Use:   "validate [path-to-theme.json]",
	Short: "Validate a custom theme JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		t, err := config.LoadTheme(path)
		if err != nil {
			fmt.Println("✗ Load error:", err)
			os.Exit(1)
		}

		if err := config.ValidateTheme(t); err != nil {
			fmt.Println("✗ Validation failed:", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Theme '%s' is valid\n", t.Meta.Name)
	},
}

var themePreviewCmd = &cobra.Command{
	Use:   "preview [theme-name]",
	Short: "Live preview of a theme before applying",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		m, err := theme.NewManager(getThemesDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		t, err := m.Load(name)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		p := preview.NewPreview(t)
		p.Run()
	},
}

func init() {
	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeApplyCmd)
	themeCmd.AddCommand(themeInfoCmd)
	themeCmd.AddCommand(themeValidateCmd)
	themeCmd.AddCommand(themePreviewCmd)
	rootCmd.AddCommand(themeCmd)
}