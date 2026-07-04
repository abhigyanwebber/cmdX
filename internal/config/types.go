// Package config defines the cmdX theme schema (the Theme struct and all
// its nested sections) along with loading and validation logic for
// theme JSON files.
package config

// Theme is the root structure of a .json theme file. Every section is
// optional except where noted by ValidateTheme, but a fully-specified
// theme covers colors, prompt, loader, progress bar, cursor, borders,
// banner, graphics effects, wallpaper, icons, and linked PNG assets.
type Theme struct {
	Meta        Meta        `json:"meta"`
	Colors      Colors      `json:"colors"`
	Prompt      Prompt      `json:"prompt"`
	Loader      Loader      `json:"loader"`
	ProgressBar ProgressBar `json:"progress_bar"`
	Cursor      Cursor      `json:"cursor"`
	Borders     Borders     `json:"borders"`
	Banner      Banner      `json:"banner"`
	Graphics    Graphics    `json:"graphics"`
	Wallpaper   Wallpaper   `json:"wallpaper"`
	Icons       IconSet     `json:"icons"`
	Assets      ThemeAssets `json:"assets"`
}

// Meta holds identifying information about a theme: its name, semver
// version, author, and a short human-readable description.
type Meta struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

// Colors defines the full semantic color palette used throughout a
// theme. Each field must be a valid hex color (e.g. "#FF00FF" or "#f0f").
type Colors struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	Accent     string `json:"accent"`
	Error      string `json:"error"`
	Success    string `json:"success"`
	Warning    string `json:"warning"`
	Muted      string `json:"muted"`
}

// Prompt configures the shell prompt: the symbol shown, segment
// separator, layout style ("single" or "multiline"), an optional ordered
// list of segments, and the overall format template.
type Prompt struct {
	Symbol    string   `json:"symbol"`
	Separator string   `json:"separator"`
	Style     string   `json:"style"`
	Segments  []string `json:"segments"`
	Format    string   `json:"format"`
}

// Loader configures the animated spinner: the sequence of frames to
// cycle through, the delay between frames in milliseconds, and the
// spinner's color.
type Loader struct {
	Frames     []string `json:"frames"`
	IntervalMs int      `json:"interval_ms"`
	Color      string   `json:"color"`
}

// ProgressBar configures the filled/empty characters, total width in
// characters, and color used when rendering progress bars.
type ProgressBar struct {
	Filled string `json:"filled"`
	Empty  string `json:"empty"`
	Width  int    `json:"width"`
	Color  string `json:"color"`
}

// Cursor configures the terminal cursor's style ("block", "bar", or
// "underline"), whether it blinks, and its color.
type Cursor struct {
	Style string `json:"style"`
	Blink bool   `json:"blink"`
	Color string `json:"color"`
}

// BorderChars defines the individual glyphs used to draw a box border:
// the four corners plus horizontal and vertical edge characters.
type BorderChars struct {
	TopLeft     string `json:"top_left"`
	TopRight    string `json:"top_right"`
	BottomLeft  string `json:"bottom_left"`
	BottomRight string `json:"bottom_right"`
	Horizontal  string `json:"horizontal"`
	Vertical    string `json:"vertical"`
}

// Borders configures the border drawing style and (for custom styles)
// the individual characters to use via Chars.
type Borders struct {
	Style string      `json:"style"`
	Chars BorderChars `json:"chars"`
}

// Banner configures the startup banner shown when a theme is injected
// into a shell: whether it's enabled, its text (supports the {user}
// placeholder), display style, and color.
type Banner struct {
	Enabled bool   `json:"enabled"`
	Text    string `json:"text"`
	Style   string `json:"style"`
	Color   string `json:"color"`
}

// Graphics holds all visual effect settings for a theme: gradients,
// dividers, text effects, and icon configuration.
type Graphics struct {
	Gradient GradientConfig `json:"gradient"`
	Divider  DividerConfig  `json:"divider"`
	Effects  EffectsConfig  `json:"effects"`
	Icons    IconsConfig    `json:"icons"`
}

// GradientConfig configures LAB-space color gradients applied to text,
// from one hex color to another, in a given direction.
type GradientConfig struct {
	Enabled   bool   `json:"enabled"`
	From      string `json:"from"`
	To        string `json:"to"`
	Direction string `json:"direction"` // horizontal, vertical
}

// DividerConfig configures the divider line style (wave, line, dots,
// stars, double) and color used between preview sections.
type DividerConfig struct {
	Style string `json:"style"` // wave, line, dots, stars, double
	Color string `json:"color"`
}

// EffectsConfig selects which text effect, if any, is applied to the
// banner and prompt (glitch, rainbow, pulse, neon, typewriter, or none).
type EffectsConfig struct {
	Banner string `json:"banner"` // glitch, rainbow, pulse, none
	Prompt string `json:"prompt"` // rainbow, none
}

// IconsConfig enables and configures inline status icons (directory,
// git branch, error, success, time) shown alongside the prompt.
type IconsConfig struct {
	Enabled   bool   `json:"enabled"`
	Directory string `json:"directory"`
	GitBranch string `json:"git_branch"`
	Error     string `json:"error"`
	Success   string `json:"success"`
	Time      string `json:"time"`
}

// Wallpaper configures a desktop/terminal background image: whether
// it's enabled, the file path, opacity, stretch mode, and alignment.
// Currently supported by Windows Terminal, iTerm2, and Kitty.
type Wallpaper struct {
	Enabled   bool    `json:"enabled"`
	Path      string  `json:"path"`
	Opacity   float64 `json:"opacity"`
	Stretch   string  `json:"stretch"`   // fill, uniform, uniformToFill, none
	Alignment string  `json:"alignment"` // center, topLeft, topRight, etc
}

// IconSet configures the full icon glyph set used across the prompt and
// status indicators: font family, per-context glyphs (git status,
// language/tool icons, error/success/warning markers, etc).
type IconSet struct {
	Enabled   bool   `json:"enabled"`
	Font      string `json:"font"` // nerd-fonts, emoji, ascii
	Directory string `json:"directory"`
	File      string `json:"file"`
	GitBranch string `json:"git_branch"`
	GitDirty  string `json:"git_dirty"`
	GitClean  string `json:"git_clean"`
	Error     string `json:"error"`
	Success   string `json:"success"`
	Warning   string `json:"warning"`
	Time      string `json:"time"`
	Package   string `json:"package"`
	Docker    string `json:"docker"`
	Python    string `json:"python"`
	NodeJS    string `json:"node"`
	Rust      string `json:"rust"`
	Go        string `json:"go"`
}

// ThemeAssets links named PNG assets (by name, resolved against the
// user's asset library) to theme slots: spinner, banner, divider, icons.
type ThemeAssets struct {
	Spinner string `json:"spinner,omitempty"`
	Banner  string `json:"banner,omitempty"`
	Divider string `json:"divider,omitempty"`
	Icons   string `json:"icons,omitempty"`
}
