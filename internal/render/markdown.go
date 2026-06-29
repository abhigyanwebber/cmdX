package render

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/x/term"
)

// termWidth returns the current terminal width, falling back to 80.
func termWidth() int {
	w, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || w < 40 {
		return 80
	}
	if w > 120 {
		return 120
	}
	return w
}

// Markdown renders a markdown string to the terminal using glamour.
// Auto-selects dark/light style based on environment and clamps to terminal width.
func Markdown(md string) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(termWidth()),
	)
	if err != nil {
		fmt.Println(md)
		return
	}

	out, err := r.Render(md)
	if err != nil {
		fmt.Println(md)
		return
	}

	fmt.Print(out)
}

// MarkdownWithStyle renders markdown using an explicit glamour style name
// ("dark", "light", "dracula", "tokyo-night", etc.).
func MarkdownWithStyle(md, style string) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath(style),
		glamour.WithWordWrap(termWidth()),
	)
	if err != nil {
		fmt.Println(md)
		return
	}

	out, err := r.Render(md)
	if err != nil {
		fmt.Println(md)
		return
	}

	fmt.Print(out)
}
