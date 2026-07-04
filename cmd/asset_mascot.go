package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abhigyanwebber/cmd-customizer/internal/assets"
	"github.com/spf13/cobra"
)

// assetMascotStateCmd is called by shell hooks after every command to
// resolve and display the current mascot state. It reads exit code,
// command, git status etc. from flags and env vars, resolves the state,
// and renders the appropriate PNG.
var assetMascotStateCmd = &cobra.Command{
	Use:    "mascot-state [asset-name]",
	Short:  "Resolve and display the current mascot state (called by shell hooks)",
	Hidden: true, // not shown in help — this is for shell hook use
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		exitCode, _ := cmd.Flags().GetInt("exit-code")
		command, _ := cmd.Flags().GetString("command")
		lastOutput, _ := cmd.Flags().GetString("output")
		idleSeconds, _ := cmd.Flags().GetInt("idle-seconds")
		gitStatus, _ := cmd.Flags().GetString("git-status")
		stateOverride, _ := cmd.Flags().GetString("state")
		modeStr, _ := cmd.Flags().GetString("mode")
		colorStr, _ := cmd.Flags().GetString("color")
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")

		if !assets.ChafaAvailable() {
			os.Exit(0) // silently skip if chafa not available
		}

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			os.Exit(0)
		}

		// if a specific state is forced via flag, skip resolution
		ctx := assets.MascotContext{
			LastExitCode: exitCode,
			LastCommand:  command,
			LastOutput:   lastOutput,
			IdleSeconds:  idleSeconds,
			GitStatus:    gitStatus,
			Env:          buildEnvSnapshot(),
		}

		var state assets.MascotState
		if stateOverride != "" {
			state = assets.MascotState(stateOverride)
		} else {
			state, err = m.ResolveMascotState(name, ctx)
			if err != nil {
				os.Exit(0)
			}
		}

		overrides := assets.RenderOverrides{
			Mode:      assets.RenderMode(modeStr),
			ColorMode: assets.ColorMode(colorStr),
			Width:     width,
			Height:    height,
		}

		// display state name for shell to optionally show
		fmt.Printf("\033[s") // save cursor position
		if err := m.PreviewMascot(name, ctx, overrides); err != nil {
			os.Exit(0)
		}
		_ = state
	},
}

// assetMascotInfoCmd shows the full state machine of a mascot asset —
// all states, their triggers, frame counts, and render overrides.
var assetMascotInfoCmd = &cobra.Command{
	Use:   "mascot-info [asset-name]",
	Short: "Show the full state machine of a mascot asset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		a, _, err := m.Get(name, assets.AssetTypeMascot)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		if a.Mascot == nil {
			fmt.Println("✗ Asset has no mascot config")
			os.Exit(1)
		}

		mc := a.Mascot
		fmt.Printf("\n  Mascot: %s\n", name)
		fmt.Printf("  Position:      %s\n", mc.Position)
		fmt.Printf("  Default state: %s\n", mc.DefaultState)
		fmt.Printf("  Max size:      %dx%d\n", mc.MaxWidth, mc.MaxHeight)
		fmt.Printf("  Global FPS:    %dms/frame\n", mc.GlobalIntervalMs)
		fmt.Printf("  Hook var:      %s\n\n", hookVarOrDefault(mc.HookVar))

		fmt.Printf("  States (%d):\n\n", len(mc.States))

		for stateName, state := range mc.States {
			frameCount := len(state.Frames)
			transCount := len(state.TransitionFrames)
			interval := state.IntervalMs
			if interval == 0 {
				interval = mc.GlobalIntervalMs
			}

			fmt.Printf("  ┌─ %s\n", stateName)
			fmt.Printf("  │  Frames:     %d", frameCount)
			if transCount > 0 {
				fmt.Printf(" (+ %d transition)", transCount)
			}
			fmt.Println()
			fmt.Printf("  │  Speed:      %dms/frame\n", interval)

			if state.RenderOverride != nil {
				ro := state.RenderOverride
				if ro.Width > 0 || ro.Height > 0 {
					fmt.Printf("  │  Size:       %dx%d (override)\n", ro.Width, ro.Height)
				}
				if ro.Tint != "" {
					fmt.Printf("  │  Tint:       %s\n", ro.Tint)
				}
			}

			if len(state.Triggers) > 0 {
				fmt.Printf("  │  Triggers:\n")
				for _, t := range state.Triggers {
					if t.Value != "" {
						fmt.Printf("  │    • %s = %s (priority %d)\n", t.Type, t.Value, t.Priority)
					} else {
						fmt.Printf("  │    • %s (priority %d)\n", t.Type, t.Priority)
					}
				}
			} else {
				fmt.Printf("  │  Triggers:   (none — manual only)\n")
			}
			fmt.Printf("  └─\n\n")
		}
	},
}

