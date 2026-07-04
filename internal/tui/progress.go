package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const progressWidth = 50

// ProgressModel shows a progress bar for download/fetch operations.
type ProgressModel struct {
	progress progress.Model
	percent  float64
	label    string
	done     bool
	primary  string
	accent   string
}

type progressTickMsg float64
type progressDoneMsg struct{}

// NewProgressModel creates a download progress bar.
func NewProgressModel(label, primary, accent string) ProgressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(progressWidth),
		progress.WithoutPercentage(),
	)
	p.FullColor = primary
	p.EmptyColor = "#333333"

	return ProgressModel{
		progress: p,
		label:    label,
		primary:  primary,
		accent:   accent,
	}
}

// Init, Update, and View satisfy tea.Model.
func (m ProgressModel) Init() tea.Cmd { return nil }

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressTickMsg:
		m.percent = float64(msg)
		if m.percent >= 1.0 {
			m.done = true
			return m, tea.Sequence(
				m.progress.SetPercent(1.0),
				tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
					return progressDoneMsg{}
				}),
			)
		}
		cmd := m.progress.SetPercent(m.percent)
		return m, cmd
	case progressDoneMsg:
		return m, tea.Quit
	case progress.FrameMsg:
		pm, cmd := m.progress.Update(msg)
		m.progress = pm.(progress.Model)
		return m, cmd
	}
	return m, nil
}

func (m ProgressModel) View() string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.primary))
	pctStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.accent))

	label := labelStyle.Render(m.label)
	pct := pctStyle.Render(fmt.Sprintf("%.0f%%", m.percent*100))

	return fmt.Sprintf("\n  %s\n  %s %s\n",
		label,
		m.progress.View(),
		pct,
	)
}

// SendProgress sends a progress update (0.0–1.0) to the running program.
// Call this from your download goroutine via p.Send(tui.progressTickMsg(0.5)).
func ProgressCmd(pct float64) tea.Msg {
	return progressTickMsg(pct)
}
