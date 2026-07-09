package fonts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	// defaultMaxDownloadBytes is the default cap on font archive
	// downloads. Most single-family Nerd Font zips are well under
	// this; full "all variants" mega-archives can exceed it, in which
	// case InstallOptions.MaxDownloadBytes must be set explicitly —
	// this is an intentional developer-freedom escape hatch rather
	// than a hard limit.
	defaultMaxDownloadBytes = 150 << 20 // 150 MiB

	downloadTimeout = 5 * time.Minute
)

var httpClient = &http.Client{Timeout: downloadTimeout}

// FontsDir returns the platform's per-user font installation directory,
// creating it if it doesn't exist.
//
//   - Windows: %LOCALAPPDATA%\Microsoft\Windows\Fonts
//   - macOS:   ~/Library/Fonts
//   - Linux:   ~/.local/share/fonts
func FontsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	var dir string
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = filepath.Join(home, "AppData", "Local")
		}
		dir = filepath.Join(localAppData, "Microsoft", "Windows", "Fonts")
	case "darwin":
		dir = filepath.Join(home, "Library", "Fonts")
	default: // linux and other unix-likes
		dir = filepath.Join(home, ".local", "share", "fonts")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("could not create fonts directory: %w", err)
	}
	return dir, nil
}

// stateFilePath returns the path to cmdX's font install tracking file.
func stateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".cmdx")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "fonts.json"), nil
}

// loadState reads the installed-fonts tracking file, returning an empty
// map if it doesn't exist yet.
func loadState() (map[string]InstalledFont, error) {
	path, err := stateFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]InstalledFont{}, nil
		}
		return nil, err
	}

	var state map[string]InstalledFont
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("could not parse font state file: %w", err)
	}
	return state, nil
}

// saveState writes the installed-fonts tracking file.
func saveState(state map[string]InstalledFont) error {
	path, err := stateFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ListInstalledFonts returns every font cmdX has installed, tracked in
// its local state file.
func ListInstalledFonts() ([]InstalledFont, error) {
	state, err := loadState()
	if err != nil {
		return nil, err
	}
	var result []InstalledFont
	for _, f := range state {
		result = append(result, f)
	}
	return result, nil
}

// IsInstalled reports whether cmdX has a tracked installation record
// for the given font name.
func IsInstalled(name string) bool {
	state, err := loadState()
	if err != nil {
		return false
	}
	_, ok := state[name]
	return ok
}

// downloadToTemp streams url to a temporary file, enforcing a maximum
// size. Returns the temp file path — caller is responsible for removing
// it once done.
func downloadToTemp(url string, maxBytes int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("could not build request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d (check the font name or URL)", resp.StatusCode)
	}

	if resp.ContentLength > 0 && resp.ContentLength > maxBytes {
		return "", fmt.Errorf("font archive is %d bytes, exceeding the %d byte limit — use InstallOptions.MaxDownloadBytes to allow larger downloads", resp.ContentLength, maxBytes)
	}

	tmp, err := os.CreateTemp("", "cmdx-font-*.zip")
	if err != nil {
		return "", fmt.Errorf("could not create temp file: %w", err)
	}
	defer tmp.Close()

	limited := io.LimitReader(resp.Body, maxBytes+1)
	written, err := io.Copy(tmp, limited)
	if err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("could not save download: %w", err)
	}
	if written > maxBytes {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("download exceeded the %d byte limit", maxBytes)
	}

	return tmp.Name(), nil
}

