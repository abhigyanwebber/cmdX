// Theme commands: list, apply, info, validate, preview, inject, remove.
// Shared helpers (detectShell, getThemesDir, loadThemeOrExit) live in
// cmd/helpers.go since they're also used by create.go, edit.go, and
// registry.go.
package main

import (
	"fmt"
	"os"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/preview"
	"github.com/abhigyanwebber/cmd-customizer/internal/theme"
	"github.com/abhigyanwebber/cmd-customizer/internal/tui"
	"github.com/spf13/cobra"
)

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

		items := buildThemeListItems(m, names)

		chosen, err := tui.RunThemeList(items, "#FF00FF", "#00FFFF", "#444444")
		if err != nil {
			printPlainThemeList(items)
			return
		}

		if chosen != "" {
			fmt.Printf("\n  Selected: %s\n", chosen)
			fmt.Printf("  Run: cmdx theme apply %s\n\n", chosen)
		}
	},
}

// buildThemeListItems loads each theme's metadata for display in the
// interactive picker, falling back to a placeholder description for
// themes that fail to load/validate rather than dropping them silently.
func buildThemeListItems(m *theme.Manager, names []string) []tui.ThemeItem {
	items := make([]tui.ThemeItem, 0, len(names))
	for _, name := range names {
		t, err := m.Load(name)
		desc := "(invalid theme)"
		author := ""
		if err == nil {
			desc = t.Meta.Description
			author = t.Meta.Author
		}
		items = append(items, tui.ThemeItem{
			Name:   name,
			Desc:   desc,
			Author: author,
		})
	}
	return items
}

// printPlainThemeList is the non-interactive fallback used when the TUI
// list picker can't run (e.g. no TTY).
func printPlainThemeList(items []tui.ThemeItem) {
	fmt.Println("\n  Available Themes:")
	for _, it := range items {
		fmt.Printf("  %-20s  %s\n", it.Name, it.Desc)
	}
	fmt.Println()
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
		_, t := loadThemeOrExit(args[0])
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
		_, t := loadThemeOrExit(args[0])
		p := preview.NewPreview(t)
		p.Run()
		showActiveFloaters()
	},
}

var themeInjectCmd = &cobra.Command{
	Use:   "inject [theme-name]",
	Short: "Inject a theme into your shell config (persists after restart)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, t := loadThemeOrExit(args[0])

		shell, shellName := detectShell()
		if shell == nil {
			fmt.Println("✗ Could not detect shell. Use --shell to specify one.")
			fmt.Println("  Options: powershell, zsh, bash")
			os.Exit(1)
		}

		fmt.Printf("  Detected shell: %s\n", shellName)

		injected, _ := shell.IsInjected()
		if injected {
			fmt.Printf("! Theme already injected. Overwriting with '%s'...\n", args[0])
		}

		if err := shell.Inject(t); err != nil {
			fmt.Println("✗ Injection failed:", err)
			os.Exit(1)
		}

		path, _ := shell.ProfilePath()
		fmt.Printf("✓ Theme '%s' injected into %s\n", args[0], path)
		fmt.Println("  Restart your terminal to see changes.")
	},
}

var themeRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove cmdx theme injection from your shell config",
	Run: func(cmd *cobra.Command, args []string) {
		shell, shellName := detectShell()
		if shell == nil {
			fmt.Println("✗ Could not detect shell.")
			os.Exit(1)
		}

		fmt.Printf("  Detected shell: %s\n", shellName)

		injected, _ := shell.IsInjected()
		if !injected {
			fmt.Println("! No cmdx theme found in shell config.")
			return
		}

		if err := shell.Remove(); err != nil {
			fmt.Println("✗ Failed to remove:", err)
			os.Exit(1)
		}

		fmt.Println("✓ cmdx theme removed from shell config.")
		fmt.Println("  Restart your terminal to revert to default.")
	},
}

func init() {
	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeApplyCmd)
	themeCmd.AddCommand(themeInfoCmd)
	themeCmd.AddCommand(themeValidateCmd)
	themeCmd.AddCommand(themePreviewCmd)
	themeCmd.AddCommand(themeInjectCmd)
	themeCmd.AddCommand(themeRemoveCmd)
	rootCmd.AddCommand(themeCmd)
}
