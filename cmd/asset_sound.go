package main

import (
	"fmt"
	"os"

	"github.com/abhigyanwebber/cmd-customizer/internal/assets"
	"github.com/spf13/cobra"
)

// assetSoundPlayCmd is called by shell hooks after every command to
// resolve and play the matching sound effect, if any. Mirrors
// mascot-state's design exactly (same context flags, same "silent
// no-op if nothing matches" behavior).
var assetSoundPlayCmd = &cobra.Command{
	Use:    "sound-play [asset-name]",
	Short:  "Resolve and play the matching sound effect (called by shell hooks)",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		exitCode, _ := cmd.Flags().GetInt("exit-code")
		command, _ := cmd.Flags().GetString("command")
		output, _ := cmd.Flags().GetString("output")
		idleSeconds, _ := cmd.Flags().GetInt("idle-seconds")
		gitStatus, _ := cmd.Flags().GetString("git-status")
		soundOverride, _ := cmd.Flags().GetString("sound")

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			os.Exit(0) // fail silently — this runs on every prompt, never block the shell
		}

		ctx := assets.MascotContext{
			LastExitCode: exitCode,
			LastCommand:  command,
			LastOutput:   output,
			IdleSeconds:  idleSeconds,
			GitStatus:    gitStatus,
			Env:          buildEnvSnapshot(),
		}

		if soundOverride != "" {
			a, assetDir, err := m.Get(name, assets.AssetTypeSound)
			if err != nil || a.Sound == nil {
				os.Exit(0)
			}
			effect, ok := a.Sound.Sounds[soundOverride]
			if !ok {
				os.Exit(0)
			}
			_ = assets.PlaySound(assetDir, name, soundOverride, effect, a.Sound)
			return
		}

		_, _ = m.PreviewSound(name, ctx)
	},
}

// assetSoundInfoCmd shows every sound effect in a sound theme asset:
// its trigger conditions, volume, cooldown, and playback mode.
var assetSoundInfoCmd = &cobra.Command{
	Use:   "sound-info [asset-name]",
	Short: "Show the sound effects and triggers in a sound theme asset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		m, err := assets.NewManager(getAssetsDir())
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		a, _, err := m.Get(name, assets.AssetTypeSound)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}
		if a.Sound == nil {
			fmt.Println("✗ Asset has no sound config")
			os.Exit(1)
		}

		s := a.Sound
		fmt.Printf("\n  Sound Theme: %s\n", name)
		fmt.Printf("  Enabled:       %v\n", s.Enabled)
		if s.GlobalVolume > 0 {
			fmt.Printf("  Global volume: %.0f%%\n", s.GlobalVolume*100)
		}
		if s.Player != "" {
			fmt.Printf("  Custom player: %s\n", s.Player)
		} else {
			fmt.Printf("  Player:        (auto-detected per platform)\n")
		}
		fmt.Printf("\n  Sounds (%d):\n\n", len(s.Sounds))

		for soundName, effect := range s.Sounds {
			fmt.Printf("  ┌─ %s\n", soundName)
			fmt.Printf("  │  File:      %s\n", effect.File)
			if effect.Volume > 0 {
				fmt.Printf("  │  Volume:    %.0f%%\n", effect.Volume*100)
			}
			if effect.CooldownMs > 0 {
				fmt.Printf("  │  Cooldown:  %dms\n", effect.CooldownMs)
			}
			fmt.Printf("  │  Mode:      %s\n", map[bool]string{true: "async (background)", false: "blocking"}[effect.Async])
			if len(effect.Triggers) > 0 {
				fmt.Printf("  │  Triggers:\n")
				for _, t := range effect.Triggers {
					if t.Value != "" {
						fmt.Printf("  │    • %s = %s (priority %d)\n", t.Type, t.Value, t.Priority)
					} else {
						fmt.Printf("  │    • %s (priority %d)\n", t.Type, t.Priority)
					}
				}
			} else {
				fmt.Printf("  │  Triggers:  (none — manual only)\n")
			}
			fmt.Printf("  └─\n\n")
		}
	},
}

// assetSoundHooksCmd prints the shell hook code for a sound theme asset.
var assetSoundHooksCmd = &cobra.Command{
	Use:   "sound-hooks [asset-name]",
	Short: "Print the shell hook code for a sound theme asset",
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

		hooks, err := m.SoundHooks(name, shell)
		if err != nil {
			fmt.Println("✗ Error:", err)
			os.Exit(1)
		}

		fmt.Printf("# Shell hooks for sound theme '%s' (%s)\n", name, shell)
		fmt.Printf("# Add this to your shell profile, or run:\n")
		fmt.Printf("#   cmdx asset use %s --as sound\n\n", name)
		fmt.Println(hooks)
	},
}

func init() {
	assetSoundPlayCmd.Flags().Int("exit-code", 0, "Exit code of the last command")
	assetSoundPlayCmd.Flags().String("command", "", "The last command run")
	assetSoundPlayCmd.Flags().String("output", "", "Stdout/stderr of the last command (for regex triggers)")
	assetSoundPlayCmd.Flags().Int("idle-seconds", 0, "Seconds since last command (for idle triggers)")
	assetSoundPlayCmd.Flags().String("git-status", "", "Git status: dirty, clean, untracked, ahead, behind")
	assetSoundPlayCmd.Flags().String("sound", "", "Force a specific sound effect by name (skips trigger resolution)")

	assetSoundHooksCmd.Flags().String("shell", "", "Shell type: bash, zsh, powershell (default: auto-detect)")

	assetCmd.AddCommand(assetSoundPlayCmd)
	assetCmd.AddCommand(assetSoundInfoCmd)
	assetCmd.AddCommand(assetSoundHooksCmd)
}
