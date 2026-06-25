package primitives

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/abhigyanwebber/cmd-customizer/internal/config"
)

// ProgressBar renders a themed progress bar
type ProgressBar struct {
	width   int
	filled  string
	empty   string
	color   string
	muted   string
}

// NewProgressBar creates a progress bar from the active theme
func NewProgressBar(t *config.Theme) *ProgressBar {
	return &ProgressBar{
		width:  t.ProgressBar.Width,
		filled: t.ProgressBar.Filled,
		empty:  t.ProgressBar.Empty,
		color:  resolveColor(t, t.ProgressBar.Color),
		muted:  t.Colors.Muted,
	}
}

// Render returns a progress bar string at the given percentage (0-100)
func (p *ProgressBar) Render(percent float64) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int(float64(p.width) * percent / 100)
	empty := p.width - filled

	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(p.color))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(p.muted))
	percentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(p.color)).Bold(true)

	bar := fmt.Sprintf("[%s%s] %s",
		filledStyle.Render(strings.Repeat(p.filled, filled)),
		emptyStyle.Render(strings.Repeat(p.empty, empty)),
		percentStyle.Render(fmt.Sprintf("%.0f%%", percent)),
	)
	return bar
}

// Animate simulates a progress bar filling up — useful for demos and previews
func (p *ProgressBar) Animate(label string, duration time.Duration) {
	steps := 40
	delay := duration / time.Duration(steps)

	for i := 0; i <= steps; i++ {
		percent := float64(i) / float64(steps) * 100
		fmt.Printf("\r%s  %s", label, p.Render(percent))
		time.Sleep(delay)
	}
	fmt.Println()
}