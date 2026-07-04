// Package assets handles loading, validating, and resolving file paths
// for PNG-based terminal assets (spinners, icons, banners, dividers) used
// by cmdX themes.
package assets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadAsset reads and parses the asset.json manifest from assetDir.
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

// ValidateAsset checks an Asset manifest for required fields and, for
// asset types that reference files on disk (spinner frames, banners,
// dividers), verifies those files actually exist under assetDir.
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

	case AssetTypeFloater:
		if a.Floater == nil {
			return fmt.Errorf("floater config is required for floater assets")
		}
		if a.Floater.File == "" {
			return fmt.Errorf("floater must specify a file")
		}
		floaterPath := filepath.Join(assetDir, a.Floater.File)
		if _, err := os.Stat(floaterPath); os.IsNotExist(err) {
			return fmt.Errorf("floater file not found: %s", a.Floater.File)
		}
		if !IsValidFloaterPosition(a.Floater.Position) {
			return fmt.Errorf("floater position must be one of: top-left, top-right, bottom-left, bottom-right (got %q)", a.Floater.Position)
		}
		if len(a.Floater.AnimateFrames) > 0 {
			if a.Floater.IntervalMs <= 0 {
				return fmt.Errorf("floater interval_ms must be greater than 0 when animate_frames is set")
			}
			for _, frame := range a.Floater.AnimateFrames {
				framePath := filepath.Join(assetDir, frame)
				if _, err := os.Stat(framePath); os.IsNotExist(err) {
					return fmt.Errorf("floater animation frame not found: %s", frame)
				}
			}
		}

	case AssetTypeMascot:
		if a.Mascot == nil {
			return fmt.Errorf("mascot config is required for mascot assets")
		}
		if err := ValidateMascot(a.Mascot, assetDir); err != nil {
			return fmt.Errorf("mascot validation failed: %w", err)
		}

	case AssetTypeStatusBar:
		if a.StatusBar == nil {
			return fmt.Errorf("status_bar config is required for status-bar assets")
		}
		if err := ValidateStatusBar(a.StatusBar); err != nil {
			return fmt.Errorf("status_bar validation failed: %w", err)
		}
	}

	return nil
}

// ResolveFramePaths returns the absolute paths of every spinner frame
// file for a, joined against assetDir. Returns an error if a is not a
// spinner asset or has no spinner configuration.
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
