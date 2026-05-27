package config

// Theme is the root structure of a .json theme file
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
}

type Meta struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

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

type Prompt struct {
	Symbol    string   `json:"symbol"`
	Separator string   `json:"separator"`
	Style     string   `json:"style"`
	Segments  []string `json:"segments"`
	Format    string   `json:"format"`
}

type Loader struct {
	Frames     []string `json:"frames"`
	IntervalMs int      `json:"interval_ms"`
	Color      string   `json:"color"`
}

type ProgressBar struct {
	Filled string `json:"filled"`
	Empty  string `json:"empty"`
	Width  int    `json:"width"`
	Color  string `json:"color"`
}

type Cursor struct {
	Style string `json:"style"`
	Blink bool   `json:"blink"`
	Color string `json:"color"`
}

type BorderChars struct {
	TopLeft     string `json:"top_left"`
	TopRight    string `json:"top_right"`
	BottomLeft  string `json:"bottom_left"`
	BottomRight string `json:"bottom_right"`
	Horizontal  string `json:"horizontal"`
	Vertical    string `json:"vertical"`
}

type Borders struct {
	Style string      `json:"style"`
	Chars BorderChars `json:"chars"`
}

type Banner struct {
	Enabled bool   `json:"enabled"`
	Text    string `json:"text"`
	Style   string `json:"style"`
	Color   string `json:"color"`
}

// Graphics holds all visual effect settings for a theme
type Graphics struct {
	Gradient GradientConfig `json:"gradient"`
	Divider  DividerConfig  `json:"divider"`
	Effects  EffectsConfig  `json:"effects"`
	Icons    IconsConfig    `json:"icons"`
}

type GradientConfig struct {
	Enabled   bool   `json:"enabled"`
	From      string `json:"from"`
	To        string `json:"to"`
	Direction string `json:"direction"` // horizontal, vertical
}

type DividerConfig struct {
	Style string `json:"style"` // wave, line, dots, stars, double
	Color string `json:"color"`
}

type EffectsConfig struct {
	Banner string `json:"banner"` // glitch, rainbow, pulse, none
	Prompt string `json:"prompt"` // rainbow, none
}

type IconsConfig struct {
	Enabled   bool   `json:"enabled"`
	Directory string `json:"directory"`
	GitBranch string `json:"git_branch"`
	Error     string `json:"error"`
	Success   string `json:"success"`
	Time      string `json:"time"`
}
