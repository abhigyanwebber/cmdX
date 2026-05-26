package preview

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/primitives"
)

// Preview runs a full live demonstration of a theme
type Preview struct {
	theme *config.Theme
}

// NewPreview creates a preview runner for the given theme
func NewPreview(t *config.Theme) *Preview {
	return &Preview{theme: t}
}

// Run executes the full theme preview sequence
func (p *Preview) Run() {
	t := p.theme
	primary := p.resolve("primary")
	accent := p.resolve("accent")
	muted := p.resolve("muted")
	success := p.resolve("success")
	errorCol := p.resolve("error")

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(primary)).
		Bold(true).
		Underline(true)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(accent)).
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(muted))

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(success))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(errorCol))

	// ── Header ────────────────────────────────────────────
	fmt.Println()
	fmt.Println(titleStyle.Render(fmt.Sprintf("  Preview: %s  ", t.Meta.Name)))
	fmt.Println(mutedStyle.Render(fmt.Sprintf("  %s", t.Meta.Description)))
	fmt.Println()
	time.Sleep(400 * time.Millisecond)

	// ── Banner ────────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Banner ]"))
	banner := primitives.NewBanner(t)
	banner.Render()
	time.Sleep(600 * time.Millisecond)

	// ── Colors ────────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Color Palette ]"))
	p.renderColorPalette()
	fmt.Println()
	time.Sleep(400 * time.Millisecond)

	// ── Prompt ────────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Prompt Style ]"))
	p.renderPromptPreview()
	fmt.Println()
	time.Sleep(400 * time.Millisecond)

	// ── Spinner ───────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Loader / Spinner ]"))
	spinner := primitives.NewSpinner(t)
	spinner.Once("Installing dependencies", 2*time.Second)
	time.Sleep(300 * time.Millisecond)

	// ── Progress Bar ──────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Progress Bar ]"))
	bar := primitives.NewProgressBar(t)
	bar.Animate("Building project", 2*time.Second)
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// ── Status Messages ───────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Status Messages ]"))
	fmt.Println(successStyle.Render("  ✓ Build successful"))
	fmt.Println(successStyle.Render("  ✓ Tests passed (42/42)"))
	fmt.Println(errorStyle.Render("  ✗ Connection refused on port 8080"))
	fmt.Println(mutedStyle.Render("  ─ Skipping optional dependency"))
	fmt.Println()
	time.Sleep(400 * time.Millisecond)

	// ── Borders ───────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Border Style ]"))
	p.renderBorderPreview()
	fmt.Println()

	// ── Footer ────────────────────────────────────────────
	fmt.Println(mutedStyle.Render(fmt.Sprintf(
		"  Run 'cmdx theme apply %s' to use this theme\n", t.Meta.Name,
	)))
}

func (p *Preview) renderColorPalette() {
	t := p.theme
	colors := []struct {
		name string
		hex  string
	}{
		{"primary", t.Colors.Primary},
		{"secondary", t.Colors.Secondary},
		{"accent", t.Colors.Accent},
		{"success", t.Colors.Success},
		{"error", t.Colors.Error},
		{"warning", t.Colors.Warning},
		{"muted", t.Colors.Muted},
	}

	for _, c := range colors {
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(c.hex)).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Render("  ")

		label := lipgloss.NewStyle().
			Foreground(lipgloss.Color(c.hex)).
			Render(fmt.Sprintf(" %-12s %s", c.name, c.hex))

		fmt.Printf("  %s%s\n", swatch, label)
	}
}

func (p *Preview) renderPromptPreview() {
	t := p.theme
	primary := p.resolve("primary")
	accent := p.resolve("accent")
	muted := p.resolve("muted")

	dirStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(accent)).Bold(true)
	gitStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(muted))
	symbolStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(primary)).Bold(true)

	dir := dirStyle.Render("~/projects/cmd-customizer")
	git := gitStyle.Render("(main)")
	symbol := symbolStyle.Render(t.Prompt.Symbol)

	if t.Prompt.Style == "multiline" {
		fmt.Printf("  %s %s\n  %s \n", dir, git, symbol)
	} else {
		fmt.Printf("  %s %s %s \n", dir, git, symbol)
	}
}

func (p *Preview) renderBorderPreview() {
	t := p.theme
	primary := p.resolve("primary")
	muted := p.resolve("muted")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(primary)).
		Foreground(lipgloss.Color(muted)).
		Padding(0, 2).
		Width(40)

	fmt.Println(boxStyle.Render(
		fmt.Sprintf("Theme: %s\nAuthor: %s\nVersion: %s",
			t.Meta.Name,
			t.Meta.Author,
			t.Meta.Version,
		),
	))
}

func (p *Preview) resolve(key string) string {
	c := p.theme.Colors
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
	return key
}