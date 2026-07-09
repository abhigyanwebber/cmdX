package assets

// AssetType defines what kind of asset this is
type AssetType string

const (
	AssetTypeSpinner AssetType = "spinner"
	AssetTypeIcon    AssetType = "icon"
	AssetTypeBanner  AssetType = "banner"
	AssetTypeDivider AssetType = "divider"
	AssetTypeFloater AssetType = "floater"
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
	Floater     *FloaterConfig    `json:"floater,omitempty"`
	Mascot      *MascotConfig     `json:"mascot,omitempty"`
	StatusBar   *StatusBarConfig  `json:"status_bar,omitempty"`
	Sound       *SoundThemeConfig `json:"sound,omitempty"`
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

// FloaterPosition identifies which corner of the terminal a floater
// asset is anchored to.
type FloaterPosition string

// Available floater corner positions.
const (
	FloaterTopLeft     FloaterPosition = "top-left"
	FloaterTopRight    FloaterPosition = "top-right"
	FloaterBottomLeft  FloaterPosition = "bottom-left"
	FloaterBottomRight FloaterPosition = "bottom-right"
)

// ValidFloaterPositions lists every accepted FloaterPosition value, used
// by ValidateAsset and by the CLI's --position flag help text.
var ValidFloaterPositions = []FloaterPosition{
	FloaterTopLeft, FloaterTopRight, FloaterBottomLeft, FloaterBottomRight,
}

// IsValidFloaterPosition reports whether pos is one of the four
// recognized corner positions.
func IsValidFloaterPosition(pos FloaterPosition) bool {
	for _, p := range ValidFloaterPositions {
		if pos == p {
			return true
		}
	}
	return false
}

// FloaterConfig defines floater-specific settings: a small decorative
// PNG anchored to one corner of the terminal, rendered alongside (not
// replacing) the normal prompt/banner content.
type FloaterConfig struct {
	File          string          `json:"file"`
	Position      FloaterPosition `json:"position"`
	MaxWidth      int             `json:"max_width"`
	MaxHeight     int             `json:"max_height"`
	MarginX       int             `json:"margin_x"`
	MarginY       int             `json:"margin_y"`
	AnimateFrames []string        `json:"animate_frames,omitempty"` // optional, for an animated floater
	IntervalMs    int             `json:"interval_ms,omitempty"`    // required if AnimateFrames is set
}

// ── Mascot ───────────────────────────────────────────────────────────────────

// AssetTypeMascot identifies a reactive character asset.
const AssetTypeMascot AssetType = "mascot"

// MascotState is a named state the mascot can be in.
// Developers define their own state names in asset.json — the built-in
// states (idle, working, success, error, warning, sleeping) are just
// conventions. Any string is valid as a custom state name.
type MascotState string

const (
	MascotStateIdle     MascotState = "idle"
	MascotStateWorking  MascotState = "working"
	MascotStateSuccess  MascotState = "success"
	MascotStateError    MascotState = "error"
	MascotStateWarning  MascotState = "warning"
	MascotStateSleeping MascotState = "sleeping"
)

// MascotTriggerType defines what kind of condition switches the mascot
// to a given state.
type MascotTriggerType string

const (
	// TriggerExitCode fires when the last command's exit code matches
	// a range or specific value. e.g. "0" for success, "1-127" for error.
	TriggerExitCode MascotTriggerType = "exit_code"

	// TriggerOutputRegex fires when the last command's stdout/stderr
	// matches a regex pattern. e.g. "error|fatal|panic" for errors.
	TriggerOutputRegex MascotTriggerType = "output_regex"

	// TriggerEnvVar fires when an environment variable matches a value.
	// e.g. CMDX_MASCOT_STATE=sleeping to force a state externally.
	TriggerEnvVar MascotTriggerType = "env_var"

	// TriggerIdleTime fires when the shell has been idle for longer
	// than the configured duration in seconds.
	TriggerIdleTime MascotTriggerType = "idle_time"

	// TriggerGitStatus fires based on the git status of the current
	// directory. Values: "dirty", "clean", "untracked", "ahead", "behind".
	TriggerGitStatus MascotTriggerType = "git_status"

	// TriggerCommand fires when the last command matches a glob pattern.
	// e.g. "go build*", "npm*", "make*" to show a working state.
	TriggerCommand MascotTriggerType = "command"

	// TriggerAlways is the default/fallback trigger — this state is
	// active whenever no other trigger matches. Typically used for "idle".
	TriggerAlways MascotTriggerType = "always"
)

// MascotTrigger defines a single condition that activates a mascot state.
type MascotTrigger struct {
	// Type is the kind of condition to check.
	Type MascotTriggerType `json:"type"`

	// Value is the trigger-specific match value:
	//   exit_code:    "0", "1", "1-127", "128+"
	//   output_regex: any valid Go regex
	//   env_var:      "VAR_NAME=value" or just "VAR_NAME" to check if set
	//   idle_time:    seconds as a string, e.g. "30"
	//   git_status:   "dirty", "clean", "untracked", "ahead", "behind"
	//   command:      glob pattern, e.g. "go build*"
	//   always:       ignored
	Value string `json:"value,omitempty"`

	// Priority controls which trigger wins when multiple match at once.
	// Higher number = higher priority. Defaults to 0.
	Priority int `json:"priority,omitempty"`
}

// MascotStateConfig defines everything about one state: which PNG files
// to display (supporting per-state animation), how to render them, and
// what conditions trigger this state.
type MascotStateConfig struct {
	// Frames is the list of PNG filenames for this state's animation.
	// A single file means the state is static. Multiple files are
	// cycled at IntervalMs speed.
	Frames []string `json:"frames"`

	// IntervalMs is the animation frame delay for this state.
	// 0 means no animation (static image or use the global default).
	IntervalMs int `json:"interval_ms,omitempty"`

	// Triggers is the list of conditions that activate this state.
	// The first trigger type that matches wins (in priority order).
	// If empty, this state is never entered automatically — it can
	// only be set manually via --state flag or env var.
	Triggers []MascotTrigger `json:"triggers,omitempty"`

	// RenderOverride allows a specific state to use a different render
	// size, mode, or color depth from the mascot's global render config.
	// Useful for making the error state render larger/brighter, or
	// making sleeping render smaller/dimmer.
	RenderOverride *MascotRenderOverride `json:"render_override,omitempty"`

	// TransitionFrames is an optional list of PNG filenames played
	// once when entering this state (before looping Frames).
	// Useful for a "flinch" animation on error or a "wake up" on success.
	TransitionFrames []string `json:"transition_frames,omitempty"`
}

// MascotRenderOverride allows a specific mascot state to override
// the global render dimensions or color treatment. Zero values mean
// "use the global default".
type MascotRenderOverride struct {
	Width     int       `json:"width,omitempty"`
	Height    int       `json:"height,omitempty"`
	ColorMode ColorMode `json:"color_mode,omitempty"`
	// Tint applies a color overlay to the rendered output.
	// Useful for tinting the mascot red on error, green on success, etc.
	// Value is a hex color: "#FF4444". Empty means no tint.
	Tint string `json:"tint,omitempty"`
}

// MascotConfig is the top-level mascot configuration block in asset.json.
type MascotConfig struct {
	// States maps state names to their configurations. Developers can
	// define any state names they want — the built-in names (idle,
	// working, success, error, warning, sleeping) are just conventions.
	States map[MascotState]MascotStateConfig `json:"states"`

	// DefaultState is the state to show when no trigger matches.
	// Defaults to "idle" if not specified.
	DefaultState MascotState `json:"default_state,omitempty"`

	// Position determines where the mascot appears in the terminal.
	// Reuses the FloaterPosition type since it's the same four-corner model.
	Position FloaterPosition `json:"position"`

	// MaxWidth and MaxHeight cap the render size regardless of the
	// global RenderConfig dimensions.
	MaxWidth  int `json:"max_width"`
	MaxHeight int `json:"max_height"`

	// MarginX and MarginY set the gap between the mascot and the
	// terminal edge in character columns/rows.
	MarginX int `json:"margin_x"`
	MarginY int `json:"margin_y"`

	// GlobalIntervalMs is the default animation speed used for any
	// state that doesn't set its own IntervalMs.
	GlobalIntervalMs int `json:"global_interval_ms,omitempty"`

	// HookVar is the name of the environment variable that the injected
	// shell hook sets to communicate the current trigger state to cmdX.
	// Defaults to "CMDX_MASCOT_TRIGGER". Advanced users can override
	// this if they have their own env var naming conventions.
	HookVar string `json:"hook_var,omitempty"`
}

// ── Status Bar ───────────────────────────────────────────────────────────────

// AssetTypeStatusBar identifies a composable terminal status bar asset.
const AssetTypeStatusBar AssetType = "status-bar"

// StatusBarPosition controls where the bar renders relative to the prompt.
type StatusBarPosition string

const (
	StatusBarAbove  StatusBarPosition = "above"  // render above the prompt line
	StatusBarBelow  StatusBarPosition = "below"  // render below the prompt line
	StatusBarInline StatusBarPosition = "inline" // replace the prompt entirely
)

// SeparatorStyle controls the glyph used between segments.
type SeparatorStyle string

const (
	SeparatorPowerline SeparatorStyle = "powerline" // ▶ / ◀ arrows
	SeparatorRounded   SeparatorStyle = "rounded"   //  /  pill shapes
	SeparatorAngular   SeparatorStyle = "angular"   // / sharp angles
	SeparatorBar       SeparatorStyle = "bar"        // │ plain pipe
	SeparatorDot       SeparatorStyle = "dot"        // · dot
	SeparatorSlash     SeparatorStyle = "slash"      // / forward slash
	SeparatorNone      SeparatorStyle = "none"       // no separator
	SeparatorCustom    SeparatorStyle = "custom"     // developer-defined glyph
)

// SegmentType identifies what a status bar segment displays.
type SegmentType string

const (
	SegmentGit         SegmentType = "git"          // branch, status, ahead/behind
	SegmentDirectory   SegmentType = "directory"    // current path (truncated)
	SegmentTime        SegmentType = "time"         // current time
	SegmentDate        SegmentType = "date"         // current date
	SegmentExitCode    SegmentType = "exit_code"    // last command exit code
	SegmentDuration    SegmentType = "duration"     // last command duration
	SegmentEnvVar      SegmentType = "env_var"      // value of an env variable
	SegmentVirtualEnv  SegmentType = "virtualenv"   // Python venv / conda env
	SegmentGoVersion   SegmentType = "go_version"   // Go toolchain version
	SegmentNodeVersion SegmentType = "node_version" // Node.js version
	SegmentKubernetes  SegmentType = "kubernetes"   // k8s context/namespace
	SegmentBattery     SegmentType = "battery"      // system battery level
	SegmentUser        SegmentType = "user"         // current username
	SegmentHost        SegmentType = "host"         // hostname
	SegmentText        SegmentType = "text"         // static custom text
	SegmentCommand     SegmentType = "command"      // output of a shell command
)

// StatusBarZone controls which side of the bar a segment appears on.
type StatusBarZone string

const (
	ZoneLeft   StatusBarZone = "left"
	ZoneCenter StatusBarZone = "center"
	ZoneRight  StatusBarZone = "right"
)

// SegmentCondition defines when a segment is shown. All conditions must
// pass for the segment to render. Empty conditions always show.
type SegmentCondition struct {
	// EnvSet shows the segment only when this env var is set.
	EnvSet string `json:"env_set,omitempty"`
	// EnvNotSet shows the segment only when this env var is NOT set.
	EnvNotSet string `json:"env_not_set,omitempty"`
	// InGitRepo shows/hides the segment based on whether the current
	// directory is inside a git repository.
	InGitRepo *bool `json:"in_git_repo,omitempty"`
	// MinWidth only shows the segment when the terminal is at least
	// this many columns wide. Prevents crowding on narrow terminals.
	MinWidth int `json:"min_width,omitempty"`
	// ExitCodeNonZero shows the segment only when last exit code != 0.
	ExitCodeNonZero bool `json:"exit_code_non_zero,omitempty"`
}

// SegmentConfig defines a single status bar segment: what it shows,
// how it looks, where it goes, and when it appears.
type SegmentConfig struct {
	// Type is what this segment displays.
	Type SegmentType `json:"type"`

	// Zone is left, center, or right.
	Zone StatusBarZone `json:"zone"`

	// Label is an optional prefix shown before the segment value.
	// e.g. " " for git branch, "⏱ " for duration.
	Label string `json:"label,omitempty"`

	// Format is a Go time format string for time/date segments,
	// or a template like "{value}" for other segment types.
	// Segment-specific placeholders: git uses {branch}, {dirty}, {ahead}, {behind}.
	Format string `json:"format,omitempty"`

	// Color is the foreground color for this segment (hex or semantic key).
	Color string `json:"color,omitempty"`

	// Background is the background color for this segment.
	// Required for Powerline/rounded/angular separator styles to render correctly.
	Background string `json:"background,omitempty"`

	// Bold renders the segment text bold.
	Bold bool `json:"bold,omitempty"`

	// Dim renders the segment text at reduced brightness.
	Dim bool `json:"dim,omitempty"`

	// MaxLength truncates the segment value to this many characters
	// (with "…" appended). 0 means no limit.
	MaxLength int `json:"max_length,omitempty"`

	// Padding adds space around the segment value inside the background block.
	Padding int `json:"padding,omitempty"`

	// Conditions is a list of show/hide rules. All must pass.
	Conditions []SegmentCondition `json:"conditions,omitempty"`

	// Command is the shell command to run for SegmentCommand type.
	// The output's first line is used as the segment value.
	Command string `json:"command,omitempty"`

	// EnvVar is the environment variable to read for SegmentEnvVar type.
	EnvVar string `json:"env_var,omitempty"`

	// Order controls the left-to-right position within the zone.
	// Lower order = closer to the start of the zone.
	Order int `json:"order,omitempty"`
}

// StatusBarConfig is the top-level config block for a status bar asset.
type StatusBarConfig struct {
	// Segments is the full list of segments that make up this bar.
	// Order within a zone is controlled by each segment's Order field.
	Segments []SegmentConfig `json:"segments"`

	// Position controls where the bar renders relative to the prompt.
	Position StatusBarPosition `json:"position"`

	// SeparatorStyle is the default separator style between segments.
	SeparatorStyle SeparatorStyle `json:"separator_style"`

	// CustomSeparator is the glyph to use when SeparatorStyle is "custom".
	// Left and Right are the glyphs used at the transition between segments
	// (pointing right for left-zone, left for right-zone).
	CustomSeparatorLeft  string `json:"custom_separator_left,omitempty"`
	CustomSeparatorRight string `json:"custom_separator_right,omitempty"`

	// Height is the number of terminal rows the bar occupies.
	// 1 is standard; 2 is used for double-line bars.
	Height int `json:"height"`

	// FillBackground fills the entire bar width with the background color
	// of the first/last segment in each zone, creating a solid bar look.
	FillBackground bool `json:"fill_background,omitempty"`

	// HideOnNarrow hides the entire bar when the terminal is narrower
	// than this many columns. 0 means always show.
	HideOnNarrow int `json:"hide_on_narrow,omitempty"`
}

// ── Sound Theme ──────────────────────────────────────────────────────────────

// AssetTypeSound identifies a sound theme asset — audio feedback tied
// to shell events. Unlike every other asset type, sound assets don't
// go through chafa at all; they shell out to a platform audio player.
const AssetTypeSound AssetType = "sound"

// SoundTrigger defines a condition that triggers a sound effect.
// Reuses the same trigger vocabulary as mascots (MascotTriggerType) so
// theme authors only need to learn one trigger syntax across both
// asset types — the same exit_code/output_regex/git_status/command
// patterns that drive a mascot's state can drive a sound effect.
type SoundTrigger struct {
	Type     MascotTriggerType `json:"type"`
	Value    string            `json:"value,omitempty"`
	Priority int               `json:"priority,omitempty"`
}

// SoundEffect is one named, triggerable sound.
type SoundEffect struct {
	// File is the audio file path relative to the asset directory.
	// WAV is the only format guaranteed to work with zero extra
	// dependencies (native support via PowerShell on Windows, afplay
	// on macOS, and paplay/aplay on most Linux distros). MP3/OGG/etc.
	// work too, but only if a Player capable of decoding them is
	// configured (e.g. ffplay via a custom Player template) — see
	// SoundThemeConfig.Player.
	File string `json:"file"`

	// Volume scales this effect's playback volume, 0.0–1.0. Not every
	// player backend supports volume control (see player.go) — when
	// unsupported, this is silently ignored rather than erroring, since
	// the sound still plays correctly, just without the effect.
	Volume float64 `json:"volume,omitempty"`

	// CooldownMs is the minimum time in milliseconds between repeated
	// plays of this specific effect, tracked persistently across CLI
	// invocations (each shell hook call is a fresh process, so this
	// can't be tracked in memory — see player.go's cooldown state
	// file). Prevents an obnoxious sound from firing on every single
	// command if its trigger condition is broad. 0 means no cooldown.
	CooldownMs int `json:"cooldown_ms,omitempty"`

	// Async controls whether the shell blocks waiting for playback to
	// finish (false — the default, useful for short confirmation
	// sounds where you want to see completion) or fires the sound in
	// the background without blocking the next prompt (true — useful
	// for longer ambient/notification sounds).
	Async bool `json:"async,omitempty"`

	// Triggers is the list of conditions that play this effect. The
	// highest-priority matching trigger across all effects wins, same
	// resolution model as mascot states.
	Triggers []SoundTrigger `json:"triggers,omitempty"`
}

// SoundThemeConfig is the top-level sound theme configuration block.
type SoundThemeConfig struct {
	Enabled bool `json:"enabled"`

	// GlobalVolume is a 0.0–1.0 multiplier applied on top of each
	// effect's own Volume, letting a user turn everything down (or up)
	// without editing every individual effect.
	GlobalVolume float64 `json:"global_volume,omitempty"`

	// Player, if set, overrides the auto-detected platform default
	// audio player entirely. This is the developer-freedom escape
	// hatch: not every format or workflow is served by the built-in
	// defaults (afplay/paplay/aplay/PowerShell SoundPlayer), so you can
	// point this at anything — ffplay for broader format support, a
	// custom wrapper script, whatever fits your setup.
	//
	// Template syntax: "%f" is replaced with the sound file's absolute
	// path. Example: "ffplay -nodisp -autoexit -loglevel quiet %f".
	// The template is split on whitespace and executed directly via
	// argv (not through a shell), so it's not vulnerable to shell
	// injection, but also doesn't support shell quoting — arguments
	// with spaces aren't supported in the template itself.
	Player string `json:"player,omitempty"`

	// Sounds maps effect name (any string you choose — "success",
	// "error", "build-complete", "coffee-break", anything) to its
	// configuration.
	Sounds map[string]SoundEffect `json:"sounds"`
}
