package assets

import (
	"fmt"
	"sort"
	"strings"
)

// ValidateStatusBar checks a StatusBarConfig for required fields and
// validates all segment configurations.
func ValidateStatusBar(sb *StatusBarConfig) error {
	if sb == nil {
		return fmt.Errorf("status_bar config is nil")
	}
	if len(sb.Segments) == 0 {
		return fmt.Errorf("status_bar must define at least one segment")
	}
	if sb.Height <= 0 {
		return fmt.Errorf("status_bar height must be greater than 0")
	}

	validPositions := map[StatusBarPosition]bool{
		StatusBarAbove: true, StatusBarBelow: true, StatusBarInline: true,
	}
	if !validPositions[sb.Position] {
		return fmt.Errorf("invalid status_bar position %q: must be above, below, or inline", sb.Position)
	}

	validSeparators := map[SeparatorStyle]bool{
		SeparatorPowerline: true, SeparatorRounded: true, SeparatorAngular: true,
		SeparatorBar: true, SeparatorDot: true, SeparatorSlash: true,
		SeparatorNone: true, SeparatorCustom: true,
	}
	if !validSeparators[sb.SeparatorStyle] {
		return fmt.Errorf("invalid separator_style %q", sb.SeparatorStyle)
	}
	if sb.SeparatorStyle == SeparatorCustom {
		if sb.CustomSeparatorLeft == "" && sb.CustomSeparatorRight == "" {
			return fmt.Errorf("separator_style 'custom' requires custom_separator_left or custom_separator_right")
		}
	}

	for i, seg := range sb.Segments {
		if err := validateSegment(seg, i); err != nil {
			return err
		}
	}

	return nil
}

// validateSegment checks one segment for required fields.
func validateSegment(seg SegmentConfig, idx int) error {
	validTypes := map[SegmentType]bool{
		SegmentGit: true, SegmentDirectory: true, SegmentTime: true,
		SegmentDate: true, SegmentExitCode: true, SegmentDuration: true,
		SegmentEnvVar: true, SegmentVirtualEnv: true, SegmentGoVersion: true,
		SegmentNodeVersion: true, SegmentKubernetes: true, SegmentBattery: true,
		SegmentUser: true, SegmentHost: true, SegmentText: true, SegmentCommand: true,
	}
	if !validTypes[seg.Type] {
		return fmt.Errorf("segment[%d]: unknown type %q", idx, seg.Type)
	}

	validZones := map[StatusBarZone]bool{ZoneLeft: true, ZoneCenter: true, ZoneRight: true}
	if !validZones[seg.Zone] {
		return fmt.Errorf("segment[%d] (%s): zone must be left, center, or right", idx, seg.Type)
	}

	if seg.Type == SegmentEnvVar && seg.EnvVar == "" {
		return fmt.Errorf("segment[%d] (env_var): env_var field is required", idx)
	}
	if seg.Type == SegmentCommand && seg.Command == "" {
		return fmt.Errorf("segment[%d] (command): command field is required", idx)
	}
	if seg.Type == SegmentText && seg.Format == "" {
		return fmt.Errorf("segment[%d] (text): format field is required for static text segments", idx)
	}

	return nil
}

// separatorGlyphs returns the left-pointing and right-pointing separator
// glyphs for a given style. Left glyph is used in left-zone transitions
// (pointing right), right glyph in right-zone transitions (pointing left).
func separatorGlyphs(style SeparatorStyle, customLeft, customRight string) (string, string) {
	switch style {
	case SeparatorPowerline:
		return "\uE0B0", "\uE0B2" // ▶ ◀ (Powerline)
	case SeparatorRounded:
		return "\uE0B4", "\uE0B6" // rounded Powerline variants
	case SeparatorAngular:
		return "\uE0B8", "\uE0BA" // angular Powerline variants
	case SeparatorBar:
		return "│", "│"
	case SeparatorDot:
		return "·", "·"
	case SeparatorSlash:
		return "/", "\\"
	case SeparatorNone:
		return " ", " "
	case SeparatorCustom:
		l := customLeft
		r := customRight
		if l == "" {
			l = " "
		}
		if r == "" {
			r = " "
		}
		return l, r
	default:
		return " ", " "
	}
}

