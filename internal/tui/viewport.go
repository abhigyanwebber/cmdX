package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewportModel wraps bubbles/viewport for scrollable content panels.
type ViewportModel struct {
	viewport viewport.Model
	ready    bool
	primary  string
	accent   string
	title    string
	content  string
}

// NewViewport creates a scrollable viewport with a title and content.
func NewViewport(title, content, primary, accent string, width, height int) ViewportModel {
	vp := viewport.New(width, height)
	vp.SetContent(content)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(primary))

	return ViewportModel{
		viewport: vp,
		ready:    true,
		primary:  primary,
		accent:   accent,
		title:    title,
		content:  content,
	}
}

// Init, Update, and View satisfy tea.Model.
func (m ViewportModel) Init() tea.Cmd { return nil }

func (m ViewportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 6
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ViewportModel) View() string {
	if !m.ready {
		return "Loading..."
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.accent)).
		Bold(true).
		Padding(0, 1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true)

	scrollPct := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.primary)).
		Render(fmt.Sprintf("%.0f%%", m.viewport.ScrollPercent()*100))

	title := titleStyle.Render(m.title) + "  " + scrollPct
	help := helpStyle.Render("[↑↓/PgUp/PgDn] Scroll  [Q] Close")

	return title + "\n" + m.viewport.View() + "\n" + help
}

// SetContent updates the viewport content.
func (m *ViewportModel) SetContent(content string) {
	m.content = content
	m.viewport.SetContent(content)
}
