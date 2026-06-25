package assets

// AssetType defines what kind of asset this is
type AssetType string

const (
	AssetTypeSpinner AssetType = "spinner"
	AssetTypeIcon    AssetType = "icon"
	AssetTypeBanner  AssetType = "banner"
	AssetTypeDivider AssetType = "divider"
)

// RenderMode defines how PNG frames are converted to terminal output
type RenderMode string

const (
	RenderModeBraille RenderMode = "braille" // ⣾⣽⣻ — smooth, detailed
	RenderModeBlocks  RenderMode = "blocks"  // █▓▒░ — bold, chunky
	RenderModeASCII   RenderMode = "ascii"   // character density mapping
	RenderModeSixel   RenderMode = "sixel"   // actual pixels (WT + iTerm2)
	RenderModeSymbols RenderMode = "symbols" // custom unicode chars
	RenderModeColor   RenderMode = "color"   // full 24-bit color blocks
)

// ColorMode defines how colors are handled during render
type ColorMode string

const (
	ColorModeNone      ColorMode = "none"      // no color
	ColorModeAnsi      ColorMode = "ansi"      // 16 colors
	ColorMode256       ColorMode = "256"       // 256 colors
	ColorModeTrueColor ColorMode = "truecolor" // 24-bit RGB
)

// SymbolSet defines which unicode symbols chafa uses for rendering
type SymbolSet string

const (
	SymbolSetAll     SymbolSet = "all"
	SymbolSetBraille SymbolSet = "braille"
	SymbolSetBlock   SymbolSet = "block"
	SymbolSetBorder  SymbolSet = "border"
	SymbolSetLegacy  SymbolSet = "legacy"
	SymbolSetAscii   SymbolSet = "ascii"
	SymbolSetNone    SymbolSet = "none"
)

// Asset is the root structure of an asset.json manifest
type Asset struct {
	Name        string         `json:"name"`
	Type        AssetType      `json:"type"`
	Version     string         `json:"version"`
	Author      string         `json:"author"`
	Description string         `json:"description"`
	License     string         `json:"license,omitempty"`
	Homepage    string         `json:"homepage,omitempty"`
	Render      RenderConfig   `json:"render"`
	Spinner     *SpinnerConfig `json:"spinner,omitempty"`
	Icon        *IconConfig    `json:"icon,omitempty"`
	Banner      *BannerConfig  `json:"banner,omitempty"`
	Divider     *DividerConfig `json:"divider,omitempty"`
}

// RenderConfig defines how chafa converts images to terminal output
type RenderConfig struct {
	Mode      RenderMode `json:"mode"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	ColorMode ColorMode  `json:"color_mode"`
	SymbolSet SymbolSet  `json:"symbol_set"`
	Dither    bool       `json:"dither"`
	Threshold float64    `json:"threshold"` // 0.0-1.0 contrast threshold
	Stretch   bool       `json:"stretch"`   // stretch to fit dimensions
	Animate   bool       `json:"animate"`   // for GIF support
	FPS       int        `json:"fps,omitempty"`
}

// SpinnerConfig defines spinner-specific settings
type SpinnerConfig struct {
	Frames     []string `json:"frames"` // PNG filenames in order
	IntervalMs int      `json:"interval_ms"`
	Reverse    bool     `json:"reverse"` // play frames in reverse
	Bounce     bool     `json:"bounce"`  // play forward then reverse
	Loop       bool     `json:"loop"`
	OnComplete string   `json:"on_complete"` // "clear", "persist", "checkmark"
}

// IconConfig defines icon-specific settings
type IconConfig struct {
	Files map[string]string `json:"files"` // e.g. "directory": "dir.png"
}

// BannerConfig defines banner-specific settings
type BannerConfig struct {
	File      string `json:"file"`
	Position  string `json:"position"` // "left", "center", "right"
	MaxWidth  int    `json:"max_width"`
	MaxHeight int    `json:"max_height"`
}

// DividerConfig defines divider-specific settings
type DividerConfig struct {
	File     string `json:"file"`
	TileMode string `json:"tile_mode"` // "repeat", "stretch", "center"
	Height   int    `json:"height"`
}
