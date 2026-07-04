// Package shells defines the cross-shell injection interface (Shell) and
// shared safety helpers used by every concrete shell implementation
// (bash, zsh, powershell) to write and remove cmdX theme configuration
// from the user's profile script.
package shells

import "github.com/abhigyanwebber/cmd-customizer/internal/config"

// Shell defines the contract every shell implementation must follow
type Shell interface {
	// Name returns the shell identifier e.g. "powershell", "zsh", "bash"
	Name() string

	// Detect checks if this shell is the active shell on the system
	Detect() bool

	// ProfilePath returns the path to the shell's config file
	ProfilePath() (string, error)

	// Inject writes the theme into the shell's config file
	Inject(t *config.Theme) error

	// Remove cleans up any cmdx injected code from the config file
	Remove() error

	// IsInjected checks if cmdx code is already present in the config
	IsInjected() (bool, error)
}

// Markers used to wrap injected code so we can find and remove it later
const (
	InjectStart = "# >>> cmdx theme start >>>"
	InjectEnd   = "# <<< cmdx theme end <<<"
)