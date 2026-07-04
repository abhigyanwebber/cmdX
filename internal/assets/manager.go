package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager handles all asset operations
type Manager struct {
	AssetsDir string
}

// NewManager creates an asset manager rooted at assetsDir, creating the
// directory if it doesn't exist.
func NewManager(assetsDir string) (*Manager, error) {
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(assetsDir, 0755); err != nil {
			return nil, fmt.Errorf("could not create assets directory: %w", err)
		}
	}
	return &Manager{AssetsDir: assetsDir}, nil
}

// List returns all assets of a given type installed under AssetsDir.
// Returns nil (not an error) if no assets of that type exist yet.
func (m *Manager) List(assetType AssetType) ([]*Asset, error) {
	subDir := filepath.Join(m.AssetsDir, string(assetType)+"s")
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(subDir)
	if err != nil {
		return nil, fmt.Errorf("could not read assets directory: %w", err)
	}

	var result []*Asset
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
		result = append(result, a)
	}
	return result, nil
}

// Get returns a specific asset by name and type, along with its
// directory on disk. Returns an error if the asset doesn't exist or
// fails to load.
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

// Import copies a source asset folder into the assets directory,
// validating the manifest and all referenced files first. Returns an
// error if the asset already exists (use Remove first) or fails
// validation.
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

// Remove deletes a named asset of the given type from the assets
// directory.
func (m *Manager) Remove(name string, assetType AssetType) error {
	assetDir := filepath.Join(m.AssetsDir, string(assetType)+"s", name)
	if _, err := os.Stat(assetDir); os.IsNotExist(err) {
		return fmt.Errorf("asset '%s' not found", name)
	}
	return os.RemoveAll(assetDir)
}

// PreviewSpinner renders and plays a spinner asset for the given
// duration, cycling through frames at the asset's configured interval.
// Supports bounce and reverse playback modes from the manifest.
func (m *Manager) PreviewSpinner(name string, duration time.Duration) error {
	return m.PreviewSpinnerWithOverrides(name, duration, RenderOverrides{})
}

// PreviewSpinnerWithOverrides is like PreviewSpinner but applies runtime
// render overrides on top of the manifest's render config. Zero-value
// override fields fall back to the manifest defaults.
func (m *Manager) PreviewSpinnerWithOverrides(name string, duration time.Duration, overrides RenderOverrides) error {
	a, assetDir, err := m.Get(name, AssetTypeSpinner)
	if err != nil {
		return err
	}

	framePaths, err := ResolveFramePaths(a, assetDir)
	if err != nil {
		return err
	}

	opts := ApplyOverrides(optionsFromConfig(a.Render), overrides)

	frames, err := RenderFrames(framePaths, opts)
	if err != nil {
		return err
	}

	if a.Spinner.Bounce {
		reversed := make([]string, len(frames))
		for i, f := range frames {
			reversed[len(frames)-1-i] = f
		}
		frames = append(frames, reversed...)
	}

	if a.Spinner.Reverse {
		for i, j := 0, len(frames)-1; i < j; i, j = i+1, j-1 {
			frames[i], frames[j] = frames[j], frames[i]
		}
	}

	interval := time.Duration(a.Spinner.IntervalMs) * time.Millisecond
	deadline := time.Now().Add(duration)
	idx := 0

	for time.Now().Before(deadline) {
		fmt.Print("\033[H\033[2J")
		fmt.Println(frames[idx%len(frames)])
		idx++
		time.Sleep(interval)
	}

	return nil
}

// RenderSpinnerFrames returns pre-rendered frames ready for animation,
// along with the configured frame interval in milliseconds.
func (m *Manager) RenderSpinnerFrames(name string) ([]string, int, error) {
	return m.RenderSpinnerFramesWithOverrides(name, RenderOverrides{})
}

// RenderSpinnerFramesWithOverrides is like RenderSpinnerFrames but applies
// runtime render overrides on top of the manifest defaults.
func (m *Manager) RenderSpinnerFramesWithOverrides(name string, overrides RenderOverrides) ([]string, int, error) {
	a, assetDir, err := m.Get(name, AssetTypeSpinner)
	if err != nil {
		return nil, 0, err
	}

	framePaths, err := ResolveFramePaths(a, assetDir)
	if err != nil {
		return nil, 0, err
	}

	opts := ApplyOverrides(optionsFromConfig(a.Render), overrides)
	frames, err := RenderFrames(framePaths, opts)
	if err != nil {
		return nil, 0, err
	}

	return frames, a.Spinner.IntervalMs, nil
}

