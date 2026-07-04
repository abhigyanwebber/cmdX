package registry

import (
	"os"
	"strings"
	"testing"
)

func sampleIndex() *Index {
	return &Index{
		UpdatedAt: "2026-01-01",
		Themes: []ThemeEntry{
			{Name: "ocean", Author: "abhigyan", Description: "calm blue waves", Tags: []string{"blue", "calm"}},
			{Name: "cyberpunk", Author: "abhigyan", Description: "neon soaked terminal", Tags: []string{"neon", "dark"}},
			{Name: "forest", Author: "someone", Description: "green canopy vibes", Tags: []string{"green", "nature"}},
		},
	}
}

func TestSearch_MatchesByName(t *testing.T) {
	idx := sampleIndex()
	results := Search(idx, "ocean")
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'ocean', got %d", len(results))
	}
	if results[0].Name != "ocean" {
		t.Errorf("expected match 'ocean', got %q", results[0].Name)
	}
}

func TestSearch_MatchesByDescription(t *testing.T) {
	idx := sampleIndex()
	results := Search(idx, "neon")
	if len(results) != 1 {
		t.Fatalf("expected 1 result matching description 'neon', got %d", len(results))
	}
}

func TestSearch_MatchesByTag(t *testing.T) {
	idx := sampleIndex()
	results := Search(idx, "nature")
	if len(results) != 1 {
		t.Fatalf("expected 1 result matching tag 'nature', got %d", len(results))
	}
	if results[0].Name != "forest" {
		t.Errorf("expected match 'forest', got %q", results[0].Name)
	}
}

func TestSearch_NoMatches(t *testing.T) {
	idx := sampleIndex()
	results := Search(idx, "nonexistent-query-xyz")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_EmptyQueryMatchesAll(t *testing.T) {
	idx := sampleIndex()
	results := Search(idx, "")
	if len(results) != len(idx.Themes) {
		t.Errorf("expected empty query to match all %d themes, got %d", len(idx.Themes), len(results))
	}
}

func TestSearch_EmptyIndex(t *testing.T) {
	idx := &Index{Themes: []ThemeEntry{}}
	results := Search(idx, "anything")
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty index, got %d", len(results))
	}
}

// Search is case-insensitive — queries match regardless of case in name,
// description, or tags.
func TestSearch_CaseInsensitiveByName(t *testing.T) {
	idx := sampleIndex()
	results := Search(idx, "OCEAN")
	if len(results) != 1 {
		t.Fatalf("expected case-insensitive match for 'OCEAN', got %d results", len(results))
	}
	if results[0].Name != "ocean" {
		t.Errorf("expected match 'ocean', got %q", results[0].Name)
	}
}

func TestSearch_CaseInsensitiveByDescription(t *testing.T) {
	idx := sampleIndex()
	results := Search(idx, "NEON SOAKED")
	if len(results) != 1 {
		t.Fatalf("expected case-insensitive description match, got %d results", len(results))
	}
}

func TestSearch_CaseInsensitiveByTag(t *testing.T) {
	idx := sampleIndex()
	results := Search(idx, "NATURE")
	if len(results) != 1 {
		t.Fatalf("expected case-insensitive tag match, got %d results", len(results))
	}
	if results[0].Name != "forest" {
		t.Errorf("expected match 'forest', got %q", results[0].Name)
	}
}

func TestSearch_MixedCaseQuery(t *testing.T) {
	idx := sampleIndex()
	results := Search(idx, "CybeRPunk")
	if len(results) != 1 {
		t.Fatalf("expected mixed-case query to match, got %d results", len(results))
	}
	if results[0].Name != "cyberpunk" {
		t.Errorf("expected match 'cyberpunk', got %q", results[0].Name)
	}
}

func TestValidThemeName_ValidNames(t *testing.T) {
	cases := []string{"ocean", "cyberpunk-2077", "my_theme", "Theme123", "a"}
	for _, name := range cases {
		if !ValidThemeName(name) {
			t.Errorf("expected %q to be a valid theme name", name)
		}
	}
}

func TestValidThemeName_RejectsPathTraversal(t *testing.T) {
	cases := []string{
		"../../../etc/passwd",
		"../../.bashrc",
		"..%2F..%2Fetc",
		"theme/../../../secrets",
	}
	for _, name := range cases {
		if ValidThemeName(name) {
			t.Errorf("expected path traversal attempt %q to be rejected", name)
		}
	}
}

func TestValidThemeName_RejectsSlashesAndDots(t *testing.T) {
	cases := []string{"a/b", "a\\b", "a.b", "a..b", "/etc/passwd", "C:\\Windows"}
	for _, name := range cases {
		if ValidThemeName(name) {
			t.Errorf("expected %q to be rejected (contains slash/dot)", name)
		}
	}
}

func TestValidThemeName_RejectsEmpty(t *testing.T) {
	if ValidThemeName("") {
		t.Error("expected empty string to be rejected")
	}
}

func TestValidThemeName_RejectsShellMetacharacters(t *testing.T) {
	cases := []string{"theme;rm -rf", "theme$(whoami)", "theme`id`", "theme|cat"}
	for _, name := range cases {
		if ValidThemeName(name) {
			t.Errorf("expected %q to be rejected (shell metacharacters)", name)
		}
	}
}

func TestFetchTheme_RejectsInvalidNameBeforeNetworkCall(t *testing.T) {
	dir := t.TempDir()
	err := FetchTheme("../../../etc/passwd", dir)
	if err == nil {
		t.Fatal("expected error for path traversal theme name, got nil")
	}
	if !strings.Contains(err.Error(), "invalid theme name") {
		t.Errorf("expected 'invalid theme name' error, got: %v", err)
	}

	// confirm nothing was written anywhere under dir or its parents
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Error("expected no files written for an invalid theme name")
	}
}
