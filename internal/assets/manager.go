package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles all asset operations
type Manager struct {
	AssetsDir string
}

// NewManager creates an asset manager
func NewManager(assetsDir string) (*Manager, error) {
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(assetsDir, 0755); err != nil {
			return nil, fmt.Errorf("could not create assets directory: %w", err)
		}
	}
	return &Manager{AssetsDir: assetsDir}, nil
}

// List returns all assets of a given type
func (m *Manager) List(assetType AssetType) ([]*Asset, error) {
	subDir := filepath.Join(m.AssetsDir, string(assetType)+"s")
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(subDir)
	if err != nil {
		return nil, fmt.Errorf("could not read assets directory: %w", err)
	}

	var assets []*Asset
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		assetDir := filepath.Join(subDir, entry.Name())
		a, err := LoadAsset(assetDir)
		if err != nil {
			fmt.Printf("  ! Skipping '%s': %s\n", entry.Name(), err)
			continue
		}
		assets = append(assets, a)
	}
	return assets, nil
}

// Get returns a specific asset by name and type
func (m *Manager) Get(name string, assetType AssetType) (*Asset, string, error) {
	assetDir := filepath.Join(m.AssetsDir, string(assetType)+"s", name)
	if _, err := os.Stat(assetDir); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("asset '%s' not found in %ss", name, assetType)
	}

	a, err := LoadAsset(assetDir)
	if err != nil {
		return nil, "", err
	}
	return a, assetDir, nil
}

// Import copies an asset folder into the assets directory
func (m *Manager) Import(sourcePath string) (*Asset, error) {
	a, err := LoadAsset(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("invalid asset: %w", err)
	}

	if err := ValidateAsset(a, sourcePath); err != nil {
		return nil, fmt.Errorf("asset validation failed: %w", err)
	}

	destDir := filepath.Join(m.AssetsDir, string(a.Type)+"s", a.Name)

	if _, err := os.Stat(destDir); err == nil {
		return nil, fmt.Errorf("asset '%s' already exists. Remove it first with 'cmdx asset remove %s'", a.Name, a.Name)
	}

	if err := copyDir(sourcePath, destDir); err != nil {
		return nil, fmt.Errorf("could not import asset: %w", err)
	}

	return a, nil
}

// PreviewSpinner renders and plays a spinner asset in the terminal
func (m *Manager) PreviewSpinner(name string, duration time.Duration) error {
	a, assetDir, err := m.Get(name, AssetTypeSpinner)
	if err != nil {
		return err
	}

	framePaths, err := ResolveFramePaths(a, assetDir)
	if err != nil {
		return err
	}

	opts := ChafaOptions{
		Mode:      a.Render.Mode,
		ColorMode: a.Render.ColorMode,
		SymbolSet: a.Render.SymbolSet,
		Width:     a.Render.Width,
		Height:    a.Render.Height,
		Dither:    a.Render.Dither,
		Stretch:   a.Render.Stretch,
		Threshold: a.Render.Threshold,
	}

	frames, err := RenderFrames(framePaths, opts)
	if err != nil {
		return err
	}

	// apply bounce mode
	if a.Spinner.Bounce {
		reversed := make([]string, len(frames))
		for i, f := range frames {
			reversed[len(frames)-1-i] = f
		}
		frames = append(frames, reversed...)
	}

	// apply reverse mode
	if a.Spinner.Reverse {
		for i, j := 0, len(frames)-1; i < j; i, j = i+1, j-1 {
			frames[i], frames[j] = frames[j], frames[i]
		}
	}

	interval := time.Duration(a.Spinner.IntervalMs) * time.Millisecond
	deadline := time.Now().Add(duration)
	idx := 0

	for time.Now().Before(deadline) {
		frame := frames[idx%len(frames)]
		fmt.Print("\033[H\033[2J") // clear screen
		fmt.Println(frame)
		idx++
		time.Sleep(interval)
	}

	return nil
}

// RenderSpinnerFrames returns pre-rendered frames ready for animation
func (m *Manager) RenderSpinnerFrames(name string) ([]string, int, error) {
	a, assetDir, err := m.Get(name, AssetTypeSpinner)
	if err != nil {
		return nil, 0, err
	}

	framePaths, err := ResolveFramePaths(a, assetDir)
	if err != nil {
		return nil, 0, err
	}

	opts := ChafaOptions{
		Mode:      a.Render.Mode,
		ColorMode: a.Render.ColorMode,
		SymbolSet: a.Render.SymbolSet,
		Width:     a.Render.Width,
		Height:    a.Render.Height,
		Dither:    a.Render.Dither,
		Stretch:   a.Render.Stretch,
		Threshold: a.Render.Threshold,
	}

	frames, err := RenderFrames(framePaths, opts)
	if err != nil {
		return nil, 0, err
	}

	return frames, a.Spinner.IntervalMs, nil
}

// Remove deletes an asset from the assets directory
func (m *Manager) Remove(name string, assetType AssetType) error {
	assetDir := filepath.Join(m.AssetsDir, string(assetType)+"s", name)
	if _, err := os.Stat(assetDir); os.IsNotExist(err) {
		return fmt.Errorf("asset '%s' not found", name)
	}
	return os.RemoveAll(assetDir)
}

// copyDir recursively copies a directory
func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, info.Mode())
	})
}