// StatusBarShellCode generates the shell function code that renders the
// status bar. This is injected into the shell profile alongside the theme.
// The generated code is pure shell — no cmdX binary call needed at prompt
// time, so it works even if cmdX isn't in PATH after injection.
func StatusBarShellCode(sb *StatusBarConfig, colors map[string]string, shell string) string {
	switch shell {
	case "bash":
		return statusBarBash(sb, colors)
	case "zsh":
		return statusBarZsh(sb, colors)
	case "powershell":
		return statusBarPowerShell(sb, colors)
	default:
		return ""
	}
}

// segmentsByZone groups and sorts segments by zone and order.
func segmentsByZone(sb *StatusBarConfig) (left, center, right []SegmentConfig) {
	for _, seg := range sb.Segments {
		switch seg.Zone {
		case ZoneLeft:
			left = append(left, seg)
		case ZoneCenter:
			center = append(center, seg)
		case ZoneRight:
			right = append(right, seg)
		}
	}
	sortSegs := func(segs []SegmentConfig) {
		sort.Slice(segs, func(i, j int) bool {
			return segs[i].Order < segs[j].Order
		})
	}
	sortSegs(left)
	sortSegs(center)
	sortSegs(right)
	return
}

// bashSegmentCode returns the bash shell code that renders one segment's value.
func bashSegmentCode(seg SegmentConfig) string {
	switch seg.Type {
	case SegmentGit:
		format := seg.Format
		if format == "" {
			format = " {branch}{dirty}"
		}
		return fmt.Sprintf(`
    __cmdx_git_branch=""
    __cmdx_git_dirty=""
    if git rev-parse --git-dir >/dev/null 2>&1; then
        __cmdx_git_branch=$(git branch --show-current 2>/dev/null || git rev-parse --short HEAD 2>/dev/null)
        [ -n "$(git status --porcelain 2>/dev/null)" ] && __cmdx_git_dirty="*"
    fi
    if [ -n "$__cmdx_git_branch" ]; then
        __seg_git="%s$__cmdx_git_branch$__cmdx_git_dirty"
    fi`, seg.Label)

	case SegmentDirectory:
		maxLen := seg.MaxLength
		if maxLen == 0 {
			maxLen = 30
		}
		label := seg.Label
		if label == "" {
			label = " "
		}
		return fmt.Sprintf(`    __seg_dir="%s$(pwd | sed "s|$HOME|~|" | awk -F/ '{if(length($0)>%d){print "…/"$NF}else{print $0}}')"`, label, maxLen)

	case SegmentTime:
		format := seg.Format
		if format == "" {
			format = "%H:%M"
		}
		return fmt.Sprintf(`    __seg_time="%s$(date +"%s")"`, seg.Label, format)

	case SegmentDate:
		format := seg.Format
		if format == "" {
			format = "%Y-%m-%d"
		}
		return fmt.Sprintf(`    __seg_date="%s$(date +"%s")"`, seg.Label, format)

	case SegmentExitCode:
		return `    __seg_exit=""
    [ $__cmdx_exit_code -ne 0 ] && __seg_exit=" ✗$__cmdx_exit_code"`

	case SegmentDuration:
		return `    __seg_dur=""
    [ -n "$__cmdx_cmd_duration" ] && __seg_dur=" ⏱${__cmdx_cmd_duration}s"`

	case SegmentUser:
		return fmt.Sprintf(`    __seg_user="%s$USER"`, seg.Label)

	case SegmentHost:
		return fmt.Sprintf(`    __seg_host="%s$(hostname -s)"`, seg.Label)

	case SegmentVirtualEnv:
		return `    __seg_venv=""
    [ -n "$VIRTUAL_ENV" ] && __seg_venv=" ($(basename $VIRTUAL_ENV))"
    [ -n "$CONDA_DEFAULT_ENV" ] && __seg_venv=" ($CONDA_DEFAULT_ENV)"`

	case SegmentEnvVar:
		return fmt.Sprintf(`    __seg_env="%s${%s}"`, seg.Label, seg.EnvVar)

	case SegmentGoVersion:
		return `    __seg_go=""
    command -v go >/dev/null 2>&1 && __seg_go=" go$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//')"`

	case SegmentNodeVersion:
		return `    __seg_node=""
    command -v node >/dev/null 2>&1 && __seg_node=" node$(node --version 2>/dev/null)"`

	case SegmentKubernetes:
		return `    __seg_k8s=""
    command -v kubectl >/dev/null 2>&1 && __seg_k8s=" ⎈$(kubectl config current-context 2>/dev/null)"`

	case SegmentBattery:
		return `    __seg_bat=""
    if [ -f /sys/class/power_supply/BAT0/capacity ]; then
        __bat=$(cat /sys/class/power_supply/BAT0/capacity)
        __seg_bat=" 🔋${__bat}%"
    fi`

	case SegmentText:
		return fmt.Sprintf(`    __seg_text="%s"`, seg.Format)

	case SegmentCommand:
		return fmt.Sprintf(`    __seg_cmd="$(%s 2>/dev/null | head -1)"`, seg.Command)

	default:
		return ""
	}
}

