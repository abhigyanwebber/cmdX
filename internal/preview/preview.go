// Package preview renders a live, animated walkthrough of a theme's visual
// elements directly in the terminal — colors, banner, prompt, loader,
// progress bar, effects, patterns, and borders — so a user can see what a
// theme looks like before applying it.
package preview

import (
	"fmt"
	"time"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/graphics"
	"github.com/abhigyanwebber/cmd-customizer/internal/primitives"
	"github.com/charmbracelet/lipgloss"
)

// Preview drives a step-by-step animated rendering of a theme.
type Preview struct {
	theme *config.Theme
}

// NewPreview creates a Preview for the given theme.
func NewPreview(t *config.Theme) *Preview {
	return &Preview{theme: t}
}

// Run walks through every visual element of the theme in sequence,
// printing each section with a short pause in between so changes are
// easy to follow.
func (p *Preview) Run() {
	t := p.theme
	primary := p.resolve("primary")
	accent := p.resolve("accent")
	muted := p.resolve("muted")
	success := p.resolve("success")
	errorCol := p.resolve("error")

	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(primary)).Bold(true).Underline(true)
	sectionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(accent)).Bold(true)
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(muted))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(success))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(errorCol))

	// ── Header ────────────────────────────────────────────
	fmt.Println()
	fmt.Println(titleStyle.Render(fmt.Sprintf("  Preview: %s  ", t.Meta.Name)))
	fmt.Println(mutedStyle.Render(fmt.Sprintf("  %s", t.Meta.Description)))
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// ── Gradient Title ────────────────────────────────────
	if t.Graphics.Gradient.Enabled {
		fmt.Println(sectionStyle.Render("[ Gradient ]"))
		grad, err := graphics.GradientText("  "+t.Meta.Name+" — "+t.Meta.Description, t.Graphics.Gradient.From, t.Graphics.Gradient.To)
		if err == nil {
			fmt.Println(grad)
		}
		fmt.Println()
		time.Sleep(300 * time.Millisecond)
	}

	// ── Banner ────────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Banner ]"))
	p.renderBannerWithEffect()
	fmt.Println()
	time.Sleep(500 * time.Millisecond)

	// ── Divider styles ────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Dividers ]"))
	p.renderDividers()
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// ── Colors ────────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Color Palette ]"))
	p.renderColorPalette()
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// ── Prompt ────────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Prompt Style ]"))
	p.renderPromptPreview()
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// ── Spinner ───────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Loader / Spinner ]"))
	spinner := primitives.NewSpinner(t)
	spinner.Once("Installing dependencies", 2*time.Second)
	time.Sleep(200 * time.Millisecond)

	// ── Progress Bar ──────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Progress Bar ]"))
	bar := primitives.NewProgressBar(t)
	bar.Animate("Building project", 2*time.Second)
	fmt.Println()
	time.Sleep(200 * time.Millisecond)

	// ── Effects ───────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Effects ]"))
	p.renderEffects()
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// ── Patterns ──────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Patterns ]"))
	p.renderPatterns()
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// ── Status Messages ───────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Status Messages ]"))
	fmt.Println(successStyle.Render("  ✓ Build successful"))
	fmt.Println(successStyle.Render("  ✓ Tests passed (42/42)"))
	fmt.Println(errorStyle.Render("  ✗ Connection refused on port 8080"))
	fmt.Println(mutedStyle.Render("  ─ Skipping optional dependency"))
	fmt.Println()
	time.Sleep(300 * time.Millisecond)

	// ── Border ────────────────────────────────────────────
	fmt.Println(sectionStyle.Render("[ Border Style ]"))
	p.renderBorderPreview()
	fmt.Println()

	// ── Footer ────────────────────────────────────────────
	fmt.Println(mutedStyle.Render(fmt.Sprintf(
		"  Run 'cmdx theme apply %s' to use this theme\n", t.Meta.Name,
	)))
}
