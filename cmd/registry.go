package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/registry"
	"github.com/abhigyanwebber/cmd-customizer/internal/render"
	"github.com/abhigyanwebber/cmd-customizer/internal/tui"
	"github.com/spf13/cobra"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Browse and download themes from the community registry",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var registryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all themes in the community registry",
	Run: func(cmd *cobra.Command, args []string) {
		var index *registry.Index

		// spinner wraps the fetch
		err := tui.RunSpinner("Fetching registry...", "#00FFFF", func() error {
			var e error
			index, e = registry.FetchIndex()
			return e
		})
		if err != nil {
			fmt.Println("✗ Error:", err)
			fmt.Println("  Make sure you have an internet connection.")
			os.Exit(1)
		}

		if len(index.Themes) == 0 {
			fmt.Println("  No themes in registry yet.")
			return
		}

		// build interactive list
		items := make([]tui.RegistryItem, len(index.Themes))
		for i, t := range index.Themes {
			items[i] = tui.RegistryItem{
				Name:    t.Name,
				Author:  t.Author,
				Desc:    t.Description,
				Version: t.Version,
				Tags:    t.Tags,
			}
		}

		// show as glamour table — interactive browse is via `registry browse`
		md := fmt.Sprintf("# Community Themes\n\n%d themes available · `cmdx registry fetch <name>` to download\n\n", len(index.Themes))
		md += "| Name | Author | Description | Tags |\n"
		md += "|------|--------|-------------|------|\n"

		for _, t := range index.Themes {
			tags := strings.Join(t.Tags, ", ")
			if tags == "" {
				tags = "—"
			}
			md += fmt.Sprintf("| **%s** | %s | %s | %s |\n",
				t.Name,
				orEmDash(t.Author),
				orEmDash(t.Description),
				tags,
			)
		}
		md += fmt.Sprintf("\n_Updated: %s_\n", index.UpdatedAt)
		render.Markdown(md)
	},
}

var registryBrowseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Interactively browse and select a theme from the registry",
	Run: func(cmd *cobra.Command, args []string) {
		var index *registry.Index

		err := tui.RunSpinner("Fetching registry...", "#00FFFF", func() error {
			var e error
			index, e = registry.FetchIndex()
			return e
		})
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		if len(index.Themes) == 0 {
			fmt.Println("  No themes in registry yet.")
			return
		}

		items := make([]tui.RegistryItem, len(index.Themes))
		for i, t := range index.Themes {
			items[i] = tui.RegistryItem{
				Name:    t.Name,
				Author:  t.Author,
				Desc:    t.Description,
				Version: t.Version,
				Tags:    t.Tags,
			}
		}

		// convert to ThemeItems for the list picker
		themeItems := make([]tui.ThemeItem, len(items))
		for i, it := range items {
			themeItems[i] = tui.ThemeItem{
				Name:   it.Name,
				Desc:   it.Description(),
				Author: it.Author,
			}
		}

		chosen, err := tui.RunThemeList(themeItems, "#FF00FF", "#00FFFF", "#444444")
		if err != nil || chosen == "" {
			return
		}

		fmt.Printf("\n  Selected: %s\n  Run: cmdx registry fetch %s\n\n", chosen, chosen)
	},
}

var registryFetchCmd = &cobra.Command{
	Use:   "fetch [theme-name]",
	Short: "Download a theme from the community registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// spinner wraps the download
		err := tui.RunSpinner(
			fmt.Sprintf("Downloading '%s'...", name),
			"#00FFFF",
			func() error {
				return registry.FetchTheme(name, getThemesDir())
			},
		)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		themePath := getThemesDir() + "/" + name + ".json"
		t, err := config.LoadTheme(themePath)
		if err != nil {
			fmt.Println("✗ Downloaded theme failed validation:", err)
			os.Exit(1)
		}

		if err := config.ValidateTheme(t); err != nil {
			fmt.Println("✗ Downloaded theme is invalid:", err)
			os.Exit(1)
		}

		md := fmt.Sprintf("# ✓ Theme Downloaded\n\n**%s** is ready to use.\n\n| Command | Action |\n|---------|--------|\n| `cmdx theme preview %s` | Preview the theme |\n| `cmdx theme info %s` | Show full theme details |\n| `cmdx theme apply %s` | Apply to your terminal |\n",
			name, name, name, name)

		render.Markdown(md)
	},
}

var registrySearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the community registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		var index *registry.Index
		err := tui.RunSpinner("Searching registry...", "#00FFFF", func() error {
			var e error
			index, e = registry.FetchIndex()
			return e
		})
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		results := registry.Search(index, query)
		if len(results) == 0 {
			render.Markdown(fmt.Sprintf("# No Results\n\nNo themes found matching **%s**.\n\nTry `cmdx registry list` to see all available themes.\n", query))
			return
		}

		md := fmt.Sprintf("# Search: \"%s\"\n\n%d result(s) found\n\n", query, len(results))
		md += "| Name | Author | Version | Description | Tags |\n"
		md += "|------|--------|---------|-------------|------|\n"

		for _, t := range results {
			tags := strings.Join(t.Tags, ", ")
			if tags == "" {
				tags = "—"
			}
			md += fmt.Sprintf("| **%s** | %s | %s | %s | %s |\n",
				t.Name,
				orEmDash(t.Author),
				orEmDash(t.Version),
				orEmDash(t.Description),
				tags,
			)
		}

		render.Markdown(md)
	},
}

func orEmDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func init() {
	registryCmd.AddCommand(registryListCmd)
	registryCmd.AddCommand(registryBrowseCmd)
	registryCmd.AddCommand(registryFetchCmd)
	registryCmd.AddCommand(registrySearchCmd)
	rootCmd.AddCommand(registryCmd)
}
