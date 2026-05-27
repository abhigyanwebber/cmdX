package wallpaper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// findWindowsTerminalSettings locates the Windows Terminal settings.json
func findWindowsTerminalSettings() (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	appData := os.Getenv("APPDATA")

	candidates := []string{
		filepath.Join(localAppData, "Packages", "Microsoft.WindowsTerminal_8wekyb3d8bbwe", "LocalState", "settings.json"),
		filepath.Join(localAppData, "Packages", "Microsoft.WindowsTerminalPreview_8wekyb3d8bbwe", "LocalState", "settings.json"),
		filepath.Join(appData, "Microsoft", "Windows Terminal", "settings.json"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("Windows Terminal settings.json not found — is Windows Terminal installed?")
}

func applyWindowsTerminal(imagePath string, opacity float64, stretch string, alignment string) error {
	settingsPath, err := findWindowsTerminalSettings()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("could not read settings.json: %w", err)
	}

	// backup original
	backupPath := settingsPath + ".cmdx.bak"
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("could not create backup: %w", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("could not parse settings.json: %w", err)
	}

	// normalize path separators for JSON
	normalizedPath := strings.ReplaceAll(imagePath, "\\", "/")

	// apply to default profile
	profiles, ok := settings["profiles"].(map[string]interface{})
	if !ok {
		profiles = make(map[string]interface{})
		settings["profiles"] = profiles
	}

	defaults, ok := profiles["defaults"].(map[string]interface{})
	if !ok {
		defaults = make(map[string]interface{})
		profiles["defaults"] = defaults
	}

	defaults["backgroundImage"] = normalizedPath
	defaults["backgroundImageOpacity"] = opacity
	defaults["backgroundImageStretchMode"] = stretch
	defaults["backgroundImageAlignment"] = alignment

	updated, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return fmt.Errorf("could not serialize settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, updated, 0644); err != nil {
		return fmt.Errorf("could not write settings.json: %w", err)
	}

	return nil
}

func removeWindowsTerminal() error {
	settingsPath, err := findWindowsTerminalSettings()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("could not read settings.json: %w", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("could not parse settings.json: %w", err)
	}

	profiles, ok := settings["profiles"].(map[string]interface{})
	if !ok {
		return nil
	}

	defaults, ok := profiles["defaults"].(map[string]interface{})
	if !ok {
		return nil
	}

	delete(defaults, "backgroundImage")
	delete(defaults, "backgroundImageOpacity")
	delete(defaults, "backgroundImageStretchMode")
	delete(defaults, "backgroundImageAlignment")

	updated, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return fmt.Errorf("could not serialize settings: %w", err)
	}

	return os.WriteFile(settingsPath, updated, 0644)
}
