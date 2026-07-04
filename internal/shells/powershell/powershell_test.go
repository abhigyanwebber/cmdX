package powershell

import (
	"os"
	"strings"
	"testing"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
	"github.com/abhigyanwebber/cmd-customizer/internal/shells"
)

func testTheme() *config.Theme {
	return &config.Theme{
		Meta: config.Meta{Name: "test-theme", Author: "tester", Description: "a sample theme"},
		Colors: config.Colors{
			Primary: "#FF00FF", Accent: "#00FFFF", Success: "#00FF88",
			Error: "#FF4444", Muted: "#444444",
		},
		Prompt: config.Prompt{Symbol: "▶"},
	}
}

func TestName_ReturnsPowershell(t *testing.T) {
	p := New()
	if p.Name() != "powershell" {
		t.Errorf("expected 'powershell', got %q", p.Name())
	}
}

func TestDetect_TrueWhenPSModulePathSet(t *testing.T) {
	original, hadOriginal := os.LookupEnv("PSModulePath")
	defer func() {
		if hadOriginal {
			os.Setenv("PSModulePath", original)
		} else {
			os.Unsetenv("PSModulePath")
		}
	}()

	os.Setenv("PSModulePath", "C:\\some\\path")
	p := New()
	if !p.Detect() {
		t.Error("expected Detect() to return true when PSModulePath is set")
	}
}

func TestDetect_FalseWhenPSModulePathUnset(t *testing.T) {
	original, hadOriginal := os.LookupEnv("PSModulePath")
	defer func() {
		if hadOriginal {
			os.Setenv("PSModulePath", original)
		}
	}()

	os.Unsetenv("PSModulePath")
	p := New()
	if p.Detect() {
		t.Error("expected Detect() to return false when PSModulePath is unset")
	}
}

func TestProfilePath_ContainsExpectedSuffix(t *testing.T) {
	p := New()
	path, err := p.ProfilePath()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(path, "Microsoft.PowerShell_profile.ps1") {
		t.Errorf("expected profile path to contain profile filename, got: %s", path)
	}
	if !strings.Contains(path, "PowerShell") {
		t.Errorf("expected profile path to contain PowerShell dir, got: %s", path)
	}
}

func TestBuildPowerShellBlock_ContainsInjectMarkers(t *testing.T) {
	block := buildPowerShellBlock(testTheme())
	if !strings.Contains(block, shells.InjectStart) {
		t.Error("expected block to contain InjectStart marker")
	}
	if !strings.Contains(block, shells.InjectEnd) {
		t.Error("expected block to contain InjectEnd marker")
	}
}

func TestBuildPowerShellBlock_ContainsThemeColors(t *testing.T) {
	th := testTheme()
	block := buildPowerShellBlock(th)
	if !strings.Contains(block, th.Colors.Primary) {
		t.Errorf("expected block to contain primary color %s", th.Colors.Primary)
	}
	if !strings.Contains(block, th.Colors.Accent) {
		t.Errorf("expected block to contain accent color %s", th.Colors.Accent)
	}
}

func TestBuildPowerShellBlock_ContainsThemeMeta(t *testing.T) {
	th := testTheme()
	block := buildPowerShellBlock(th)
	if !strings.Contains(block, th.Meta.Name) {
		t.Errorf("expected block to contain theme name %s", th.Meta.Name)
	}
	if !strings.Contains(block, th.Meta.Author) {
		t.Errorf("expected block to contain author %s", th.Meta.Author)
	}
}

func TestRemoveInjection_StripsMarkedBlock(t *testing.T) {
	content := "existing content\n" + shells.InjectStart + "\ninjected stuff\n" + shells.InjectEnd + "\nmore content"
	cleaned := removeInjection(content)

	if strings.Contains(cleaned, shells.InjectStart) {
		t.Error("expected InjectStart marker to be removed")
	}
	if strings.Contains(cleaned, "injected stuff") {
		t.Error("expected injected content to be removed")
	}
	if !strings.Contains(cleaned, "existing content") {
		t.Error("expected surrounding content to be preserved")
	}
	if !strings.Contains(cleaned, "more content") {
		t.Error("expected trailing content to be preserved")
	}
}

func TestRemoveInjection_NoMarkersPresent(t *testing.T) {
	content := "just plain profile content, no cmdx here"
	cleaned := removeInjection(content)
	if cleaned != content {
		t.Errorf("expected content unchanged when no markers present, got: %q", cleaned)
	}
}

func TestRemoveInjection_MultipleInjectedBlocks(t *testing.T) {
	content := "a\n" + shells.InjectStart + "\nfirst\n" + shells.InjectEnd +
		"\nb\n" + shells.InjectStart + "\nsecond\n" + shells.InjectEnd + "\nc"
	cleaned := removeInjection(content)

	if strings.Contains(cleaned, "first") || strings.Contains(cleaned, "second") {
		t.Error("expected all injected blocks to be removed")
	}
}

func TestIsInjected_TrueAfterInject(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	p := New()
	if err := p.Inject(testTheme()); err != nil {
		t.Fatalf("inject failed: %v", err)
	}

	injected, err := p.IsInjected()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !injected {
		t.Error("expected IsInjected to return true after Inject")
	}
}

func TestIsInjected_FalseWithNoProfile(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	p := New()
	injected, err := p.IsInjected()
	if err != nil {
		t.Fatalf("expected no error (missing file should return false,nil), got: %v", err)
	}
	if injected {
		t.Error("expected IsInjected false when no profile exists")
	}
}

func TestInjectThenRemove_RoundTrip(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	p := New()
	if err := p.Inject(testTheme()); err != nil {
		t.Fatalf("inject failed: %v", err)
	}

	if err := p.Remove(); err != nil {
		t.Fatalf("remove failed: %v", err)
	}

	injected, _ := p.IsInjected()
	if injected {
		t.Error("expected theme to be removed after Remove()")
	}
}

func TestInject_OverwritesPreviousInjection(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("HOME", tmpHome)

	p := New()
	th1 := testTheme()
	th1.Meta.Name = "first-theme"
	p.Inject(th1)

	th2 := testTheme()
	th2.Meta.Name = "second-theme"
	p.Inject(th2)

	path, _ := p.ProfilePath()
	data, _ := os.ReadFile(path)
	content := string(data)

	if strings.Contains(content, "first-theme") {
		t.Error("expected first theme injection to be replaced")
	}
	if !strings.Contains(content, "second-theme") {
		t.Error("expected second theme injection to be present")
	}
}