// PreviewBanner renders and prints a banner asset to the terminal.
func (m *Manager) PreviewBanner(name string) error {
	return m.PreviewBannerWithOverrides(name, RenderOverrides{})
}

// PreviewBannerWithOverrides is like PreviewBanner but applies runtime
// render overrides.
func (m *Manager) PreviewBannerWithOverrides(name string, overrides RenderOverrides) error {
	a, assetDir, err := m.Get(name, AssetTypeBanner)
	if err != nil {
		return err
	}

	if a.Banner == nil {
		return fmt.Errorf("asset '%s' has no banner config", name)
	}

	bannerPath := filepath.Join(assetDir, a.Banner.File)
	opts := ApplyOverrides(optionsFromConfig(a.Render), overrides)

	rendered, err := Render(bannerPath, opts)
	if err != nil {
		return fmt.Errorf("could not render banner: %w", err)
	}

	fmt.Println(rendered)
	return nil
}

// PreviewDivider renders and prints a divider asset to the terminal.
func (m *Manager) PreviewDivider(name string) error {
	return m.PreviewDividerWithOverrides(name, RenderOverrides{})
}

// PreviewDividerWithOverrides is like PreviewDivider but applies runtime
// render overrides.
func (m *Manager) PreviewDividerWithOverrides(name string, overrides RenderOverrides) error {
	a, assetDir, err := m.Get(name, AssetTypeDivider)
	if err != nil {
		return err
	}

	if a.Divider == nil {
		return fmt.Errorf("asset '%s' has no divider config", name)
	}

	dividerPath := filepath.Join(assetDir, a.Divider.File)
	opts := ApplyOverrides(optionsFromConfig(a.Render), overrides)

	rendered, err := Render(dividerPath, opts)
	if err != nil {
		return fmt.Errorf("could not render divider: %w", err)
	}

	fmt.Println(rendered)
	return nil
}

// PreviewFloater renders and displays a floater asset at its configured
// corner position.
func (m *Manager) PreviewFloater(name string) error {
	return m.PreviewFloaterWithOverrides(name, "", RenderOverrides{})
}

// PreviewFloaterWithOverrides is like PreviewFloater but allows
// overriding both the corner position and any render options at runtime.
// An empty positionOverride means "use the manifest's configured position".
func (m *Manager) PreviewFloaterWithOverrides(name string, positionOverride FloaterPosition, overrides RenderOverrides) error {
	a, assetDir, err := m.Get(name, AssetTypeFloater)
	if err != nil {
		return err
	}

	if a.Floater == nil {
		return fmt.Errorf("asset '%s' has no floater config", name)
	}

	pos := a.Floater.Position
	if positionOverride != "" {
		pos = positionOverride
	}

	floaterPath := filepath.Join(assetDir, a.Floater.File)
	opts := ApplyOverrides(optionsFromConfig(a.Render), overrides)

	rendered, err := Render(floaterPath, opts)
	if err != nil {
		return fmt.Errorf("could not render floater: %w", err)
	}

	fmt.Printf("  [%s]\n%s\n", pos, rendered)
	return nil
}

// RenderFloaterFrames returns pre-rendered animation frames for an
// animated floater, along with frame interval and position.
func (m *Manager) RenderFloaterFrames(name string) ([]string, int, FloaterPosition, error) {
	return m.RenderFloaterFramesWithOverrides(name, "", RenderOverrides{})
}

