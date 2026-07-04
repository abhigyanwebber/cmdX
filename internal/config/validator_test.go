package config

import (
	"strings"
	"testing"
)

func validTheme() *Theme {
	return &Theme{
		Meta: Meta{Name: "test", Version: "1.0.0"},
		Colors: Colors{
			Primary:    "#FF00FF",
			Secondary:  "#00FFFF",
			Background: "#0D0D0D",
			Foreground: "#FFFFFF",
			Accent:     "#FFD700",
			Error:      "#FF4444",
			Success:    "#00FF88",
			Warning:    "#FFA500",
			Muted:      "#444444",
		},
		Prompt: Prompt{
			Symbol: "▶",
			Format: "{user}@{dir}",
			Style:  "single",
		},
		Loader: Loader{
			Frames:     []string{"◐", "◓", "◑", "◒"},
			IntervalMs: 100,
		},
		ProgressBar: ProgressBar{
			Filled: "█",
			Empty:  "░",
			Width:  20,
		},
		Cursor: Cursor{Style: "block"},
	}
}

func TestValidateTheme_Valid(t *testing.T) {
	if err := ValidateTheme(validTheme()); err != nil {
		t.Fatalf("expected valid theme to pass, got error: %v", err)
	}
}

func TestValidateMeta_MissingName(t *testing.T) {
	th := validTheme()
	th.Meta.Name = ""
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for missing name, got nil")
	} else if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("expected 'name is required' error, got: %v", err)
	}
}

func TestValidateMeta_MissingVersion(t *testing.T) {
	th := validTheme()
	th.Meta.Version = ""
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
}

func TestValidateColors_MissingField(t *testing.T) {
	th := validTheme()
	th.Colors.Primary = ""
	err := ValidateTheme(th)
	if err == nil {
		t.Fatal("expected error for missing primary color, got nil")
	}
	if !strings.Contains(err.Error(), "primary color is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidateColors_InvalidHex(t *testing.T) {
	cases := []string{"red", "#GG0000", "#12345", "FF00FF", "#1234567"}
	for _, hex := range cases {
		th := validTheme()
		th.Colors.Primary = hex
		if err := ValidateTheme(th); err == nil {
			t.Errorf("expected error for invalid hex %q, got nil", hex)
		}
	}
}

func TestValidateColors_ValidHexFormats(t *testing.T) {
	cases := []string{"#FF00FF", "#fff", "#000000", "#AbC123"}
	for _, hex := range cases {
		th := validTheme()
		th.Colors.Primary = hex
		if err := ValidateTheme(th); err != nil {
			t.Errorf("expected hex %q to be valid, got error: %v", hex, err)
		}
	}
}

func TestValidatePrompt_MissingSymbol(t *testing.T) {
	th := validTheme()
	th.Prompt.Symbol = ""
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for missing prompt symbol")
	}
}

func TestValidatePrompt_MissingFormat(t *testing.T) {
	th := validTheme()
	th.Prompt.Format = ""
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for missing prompt format")
	}
}

func TestValidatePrompt_InvalidStyle(t *testing.T) {
	th := validTheme()
	th.Prompt.Style = "fancy"
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for invalid prompt style")
	}
}

func TestValidatePrompt_ValidStyles(t *testing.T) {
	for _, style := range []string{"single", "multiline"} {
		th := validTheme()
		th.Prompt.Style = style
		if err := ValidateTheme(th); err != nil {
			t.Errorf("expected style %q to be valid, got error: %v", style, err)
		}
	}
}

func TestValidateLoader_EmptyFrames(t *testing.T) {
	th := validTheme()
	th.Loader.Frames = []string{}
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for empty loader frames")
	}
}

func TestValidateLoader_ZeroInterval(t *testing.T) {
	th := validTheme()
	th.Loader.IntervalMs = 0
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for zero interval_ms")
	}
}

func TestValidateLoader_NegativeInterval(t *testing.T) {
	th := validTheme()
	th.Loader.IntervalMs = -50
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for negative interval_ms")
	}
}

func TestValidateProgressBar_MissingFilled(t *testing.T) {
	th := validTheme()
	th.ProgressBar.Filled = ""
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for missing filled char")
	}
}

func TestValidateProgressBar_MissingEmpty(t *testing.T) {
	th := validTheme()
	th.ProgressBar.Empty = ""
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for missing empty char")
	}
}

func TestValidateProgressBar_ZeroWidth(t *testing.T) {
	th := validTheme()
	th.ProgressBar.Width = 0
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for zero width")
	}
}

func TestValidateCursor_InvalidStyle(t *testing.T) {
	th := validTheme()
	th.Cursor.Style = "wiggle"
	if err := ValidateTheme(th); err == nil {
		t.Fatal("expected error for invalid cursor style")
	}
}

func TestValidateCursor_ValidStyles(t *testing.T) {
	for _, style := range []string{"block", "bar", "underline"} {
		th := validTheme()
		th.Cursor.Style = style
		if err := ValidateTheme(th); err != nil {
			t.Errorf("expected cursor style %q to be valid, got error: %v", style, err)
		}
	}
}

func TestValidateTheme_ErrorWrapping(t *testing.T) {
	th := validTheme()
	th.Meta.Name = ""
	err := ValidateTheme(th)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.HasPrefix(err.Error(), "meta:") {
		t.Errorf("expected error to be wrapped with 'meta:' prefix, got: %v", err)
	}
}
