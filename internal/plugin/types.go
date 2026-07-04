// Package plugin defines the manifest schema for cmdX plugins and
// provides discovery, loading, and validation of plugins on disk.
// A plugin extends cmdX's built-in spinners, banners, and prompt
// templates without modifying the core codebase.
package plugin

// Plugin is the root structure of a plugin manifest (plugin.json),
// contributing optional sets of named spinners, banners, and prompts
// that themes and the CLI can reference by name.
type Plugin struct {
	Meta     PluginMeta       `json:"meta"`
	Spinners []SpinnerSet     `json:"spinners,omitempty"`
	Banners  []BannerTemplate `json:"banners,omitempty"`
	Prompts  []PromptTemplate `json:"prompts,omitempty"`
}

// PluginMeta holds identifying information about a plugin: its name,
// semver version, author, description, and an optional homepage URL.
type PluginMeta struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Homepage    string `json:"homepage,omitempty"`
}

// SpinnerSet is a named, reusable spinner animation contributed by a
// plugin: a sequence of frames, the delay between them in milliseconds,
// and an optional color.
type SpinnerSet struct {
	Name       string   `json:"name"`
	Frames     []string `json:"frames"`
	IntervalMs int      `json:"interval_ms"`
	Color      string   `json:"color,omitempty"`
}

// BannerTemplate is a named, reusable startup banner contributed by a
// plugin: display text, style, and an optional color.
type BannerTemplate struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Style string `json:"style"`
	Color string `json:"color,omitempty"`
}

// PromptTemplate is a named, reusable prompt layout contributed by a
// plugin, mirroring the fields of config.Prompt.
type PromptTemplate struct {
	Name      string   `json:"name"`
	Symbol    string   `json:"symbol"`
	Separator string   `json:"separator"`
	Style     string   `json:"style"`
	Segments  []string `json:"segments"`
	Format    string   `json:"format"`
}
