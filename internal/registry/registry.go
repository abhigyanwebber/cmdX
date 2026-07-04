// Package registry handles communication with the community theme
// registry hosted at github.com/abhigyanwebber/cmdX-themes — fetching
// the theme index, downloading individual themes, and searching.
package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	// RegistryBase is the root raw-content URL for the community themes repo.
	RegistryBase = "https://raw.githubusercontent.com/abhigyanwebber/cmdX-themes/main"
	// RegistryIndex is the URL of the JSON index listing all available themes.
	RegistryIndex = RegistryBase + "/index.json"

	// requestTimeout bounds every individual HTTP request made to the
	// registry, independent of the caller's own context (if any).
	requestTimeout = 15 * time.Second

	// maxDownloadBytes caps the size of any single response body read
	// into memory (theme JSON or the registry index itself). This
	// protects against a compromised or malicious registry serving an
	// oversized response that could exhaust client memory — theme files
	// are small JSON documents and have no legitimate reason to
	// approach this size.
	maxDownloadBytes = 1 << 20 // 1 MiB
)

// themeNameValidator restricts registry theme names to a safe character
// set: lowercase/uppercase letters, digits, hyphens, and underscores.
// This is the same allowlist enforced by the interactive theme creator
// for locally-created themes (see cmd/create.go), applied here for
// theme names obtained from network input (the registry index, or a
// user-supplied "cmdx registry fetch <name>" argument).
//
// Without this, a name containing path traversal sequences (e.g.
// "../../../.bashrc") could escape the themes directory in FetchTheme's
// filepath.Join, allowing a malicious registry entry — or a typo'd CLI
// argument forwarded into a script — to overwrite arbitrary files on
// the user's machine.
var themeNameValidator = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidThemeName reports whether name is safe to use as a filename and
// as a path segment in a registry URL.
func ValidThemeName(name string) bool {
	return name != "" && themeNameValidator.MatchString(name)
}

// ThemeEntry represents a single theme's metadata as listed in the
// registry index.
type ThemeEntry struct {
	Name        string   `json:"name"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
	URL         string   `json:"url"`
}

// Index is the full registry index — every theme currently published to
// the community registry.
type Index struct {
	UpdatedAt string       `json:"updated_at"`
	Themes    []ThemeEntry `json:"themes"`
}

// httpClient is shared across all registry requests. Per-request timeouts
// are additionally enforced via context in each call site below.
var httpClient = &http.Client{
	Timeout: requestTimeout,
}

// doGet issues a GET request against url with a bounded context timeout,
// centralizing context creation so every registry HTTP call is
// consistently cancellable rather than relying solely on the client's
// blanket timeout.
func doGet(url string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not build request: %w", err)
	}

	return httpClient.Do(req)
}

// readBodyLimited reads resp.Body up to maxDownloadBytes and returns an
// error if the response is larger than that, rather than buffering an
// unbounded amount of attacker-controlled data into memory.
func readBodyLimited(resp *http.Response) ([]byte, error) {
	limited := io.LimitReader(resp.Body, maxDownloadBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if len(data) > maxDownloadBytes {
		return nil, fmt.Errorf("response exceeded maximum allowed size of %d bytes", maxDownloadBytes)
	}
	return data, nil
}

// FetchIndex downloads and parses the registry index.
func FetchIndex() (*Index, error) {
	resp, err := doGet(RegistryIndex)
	if err != nil {
		return nil, fmt.Errorf("could not reach registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned status %d", resp.StatusCode)
	}

	data, err := readBodyLimited(resp)
	if err != nil {
		return nil, fmt.Errorf("could not read registry index: %w", err)
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("could not parse registry index: %w", err)
	}

	return &index, nil
}

// FetchTheme downloads a theme JSON by name from the registry and saves
// it into themesDir as "<name>.json".
//
// name is validated against ValidThemeName before being used in either
// the download URL or the output filesystem path, preventing path
// traversal (e.g. "../../.bashrc") and URL injection. The downloaded
// body is size-capped and parsed as JSON before being written to disk,
// so a non-JSON or oversized response never reaches the themes
// directory as a ".json" file.
func FetchTheme(name string, themesDir string) error {
	if !ValidThemeName(name) {
		return fmt.Errorf("invalid theme name %q: only letters, digits, hyphens, and underscores are allowed", name)
	}

	url := fmt.Sprintf("%s/themes/%s.json", RegistryBase, name)

	resp, err := doGet(url)
	if err != nil {
		return fmt.Errorf("could not download theme: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("theme '%s' not found in registry", name)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	data, err := readBodyLimited(resp)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	// Reject anything that isn't even structurally valid JSON before it
	// ever touches disk. Full schema validation (config.ValidateTheme)
	// still happens in the caller after this — this is just a sanity
	// gate against writing garbage/non-JSON to the themes directory.
	if !json.Valid(data) {
		return fmt.Errorf("downloaded theme is not valid JSON")
	}

	outPath := filepath.Join(themesDir, name+".json")
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return fmt.Errorf("could not save theme: %w", err)
	}

	return nil
}

// Search filters the index by name, description, or tag.
// Matching is case-insensitive.
func Search(index *Index, query string) []ThemeEntry {
	q := strings.ToLower(query)
	var results []ThemeEntry
	for _, t := range index.Themes {
		if strings.Contains(strings.ToLower(t.Name), q) || strings.Contains(strings.ToLower(t.Description), q) {
			results = append(results, t)
			continue
		}
		for _, tag := range t.Tags {
			if strings.Contains(strings.ToLower(tag), q) {
				results = append(results, t)
				break
			}
		}
	}
	return results
}
