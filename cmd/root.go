package main

import (
	"fmt"
	"os"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmdx",
	Short: "cmd-customizer — take control of your terminal",
	Long: `
 ██████╗███╗   ███╗██████╗ ██╗  ██╗
██╔════╝████╗ ████║██╔══██╗╚██╗██╔╝
██║     ██╔████╔██║██║  ██║ ╚███╔╝ 
██║     ██║╚██╔╝██║██║  ██║ ██╔██╗ 
╚██████╗██║ ╚═╝ ██║██████╔╝██╔╝ ██╗
 ╚═════╝╚═╝     ╚═╝╚═════╝ ╚═╝  ╚═╝

cmd-customizer — break free from the boring terminal.
Themes, prompts, spinners, cursors — all yours to control.

GitHub: https://github.com/abhigyanwebber/cmdX`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// initialize global config on every command
		config.InitGlobalConfig()
	},

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute runs the root cobra command, printing any error to stderr and
// exiting with a non-zero status on failure. Called from main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
