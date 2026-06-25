package plugin

// Plugin is the root structure of a plugin manifest
type Plugin struct {
	Meta     PluginMeta       `json:"meta"`
	Spinners []SpinnerSet     `json:"spinners,omitempty"`
	Banners  []BannerTemplate `json:"banners,omitempty"`
	Prompts  []PromptTemplate `json:"prompts,omitempty"`
}

type PluginMeta struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Homepage    string `json:"homepage,omitempty"`
}

type SpinnerSet struct {
	Name       string   `json:"name"`
	Frames     []string `json:"frames"`
	IntervalMs int      `json:"interval_ms"`
	Color      string   `json:"color,omitempty"`
}

type BannerTemplate struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Style string `json:"style"`
	Color string `json:"color,omitempty"`
}

type PromptTemplate struct {
	Name      string   `json:"name"`
	Symbol    string   `json:"symbol"`
	Separator string   `json:"separator"`
	Style     string   `json:"style"`
	Segments  []string `json:"segments"`
	Format    string   `json:"format"`
}
