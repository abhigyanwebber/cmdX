package shells

import "strings"

// dangerousShellChars are characters that have special meaning in
// bash/zsh/PowerShell when they appear inside a double-quoted string or
// an unquoted command position — backticks and $() trigger command
// substitution, semicolons and pipes chain commands, and quotes/escapes
// can break out of the surrounding string literal entirely.
var dangerousShellChars = []string{
	"`", "$", "\"", "'", "\\", ";", "|", "&", "\n", "\r",
}

// SanitizeForShell strips characters that could break out of the quoted
// string context a theme field is embedded into, or trigger command
// substitution / command chaining, when that field's value is written
// into a shell profile script.
//
// Theme JSON files are untrusted input — they may be downloaded from the
// community registry or shared by other users — so every text field that
// gets embedded into an injected shell block (banner text, theme name,
// author, prompt symbol, etc.) must pass through this function first.
// Without this, a malicious theme's Banner.Text field could inject
// arbitrary shell commands that execute on every new terminal session.
func SanitizeForShell(s string) string {
	for _, ch := range dangerousShellChars {
		s = strings.ReplaceAll(s, ch, "")
	}
	// also strip the injection markers themselves, so a malicious theme
	// can't forge a fake end-marker and break out of its own block
	s = strings.ReplaceAll(s, InjectStart, "")
	s = strings.ReplaceAll(s, InjectEnd, "")
	return s
}
