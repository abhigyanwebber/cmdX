package wallpaper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ── iTerm2 ────────────────────────────────────────────────

func applyITerm2(imagePath string, opacity float64) error {
	prefsPath := filepath.Join(os.Getenv("HOME"), "Library", "Preferences", "com.googlecode.iterm2.plist")

	if _, err := os.Stat(prefsPath); os.IsNotExist(err) {
		return fmt.Errorf("iTerm2 preferences not found — is iTerm2 installed?")
	}

	// iTerm2 uses a plist — we write an AppleScript to set the background
	script := fmt.Sprintf(`
tell application "iTerm2"
    tell current session of current window
        set background image to "%s"
    end tell
end tell`, imagePath)

	scriptPath := filepath.Join(os.TempDir(), "cmdx_iterm.scpt")
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return fmt.Errorf("could not write AppleScript: %w", err)
	}

	fmt.Printf("✓ Run this to apply in iTerm2:\n  osascript %s\n", scriptPath)
	return nil
}

func removeITerm2() error {
	script := `
tell application "iTerm2"
    tell current session of current window
        set background image to ""
    end tell
end tell`

	scriptPath := filepath.Join(os.TempDir(), "cmdx_iterm_remove.scpt")
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return fmt.Errorf("could not write AppleScript: %w", err)
	}

	fmt.Printf("✓ Run this to remove in iTerm2:\n  osascript %s\n", scriptPath)
	return nil
}

// ── Kitty ─────────────────────────────────────────────────

func applyKitty(imagePath string, opacity float64) error {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "kitty", "kitty.conf")

	existing := ""
	if data, err := os.ReadFile(configPath); err == nil {
		existing = string(data)
	}

	existing = removeKittyWallpaper(existing)

	block := fmt.Sprintf(`
# cmdx wallpaper
background_image %s
background_image_layout scaled
background_opacity %.2f
`, imagePath, opacity)

	final := existing + block
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("could not create kitty config dir: %w", err)
	}

	return os.WriteFile(configPath, []byte(final), 0644)
}

func removeKitty() error {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "kitty", "kitty.conf")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	cleaned := removeKittyWallpaper(string(data))
	return os.WriteFile(configPath, []byte(cleaned), 0644)
}

func removeKittyWallpaper(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	skip := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "# cmdx wallpaper" {
			skip = true
			continue
		}
		if skip && (strings.HasPrefix(line, "background_image") ||
			strings.HasPrefix(line, "background_opacity")) {
			continue
		}
		skip = false
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
