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