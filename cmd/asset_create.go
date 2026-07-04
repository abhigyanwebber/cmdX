package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abhigyanwebber/cmd-customizer/internal/assets"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var assetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Scaffold a new asset with an interactive wizard",
	Long: `Scaffold a new asset folder with a pre-filled asset.json manifest.

Supports all asset types:
  floater  — decorative corner PNG (top-left, top-right, bottom-left, bottom-right)
  spinner  — animated PNG sequence for loading states
  banner   — startup banner graphic shown when the theme is injected
  divider  — horizontal separator graphic
  icon     — named icon set (directory, git, error, success, etc.)

The wizard generates a ready-to-edit asset.json and a placeholder folder
structure. Drop your PNG files in, adjust the manifest, then import with:
  cmdx asset import ./my-asset/`,
	Run: func(cmd *cobra.Command, args []string) {
		var assetType string
		var name string
		var author string
		var description string
		var renderMode string
		var colorMode string
		var width int
		var height int

		// ── Step 1: Type ──────────────────────────────────────
		typeForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Asset type").
					Description("What kind of asset are you creating?").
					Options(
						huh.NewOption("Floater — decorative corner PNG", "floater"),
						huh.NewOption("Mascot — reactive character with state machine", "mascot"),
						huh.NewOption("Spinner — animated loading sequence", "spinner"),
						huh.NewOption("Banner — startup graphic", "banner"),
						huh.NewOption("Divider — horizontal separator", "divider"),
						huh.NewOption("Icon set — named icon collection", "icon"),
					).
					Value(&assetType),
			),
		)

		if err := typeForm.Run(); err != nil {
			fmt.Println("✗ Cancelled")
			return
		}

		// ── Step 2: Identity ──────────────────────────────────
		metaForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Asset name").
					Description("Lowercase, hyphens only (e.g. 'corner-cat', 'my-spinner')").
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("name is required")
						}
						if strings.ContainsAny(s, " /\\:.") {
							return fmt.Errorf("name must not contain spaces, slashes, or dots")
						}
						return nil
					}).
					Value(&name),
				huh.NewInput().
					Title("Author").
					Description("Your name or GitHub handle").
					Value(&author),
				huh.NewInput().
					Title("Description").
					Description("One line describing what this asset looks like or does").
					Value(&description),
			),
		)

		if err := metaForm.Run(); err != nil {
			fmt.Println("✗ Cancelled")
			return
		}

		// ── Step 3: Render config ────────────────────────────
		renderForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Render mode").
					Description("How should chafa convert your PNG to terminal output?").
					Options(
						huh.NewOption("braille  — fine detail, smooth edges (recommended)", "braille"),
						huh.NewOption("blocks   — bold, chunky, retro look", "blocks"),
						huh.NewOption("ascii    — no color, pure characters", "ascii"),
						huh.NewOption("sixel    — actual pixels (Windows Terminal / iTerm2 only)", "sixel"),
						huh.NewOption("color    — full 24-bit color blocks", "color"),
					).
					Value(&renderMode),
				huh.NewSelect[string]().
					Title("Color depth").
					Options(
						huh.NewOption("truecolor — full 24-bit RGB (recommended)", "truecolor"),
						huh.NewOption("256       — xterm 256 color palette", "256"),
						huh.NewOption("ansi      — 16-color ANSI palette", "ansi"),
						huh.NewOption("none      — no color (monochrome)", "none"),
					).
					Value(&colorMode),
			),
		)

		if err := renderForm.Run(); err != nil {
			fmt.Println("✗ Cancelled")
			return
		}

		// ── Step 4: Size (type-appropriate defaults) ─────────
		defaultWidth, defaultHeight := defaultSizeForType(assets.AssetType(assetType), assets.RenderMode(renderMode))
		width = defaultWidth
		height = defaultHeight

		// ── Generate manifest ─────────────────────────────────
		a := buildManifest(assetType, name, author, description, renderMode, colorMode, width, height)

		// ── Write to disk ─────────────────────────────────────
		outDir := "./" + name
		if err := os.MkdirAll(outDir, 0755); err != nil {
			fmt.Printf("✗ Could not create directory: %v\n", err)
			os.Exit(1)
		}

		data, err := json.MarshalIndent(a, "", "  ")
		if err != nil {
			fmt.Printf("✗ Could not serialize manifest: %v\n", err)
			os.Exit(1)
		}

		manifestPath := filepath.Join(outDir, "asset.json")
		if err := os.WriteFile(manifestPath, data, 0644); err != nil {
			fmt.Printf("✗ Could not write asset.json: %v\n", err)
			os.Exit(1)
		}

		// Write placeholder README
		readme := buildAssetReadme(assetType, name, author, description)
		os.WriteFile(filepath.Join(outDir, "README.md"), []byte(readme), 0644)

		fmt.Printf("\n✓ Asset scaffolded at: %s\n\n", outDir)
		fmt.Println("  Next steps:")
		printNextSteps(assetType, name, outDir)
	},
}