// assetMascotHooksCmd prints the shell hook code for a mascot to stdout
// so users can inspect or manually add it to their profile.
var assetMascotHooksCmd = &cobra.Command{
	Use:   "mascot-hooks [asset-name]",
	Short: "Print the shell hook code for a mascot asset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		shell, _ := cmd.Flags().GetString("shell")

		if shell == "" {
			sh, shellName := detectShell()
			if sh == nil {
				fmt.Println("✗ Could not detect shell. Use --shell to specify one.")
				fmt.Println("  Options: bash, zsh, powershell")
				os.Exit(1)
			}
			shell = shellName
		}

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		hooks, err := m.MascotHooks(name, shell)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		fmt.Printf("# Shell hooks for mascot '%s' (%s)\n", name, shell)
		fmt.Printf("# Add this to your shell profile, or run:\n")
		fmt.Printf("#   cmdx asset use %s --as mascot\n\n", name)
		fmt.Println(hooks)
	},
}

// buildEnvSnapshot captures relevant env vars for trigger resolution.
func buildEnvSnapshot() map[string]string {
	env := map[string]string{}
	relevant := []string{
		"CMDX_MASCOT_TRIGGER", "CMDX_MASCOT_STATE",
		"GIT_BRANCH", "VIRTUAL_ENV", "CONDA_DEFAULT_ENV",
		"NODE_ENV", "GOPATH", "CARGO_HOME",
	}
	for _, key := range relevant {
		if val := os.Getenv(key); val != "" {
			env[key] = val
		}
	}
	return env
}

func hookVarOrDefault(v string) string {
	if v == "" {
		return "CMDX_MASCOT_TRIGGER"
	}
	return v
}

func init() {
	// mascot-state flags (for shell hooks)
	assetMascotStateCmd.Flags().Int("exit-code", 0, "Exit code of the last command")
	assetMascotStateCmd.Flags().String("command", "", "The last command run")
	assetMascotStateCmd.Flags().String("output", "", "Stdout/stderr of the last command (for regex triggers)")
	assetMascotStateCmd.Flags().Int("idle-seconds", 0, "Seconds since last command (for idle triggers)")
	assetMascotStateCmd.Flags().String("git-status", "", "Git status: dirty, clean, untracked, ahead, behind")
	assetMascotStateCmd.Flags().String("state", "", "Force a specific state (skips trigger resolution)")
	assetMascotStateCmd.Flags().String("mode", "", "Override render mode")
	assetMascotStateCmd.Flags().String("color", "", "Override color depth")
	assetMascotStateCmd.Flags().Int("width", 0, "Override width")
	assetMascotStateCmd.Flags().Int("height", 0, "Override height")

	// mascot-hooks flags
	assetMascotHooksCmd.Flags().String("shell", "", "Shell type: bash, zsh, powershell (default: auto-detect)")

	assetCmd.AddCommand(assetMascotStateCmd)
	assetCmd.AddCommand(assetMascotInfoCmd)
	assetCmd.AddCommand(assetMascotHooksCmd)
}

// exitCodeFromStr is a helper used by shell hooks when passing exit codes as strings.
func exitCodeFromStr(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
