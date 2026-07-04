package graphics

import (
	"strings"
	"testing"
)

func TestGlitchText_PreservesLength(t *testing.T) {
	input := "hello world"
	out := GlitchText(input, 0.5)
	if len([]rune(out)) != len([]rune(input)) {
		t.Errorf("expected output length %d to match input length, got %d",
			len([]rune(input)), len([]rune(out)))
	}
}

func TestGlitchText_ZeroIntensityPreservesText(t *testing.T) {
	input := "stable text"
	out := GlitchText(input, 0.0)
	if out != input {
		t.Errorf("expected zero intensity to preserve text exactly, got: %q", out)
	}
}

func TestGlitchText_PreservesSpaces(t *testing.T) {
	input := "a b c"
	out := GlitchText(input, 1.0) // max intensity
	runes := []rune(out)
	if runes[1] != ' ' || runes[3] != ' ' {
		t.Errorf("expected spaces to be preserved even at max intensity, got: %q", out)
	}
}

func TestGlitchText_FullIntensityChangesNonSpaceChars(t *testing.T) {
	input := "aaaaaaaaaa"
	out := GlitchText(input, 1.0)
	if out == input {
		t.Error("expected full intensity glitch to alter text, got identical output")
	}
}

func TestGlitchFrames_ReturnsRequestedCount(t *testing.T) {
	frames := GlitchFrames("test", 5)
	if len(frames) != 5 {
		t.Errorf("expected 5 frames, got %d", len(frames))
	}
}

func TestGlitchFrames_ZeroCount(t *testing.T) {
	frames := GlitchFrames("test", 0)
	if len(frames) != 0 {
		t.Errorf("expected 0 frames, got %d", len(frames))
	}
}

func TestGlitchFrames_EachFrameSameLength(t *testing.T) {
	input := "consistent"
	frames := GlitchFrames(input, 3)
	for i, f := range frames {
		if len([]rune(f)) != len([]rune(input)) {
			t.Errorf("frame %d: expected length %d, got %d", i, len([]rune(input)), len([]rune(f)))
		}
	}
}

func TestPulseText_WrapsWithAnsiCodes(t *testing.T) {
	out := PulseText("hi", 0)
	if !strings.Contains(out, "hi") {
		t.Errorf("expected text preserved, got: %q", out)
	}
	if !strings.HasSuffix(out, "\033[0m") {
		t.Errorf("expected reset code at end, got: %q", out)
	}
}

func TestPulseText_CyclesThroughFrames(t *testing.T) {
	// styles array has 5 entries — frame 5 should wrap to style[0]
	out0 := PulseText("x", 0)
	out5 := PulseText("x", 5)
	if !strings.Contains(out0, "x") || !strings.Contains(out5, "x") {
		t.Error("expected text preserved across all frame indices")
	}
}

func TestNeonText_ValidColor(t *testing.T) {
	out, err := NeonText("glow", "#00FFFF")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	// each character is individually ANSI-wrapped, so check each rune is present
	// rather than expecting "glow" as a contiguous substring
	for _, ch := range "glow" {
		if !strings.ContainsRune(out, ch) {
			t.Errorf("expected character %q to be present in output, got: %q", ch, out)
		}
	}
	if !strings.HasPrefix(out, "\033[1m") {
		t.Errorf("expected bold prefix, got: %q", out)
	}
}

func TestNeonText_InvalidColor(t *testing.T) {
	_, err := NeonText("test", "not-a-hex-color")
	if err == nil {
		t.Fatal("expected error for invalid color, got nil")
	}
}

func TestNeonText_SingleCharacter(t *testing.T) {
	out, err := NeonText("X", "#FF00FF")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "X") {
		t.Errorf("expected character preserved, got: %q", out)
	}
}

func TestNeonText_EmptyString(t *testing.T) {
	out, err := NeonText("", "#FF00FF")
	if err != nil {
		t.Fatalf("expected no error for empty string, got: %v", err)
	}
	// should still have bold wrapper even with empty content
	if !strings.Contains(out, "\033[1m") {
		t.Errorf("expected bold wrapper present, got: %q", out)
	}
}