func defaultSizeForType(assetType assets.AssetType, mode assets.RenderMode) (int, int) {
	switch assetType {
	case assets.AssetTypeFloater:
		return 12, 8
	case assets.AssetTypeBanner:
		if mode == assets.RenderModeSixel {
			return 60, 12
		}
		return 40, 6
	case assets.AssetTypeDivider:
		return 60, 2
	case assets.AssetTypeSpinner:
		return 8, 4
	case assets.AssetTypeIcon:
		return 4, 2
	default:
		return 16, 8
	}
}

func buildManifest(assetType, name, author, description, renderMode, colorMode string, width, height int) map[string]interface{} {
	base := map[string]interface{}{
		"name":        name,
		"type":        assetType,
		"version":     "1.0.0",
		"author":      author,
		"description": description,
		"render": map[string]interface{}{
			"mode":       renderMode,
			"color_mode": colorMode,
			"width":      width,
			"height":     height,
			"dither":     false,
			"stretch":    false,
			"threshold":  0.5,
		},
	}

	switch assetType {
	case "floater":
		base["floater"] = map[string]interface{}{
			"file":           "floater.png",
			"position":       "top-right",
			"max_width":      width,
			"max_height":     height,
			"margin_x":       1,
			"margin_y":       1,
			"animate_frames": []string{},
			"interval_ms":    0,
		}
	case "mascot":
		base["mascot"] = map[string]interface{}{
			"position":           "bottom-right",
			"max_width":          width,
			"max_height":         height,
			"margin_x":           2,
			"margin_y":           1,
			"default_state":      "idle",
			"global_interval_ms": 120,
			"hook_var":           "CMDX_MASCOT_TRIGGER",
			"states": map[string]interface{}{
				"idle": map[string]interface{}{
					"frames":      []string{"idle_01.png", "idle_02.png"},
					"interval_ms": 200,
					"triggers": []map[string]interface{}{
						{"type": "always", "priority": 0},
					},
				},
				"working": map[string]interface{}{
					"frames":      []string{"working_01.png", "working_02.png", "working_03.png"},
					"interval_ms": 80,
					"triggers": []map[string]interface{}{
						{"type": "command", "value": "go build*", "priority": 10},
						{"type": "command", "value": "npm*", "priority": 10},
						{"type": "command", "value": "make*", "priority": 10},
					},
				},
				"success": map[string]interface{}{
					"frames":           []string{"success_01.png"},
					"transition_frames": []string{"success_transition.png"},
					"triggers": []map[string]interface{}{
						{"type": "exit_code", "value": "0", "priority": 20},
					},
					"render_override": map[string]interface{}{
						"tint": "#00FF88",
					},
				},
				"error": map[string]interface{}{
					"frames":           []string{"error_01.png", "error_02.png"},
					"transition_frames": []string{"error_transition.png"},
					"triggers": []map[string]interface{}{
						{"type": "exit_code", "value": "1-127", "priority": 20},
						{"type": "output_regex", "value": "error|fatal|panic|FAILED", "priority": 15},
					},
					"render_override": map[string]interface{}{
						"tint": "#FF4444",
					},
				},
				"warning": map[string]interface{}{
					"frames": []string{"warning_01.png"},
					"triggers": []map[string]interface{}{
						{"type": "exit_code", "value": "128+", "priority": 18},
						{"type": "git_status", "value": "dirty", "priority": 5},
					},
					"render_override": map[string]interface{}{
						"tint": "#FFA500",
					},
				},
				"sleeping": map[string]interface{}{
					"frames":      []string{"sleep_01.png", "sleep_02.png"},
					"interval_ms": 800,
					"triggers": []map[string]interface{}{
						{"type": "idle_time", "value": "30", "priority": 1},
					},
					"render_override": map[string]interface{}{
						"width":  width / 2,
						"height": height / 2,
					},
				},
			},
		}
	case "spinner":
		base["spinner"] = map[string]interface{}{
			"frames":      []string{"frame_01.png", "frame_02.png", "frame_03.png", "frame_04.png"},
			"interval_ms": 100,
			"reverse":     false,
			"bounce":      false,
			"loop":        true,
			"on_complete": "persist",
		}
	case "banner":
		base["banner"] = map[string]interface{}{
			"file":       "banner.png",
			"position":   "center",
			"max_width":  width,
			"max_height": height,
		}
	case "divider":
		base["divider"] = map[string]interface{}{
			"file":      "divider.png",
			"tile_mode": "repeat",
			"height":    height,
		}
	case "icon":
		base["icon"] = map[string]interface{}{
			"files": map[string]string{
				"directory": "dir.png",
				"file":      "file.png",
				"git":       "git.png",
				"error":     "error.png",
				"success":   "success.png",
				"warning":   "warning.png",
			},
		}
	}

	return base
}

