package assets

import (
	"regexp"
	"strings"
	"testing"
)

func buildValidStatusBar() *StatusBarConfig {
	inGitRepo := true
	return &StatusBarConfig{
		Position:       StatusBarAbove,
		SeparatorStyle: SeparatorBar,
		Height:         1,
		Segments: []SegmentConfig{
			{Type: SegmentDirectory, Zone: ZoneLeft, Color: "#00FFFF", Label: " ", Order: 1},
			{Type: SegmentGit, Zone: ZoneLeft, Color: "#FF00FF", Label: " ", Order: 2,
				Conditions: []SegmentCondition{{InGitRepo: &inGitRepo}}},
			{Type: SegmentTime, Zone: ZoneRight, Color: "#FFD700", Format: "%H:%M", Order: 1},
			{Type: SegmentExitCode, Zone: ZoneRight, Color: "#FF4444", Order: 2},
		},
	}
}

// ── ValidateStatusBar ─────────────────────────────────────────────────────────

func TestValidateStatusBar_Valid(t *testing.T) {
	if err := ValidateStatusBar(buildValidStatusBar()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateStatusBar_Nil(t *testing.T) {
	if err := ValidateStatusBar(nil); err == nil {
		t.Fatal("expected error for nil config, got nil")
	}
}

func TestValidateStatusBar_NoSegments(t *testing.T) {
	sb := buildValidStatusBar()
	sb.Segments = nil
	if err := ValidateStatusBar(sb); err == nil {
		t.Fatal("expected error for zero segments, got nil")
	}
}

func TestValidateStatusBar_ZeroHeight(t *testing.T) {
	sb := buildValidStatusBar()
	sb.Height = 0
	if err := ValidateStatusBar(sb); err == nil {
		t.Fatal("expected error for zero height, got nil")
	}
}

func TestValidateStatusBar_InvalidPosition(t *testing.T) {
	sb := buildValidStatusBar()
	sb.Position = "floating"
	if err := ValidateStatusBar(sb); err == nil {
		t.Fatal("expected error for invalid position, got nil")
	}
}

func TestValidateStatusBar_AllValidPositions(t *testing.T) {
	for _, pos := range []StatusBarPosition{StatusBarAbove, StatusBarBelow, StatusBarInline} {
		sb := buildValidStatusBar()
		sb.Position = pos
		if err := ValidateStatusBar(sb); err != nil {
			t.Errorf("expected position %q to be valid, got: %v", pos, err)
		}
	}
}

func TestValidateStatusBar_InvalidSeparator(t *testing.T) {
	sb := buildValidStatusBar()
	sb.SeparatorStyle = "zigzag"
	if err := ValidateStatusBar(sb); err == nil {
		t.Fatal("expected error for invalid separator style, got nil")
	}
}

func TestValidateStatusBar_AllValidSeparators(t *testing.T) {
	separators := []SeparatorStyle{
		SeparatorPowerline, SeparatorRounded, SeparatorAngular,
		SeparatorBar, SeparatorDot, SeparatorSlash, SeparatorNone,
	}
	for _, sep := range separators {
		sb := buildValidStatusBar()
		sb.SeparatorStyle = sep
		if err := ValidateStatusBar(sb); err != nil {
			t.Errorf("expected separator %q to be valid, got: %v", sep, err)
		}
	}
}

func TestValidateStatusBar_CustomSeparatorNeedsGlyphs(t *testing.T) {
	sb := buildValidStatusBar()
	sb.SeparatorStyle = SeparatorCustom
	// no custom glyphs set → error
	if err := ValidateStatusBar(sb); err == nil {
		t.Fatal("expected error for custom separator without glyphs, got nil")
	}
}

func TestValidateStatusBar_CustomSeparatorWithGlyphs(t *testing.T) {
	sb := buildValidStatusBar()
	sb.SeparatorStyle = SeparatorCustom
	sb.CustomSeparatorLeft = ">"
	if err := ValidateStatusBar(sb); err != nil {
		t.Fatalf("expected no error with custom left glyph set, got: %v", err)
	}
}

// ── validateSegment ───────────────────────────────────────────────────────────

func TestValidateSegment_AllSegmentTypes(t *testing.T) {
	types := []SegmentType{
		SegmentGit, SegmentDirectory, SegmentTime, SegmentDate,
		SegmentExitCode, SegmentDuration, SegmentVirtualEnv,
		SegmentGoVersion, SegmentNodeVersion, SegmentKubernetes,
		SegmentBattery, SegmentUser, SegmentHost,
	}
	for _, st := range types {
		seg := SegmentConfig{Type: st, Zone: ZoneLeft}
		if err := validateSegment(seg, 0); err != nil {
			t.Errorf("expected segment type %q to be valid, got: %v", st, err)
		}
	}
}

func TestValidateSegment_EnvVarRequiresEnvVarField(t *testing.T) {
	seg := SegmentConfig{Type: SegmentEnvVar, Zone: ZoneLeft}
	if err := validateSegment(seg, 0); err == nil {
		t.Fatal("expected error for env_var segment with no env_var field, got nil")
	}
}

func TestValidateSegment_CommandRequiresCommandField(t *testing.T) {
	seg := SegmentConfig{Type: SegmentCommand, Zone: ZoneLeft}
	if err := validateSegment(seg, 0); err == nil {
		t.Fatal("expected error for command segment with no command field, got nil")
	}
}

func TestValidateSegment_TextRequiresFormat(t *testing.T) {
	seg := SegmentConfig{Type: SegmentText, Zone: ZoneLeft}
	if err := validateSegment(seg, 0); err == nil {
		t.Fatal("expected error for text segment with no format, got nil")
	}
}

func TestValidateSegment_TextWithFormatValid(t *testing.T) {
	seg := SegmentConfig{Type: SegmentText, Zone: ZoneLeft, Format: "hello"}
	if err := validateSegment(seg, 0); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateSegment_InvalidZone(t *testing.T) {
	seg := SegmentConfig{Type: SegmentGit, Zone: "middle"}
	if err := validateSegment(seg, 0); err == nil {
		t.Fatal("expected error for invalid zone, got nil")
	}
}

func TestValidateSegment_AllValidZones(t *testing.T) {
	for _, z := range []StatusBarZone{ZoneLeft, ZoneCenter, ZoneRight} {
		seg := SegmentConfig{Type: SegmentGit, Zone: z}
		if err := validateSegment(seg, 0); err != nil {
			t.Errorf("expected zone %q to be valid, got: %v", z, err)
		}
	}
}

func TestValidateSegment_UnknownType(t *testing.T) {
	seg := SegmentConfig{Type: "moon_phase", Zone: ZoneLeft}
	if err := validateSegment(seg, 0); err == nil {
		t.Fatal("expected error for unknown segment type, got nil")
	}
}

// ── separatorGlyphs ───────────────────────────────────────────────────────────

func TestSeparatorGlyphs_PowerlineNotEmpty(t *testing.T) {
	l, r := separatorGlyphs(SeparatorPowerline, "", "")
	if l == "" || r == "" {
		t.Error("expected non-empty powerline glyphs")
	}
	if l == r {
		t.Error("expected powerline left and right glyphs to differ")
	}
}

func TestSeparatorGlyphs_NoneReturnsSpace(t *testing.T) {
	l, r := separatorGlyphs(SeparatorNone, "", "")
	if l != " " || r != " " {
		t.Errorf("expected spaces for none separator, got %q and %q", l, r)
	}
}

func TestSeparatorGlyphs_CustomUsesProvidedGlyphs(t *testing.T) {
	l, r := separatorGlyphs(SeparatorCustom, ">>", "<<")
	if l != ">>" || r != "<<" {
		t.Errorf("expected custom glyphs '>>' and '<<', got %q and %q", l, r)
	}
}

func TestSeparatorGlyphs_CustomFallsBackToSpace(t *testing.T) {
	l, r := separatorGlyphs(SeparatorCustom, "", "")
	if l != " " || r != " " {
		t.Errorf("expected space fallback for empty custom glyphs, got %q and %q", l, r)
	}
}

// ── StatusBarShellCode ────────────────────────────────────────────────────────

func TestStatusBarShellCode_BashNotEmpty(t *testing.T) {
	sb := buildValidStatusBar()
	colors := map[string]string{"primary": "#FF00FF"}
	code := StatusBarShellCode(sb, colors, "bash")
	if code == "" {
		t.Fatal("expected non-empty bash code, got empty")
	}
}

func TestStatusBarShellCode_ZshNotEmpty(t *testing.T) {
	sb := buildValidStatusBar()
	code := StatusBarShellCode(sb, nil, "zsh")
	if code == "" {
		t.Fatal("expected non-empty zsh code, got empty")
	}
}

func TestStatusBarShellCode_PowerShellNotEmpty(t *testing.T) {
	sb := buildValidStatusBar()
	code := StatusBarShellCode(sb, nil, "powershell")
	if code == "" {
		t.Fatal("expected non-empty powershell code, got empty")
	}
}

func TestStatusBarShellCode_UnsupportedShellEmpty(t *testing.T) {
	sb := buildValidStatusBar()
	code := StatusBarShellCode(sb, nil, "fish")
	if code != "" {
		t.Error("expected empty code for unsupported shell 'fish'")
	}
}

func TestStatusBarShellCode_BashContainsPromptCommand(t *testing.T) {
	sb := buildValidStatusBar()
	sb.Position = StatusBarAbove
	code := StatusBarShellCode(sb, nil, "bash")
	if !strings.Contains(code, "PROMPT_COMMAND") {
		t.Error("expected bash status bar above to use PROMPT_COMMAND")
	}
}

func TestStatusBarShellCode_ZshContainsAddZshHook(t *testing.T) {
	sb := buildValidStatusBar()
	code := StatusBarShellCode(sb, nil, "zsh")
	if !strings.Contains(code, "add-zsh-hook") {
		t.Error("expected zsh status bar to use add-zsh-hook")
	}
}

func TestStatusBarShellCode_PowerShellContainsLastExitCode(t *testing.T) {
	sb := buildValidStatusBar()
	code := StatusBarShellCode(sb, nil, "powershell")
	if !strings.Contains(code, "LASTEXITCODE") {
		t.Error("expected powershell status bar to check $LASTEXITCODE")
	}
}

func TestStatusBarShellCode_GitSegmentGeneratesGitCode(t *testing.T) {
	sb := buildValidStatusBar()
	code := StatusBarShellCode(sb, nil, "bash")
	if !strings.Contains(code, "git") {
		t.Error("expected bash status bar with git segment to include git commands")
	}
}

func TestStatusBarShellCode_BashAboveVsBelowDiffers(t *testing.T) {
	sbAbove := buildValidStatusBar()
	sbAbove.Position = StatusBarAbove
	sbBelow := buildValidStatusBar()
	sbBelow.Position = StatusBarBelow

	codeAbove := StatusBarShellCode(sbAbove, nil, "bash")
	codeBelow := StatusBarShellCode(sbBelow, nil, "bash")

	if codeAbove == codeBelow {
		t.Error("expected above and below positions to generate different shell code")
	}
}

// ── segmentsByZone ────────────────────────────────────────────────────────────

func TestSegmentsByZone_CorrectCounts(t *testing.T) {
	sb := buildValidStatusBar()
	left, center, right := segmentsByZone(sb)
	if len(left) != 2 {
		t.Errorf("expected 2 left segments, got %d", len(left))
	}
	if len(center) != 0 {
		t.Errorf("expected 0 center segments, got %d", len(center))
	}
	if len(right) != 2 {
		t.Errorf("expected 2 right segments, got %d", len(right))
	}
}

func TestSegmentsByZone_SortedByOrder(t *testing.T) {
	sb := &StatusBarConfig{
		Position:       StatusBarAbove,
		SeparatorStyle: SeparatorNone,
		Height:         1,
		Segments: []SegmentConfig{
			{Type: SegmentTime, Zone: ZoneLeft, Order: 3},
			{Type: SegmentDirectory, Zone: ZoneLeft, Order: 1},
			{Type: SegmentGit, Zone: ZoneLeft, Order: 2},
		},
	}
	left, _, _ := segmentsByZone(sb)
	if left[0].Type != SegmentDirectory || left[1].Type != SegmentGit || left[2].Type != SegmentTime {
		t.Error("expected segments sorted by Order: directory(1) < git(2) < time(3)")
	}
}

// ── resolveSegColor ───────────────────────────────────────────────────────────

func TestResolveSegColor_HexPassthrough(t *testing.T) {
	result := resolveSegColor("#FF00FF", nil)
	if result != "#FF00FF" {
		t.Errorf("expected hex to pass through unchanged, got %q", result)
	}
}

func TestResolveSegColor_SemanticKeyResolved(t *testing.T) {
	palette := map[string]string{"primary": "#00FFFF"}
	result := resolveSegColor("primary", palette)
	if result != "#00FFFF" {
		t.Errorf("expected 'primary' to resolve to '#00FFFF', got %q", result)
	}
}

func TestResolveSegColor_EmptyReturnsEmpty(t *testing.T) {
	if result := resolveSegColor("", nil); result != "" {
		t.Errorf("expected empty color to stay empty, got %q", result)
	}
}

func TestResolveSegColor_UnknownKeyPassthrough(t *testing.T) {
	result := resolveSegColor("unknown-key", map[string]string{})
	if result != "unknown-key" {
		t.Errorf("expected unknown key to pass through, got %q", result)
	}
}

// ── Compute/assembly variable-name consistency ─────────────────────────────
//
// Regression coverage for a real, shipped bug: the per-segment value
// computation and the final bar assembly used two independently-written,
// disconnected variable-naming schemes. Every segment beyond the first
// in a zone silently referenced a variable that was never assigned,
// so the rendered status bar was empty for any realistic multi-segment
// bar. Prior tests did call StatusBarShellCode, but only checked
// superficial properties (non-empty output, presence of expected
// substrings like "PROMPT_COMMAND") — none inspected whether the
// variables the assembly step checks were ever actually assigned by
// the compute step, which is exactly where this bug lived.
//
// assignedVars finds every `name="..."` / `name=""` / `local name="..."`
// assignment target in generated shell code.
func assignedVars(code string) map[string]bool {
	found := map[string]bool{}
	re := regexp.MustCompile(`(?m)^\s*(?:local\s+)?(\$?__seg_\w+)=`)
	for _, m := range re.FindAllStringSubmatch(code, -1) {
		found[strings.TrimPrefix(m[1], "$")] = true
	}
	return found
}

// checkedVars finds every variable referenced in a `-n "$name"` /
// `-n \"$name\"` truthiness check — the assembly step's gate for
// whether to include a segment in the rendered bar.
func checkedVars(code string) []string {
	re := regexp.MustCompile(`-n\s+\\?"\$(__seg_\w+)\\?"`)
	var found []string
	for _, m := range re.FindAllStringSubmatch(code, -1) {
		found = append(found, m[1])
	}
	return found
}

func TestStatusBarShellCode_BashAssemblyVarsAreAllAssigned(t *testing.T) {
	code := StatusBarShellCode(buildValidStatusBar(), nil, "bash")
	assigned := assignedVars(code)
	for _, v := range checkedVars(code) {
		if !assigned[v] {
			t.Errorf("bash assembly checks $%s but no compute step ever assigns it — the bar would silently omit this segment forever; full generated code:\n%s", v, code)
		}
	}
	if len(checkedVars(code)) == 0 {
		t.Fatal("test fixture produced no assembly checks at all — fixture or regex is broken, this test isn't actually testing anything")
	}
}

func TestStatusBarShellCode_ZshAssemblyVarsAreAllAssigned(t *testing.T) {
	code := StatusBarShellCode(buildValidStatusBar(), nil, "zsh")
	assigned := assignedVars(code)
	for _, v := range checkedVars(code) {
		if !assigned[v] {
			t.Errorf("zsh assembly checks $%s but no compute step ever assigns it — the bar would silently omit this segment forever; full generated code:\n%s", v, code)
		}
	}
	if len(checkedVars(code)) == 0 {
		t.Fatal("test fixture produced no assembly checks at all — fixture or regex is broken, this test isn't actually testing anything")
	}
}

// TestBashSegVarName_UniquePerZoneRelativeIndex guards the specific
// mechanism the bug above hinged on: two segments of the same type at
// different zone-relative indices must get different variable names,
// and the same (type, index) pair must always produce the same name
// (compute and assembly only stay in sync if this holds).
func TestBashSegVarName_UniquePerZoneRelativeIndex(t *testing.T) {
	seg := SegmentConfig{Type: SegmentText}
	if bashSegVarName(seg, 0) == bashSegVarName(seg, 1) {
		t.Error("expected different indices to produce different variable names")
	}
	if bashSegVarName(seg, 2) != bashSegVarName(seg, 2) {
		t.Error("expected the same (type, index) pair to always produce the same name")
	}
}
