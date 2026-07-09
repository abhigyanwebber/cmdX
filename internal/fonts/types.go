// Package fonts handles discovery, download, and cross-platform
// installation of terminal fonts — primarily Nerd Fonts (patched with
// icon glyphs for the icons/status-bar segment types), but also any
// arbitrary font a developer wants to install via a direct URL.
package fonts

import "time"

// FontEntry describes one font in the curated catalog.
type FontEntry struct {
	// Name is the identifier used on the CLI (e.g. "firacode").
	Name string

	// DisplayName is the human-readable name (e.g. "FiraCode Nerd Font").
	DisplayName string

	// Description is a one-line summary of the font's style.
	Description string

	// ZipName is the filename (without extension) of this font's
	// release zip in the Nerd Fonts GitHub releases.
	ZipName string

	// License is the font's license identifier (e.g. "MIT", "OFL-1.1").
	License string

	// Monospace indicates whether this is a monospace font suitable
	// for terminal use (all curated entries are, but kept explicit
	// for future non-monospace additions).
	Monospace bool
}

// InstalledFont is a record of a font cmdX has installed, tracked so it
// can be cleanly removed later without guessing which files belong to it.
type InstalledFont struct {
	Name        string    `json:"name"`
	Source      string    `json:"source"` // curated catalog name or custom URL
	Files       []string  `json:"files"`  // absolute paths of installed font files
	InstalledAt time.Time `json:"installed_at"`
}

// InstallOptions controls how InstallFont behaves.
type InstallOptions struct {
	// URL, if set, overrides the curated catalog entirely — install
	// any font from any direct zip/font file URL. This is the escape
	// hatch for developer freedom: not every font is in the curated
	// list, and this lets you install anything.
	URL string

	// Variant filters which font weight/style files get installed
	// (e.g. "Regular", "Bold", "Italic", "Mono"). Empty means install
	// every file found in the archive. Matching is a case-insensitive
	// substring match against each filename.
	Variant string

	// AllVariants explicitly installs every variant, ignoring Variant.
	// This is the default if both are unset.
	AllVariants bool

	// Force reinstalls even if the font is already tracked as installed.
	Force bool

	// DryRun prints what would be installed without writing any files.
	DryRun bool

	// MaxDownloadBytes overrides the default download size cap. Some
	// full Nerd Font family archives exceed the conservative default;
	// set this explicitly to opt into a larger download rather than
	// silently allowing unbounded downloads.
	MaxDownloadBytes int64
}
