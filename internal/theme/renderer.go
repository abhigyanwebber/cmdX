package theme

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/abhigyanwebber/cmd-customizer/internal/config"
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

	// replace any unreplaced tokens with empty string
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

	_ = chars // lipgloss handles border chars internally for now
	return style.Render(content)
}

// RenderThemeInfo prints a summary card of the theme
func (r *Renderer) RenderThemeInfo() {
	t := r.Theme
	primary := r.resolveColor("primary")
	accent := r.resolveColor("accent")
	muted := r.resolveColor("muted")

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(primary)).
		Bold(true).
		Underline(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(muted))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(accent))

	fmt.Println(titleStyle.Render("  " + t.Meta.Name + "  "))
	fmt.Println(labelStyle.Render("Author:      ") + valueStyle.Render(t.Meta.Author))
	fmt.Println(labelStyle.Render("Version:     ") + valueStyle.Render(t.Meta.Version))
	fmt.Println(labelStyle.Render("Description: ") + valueStyle.Render(t.Meta.Description))
	fmt.Println(labelStyle.Render("Cursor:      ") + valueStyle.Render(t.Cursor.Style))
	fmt.Println(labelStyle.Render("Prompt:      ") + valueStyle.Render(t.Prompt.Format))
}

// resolveColor maps a color key to its hex value
func (r *Renderer) resolveColor(key string) string {
	c := r.Theme.Colors
	colorMap := map[string]string{
		"primary":    c.Primary,
		"secondary":  c.Secondary,
		"background": c.Background,
		"foreground": c.Foreground,
		"accent":     c.Accent,
		"error":      c.Error,
		"success":    c.Success,
		"warning":    c.Warning,
		"muted":      c.Muted,
	}
	if val, ok := colorMap[key]; ok {
		return val
	}
	return key // assume it's already a hex value
}