// RenderFloaterFramesWithOverrides is like RenderFloaterFrames but
// applies runtime render overrides and an optional position override.
func (m *Manager) RenderFloaterFramesWithOverrides(name string, positionOverride FloaterPosition, overrides RenderOverrides) ([]string, int, FloaterPosition, error) {
	a, assetDir, err := m.Get(name, AssetTypeFloater)
	if err != nil {
		return nil, 0, "", err
	}

	if a.Floater == nil {
		return nil, 0, "", fmt.Errorf("asset '%s' has no floater config", name)
	}
	if len(a.Floater.AnimateFrames) == 0 {
		return nil, 0, "", fmt.Errorf("floater '%s' has no animate_frames configured", name)
	}

	pos := a.Floater.Position
	if positionOverride != "" {
		pos = positionOverride
	}

	framePaths := make([]string, len(a.Floater.AnimateFrames))
	for i, frame := range a.Floater.AnimateFrames {
		framePaths[i] = filepath.Join(assetDir, frame)
	}

	opts := ApplyOverrides(optionsFromConfig(a.Render), overrides)
	frames, err := RenderFrames(framePaths, opts)
	if err != nil {
		return nil, 0, "", err
	}

	return frames, a.Floater.IntervalMs, pos, nil
}

// PreviewIcon renders and displays a single icon from an icon asset.
func (m *Manager) PreviewIcon(name string, iconKey string) error {
	return m.PreviewIconWithOverrides(name, iconKey, RenderOverrides{})
}

// PreviewIconWithOverrides is like PreviewIcon but applies runtime
// render overrides.
func (m *Manager) PreviewIconWithOverrides(name string, iconKey string, overrides RenderOverrides) error {
	a, assetDir, err := m.Get(name, AssetTypeIcon)
	if err != nil {
		return err
	}

	if a.Icon == nil {
		return fmt.Errorf("asset '%s' has no icon config", name)
	}

	file, ok := a.Icon.Files[iconKey]
	if !ok {
		return fmt.Errorf("icon key '%s' not found in asset '%s'", iconKey, name)
	}

	iconPath := filepath.Join(assetDir, file)
	opts := ApplyOverrides(optionsFromConfig(a.Render), overrides)

	rendered, err := Render(iconPath, opts)
	if err != nil {
		return fmt.Errorf("could not render icon: %w", err)
	}

	fmt.Printf("  %s: %s", iconKey, rendered)
	return nil
}

// RenderRaw renders any image file directly through the chafa pipeline
// with the given options, without requiring an asset manifest. This is
// the low-level API for developers who want to test images or build
// their own tooling on top of cmdX's render pipeline.
func RenderRaw(imagePath string, opts ChafaOptions) (string, error) {
	return Render(imagePath, opts)
}

// copyDir recursively copies a directory tree.
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

// PreviewMascot renders and displays the mascot in its current resolved
// state based on the provided context.
func (m *Manager) PreviewMascot(name string, ctx MascotContext, overrides RenderOverrides) error {
	a, assetDir, err := m.Get(name, AssetTypeMascot)
	if err != nil {
		return err
	}
	if a.Mascot == nil {
		return fmt.Errorf("asset '%s' has no mascot config", name)
	}

	state := ResolveState(a.Mascot, ctx)
	frames, transition, intervalMs, err := RenderMascotState(a, assetDir, state, overrides)
	if err != nil {
		return fmt.Errorf("could not render mascot state %q: %w", state, err)
	}

	fmt.Printf("  [mascot: %s | state: %s]\n\n", name, state)

	// play transition frames once if present
	for _, f := range transition {
		fmt.Println(f)
	}

	// display first frame (for preview — not a live animation loop)
	if len(frames) > 0 {
		fmt.Println(frames[0])
	}

	_ = intervalMs
	return nil
}

// ResolveMascotState resolves and returns the current mascot state name
// without rendering anything. Useful for shell hooks and state display.
func (m *Manager) ResolveMascotState(name string, ctx MascotContext) (MascotState, error) {
	a, _, err := m.Get(name, AssetTypeMascot)
	if err != nil {
		return "", err
	}
	if a.Mascot == nil {
		return "", fmt.Errorf("asset '%s' has no mascot config", name)
	}
	return ResolveState(a.Mascot, ctx), nil
}

// MascotHooks returns the shell hook code for a mascot asset, ready to
// be injected into the user's shell profile.
func (m *Manager) MascotHooks(name string, shell string) (string, error) {
	a, _, err := m.Get(name, AssetTypeMascot)
	if err != nil {
		return "", err
	}
	if a.Mascot == nil {
		return "", fmt.Errorf("asset '%s' has no mascot config", name)
	}

	hookVar := a.Mascot.HookVar
	if hookVar == "" {
		hookVar = "CMDX_MASCOT_TRIGGER"
	}

	hooks := MascotShellHooks(name, hookVar, shell)
	if hooks == "" {
		return "", fmt.Errorf("unsupported shell %q for mascot hooks", shell)
	}
	return hooks, nil
}

