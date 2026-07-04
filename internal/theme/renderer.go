package theme

import (
	"fmt"
	"strings"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/render"
	"github.com/charmbracelet/lipgloss"
)

// Renderer handles all visual output for a theme
type Renderer struct {
	Theme *config.Theme
}

// NewRenderer creates a renderer for the given theme
func NewRenderer(t *config.Theme) *Renderer {
	return &Renderer{Theme: t}
}

// RenderBanner prints the startup banner
func (r *Renderer) RenderBanner(username string) {
	if !r.Theme.Banner.Enabled {
		return
	}

	color := r.resolveColor(r.Theme.Banner.Color)
	text := strings.ReplaceAll(r.Theme.Banner.Text, "{user}", username)

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		PaddingTop(1).
		PaddingBottom(1)

	fmt.Println(style.Render(text))
}

// RenderPrompt builds and prints the prompt string
func (r *Renderer) RenderPrompt(data map[string]string) string {
	p := r.Theme.Prompt
	primary := r.resolveColor("primary")

	result := p.Format
	for key, val := range data {
		result = strings.ReplaceAll(result, "{"+key+"}", val)
	}

	for _, seg := range []string{"{user}", "{dir}", "{git}", "{time}", "{symbol}"} {
		result = strings.ReplaceAll(result, seg, "")
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(primary)).
		Bold(true)

	return style.Render(result)
}

// RenderProgressBar renders a styled progress bar
func (r *Renderer) RenderProgressBar(percent float64) string {
	bar := r.Theme.ProgressBar
	color := r.resolveColor(bar.Color)

	filled := int(float64(bar.Width) * percent / 100)
	if filled > bar.Width {
		filled = bar.Width
	}
	empty := bar.Width - filled

	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(r.Theme.Colors.Muted))

	return fmt.Sprintf("[%s%s] %.0f%%",
		filledStyle.Render(strings.Repeat(bar.Filled, filled)),
		emptyStyle.Render(strings.Repeat(bar.Empty, empty)),
		percent,
	)
}

// RenderBorder draws a box with the theme's border style
func (r *Renderer) RenderBorder(content string) string {
	color := r.resolveColor("primary")
	chars := r.Theme.Borders.Chars

	style := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(color)).
		Padding(0, 1)

	_ = chars
	return style.Render(content)
}

// RenderThemeInfo prints a rich markdown summary card of the theme via glamour.
func (r *Renderer) RenderThemeInfo() {
	t := r.Theme

	swatch := func(hex string) string {
		if hex == "" {
			return "—"
		}
		return lipgloss.NewStyle().
			Background(lipgloss.Color(hex)).
			Foreground(lipgloss.Color(hex)).
			Render("  ") + " `" + hex + "`"
	}

	loaderFrames := strings.Join(t.Loader.Frames, " ")
	if len(loaderFrames) > 40 {
		loaderFrames = loaderFrames[:40] + "…"
	}

	tags := buildTags(t)

	md := fmt.Sprintf(`# %s

> %s

| Field      | Value |
|------------|-------|
| Author     | %s |
| Version    | %s |
| Cursor     | %s |
| Prompt     | %s |
| Loader     | %s (%dms) |
| Progress   | %s / %s (%d wide) |
| Border     | %s |
| Effects    | banner: %s · prompt: %s |
| Gradient   | %s → %s (%s) |
| Wallpaper  | %s |
| Assets     | spinner: %s · banner: %s · divider: %s |

## Color Palette

| Role       | Color |
|------------|-------|
| Primary    | %s |
| Secondary  | %s |
| Accent     | %s |
| Background | %s |
| Foreground | %s |
| Error      | %s |
| Success    | %s |
| Warning    | %s |
| Muted      | %s |

%s
`,
		t.Meta.Name,
		t.Meta.Description,
		orDash(t.Meta.Author),
		orDash(t.Meta.Version),
		orDash(t.Cursor.Style),
		orDash(t.Prompt.Format),
		orDash(loaderFrames), t.Loader.IntervalMs,
		orDash(t.ProgressBar.Filled), orDash(t.ProgressBar.Empty), t.ProgressBar.Width,
		orDash(t.Borders.Style),
		orDash(t.Graphics.Effects.Banner), orDash(t.Graphics.Effects.Prompt),
		orDash(t.Graphics.Gradient.From), orDash(t.Graphics.Gradient.To), orDash(t.Graphics.Gradient.Direction),
		wallpaperStatus(t),
		orDash(t.Assets.Spinner), orDash(t.Assets.Banner), orDash(t.Assets.Divider),
		swatch(t.Colors.Primary),
		swatch(t.Colors.Secondary),
		swatch(t.Colors.Accent),
		swatch(t.Colors.Background),
		swatch(t.Colors.Foreground),
		swatch(t.Colors.Error),
		swatch(t.Colors.Success),
		swatch(t.Colors.Warning),
		swatch(t.Colors.Muted),
		tags,
	)

	render.Markdown(md)
}

func buildTags(t *config.Theme) string {
	var tags []string
	if t.Graphics.Gradient.Enabled {
		tags = append(tags, "`gradient`")
	}
	if t.Banner.Enabled {
		tags = append(tags, "`banner`")
	}
	if t.Wallpaper.Enabled {
		tags = append(tags, "`wallpaper`")
	}
	if t.Icons.Enabled {
		tags = append(tags, "`icons`")
	}
	if t.Assets.Spinner != "" {
		tags = append(tags, "`png-spinner`")
	}
	if len(tags) == 0 {
		return ""
	}
	return "**Features:** " + strings.Join(tags, " · ")
}

func wallpaperStatus(t *config.Theme) string {
	if !t.Wallpaper.Enabled {
		return "disabled"
	}
	return fmt.Sprintf("enabled (opacity %.0f%%)", t.Wallpaper.Opacity*100)
}

func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

// resolveColor maps a color key to its hex value
func (r *Renderer) resolveColor(key string) string {
	return config.ResolveColor(r.Theme.Colors, key)
}
