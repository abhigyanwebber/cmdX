package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerModel wraps bubbles/spinner for loading states in CLI operations.
type SpinnerModel struct {
	spinner  spinner.Model
	message  string
	done     bool
	err      error
	primary  string
	task     func() error
	resultCh chan error
}

type spinnerDoneMsg struct{ err error }

// NewSpinner creates a spinner with a message and a background task to run.
func NewSpinner(message, primary string, task func() error) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(primary))

	return SpinnerModel{
		spinner:  s,
		message:  message,
		primary:  primary,
		task:     task,
		resultCh: make(chan error, 1),
	}
}

// Init starts the spinner animation and kicks off the background task
// concurrently. Satisfies tea.Model.
func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			err := m.task()
			return spinnerDoneMsg{err: err}
		},
	)
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinnerDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m SpinnerModel) View() string {
	if m.done {
		return ""
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(m.primary))
	return fmt.Sprintf("\n  %s %s\n", m.spinner.View(), style.Render(m.message))
}

// Err returns the error from the background task, if any.
func (m SpinnerModel) Err() error { return m.err }

// RunSpinner runs a task with a spinner and returns any error from the task.
func RunSpinner(message, primary string, task func() error) error {
	m := NewSpinner(message, primary, task)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return err
	}
	return result.(SpinnerModel).Err()
}