// InstallFont downloads and installs a font by curated catalog name
// (e.g. "firacode") or, if opts.URL is set, from any direct archive URL.
// Returns the list of installed file paths.
func InstallFont(name string, opts InstallOptions) ([]InstalledFont, error) {
	var downloadURL string
	var displayName string
	var source string

	if opts.URL != "" {
		downloadURL = opts.URL
		displayName = name
		source = opts.URL
	} else {
		entry := Find(name)
		if entry == nil {
			return nil, fmt.Errorf("font %q not found in catalog. Run 'cmdx font list' to see available fonts, or use --url for a custom font", name)
		}
		downloadURL = entry.DownloadURL()
		displayName = entry.DisplayName
		source = name
	}

	if IsInstalled(name) && !opts.Force {
		return nil, fmt.Errorf("font %q is already installed. Use --force to reinstall", name)
	}

	maxBytes := int64(defaultMaxDownloadBytes)
	if opts.MaxDownloadBytes > 0 {
		maxBytes = opts.MaxDownloadBytes
	}

	if opts.DryRun {
		fmt.Printf("  [dry-run] Would download: %s\n", downloadURL)
		fmt.Printf("  [dry-run] Would install to font directory (variant filter: %s)\n", variantLabel(opts))
		return nil, nil
	}

	zipPath, err := downloadToTemp(downloadURL, maxBytes)
	if err != nil {
		return nil, err
	}
	defer os.Remove(zipPath)

	fontsDir, err := FontsDir()
	if err != nil {
		return nil, err
	}

	extracted, err := extractFontFiles(zipPath, fontsDir, opts.Variant, opts.AllVariants || opts.Variant == "")
	if err != nil {
		return nil, err
	}

	// register with the OS (no-op on macOS/Linux)
	for _, path := range extracted {
		if err := registerFont(filepath.Base(path), path); err != nil {
			// non-fatal — font files are already installed, registry
			// registration is a nice-to-have for full app compatibility
			fmt.Printf("  ! Warning: could not register %s in Windows font registry: %v\n", filepath.Base(path), err)
		}
	}

	refreshFontCache()

	record := InstalledFont{
		Name:        name,
		Source:      source,
		Files:       extracted,
		InstalledAt: time.Now(),
	}

	state, err := loadState()
	if err != nil {
		return nil, err
	}
	state[name] = record
	if err := saveState(state); err != nil {
		return nil, fmt.Errorf("font installed but could not save state (removal tracking may not work): %w", err)
	}

	_ = displayName
	return []InstalledFont{record}, nil
}

// RemoveFont deletes a previously cmdX-installed font's files and
// tracking record.
func RemoveFont(name string) error {
	state, err := loadState()
	if err != nil {
		return err
	}

	record, ok := state[name]
	if !ok {
		return fmt.Errorf("font %q is not tracked as installed by cmdX (it may have been installed manually)", name)
	}

	for _, path := range record.Files {
		if err := unregisterFont(filepath.Base(path)); err != nil {
			fmt.Printf("  ! Warning: could not unregister %s: %v\n", filepath.Base(path), err)
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			fmt.Printf("  ! Warning: could not remove %s: %v\n", path, err)
		}
	}

	delete(state, name)
	if err := saveState(state); err != nil {
		return fmt.Errorf("removed font files but could not update state: %w", err)
	}

	refreshFontCache()
	return nil
}

// refreshFontCache runs the platform's font cache refresh command, if
// applicable. Best-effort — failures are silent since font cache
// refresh isn't strictly required for the font to become available
// (most terminal apps re-scan on next launch regardless).
func refreshFontCache() {
	if runtime.GOOS != "linux" {
		return
	}
	if _, err := exec.LookPath("fc-cache"); err != nil {
		return
	}
	dir, err := FontsDir()
	if err != nil {
		return
	}
	_ = exec.Command("fc-cache", "-f", dir).Run()
}

// variantLabel formats the variant selection for dry-run output.
func variantLabel(opts InstallOptions) string {
	if opts.AllVariants || opts.Variant == "" {
		return "all variants"
	}
	return opts.Variant
}

// RegistrySupported reports whether font registry registration runs on
// this platform (Windows only — macOS/Linux discover fonts by directory
// scan alone).
func RegistrySupported() bool {
	return registrySupported()
}