func printNextSteps(assetType, name, outDir string) {
	switch assetType {
	case "mascot":
		fmt.Printf("  1. Add PNG files for each state in %s/\n", outDir)
		fmt.Printf("     idle_01.png, working_01.png, success_01.png, error_01.png ...\n")
		fmt.Printf("  2. Edit asset.json — customize triggers, tints, render overrides\n")
		fmt.Printf("  3. cmdx asset import %s/ && cmdx asset use %s --as mascot\n", outDir, name)
		fmt.Printf("  4. cmdx asset mascot-hooks %s  # print shell hook code\n", name)
		fmt.Printf("  5. cmdx asset preview %s --state error  # preview a specific state\n", name)
	case "floater":
		fmt.Printf("  1. Add your PNG file as: %s/floater.png\n", outDir)
		fmt.Printf("  2. Edit %s/asset.json — set position (top-left, top-right, bottom-left, bottom-right)\n", outDir)
		fmt.Printf("  3. Optional: add animate_frames for an animated floater\n")
	case "spinner":
		fmt.Printf("  1. Add your PNG frames as: %s/frame_01.png, frame_02.png, ...\n", outDir)
		fmt.Printf("  2. Edit %s/asset.json — list frame filenames in order\n", outDir)
		fmt.Printf("  3. Adjust interval_ms to control animation speed\n")
	case "banner":
		fmt.Printf("  1. Add your banner PNG as: %s/banner.png\n", outDir)
		fmt.Printf("  2. Edit %s/asset.json — adjust max_width and position\n", outDir)
	case "divider":
		fmt.Printf("  1. Add your divider PNG as: %s/divider.png\n", outDir)
		fmt.Printf("  2. Edit %s/asset.json — set tile_mode (repeat, stretch, center)\n", outDir)
	case "icon":
		fmt.Printf("  1. Add PNG files for each icon key in %s/\n", outDir)
		fmt.Printf("  2. Edit %s/asset.json — map icon keys to filenames\n", outDir)
	}
	fmt.Printf("\n  Test render:  cmdx asset render %s/<your-file>.png\n", outDir)
	fmt.Printf("  Validate:     cmdx asset validate %s\n", outDir)
	fmt.Printf("  Import:       cmdx asset import %s\n", outDir)
	fmt.Printf("  Preview:      cmdx asset preview %s\n\n", name)
}

func buildAssetReadme(assetType, name, author, description string) string {
	codeBlock := "```"
	return fmt.Sprintf("# %s\n\n> %s\n\n**Type:** %s  \n**Author:** %s  \n**Compatible with:** cmdX v1.0.0+\n\n## Installation\n\n%sbash\ncmdx asset import ./%s\ncmdx asset use %s --as %s\n%s\n\n## Customization\n\nOverride render options at runtime without editing asset.json:\n\n%sbash\n# Preview with different render mode\ncmdx asset preview %s --mode sixel --width 20\n\n# Preview in a different corner position (floaters only)\ncmdx asset preview %s --position bottom-left\n%s\n\n## asset.json\n\nSee the generated asset.json for the full configuration schema.\nAll options can be overridden at the CLI level using flags on 'cmdx asset preview'.\n",
		name, description, assetType, author,
		codeBlock, name, name, assetType, codeBlock,
		codeBlock, name, name, codeBlock,
	)
}

func init() {
	assetCmd.AddCommand(assetCreateCmd)
}