// PreviewStatusBar prints the shell code for a status bar asset — since
// a status bar is a shell function, not a rendered image, preview shows
// the generated code and a visual mockup of the segments.
func (m *Manager) PreviewStatusBar(name string, shell string, colors map[string]string) error {
	a, _, err := m.Get(name, AssetTypeStatusBar)
	if err != nil {
		return err
	}
	if a.StatusBar == nil {
		return fmt.Errorf("asset '%s' has no status_bar config", name)
	}

	sb := a.StatusBar

	// print visual mockup
	fmt.Printf("\n  Status Bar: %s\n", name)
	fmt.Printf("  Position:   %s\n", sb.Position)
	fmt.Printf("  Separator:  %s\n", sb.SeparatorStyle)
	fmt.Printf("  Segments:   %d total\n\n", len(sb.Segments))

	left, center, right := segmentsByZone(sb)
	sepL, _ := separatorGlyphs(sb.SeparatorStyle, sb.CustomSeparatorLeft, sb.CustomSeparatorRight)

	// visual mockup line
	var leftParts, centerParts, rightParts []string
	for _, seg := range left {
		leftParts = append(leftParts, segmentMockup(seg))
	}
	for _, seg := range center {
		centerParts = append(centerParts, segmentMockup(seg))
	}
	for _, seg := range right {
		rightParts = append(rightParts, segmentMockup(seg))
	}

	leftStr := strings.Join(leftParts, sepL)
	centerStr := strings.Join(centerParts, sepL)
	rightStr := strings.Join(rightParts, sepL)

	fmt.Printf("  ┌─ Visual Mockup ─────────────────────────────────────────────┐\n")
	fmt.Printf("  │ %-25s %-15s %20s │\n", leftStr, centerStr, rightStr)
	fmt.Printf("  └─────────────────────────────────────────────────────────────┘\n\n")

	// print generated shell code
	if shell != "" {
		fmt.Printf("  Generated %s code:\n\n", shell)
		code := StatusBarShellCode(sb, colors, shell)
		if code == "" {
			fmt.Printf("  (unsupported shell: %s)\n", shell)
		} else {
			fmt.Println(code)
		}
	} else {
		fmt.Printf("  Run with --shell bash|zsh|powershell to see generated hook code.\n")
	}

	return nil
}

// StatusBarHooks returns the shell hook code for a status bar asset,
// ready to be injected into the user's shell profile.
func (m *Manager) StatusBarHooks(name string, shell string, colors map[string]string) (string, error) {
	a, _, err := m.Get(name, AssetTypeStatusBar)
	if err != nil {
		return "", err
	}
	if a.StatusBar == nil {
		return "", fmt.Errorf("asset '%s' has no status_bar config", name)
	}

	code := StatusBarShellCode(a.StatusBar, colors, shell)
	if code == "" {
		return "", fmt.Errorf("unsupported shell %q for status bar hooks", shell)
	}
	return code, nil
}

// segmentMockup returns a human-readable placeholder for a segment in the preview.
func segmentMockup(seg SegmentConfig) string {
	label := seg.Label
	switch seg.Type {
	case SegmentGit:
		return label + "main*"
	case SegmentDirectory:
		return label + "~/projects"
	case SegmentTime:
		return label + "14:32"
	case SegmentDate:
		return label + "2026-07-04"
	case SegmentExitCode:
		return label + "✗127"
	case SegmentDuration:
		return label + "⏱2.3s"
	case SegmentUser:
		return label + "$USER"
	case SegmentHost:
		return label + "hostname"
	case SegmentVirtualEnv:
		return label + "(venv)"
	case SegmentGoVersion:
		return label + "go1.22"
	case SegmentNodeVersion:
		return label + "node20"
	case SegmentKubernetes:
		return label + "⎈prod"
	case SegmentBattery:
		return label + "🔋85%"
	case SegmentText:
		return seg.Format
	case SegmentEnvVar:
		return label + "$" + seg.EnvVar
	case SegmentCommand:
		return label + "(cmd)"
	default:
		return label + string(seg.Type)
	}
}
