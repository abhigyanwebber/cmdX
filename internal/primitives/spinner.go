package primitives

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/abhigyanwebber/cmd-customizer/internal/config"
)

// Spinner renders animated loading indicators
type Spinner struct {
	frames   []string
	interval time.Duration
	color    string
	index    int
	stop     chan struct{}
}

// NewSpinner creates a spinner from the active theme
func NewSpinner(t *config.Theme) *Spinner {
	return &Spinner{
		frames:   t.Loader.Frames,
		interval: time.Duration(t.Loader.IntervalMs) * time.Millisecond,
		color:    resolveColor(t, t.Loader.Color),
		stop:     make(chan struct{}),
	}
}

// Start begins spinning with a message, runs until Stop() is called
func (s *Spinner) Start(message string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(s.color))

	go func() {
		for {
			select {
			case <-s.stop:
				fmt.Printf("\r%s\n", "                              ")
				return
			default:
				frame := s.frames[s.index%len(s.frames)]
				fmt.Printf("\r%s %s", style.Render(frame), message)
				s.index++
				time.Sleep(s.interval)
			}
		}
	}()
}

// Stop halts the spinner and prints a completion message
func (s *Spinner) Stop(finalMessage string) {
	s.stop <- struct{}{}
	time.Sleep(50 * time.Millisecond)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(s.color))
	fmt.Printf("\r%s %s\n", style.Render("✓"), finalMessage)
}

// StopWithError halts the spinner and prints an error message
func (s *Spinner) StopWithError(errMessage string) {
	s.stop <- struct{}{}
	time.Sleep(50 * time.Millisecond)
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444"))
	fmt.Printf("\r%s %s\n", errStyle.Render("✗"), errMessage)
}

// Once runs a spinner for a fixed duration — useful for demos
func (s *Spinner) Once(message string, duration time.Duration) {
	s.Start(message)
	time.Sleep(duration)
	s.Stop(message + " done")
}