// bashColorCode returns ANSI escape code for a hex color in bash PS1.
// Uses \[\033[38;2;R;G;Bm\] format for truecolor, falling back to the
// color string directly if it looks like a predefined name.
func bashColorCode(hexOrName string, bg bool) string {
	if hexOrName == "" {
		return ""
	}
	hex := strings.TrimPrefix(hexOrName, "#")
	if len(hex) == 6 {
		var r, g, b int
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		if bg {
			return fmt.Sprintf(`\[\033[48;2;%d;%d;%dm\]`, r, g, b)
		}
		return fmt.Sprintf(`\[\033[38;2;%d;%d;%dm\]`, r, g, b)
	}
	return ""
}

const bashReset = `\[\033[0m\]`

// statusBarBash generates a bash PS1 function for the status bar.
func statusBarBash(sb *StatusBarConfig, colors map[string]string) string {
	left, center, right := segmentsByZone(sb)
	sepL, sepR := separatorGlyphs(sb.SeparatorStyle, sb.CustomSeparatorLeft, sb.CustomSeparatorRight)

	var b strings.Builder

	b.WriteString(`
# cmdX status bar — computed segment values
__cmdx_exit_code=$?
__cmdx_cmd_duration=""
`)

	// segment value computations
	for _, seg := range append(append(left, center...), right...) {
		b.WriteString(bashSegmentCode(seg))
		b.WriteString("\n")
	}

	// assemble bar
	b.WriteString("\n    # Assemble status bar\n    __cmdx_bar=\"\"\n")

	for i, seg := range left {
		varName := bashSegVarName(seg, i)
		fg := bashColorCode(resolveSegColor(seg.Color, colors), false)
		bg := bashColorCode(resolveSegColor(seg.Background, colors), true)
		bold := ""
		if seg.Bold {
			bold = `\[\033[1m\]`
		}
		b.WriteString(fmt.Sprintf(`    [ -n "$%s" ] && __cmdx_bar="${__cmdx_bar}%s%s%s$%s%s%s"`,
			varName, bg, fg, bold, varName, bashReset, sepL))
		b.WriteString("\n")
	}

	for i, seg := range right {
		varName := bashSegVarName(seg, i+100)
		fg := bashColorCode(resolveSegColor(seg.Color, colors), false)
		bg := bashColorCode(resolveSegColor(seg.Background, colors), true)
		b.WriteString(fmt.Sprintf(`    [ -n "$%s" ] && __cmdx_bar_right="${__cmdx_bar_right}%s%s$%s%s%s"`,
			varName, bg, fg, varName, bashReset, sepR))
		b.WriteString("\n")
	}

	_ = center
	_ = sepR

	switch sb.Position {
	case StatusBarAbove:
		return fmt.Sprintf("# cmdX status bar\n__cmdx_statusbar() {\n%s\n    echo -e \"$__cmdx_bar\"\n}\nPROMPT_COMMAND=\"__cmdx_statusbar${PROMPT_COMMAND:+;$PROMPT_COMMAND}\"", b.String())
	case StatusBarBelow:
		return fmt.Sprintf("# cmdX status bar\n__cmdx_statusbar() {\n%s\n}\nPROMPT_COMMAND=\"__cmdx_statusbar${PROMPT_COMMAND:+;$PROMPT_COMMAND}\"", b.String())
	default: // inline
		return fmt.Sprintf("# cmdX status bar (inline)\n__cmdx_statusbar() {\n%s\n    echo -n \"$__cmdx_bar \"\n}\nPS1='$(__cmdx_statusbar)'", b.String())
	}
}

