package tui

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsInteractive reports whether both stdin and stdout are attached to a
// real terminal. RunSpinner and RunThemeList use this to decide whether
// launching a Bubble Tea program is safe: tea.NewProgram's Run() blocks
// on its internal terminal/input initialization when stdin isn't a real
// TTY (piped, redirected, or a non-interactive automation context), and
// — critically — it does this by hanging rather than returning an error,
// so callers cannot simply check p.Run()'s error to detect this case.
// Checking before ever calling tea.NewProgram avoids the hang entirely
// and costs nothing when a real terminal is present, since this function
// returns true immediately in that case.
func IsInteractive() bool {
	return isTerminal(os.Stdin.Fd()) && isTerminal(os.Stdout.Fd())
}

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}
