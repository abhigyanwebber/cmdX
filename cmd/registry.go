package main

import (
	"fmt"
	"os"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/registry"
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
		fmt.Println("\n  Fetching registry...")

		index, err := registry.FetchIndex()
		if err != nil {
			fmt.Println("✗ Error:", err)
			fmt.Println("  Make sure you have an internet connection.")
			os.Exit(1)
		}

		if len(index.Themes) == 0 {
			fmt.Println("  No themes in registry yet.")
			return
		}

		fmt.Printf("\n  Community Themes (%d available):\n\n", len(index.Themes))
		for _, t := range index.Themes {
			tags := ""
			for i, tag := range t.Tags {
				if i > 0 {
					tags += ", "
				}
				tags += tag
			}
			fmt.Printf("  %-20s  %s\n", t.Name, t.Description)
			fmt.Printf("  %-20s  by %s  [%s]\n\n", "", t.Author, tags)
		}
	},
}

var registryFetchCmd = &cobra.Command{
	Use:   "fetch [theme-name]",
	Short: "Download a theme from the community registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		fmt.Printf("  Downloading theme '%s'...\n", name)

		if err := registry.FetchTheme(name, getThemesDir()); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		// validate the downloaded theme
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

		fmt.Printf("✓ Theme '%s' downloaded successfully\n", name)
		fmt.Printf("  Run 'cmdx theme preview %s' to preview it\n", name)
		fmt.Printf("  Run 'cmdx theme apply %s' to apply it\n", name)
	},
}

var registrySearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the community registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		fmt.Printf("\n  Searching registry for '%s'...\n\n", query)

		index, err := registry.FetchIndex()
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		results := registry.Search(index, query)
		if len(results) == 0 {
			fmt.Printf("  No themes found for '%s'\n\n", query)
			return
		}

		fmt.Printf("  Found %d result(s):\n\n", len(results))
		for _, t := range results {
			fmt.Printf("  %-20s  %s\n", t.Name, t.Description)
			fmt.Printf("  %-20s  by %s  v%s\n\n", "", t.Author, t.Version)
		}
	},
}

func init() {
	registryCmd.AddCommand(registryListCmd)
	registryCmd.AddCommand(registryFetchCmd)
	registryCmd.AddCommand(registrySearchCmd)
	rootCmd.AddCommand(registryCmd)
}
