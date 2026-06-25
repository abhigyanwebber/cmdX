package primitives

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/abhigyanwebber/cmd-customizer/internal/config"
)

// Banner renders the startup banner for a theme
type Banner struct {
	theme *config.Theme
}

// NewBanner creates a banner renderer from the active theme
func NewBanner(t *config.Theme) *Banner {
	return &Banner{theme: t}
}

// Render prints the themed startup banner
func (b *Banner) Render() {
	if !b.theme.Banner.Enabled {
		return
	}

	username := getUsername()
	color := resolveColor(b.theme, b.theme.Banner.Color)
	text := strings.ReplaceAll(b.theme.Banner.Text, "{user}", username)

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		PaddingTop(1).
		PaddingBottom(1).
		PaddingLeft(2)

	borderColor := resolveColor(b.theme, "primary")
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 2)

	fmt.Println(style.Render(boxStyle.Render(text)))
}

// RenderWithText prints the banner with custom text override
func (b *Banner) RenderWithText(text string) {
	color := resolveColor(b.theme, b.theme.Banner.Color)
	borderColor := resolveColor(b.theme, "primary")

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 2)

	fmt.Println(textStyle.Render(boxStyle.Render(text)))
}

func getUsername() string {
	u, err := user.Current()
	if err != nil {
		return "developer"
	}
	return u.Username
}