package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var themeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new theme interactively",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n  cmdX Theme Creator")

		// ── Step 1: Meta ──────────────────────────────
		var name, author, description string

		metaForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Theme name").
					Description("Unique identifier, no spaces (e.g. my-theme)").
					Placeholder("my-theme").
					Value(&name).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("name is required")
						}
						if strings.Contains(s, " ") {
							return fmt.Errorf("name cannot contain spaces — use hyphens (e.g. my-theme)")
						}
						return nil
					}),

				huh.NewInput().
					Title("Author").
					Placeholder("your-username").
					Value(&author),

				huh.NewInput().
					Title("Description").
					Placeholder("A short description of your theme").
					Value(&description),
			),
		)

		if err := metaForm.Run(); err != nil {
			fmt.Println("✗ Cancelled")
			return
		}

		// ── Step 2: Colors ────────────────────────────
		var primary, secondary, background, foreground, accent, errorCol, success, warning, muted string

		primary = "#00BFFF"
		secondary = "#6A9FB5"
		background = "#1E1E1E"
		foreground = "#D4D4D4"
		accent = "#569CD6"
		errorCol = "#F44747"
		success = "#4EC9B0"
		warning = "#CE9178"
		muted = "#555555"

		colorForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Primary color (hex)").Value(&primary).
					Validate(validateHex),
				huh.NewInput().Title("Secondary color (hex)").Value(&secondary).
					Validate(validateHex),
				huh.NewInput().Title("Background color (hex)").Value(&background).
					Validate(validateHex),
				huh.NewInput().Title("Foreground color (hex)").Value(&foreground).
					Validate(validateHex),
				huh.NewInput().Title("Accent color (hex)").Value(&accent).
					Validate(validateHex),
			),
			huh.NewGroup(
				huh.NewInput().Title("Error color (hex)").Value(&errorCol).
					Validate(validateHex),
				huh.NewInput().Title("Success color (hex)").Value(&success).
					Validate(validateHex),
				huh.NewInput().Title("Warning color (hex)").Value(&warning).
					Validate(validateHex),
				huh.NewInput().Title("Muted color (hex)").Value(&muted).
					Validate(validateHex),
			),
		)

		if err := colorForm.Run(); err != nil {
			fmt.Println("✗ Cancelled")
			return
		}

		// ── Step 3: Prompt ────────────────────────────
		var promptSymbol, promptStyle string
		promptSymbol = "❯"
		promptStyle = "single"

		promptForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Prompt symbol").
					Description("Character at the end of your prompt").
					Value(&promptSymbol),

				huh.NewSelect[string]().
					Title("Prompt style").
					Options(
						huh.NewOption("Single line", "single"),
						huh.NewOption("Multiline", "multiline"),
					).
					Value(&promptStyle),
			),
		)

		if err := promptForm.Run(); err != nil {
			fmt.Println("✗ Cancelled")
			return
		}

		// ── Step 4: Loader ────────────────────────────
		var loaderStyle, intervalStr string
		loaderStyle = "braille"
		intervalStr = "80"

		loaderForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Loader style").
					Options(
						huh.NewOption("Braille (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)", "braille"),
						huh.NewOption("Circle (◐◓◑◒)", "circle"),
						huh.NewOption("Classic (|/-\\)", "classic"),
						huh.NewOption("Dots (... ....)", "dots"),
						huh.NewOption("Stars (✦✧)", "stars"),
					).
					Value(&loaderStyle),

				huh.NewInput().
					Title("Animation speed (ms)").
					Description("Lower = faster (e.g. 80)").
					Value(&intervalStr).
					Validate(func(s string) error {
						v, err := strconv.Atoi(s)
						if err != nil || v <= 0 {
							return fmt.Errorf("must be a positive number")
						}
						return nil
					}),
			),
		)

		if err := loaderForm.Run(); err != nil {
			fmt.Println("✗ Cancelled")
			return
		}

		// ── Step 5: Cursor & Effects ──────────────────
		var cursorStyle, bannerEffect, dividerStyle string
		var bannerEnabled bool
		var bannerText string

		cursorStyle = "bar"
		bannerEffect = "none"
		dividerStyle = "line"
		bannerEnabled = true
		bannerText = "Welcome, {user}"

		styleForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Cursor style").
					Options(
						huh.NewOption("Bar |", "bar"),
						huh.NewOption("Block █", "block"),
						huh.NewOption("Underline _", "underline"),
					).
					Value(&cursorStyle),

				huh.NewSelect[string]().
					Title("Divider style").
					Options(
						huh.NewOption("Line ─────", "line"),
						huh.NewOption("Wave ~-~-~", "wave"),
						huh.NewOption("Dots · · ·", "dots"),
						huh.NewOption("Zigzag /\\/\\", "zigzag"),
						huh.NewOption("Stars ✦ ✦ ✦", "stars"),
						huh.NewOption("Double ═════", "double"),
					).
					Value(&dividerStyle),
			),
			huh.NewGroup(
				huh.NewConfirm().
					Title("Enable startup banner?").
					Value(&bannerEnabled),

				huh.NewInput().
					Title("Banner text").
					Description("Use {user} for username").
					Value(&bannerText),

				huh.NewSelect[string]().
					Title("Banner effect").
					Options(
						huh.NewOption("None", "none"),
						huh.NewOption("Neon glow", "neon"),
						huh.NewOption("Glitch", "glitch"),
						huh.NewOption("Rainbow", "rainbow"),
						huh.NewOption("Typewriter", "typewriter"),
					).
					Value(&bannerEffect),
			),
		)

		if err := styleForm.Run(); err != nil {
			fmt.Println("✗ Cancelled")
			return
		}

		// ── Step 6: Gradient ──────────────────────────
		var gradientEnabled bool
		var gradientFrom, gradientTo string
		gradientEnabled = true
		gradientFrom = primary
		gradientTo = secondary

		gradientForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Enable gradient effects?").
					Value(&gradientEnabled),

				huh.NewInput().
					Title("Gradient from (hex)").
					Value(&gradientFrom).
					Validate(validateHex),

				huh.NewInput().
					Title("Gradient to (hex)").
					Value(&gradientTo).
					Validate(validateHex),
			),
		)

		if err := gradientForm.Run(); err != nil {
			fmt.Println("✗ Cancelled")
			return
		}

		// ── Build Theme ───────────────────────────────
		intervalMs, _ := strconv.Atoi(intervalStr)
		loaderFrames := getLoaderFrames(loaderStyle)

		promptFormat := "{dir} {git} {symbol} "
		if promptStyle == "multiline" {
			promptFormat = "{user} @ {dir} {git}\n{symbol} "
		}

		theme := config.Theme{
			Meta: config.Meta{
				Name:        name,
				Version:     "1.0.0",
				Author:      author,
				Description: description,
			},
			Colors: config.Colors{
				Primary:    primary,
				Secondary:  secondary,
				Background: background,
				Foreground: foreground,
				Accent:     accent,
				Error:      errorCol,
				Success:    success,
				Warning:    warning,
				Muted:      muted,
			},
			Prompt: config.Prompt{
				Symbol:    promptSymbol,
				Separator: "›",
				Style:     promptStyle,
				Segments:  []string{"directory", "git"},
				Format:    promptFormat,
			},
			Loader: config.Loader{
				Frames:     loaderFrames,
				IntervalMs: intervalMs,
				Color:      "primary",
			},
			ProgressBar: config.ProgressBar{
				Filled: "█",
				Empty:  "░",
				Width:  40,
				Color:  "primary",
			},
			Cursor: config.Cursor{
				Style: cursorStyle,
				Blink: true,
				Color: "primary",
			},
			Borders: config.Borders{
				Style: "rounded",
				Chars: config.BorderChars{
					TopLeft:     "╭",
					TopRight:    "╮",
					BottomLeft:  "╰",
					BottomRight: "╯",
					Horizontal:  "─",
					Vertical:    "│",
				},
			},
			Banner: config.Banner{
				Enabled: bannerEnabled,
				Text:    bannerText,
				Style:   "plain",
				Color:   "primary",
			},
			Graphics: config.Graphics{
				Gradient: config.GradientConfig{
					Enabled:   gradientEnabled,
					From:      gradientFrom,
					To:        gradientTo,
					Direction: "horizontal",
				},
				Divider: config.DividerConfig{
					Style: dividerStyle,
					Color: primary,
				},
				Effects: config.EffectsConfig{
					Banner: bannerEffect,
					Prompt: "none",
				},
				Icons: config.IconsConfig{
					Enabled:   true,
					Directory: "󰉋",
					GitBranch: "",
					Error:     "",
					Success:   "",
					Time:      "",
				},
			},
		}

		// validate
		if err := config.ValidateTheme(&theme); err != nil {
			fmt.Println("✗ Theme validation failed:", err)
			os.Exit(1)
		}

		// save
		outPath := filepath.Join(getThemesDir(), name+".json")
		data, err := json.MarshalIndent(theme, "", "  ")
		if err != nil {
			fmt.Println("✗ Could not serialize theme:", err)
			os.Exit(1)
		}

		if err := os.WriteFile(outPath, data, 0644); err != nil {
			fmt.Println("✗ Could not save theme:", err)
			os.Exit(1)
		}

		fmt.Printf("\n✓ Theme '%s' created at %s\n", name, outPath)
		fmt.Printf("  Preview: cmdx theme preview %s\n", name)
		fmt.Printf("  Apply:   cmdx theme apply %s\n\n", name)
	},
}

func validateHex(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("color is required")
	}
	if s[0] != '#' || (len(s) != 7 && len(s) != 4) {
		return fmt.Errorf("must be a valid hex color (e.g. #FF00FF)")
	}
	return nil
}

func getLoaderFrames(style string) []string {
	switch style {
	case "braille":
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	case "circle":
		return []string{"◐", "◓", "◑", "◒"}
	case "classic":
		return []string{"|", "/", "-", "\\"}
	case "dots":
		return []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	case "stars":
		return []string{"✦", "✧", "✦", "✧", "·", "✦"}
	default:
		return []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	}
}

func init() {
	themeCmd.AddCommand(themeCreateCmd)
}