// statusBarZsh generates a zsh PROMPT function for the status bar.
func statusBarZsh(sb *StatusBarConfig, colors map[string]string) string {
	left, _, right := segmentsByZone(sb)
	sepL, sepR := separatorGlyphs(sb.SeparatorStyle, sb.CustomSeparatorLeft, sb.CustomSeparatorRight)

	var b strings.Builder
	b.WriteString(`
# cmdX status bar
__cmdx_statusbar() {
    local __cmdx_bar=""
    local __cmdx_bar_right=""
`)

	for i, seg := range left {
		b.WriteString(zshSegmentCode(seg))
		varName := bashSegVarName(seg, i)
		color := resolveSegColor(seg.Color, colors)
		if color != "" && strings.HasPrefix(color, "#") {
			color = color[1:]
		}
		b.WriteString(fmt.Sprintf("    [[ -n \"$%s\" ]] && __cmdx_bar+=\"%%F{#%s}$%s%%f%s\"\n", varName, color, varName, sepL))
	}

	for i, seg := range right {
		b.WriteString(zshSegmentCode(seg))
		varName := bashSegVarName(seg, i+100)
		color := resolveSegColor(seg.Color, colors)
		if color != "" && strings.HasPrefix(color, "#") {
			color = color[1:]
		}
		b.WriteString(fmt.Sprintf("    [[ -n \"$%s\" ]] && __cmdx_bar_right+=\"%s%%F{#%s}$%s%%f\"\n", varName, sepR, color, varName))
	}

	switch sb.Position {
	case StatusBarAbove:
		b.WriteString("    print -P \"$__cmdx_bar\"\n}\nadd-zsh-hook precmd __cmdx_statusbar")
	case StatusBarBelow:
		b.WriteString("    RPROMPT=\"$__cmdx_bar_right\"\n    print -P \"$__cmdx_bar\"\n}\nadd-zsh-hook precmd __cmdx_statusbar")
	default:
		b.WriteString("    echo -n \"$__cmdx_bar \"\n}\nPROMPT='$(__cmdx_statusbar)'")
	}

	b.WriteString("\n")
	return b.String()
}

// statusBarPowerShell generates PowerShell prompt code for the status bar.
func statusBarPowerShell(sb *StatusBarConfig, colors map[string]string) string {
	left, _, _ := segmentsByZone(sb)
	sepL, _ := separatorGlyphs(sb.SeparatorStyle, sb.CustomSeparatorLeft, sb.CustomSeparatorRight)

	var b strings.Builder
	b.WriteString(`
# cmdX status bar
function __cmdx_statusbar {
    $exitCode = $LASTEXITCODE
    $bar = ""
`)

	for _, seg := range left {
		b.WriteString(psSegmentCode(seg))
		color := resolveSegColor(seg.Color, colors)
		label := seg.Label
		b.WriteString(fmt.Sprintf(`
    if ($__seg -ne "") {
        Write-Host -NoNewline "%s$__seg%s" -ForegroundColor %s
    }`, label, sepL, psColor(color)))
	}

	switch sb.Position {
	case StatusBarAbove:
		b.WriteString(`
    Write-Host ""  # newline after bar
}
$__cmdx_orig_prompt = $function:prompt
function prompt {
    __cmdx_statusbar
    & $__cmdx_orig_prompt
}`)
	default:
		b.WriteString(`
    Write-Host ""
}
function prompt {
    __cmdx_statusbar
    return " "
}`)
	}

	b.WriteString("\n")
	return b.String()
}

