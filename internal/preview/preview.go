package preview

import (
	"fmt"
	"strings"
	"time"

	"os/user"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/graphics"
	"github.com/abhigyanwebber/cmd-customizer/internal/primitives"
	"github.com/charmbracelet/lipgloss"
)

type Preview struct {
	theme *config.Theme
}

func NewPreview(t *config.Theme) *Preview {
	return &Preview{theme: t}
}

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

func (p *Preview) renderBannerWithEffect() {
	t := p.theme

	// at the top of renderBannerWithEffect
	u, _ := user.Current()
	username := "developer"
	if u != nil {
		username = u.Username
	}
	bannerText := strings.ReplaceAll(t.Banner.Text, "{user}", username)

	switch t.Graphics.Effects.Banner {
	case "glitch":
		frames := graphics.GlitchFrames("  "+bannerText, 6)
		for _, frame := range frames {
			color, _ := graphics.ParseHex(p.resolve(t.Banner.Color))
			fmt.Printf("\r%s", graphics.AnsiColor(graphics.ToRGB(color), frame))
			time.Sleep(80 * time.Millisecond)
		}
		fmt.Println()

	case "rainbow":
		fmt.Println(graphics.RainbowText("  " + bannerText))

	case "neon":
		neon, err := graphics.NeonText("  "+bannerText, p.resolve("primary"))
		if err == nil {
			fmt.Println(neon)
		}

	case "typewriter":
		graphics.TypewriterPrint("  "+bannerText, 40)

	default:
		banner := primitives.NewBanner(t)
		banner.Render()
	}
}

func (p *Preview) renderDividers() {
	t := p.theme
	from := t.Colors.Primary
	to := t.Colors.Secondary
	width := 50

	styles := []graphics.DividerStyle{
		graphics.StyleLine,
		graphics.StyleWave,
		graphics.StyleDots,
		graphics.StyleZigzag,
		graphics.StyleStars,
		graphics.StyleDouble,
	}

	names := []string{"line", "wave", "dots", "zigzag", "stars", "double"}

	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Colors.Muted))

	for i, style := range styles {
		div, err := graphics.GradientDivider(style, width, from, to)
		if err != nil {
			div = graphics.Divider(style, width, from)
		}
		fmt.Printf("  %s  %s\n", div, mutedStyle.Render(names[i]))
	}
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

func (p *Preview) renderEffects() {
	t := p.theme
	primary := p.resolve("primary")

	// Rainbow
	fmt.Println(graphics.RainbowText("  The quick brown fox jumps over the lazy dog"))

	// Glitch
	glitched := graphics.GlitchText("  SYSTEM ERROR: unexpected token at line 42", 0.2)
	color, _ := graphics.ParseHex(t.Colors.Error)
	fmt.Println(graphics.AnsiColor(graphics.ToRGB(color), glitched))

	// Neon
	neon, err := graphics.NeonText("  NEON GLOW ACTIVE", primary)
	if err == nil {
		fmt.Println(neon)
	}

	// Gradient text
	grad, err := graphics.GradientText("  Gradient text rendering engine", t.Colors.Primary, t.Colors.Secondary)
	if err == nil {
		fmt.Println(grad)
	}
}

func (p *Preview) renderPatterns() {
	t := p.theme
	col := t.Colors.Muted

	patterns := []struct {
		name    string
		pattern graphics.Pattern
	}{
		{"dots", graphics.PatternDots},
		{"circuit", graphics.PatternCircuit},
		{"stars", graphics.PatternStars},
	}

	for _, pat := range patterns {
		rendered := graphics.RenderPattern(pat.pattern, 40, 2, col)
		label := lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Muted)).
			Render("  " + pat.name)
		fmt.Println(label)
		lines := rendered
		for _, ch := range lines {
			if ch == '\n' {
				fmt.Print("  ")
			}
			fmt.Print(string(ch))
		}
		fmt.Println()
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
