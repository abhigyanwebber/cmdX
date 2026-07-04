package graphics

import (
	"strings"
	"testing"
)

func TestParseHex_Valid(t *testing.T) {
	c, err := ParseHex("#FF00FF")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	r, g, b := c.RGB255()
	if r != 255 || g != 0 || b != 255 {
		t.Errorf("expected RGB(255,0,255), got (%d,%d,%d)", r, g, b)
	}
}

func TestParseHex_Invalid(t *testing.T) {
	_, err := ParseHex("not-a-color")
	if err == nil {
		t.Fatal("expected error for invalid hex, got nil")
	}
}

func TestParseHex_ShortForm(t *testing.T) {
	_, err := ParseHex("#f0f")
	if err != nil {
		t.Fatalf("expected short hex form to parse, got error: %v", err)
	}
}

func TestToRGB_Conversion(t *testing.T) {
	c, _ := ParseHex("#00FF00")
	rgb := ToRGB(c)
	if rgb.R != 0 || rgb.G != 255 || rgb.B != 0 {
		t.Errorf("expected RGB(0,255,0), got (%d,%d,%d)", rgb.R, rgb.G, rgb.B)
	}
}

func TestAnsiColor_ContainsEscapeCode(t *testing.T) {
	out := AnsiColor(RGB{255, 0, 0}, "hello")
	if !strings.Contains(out, "\033[38;2;255;0;0m") {
		t.Errorf("expected ANSI foreground escape code, got: %q", out)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected text to be preserved, got: %q", out)
	}
	if !strings.HasSuffix(out, "\033[0m") {
		t.Errorf("expected reset code at end, got: %q", out)
	}
}

func TestAnsiBgColor_ContainsEscapeCode(t *testing.T) {
	out := AnsiBgColor(RGB{10, 20, 30}, "x")
	if !strings.Contains(out, "\033[48;2;10;20;30m") {
		t.Errorf("expected ANSI background escape code, got: %q", out)
	}
}

func TestAnsiColorHex_ValidHex(t *testing.T) {
	out := AnsiColorHex("#FF0000", "test")
	if !strings.Contains(out, "test") {
		t.Errorf("expected text preserved, got: %q", out)
	}
	if !strings.Contains(out, "255;0;0") {
		t.Errorf("expected red RGB values in escape code, got: %q", out)
	}
}

func TestAnsiColorHex_InvalidHexFallsBackToPlainText(t *testing.T) {
	out := AnsiColorHex("garbage", "fallback")
	if out != "fallback" {
		t.Errorf("expected plain text fallback on invalid hex, got: %q", out)
	}
}

func TestGradientText_EmptyString(t *testing.T) {
	out, err := GradientText("", "#FF0000", "#0000FF")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output for empty input, got: %q", out)
	}
}

func TestGradientText_SingleChar(t *testing.T) {
	out, err := GradientText("X", "#FF0000", "#0000FF")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "X") {
		t.Errorf("expected character preserved, got: %q", out)
	}
}

func TestGradientText_InvalidFromColor(t *testing.T) {
	_, err := GradientText("test", "invalid", "#0000FF")
	if err == nil {
		t.Fatal("expected error for invalid 'from' color, got nil")
	}
}

func TestGradientText_InvalidToColor(t *testing.T) {
	_, err := GradientText("test", "#FF0000", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid 'to' color, got nil")
	}
}

func TestGradientText_MultiCharProducesColorPerChar(t *testing.T) {
	out, err := GradientText("ABC", "#FF0000", "#0000FF")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	// each char should get its own escape sequence -> 3 reset codes minimum
	count := strings.Count(out, "\033[0m")
	if count != 3 {
		t.Errorf("expected 3 ANSI reset codes (one per char), got %d", count)
	}
}

func TestGradientLine_ProducesCorrectWidth(t *testing.T) {
	out, err := GradientLine("─", 5, "#FF0000", "#0000FF")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	runeCount := strings.Count(out, "─")
	if runeCount != 5 {
		t.Errorf("expected 5 repeated chars, got %d", runeCount)
	}
}

func TestRainbowText_PreservesCharacters(t *testing.T) {
	out := RainbowText("hi")
	if !strings.Contains(out, "h") || !strings.Contains(out, "i") {
		t.Errorf("expected original characters preserved, got: %q", out)
	}
}

func TestComplementaryColor_Valid(t *testing.T) {
	comp, err := ComplementaryColor("#FF0000")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if comp == "" {
		t.Error("expected non-empty complementary color")
	}
	if comp == "#ff0000" {
		t.Error("expected complementary color to differ from input")
	}
}

func TestComplementaryColor_InvalidInput(t *testing.T) {
	_, err := ComplementaryColor("not-hex")
	if err == nil {
		t.Fatal("expected error for invalid hex input, got nil")
	}
}

func TestAnalogousColors_ReturnsDistinctColors(t *testing.T) {
	c1, c2, err := AnalogousColors("#FF0000")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if c1 == c2 {
		t.Error("expected two distinct analogous colors")
	}
}

func TestTriadicColors_ReturnsDistinctColors(t *testing.T) {
	c1, c2, err := TriadicColors("#00FF00")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if c1 == c2 {
		t.Error("expected two distinct triadic colors")
	}
}

func TestLightenColor_IncreasesLightness(t *testing.T) {
	lightened, err := LightenColor("#404040", 0.3)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	orig, _ := ParseHex("#404040")
	light, _ := ParseHex(lightened)
	_, _, vOrig := orig.Hsv()
	_, _, vLight := light.Hsv()
	if vLight <= vOrig {
		t.Errorf("expected lightened color to have higher value, orig=%f light=%f", vOrig, vLight)
	}
}

func TestDarkenColor_DecreasesLightness(t *testing.T) {
	darkened, err := DarkenColor("#C0C0C0", 0.3)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	orig, _ := ParseHex("#C0C0C0")
	dark, _ := ParseHex(darkened)
	_, _, vOrig := orig.Hsv()
	_, _, vDark := dark.Hsv()
	if vDark >= vOrig {
		t.Errorf("expected darkened color to have lower value, orig=%f dark=%f", vOrig, vDark)
	}
}

func TestLightenColor_ClampsAtMax(t *testing.T) {
	// already near max — should clamp at 1.0, not error or overflow
	result, err := LightenColor("#FFFFFF", 0.9)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	c, _ := ParseHex(result)
	_, _, v := c.Hsv()
	if v > 1.0 {
		t.Errorf("expected value clamped to <= 1.0, got %f", v)
	}
}

func TestDarkenColor_ClampsAtMin(t *testing.T) {
	result, err := DarkenColor("#000000", 0.9)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	c, _ := ParseHex(result)
	_, _, v := c.Hsv()
	if v < 0.0 {
		t.Errorf("expected value clamped to >= 0.0, got %f", v)
	}
}
