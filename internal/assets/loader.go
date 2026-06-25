package assets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadAsset(assetDir string) (*Asset, error) {
	manifestPath := filepath.Join(assetDir, "asset.json")

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("could not read asset.json: %w", err)
	}

	var a Asset
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("invalid asset.json: %w", err)
	}

	return &a, nil
}

func ValidateAsset(a *Asset, assetDir string) error {
	if a.Name == "" {
		return fmt.Errorf("name is required")
	}
	if a.Type == "" {
		return fmt.Errorf("type is required")
	}
	if a.Version == "" {
		return fmt.Errorf("version is required")
	}

	switch a.Type {
	case AssetTypeSpinner:
		if a.Spinner == nil {
			return fmt.Errorf("spinner config is required for spinner assets")
		}
		if len(a.Spinner.Frames) == 0 {
			return fmt.Errorf("spinner must have at least one frame")
		}
		if a.Spinner.IntervalMs <= 0 {
			return fmt.Errorf("interval_ms must be greater than 0")
		}
		for _, frame := range a.Spinner.Frames {
			framePath := filepath.Join(assetDir, frame)
			if _, err := os.Stat(framePath); os.IsNotExist(err) {
				return fmt.Errorf("frame file not found: %s", frame)
			}
		}

	case AssetTypeIcon:
		if a.Icon == nil {
			return fmt.Errorf("icon config is required for icon assets")
		}
		if len(a.Icon.Files) == 0 {
			return fmt.Errorf("icon asset must define at least one icon file")
		}

	case AssetTypeBanner:
		if a.Banner == nil {
			return fmt.Errorf("banner config is required for banner assets")
		}
		bannerPath := filepath.Join(assetDir, a.Banner.File)
		if _, err := os.Stat(bannerPath); os.IsNotExist(err) {
			return fmt.Errorf("banner file not found: %s", a.Banner.File)
		}

	case AssetTypeDivider:
		if a.Divider == nil {
			return fmt.Errorf("divider config is required for divider assets")
		}
		dividerPath := filepath.Join(assetDir, a.Divider.File)
		if _, err := os.Stat(dividerPath); os.IsNotExist(err) {
			return fmt.Errorf("divider file not found: %s", a.Divider.File)
		}
	}

	return nil
}

func ResolveFramePaths(a *Asset, assetDir string) ([]string, error) {
	if a.Type != AssetTypeSpinner || a.Spinner == nil {
		return nil, fmt.Errorf("asset is not a spinner")
	}

	paths := make([]string, len(a.Spinner.Frames))
	for i, frame := range a.Spinner.Frames {
		paths[i] = filepath.Join(assetDir, frame)
	}
	return paths, nil
}
