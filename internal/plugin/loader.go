package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadPlugin reads and parses a plugin.json from a plugin directory
func LoadPlugin(pluginDir string) (*Plugin, error) {
	manifestPath := filepath.Join(pluginDir, "plugin.json")

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("could not read plugin manifest: %w", err)
	}

	var p Plugin
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("invalid plugin manifest: %w", err)
	}

	return &p, nil
}

// LoadSpinners reads spinners.json from a plugin directory if it exists
func LoadSpinners(pluginDir string) ([]SpinnerSet, error) {
	path := filepath.Join(pluginDir, "spinners.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not read spinners.json: %w", err)
	}

	var spinners []SpinnerSet
	if err := json.Unmarshal(data, &spinners); err != nil {
		return nil, fmt.Errorf("invalid spinners.json: %w", err)
	}

	return spinners, nil
}

// LoadBanners reads banners.json from a plugin directory if it exists
func LoadBanners(pluginDir string) ([]BannerTemplate, error) {
	path := filepath.Join(pluginDir, "banners.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not read banners.json: %w", err)
	}

	var banners []BannerTemplate
	if err := json.Unmarshal(data, &banners); err != nil {
		return nil, fmt.Errorf("invalid banners.json: %w", err)
	}

	return banners, nil
}

// LoadPrompts reads prompts.json from a plugin directory if it exists
func LoadPrompts(pluginDir string) ([]PromptTemplate, error) {
	path := filepath.Join(pluginDir, "prompts.json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not read prompts.json: %w", err)
	}

	var prompts []PromptTemplate
	if err := json.Unmarshal(data, &prompts); err != nil {
		return nil, fmt.Errorf("invalid prompts.json: %w", err)
	}

	return prompts, nil
}
