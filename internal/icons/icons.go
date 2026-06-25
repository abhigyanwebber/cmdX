package icons

import "github.com/abhigyanwebber/cmd-customizer/internal/config"

// Set holds all resolved icons for the current theme
type Set struct {
	Directory string
	File      string
	GitBranch string
	GitDirty  string
	GitClean  string
	Error     string
	Success   string
	Warning   string
	Time      string
	Package   string
	Docker    string
	Python    string
	NodeJS    string
	Rust      string
	Go        string
}

// NerdFonts is the default Nerd Fonts icon set
var NerdFonts = Set{
	Directory: "󰉋",
	File:      "󰈔",
	GitBranch: "",
	GitDirty:  "✗",
	GitClean:  "✓",
	Error:     "",
	Success:   "",
	Warning:   "",
	Time:      "",
	Package:   "󰏗",
	Docker:    "",
	Python:    "",
	NodeJS:    "",
	Rust:      "",
	Go:        "",
}

// Emoji is a fallback icon set using unicode emoji
var Emoji = Set{
	Directory: "📁",
	File:      "📄",
	GitBranch: "🌿",
	GitDirty:  "✗",
	GitClean:  "✓",
	Error:     "❌",
	Success:   "✅",
	Warning:   "⚠️",
	Time:      "🕐",
	Package:   "📦",
	Docker:    "🐳",
	Python:    "🐍",
	NodeJS:    "💚",
	Rust:      "🦀",
	Go:        "🐹",
}

// ASCII is a plain ASCII fallback for terminals without unicode support
var ASCII = Set{
	Directory: "/",
	File:      "-",
	GitBranch: "git:",
	GitDirty:  "x",
	GitClean:  "v",
	Error:     "ERR",
	Success:   "OK",
	Warning:   "(!)",
	Time:      "~",
	Package:   "pkg",
	Docker:    "dkr",
	Python:    "py",
	NodeJS:    "js",
	Rust:      "rs",
	Go:        "go",
}

// Resolve returns the right icon set based on theme config
func Resolve(t *config.Theme) Set {
	if !t.Icons.Enabled {
		return ASCII
	}

	switch t.Icons.Font {
	case "emoji":
		return Emoji
	case "ascii":
		return ASCII
	default:
		// use theme overrides if provided, otherwise nerd fonts defaults
		s := NerdFonts
		if t.Icons.Directory != "" {
			s.Directory = t.Icons.Directory
		}
		if t.Icons.GitBranch != "" {
			s.GitBranch = t.Icons.GitBranch
		}
		if t.Icons.Error != "" {
			s.Error = t.Icons.Error
		}
		if t.Icons.Success != "" {
			s.Success = t.Icons.Success
		}
		return s
	}
}

// FormatDirectory formats a directory path with the right icon
func FormatDirectory(path string, s Set) string {
	return s.Directory + " " + path
}

// FormatGitBranch formats a git branch with the right icon
func FormatGitBranch(branch string, dirty bool, s Set) string {
	if branch == "" {
		return ""
	}
	indicator := s.GitClean
	if dirty {
		indicator = s.GitDirty
	}
	return s.GitBranch + " " + branch + " " + indicator
}
