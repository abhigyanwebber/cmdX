package assets

// RenderOverrides allows callers to override individual render options
// from an asset's manifest at runtime, without modifying the manifest
// itself. Zero values mean "use the manifest default" for every field —
// overrides are applied selectively, not wholesale.
//
// This is the core API that gives developers runtime control over the
// chafa render pipeline. The CLI exposes these as flags on commands like
// "cmdx asset preview" and "cmdx asset render"; external tooling and
// scripts can also call ApplyOverrides directly via the assets package.
type RenderOverrides struct {
	// Mode overrides the render mode (braille/blocks/ascii/sixel/color).
	// Empty string means "use manifest default".
	Mode RenderMode

	// ColorMode overrides the color depth. Empty means manifest default.
	ColorMode ColorMode

	// SymbolSet overrides the chafa symbol set used for rendering.
	// Only applies when Mode is not one of the named modes that imply
	// a specific symbol set (braille/blocks/ascii). Empty means default.
	SymbolSet SymbolSet

	// Width overrides the output width in terminal columns. 0 means default.
	Width int

	// Height overrides the output height in terminal rows. 0 means default.
	Height int

	// Dither, if non-nil, overrides the dithering setting.
	// Pointer so that false can be distinguished from "not set".
	Dither *bool

	// Stretch, if non-nil, overrides stretch-to-fit.
	Stretch *bool

	// Threshold overrides the alpha/contrast threshold (0.0–1.0).
	// Negative value means "use manifest default".
	Threshold float64
}

// ApplyOverrides merges a RenderOverrides on top of a base ChafaOptions,
// returning the merged options. Only non-zero override fields are applied;
// base fields are kept for anything the override leaves unset.
//
// This lets "cmdx asset preview my-asset --mode sixel --width 32" work
// without altering my-asset's asset.json: the CLI builds a RenderOverrides
// from its flags and calls ApplyOverrides(baseFromManifest, overrides).
func ApplyOverrides(base ChafaOptions, overrides RenderOverrides) ChafaOptions {
	result := base

	if overrides.Mode != "" {
		result.Mode = overrides.Mode
	}
	if overrides.ColorMode != "" {
		result.ColorMode = overrides.ColorMode
	}
	if overrides.SymbolSet != "" {
		result.SymbolSet = overrides.SymbolSet
	}
	if overrides.Width > 0 {
		result.Width = overrides.Width
	}
	if overrides.Height > 0 {
		result.Height = overrides.Height
	}
	if overrides.Dither != nil {
		result.Dither = *overrides.Dither
	}
	if overrides.Stretch != nil {
		result.Stretch = *overrides.Stretch
	}
	if overrides.Threshold > 0 {
		result.Threshold = overrides.Threshold
	}

	return result
}

// OverridesFromFlags builds a RenderOverrides from the common render
// flag values that "cmdx asset preview" and "cmdx asset render" expose.
// Any flag left at its zero value is treated as "not set" and won't
// override the manifest default.
func OverridesFromFlags(mode, color, symbols string, width, height int, dither, stretch bool, threshold float64, ditherSet, stretchSet bool) RenderOverrides {
	o := RenderOverrides{
		Mode:      RenderMode(mode),
		ColorMode: ColorMode(color),
		SymbolSet: SymbolSet(symbols),
		Width:     width,
		Height:    height,
		Threshold: threshold,
	}
	if ditherSet {
		o.Dither = &dither
	}
	if stretchSet {
		o.Stretch = &stretch
	}
	return o
}

// optionsFromConfig converts a RenderConfig (the asset.json schema type)
// to a ChafaOptions (the runtime render call type). This is the bridge
// between what's stored in the manifest and what gets passed to chafa.
func optionsFromConfig(r RenderConfig) ChafaOptions {
	return ChafaOptions{
		Mode:      r.Mode,
		ColorMode: r.ColorMode,
		SymbolSet: r.SymbolSet,
		Width:     r.Width,
		Height:    r.Height,
		Dither:    r.Dither,
		Stretch:   r.Stretch,
		Threshold: r.Threshold,
		Animate:   r.Animate,
		FPS:       r.FPS,
	}
}
