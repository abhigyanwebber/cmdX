package preview

import (
	"fmt"
	"os/user"
	"strings"
	"time"

	"github.com/abhigyanwebber/cmd-customizer/internal/graphics"
	"github.com/abhigyanwebber/cmd-customizer/internal/primitives"
	"github.com/charmbracelet/lipgloss"
)

// renderBannerWithEffect renders the theme's banner using whichever
// graphics effect (glitch, rainbow, neon, typewriter) the theme specifies,
// falling back to the plain primitives.Banner renderer otherwise.
func (p *Preview) renderBannerWithEffect() {
	t := p.theme

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

// renderDividers prints one line of every available divider style so the
// user can compare them side by side.
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

// renderColorPalette prints a swatch + hex label for every named color
// in the theme.
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

// renderPromptPreview shows what the theme's configured prompt looks like
// against a sample working directory and git branch.
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

// renderEffects demonstrates the graphics effects engine (rainbow, glitch,
// neon, gradient text) against sample strings, independent of what the
// theme's own banner/prompt effect settings are.
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

// renderPatterns prints a sample of each available background pattern.
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
		for _, ch := range rendered {
			if ch == '\n' {
				fmt.Print("  ")
			}
			fmt.Print(string(ch))
		}
		fmt.Println()
	}
}

// renderBorderPreview shows a sample bordered box using the theme's
// primary color, containing the theme's own metadata.
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
