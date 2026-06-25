package main

import (
	"fmt"
	"os"

	"github.com/abhigyanwebber/cmd-customizer/internal/wallpaper"
	"github.com/spf13/cobra"
)

var wallpaperCmd = &cobra.Command{
	Use:   "wallpaper",
	Short: "Manage terminal background wallpaper",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var wallpaperSetCmd = &cobra.Command{
	Use:   "set [image-path]",
	Short: "Set a background image for your terminal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]

		opacity, _ := cmd.Flags().GetFloat64("opacity")
		stretch, _ := cmd.Flags().GetString("stretch")
		alignment, _ := cmd.Flags().GetString("alignment")

		engine := wallpaper.NewEngine()

		fmt.Printf("  Terminal: %s\n", engine.DetectedTerminal())
		fmt.Printf("  Image:    %s\n", imagePath)
		fmt.Printf("  Opacity:  %.2f\n", opacity)
		fmt.Println()

		if err := engine.Apply(imagePath, opacity, stretch, alignment); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		fmt.Println("✓ Wallpaper applied successfully")
		fmt.Println("  Restart your terminal to see changes.")
	},
}

var wallpaperRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove the terminal background wallpaper",
	Run: func(cmd *cobra.Command, args []string) {
		engine := wallpaper.NewEngine()

		fmt.Printf("  Terminal: %s\n", engine.DetectedTerminal())

		if err := engine.Remove(); err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		fmt.Println("✓ Wallpaper removed")
		fmt.Println("  Restart your terminal to see changes.")
	},
}

var wallpaperInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show wallpaper support info for your terminal",
	Run: func(cmd *cobra.Command, args []string) {
		engine := wallpaper.NewEngine()

		fmt.Println()
		fmt.Printf("  OS:       %s\n", engine.OS)
		fmt.Printf("  Terminal: %s\n", engine.DetectedTerminal())
		fmt.Println()
		fmt.Println("  Supported features:")
		fmt.Println("    ✓ Background image")
		fmt.Println("    ✓ Opacity control")
		fmt.Println("    ✓ Stretch modes (fill, uniform, uniformToFill, none)")
		fmt.Println("    ✓ Alignment control")
		fmt.Println()
	},
}

func init() {
	wallpaperSetCmd.Flags().Float64P("opacity", "o", 0.3, "Background opacity (0.0 to 1.0)")
	wallpaperSetCmd.Flags().StringP("stretch", "s", "uniformToFill", "Stretch mode: fill, uniform, uniformToFill, none")
	wallpaperSetCmd.Flags().StringP("alignment", "a", "center", "Alignment: center, topLeft, topRight, bottomLeft, bottomRight")

	wallpaperCmd.AddCommand(wallpaperSetCmd)
	wallpaperCmd.AddCommand(wallpaperRemoveCmd)
	wallpaperCmd.AddCommand(wallpaperInfoCmd)
	rootCmd.AddCommand(wallpaperCmd)
}