// zshSegmentCode returns the zsh code that computes a segment's value.
func zshSegmentCode(seg SegmentConfig) string {
	varName := bashSegVarName(seg, 0)
	switch seg.Type {
	case SegmentGit:
		return fmt.Sprintf(`
    local %s=""
    if git rev-parse --git-dir &>/dev/null; then
        local __branch=$(git branch --show-current 2>/dev/null || git rev-parse --short HEAD 2>/dev/null)
        local __dirty=""
        [[ -n $(git status --porcelain 2>/dev/null) ]] && __dirty="*"
        %s="%s${__branch}${__dirty}"
    fi`, varName, varName, seg.Label)
	case SegmentDirectory:
		return fmt.Sprintf(`    local %s="%s${PWD/#$HOME/~}"`, varName, seg.Label)
	case SegmentTime:
		format := seg.Format
		if format == "" {
			format = "%H:%M"
		}
		return fmt.Sprintf(`    local %s="%s$(date +%s)"`, varName, seg.Label, format)
	case SegmentUser:
		return fmt.Sprintf(`    local %s="%s$USER"`, varName, seg.Label)
	case SegmentVirtualEnv:
		return fmt.Sprintf(`    local %s=""
    [[ -n "$VIRTUAL_ENV" ]] && %s=" ($(basename $VIRTUAL_ENV))"`, varName, varName)
	case SegmentExitCode:
		return fmt.Sprintf(`    local %s=""
    [[ $? -ne 0 ]] && %s=" ✗$?"`, varName, varName)
	case SegmentText:
		return fmt.Sprintf(`    local %s="%s"`, varName, seg.Format)
	case SegmentEnvVar:
		return fmt.Sprintf(`    local %s="%s${%s}"`, varName, seg.Label, seg.EnvVar)
	default:
		return fmt.Sprintf(`    local %s=""`, varName)
	}
}

// psSegmentCode returns PowerShell code that computes a segment value into $__seg.
func psSegmentCode(seg SegmentConfig) string {
	switch seg.Type {
	case SegmentGit:
		return `
    $__seg = ""
    if (Test-Path .git) {
        $__branch = git branch --show-current 2>$null
        $__dirty = if ((git status --porcelain 2>$null) -ne "") { "*" } else { "" }
        $__seg = "$__branch$__dirty"
    }`
	case SegmentDirectory:
		return `
    $__seg = Split-Path -Leaf (Get-Location)`
	case SegmentTime:
		return `
    $__seg = Get-Date -Format "HH:mm"`
	case SegmentUser:
		return `
    $__seg = $env:USERNAME`
	case SegmentExitCode:
		return `
    $__seg = if ($exitCode -ne 0) { " ✗$exitCode" } else { "" }`
	case SegmentVirtualEnv:
		return `
    $__seg = if ($env:VIRTUAL_ENV) { " ($(Split-Path -Leaf $env:VIRTUAL_ENV))" } else { "" }`
	case SegmentText:
		return fmt.Sprintf(`
    $__seg = "%s"`, seg.Format)
	case SegmentEnvVar:
		return fmt.Sprintf(`
    $__seg = $env:%s`, seg.EnvVar)
	default:
		return `
    $__seg = ""`
	}
}

// bashSegVarName returns a unique bash variable name for a segment.
func bashSegVarName(seg SegmentConfig, idx int) string {
	return fmt.Sprintf("__seg_%s_%d", strings.ReplaceAll(string(seg.Type), "_", ""), idx)
}

// resolveSegColor resolves a segment color to a hex value, checking the
// theme's color palette if a semantic key is provided.
func resolveSegColor(color string, palette map[string]string) string {
	if color == "" {
		return ""
	}
	if strings.HasPrefix(color, "#") {
		return color
	}
	if v, ok := palette[color]; ok {
		return v
	}
	return color
}

// psColor converts a hex color to a PowerShell ConsoleColor name.
// This is a best-effort approximation — PowerShell's 16-color console
// doesn't support truecolor, so we map to the nearest named color.
func psColor(hex string) string {
	if hex == "" {
		return "White"
	}
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "White"
	}
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)

	// Very rough nearest-color mapping
	switch {
	case r > 200 && g < 100 && b < 100:
		return "Red"
	case r < 100 && g > 200 && b < 100:
		return "Green"
	case r < 100 && g < 100 && b > 200:
		return "Blue"
	case r > 200 && g > 200 && b < 100:
		return "Yellow"
	case r > 200 && g < 100 && b > 200:
		return "Magenta"
	case r < 100 && g > 200 && b > 200:
		return "Cyan"
	case r > 200 && g > 200 && b > 200:
		return "White"
	case r < 80 && g < 80 && b < 80:
		return "DarkGray"
	default:
		return "Gray"
	}
}
