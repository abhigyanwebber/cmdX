package fonts

import "strings"

// NerdFontsVersion pins the Nerd Fonts release tag used to build
// download URLs for curated catalog entries. Bump this constant when a
// newer Nerd Fonts release is out. Developers who need a different
// version immediately (without waiting for a cmdX update) can override
// it per-install with InstallOptions.URL pointing at any release asset.
const NerdFontsVersion = "v3.2.1"

// nerdFontsBaseURL is the GitHub releases base for the Nerd Fonts project.
const nerdFontsBaseURL = "https://github.com/ryanoasis/nerd-fonts/releases/download"

// Catalog is the curated list of popular terminal-friendly Nerd Fonts.
// This is a convenience starting point, not a restriction — any font
// can be installed via InstallOptions.URL regardless of whether it's
// listed here.
var Catalog = []FontEntry{
	{
		Name:        "firacode",
		DisplayName: "FiraCode Nerd Font",
		Description: "Popular ligature-rich monospace font, clean and modern",
		ZipName:     "FiraCode",
		License:     "OFL-1.1",
		Monospace:   true,
	},
	{
		Name:        "jetbrainsmono",
		DisplayName: "JetBrainsMono Nerd Font",
		Description: "Designed for code readability by JetBrains, wide character set",
		ZipName:     "JetBrainsMono",
		License:     "OFL-1.1",
		Monospace:   true,
	},
	{
		Name:        "hack",
		DisplayName: "Hack Nerd Font",
		Description: "Bitmap-inspired, highly legible at small sizes",
		ZipName:     "Hack",
		License:     "MIT",
		Monospace:   true,
	},
	{
		Name:        "cascadiacode",
		DisplayName: "CascadiaCode Nerd Font",
		Description: "Microsoft's terminal font, ships with Windows Terminal",
		ZipName:     "CascadiaCode",
		License:     "OFL-1.1",
		Monospace:   true,
	},
	{
		Name:        "meslo",
		DisplayName: "Meslo Nerd Font",
		Description: "Customized Menlo, a favorite for macOS terminal setups",
		ZipName:     "Meslo",
		License:     "Apache-2.0",
		Monospace:   true,
	},
	{
		Name:        "sourcecodepro",
		DisplayName: "SourceCodePro Nerd Font",
		Description: "Adobe's monospace font, clean and professional",
		ZipName:     "SourceCodePro",
		License:     "OFL-1.1",
		Monospace:   true,
	},
	{
		Name:        "ubuntumono",
		DisplayName: "UbuntuMono Nerd Font",
		Description: "Ubuntu's default monospace font",
		ZipName:     "UbuntuMono",
		License:     "UFL-1.0",
		Monospace:   true,
	},
	{
		Name:        "iosevka",
		DisplayName: "Iosevka Nerd Font",
		Description: "Narrow, space-efficient font good for dense terminals",
		ZipName:     "Iosevka",
		License:     "OFL-1.1",
		Monospace:   true,
	},
	{
		Name:        "robotomono",
		DisplayName: "RobotoMono Nerd Font",
		Description: "Google's monospace companion to Roboto",
		ZipName:     "RobotoMono",
		License:     "Apache-2.0",
		Monospace:   true,
	},
	{
		Name:        "dejavusansmono",
		DisplayName: "DejaVuSansMono Nerd Font",
		Description: "Broad Unicode coverage, good for symbol-heavy prompts",
		ZipName:     "DejaVuSansMono",
		License:     "Bitstream Vera",
		Monospace:   true,
	},
}

// Find returns the catalog entry matching name (case-sensitive, as
// listed), or nil if not found.
func Find(name string) *FontEntry {
	for i := range Catalog {
		if Catalog[i].Name == name {
			return &Catalog[i]
		}
	}
	return nil
}

// Search returns catalog entries whose name or description contains
// query (case-insensitive).
func Search(query string) []FontEntry {
	q := strings.ToLower(query)
	var results []FontEntry
	for _, f := range Catalog {
		if strings.Contains(strings.ToLower(f.Name), q) ||
			strings.Contains(strings.ToLower(f.Description), q) ||
			strings.Contains(strings.ToLower(f.DisplayName), q) {
			results = append(results, f)
		}
	}
	return results
}

// DownloadURL builds the Nerd Fonts release URL for a catalog entry,
// using the pinned NerdFontsVersion.
func (f FontEntry) DownloadURL() string {
	return nerdFontsBaseURL + "/" + NerdFontsVersion + "/" + f.ZipName + ".zip"
}
