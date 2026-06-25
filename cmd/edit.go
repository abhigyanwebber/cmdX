package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/editor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [theme-name]",
	Short: "Open the interactive TUI theme editor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		themePath := filepath.Join(getThemesDir(), name+".json")

		if _, err := os.Stat(themePath); os.IsNotExist(err) {
			fmt.Printf("✗ Theme '%s' not found.\n", name)
			fmt.Println("  Run 'cmdx theme list' to see available themes.")
			os.Exit(1)
		}

		t, err := config.LoadTheme(themePath)
		if err != nil {
			fmt.Println("✗ Could not load theme:", err)
			os.Exit(1)
		}

		m := editor.NewModel(t, themePath)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Println("✗ Editor error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
