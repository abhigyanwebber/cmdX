package fonts

import (
	"strings"
	"testing"
)

func TestFind_ExistingFont(t *testing.T) {
	f := Find("firacode")
	if f == nil {
		t.Fatal("expected to find 'firacode' in catalog, got nil")
	}
	if f.Name != "firacode" {
		t.Errorf("expected name 'firacode', got %q", f.Name)
	}
}

func TestFind_NonexistentFont(t *testing.T) {
	if f := Find("not-a-real-font-xyz"); f != nil {
		t.Errorf("expected nil for nonexistent font, got %+v", f)
	}
}

func TestSearch_MatchesByName(t *testing.T) {
	results := Search("fira")
	if len(results) == 0 {
		t.Fatal("expected at least one match for 'fira'")
	}
	found := false
	for _, r := range results {
		if r.Name == "firacode" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'firacode' in search results for 'fira'")
	}
}

func TestSearch_MatchesByDescription(t *testing.T) {
	results := Search("ligature")
	if len(results) == 0 {
		t.Fatal("expected at least one match for 'ligature' (FiraCode's description)")
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	lower := Search("hack")
	upper := Search("HACK")
	if len(lower) != len(upper) {
		t.Errorf("expected case-insensitive search to return same count, got %d vs %d", len(lower), len(upper))
	}
	if len(lower) == 0 {
		t.Fatal("expected at least one match for 'hack'")
	}
}

func TestSearch_NoMatches(t *testing.T) {
	results := Search("zzz-nonexistent-query-zzz")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestDownloadURL_ContainsVersionAndZipName(t *testing.T) {
	f := Find("jetbrainsmono")
	if f == nil {
		t.Fatal("expected to find jetbrainsmono")
	}
	url := f.DownloadURL()
	if !strings.Contains(url, NerdFontsVersion) {
		t.Errorf("expected URL to contain version %q, got %q", NerdFontsVersion, url)
	}
	if !strings.Contains(url, f.ZipName) {
		t.Errorf("expected URL to contain zip name %q, got %q", f.ZipName, url)
	}
	if !strings.HasSuffix(url, ".zip") {
		t.Errorf("expected URL to end in .zip, got %q", url)
	}
}

func TestCatalog_AllEntriesHaveRequiredFields(t *testing.T) {
	for _, f := range Catalog {
		if f.Name == "" {
			t.Errorf("catalog entry with empty Name: %+v", f)
		}
		if f.DisplayName == "" {
			t.Errorf("catalog entry %q has empty DisplayName", f.Name)
		}
		if f.ZipName == "" {
			t.Errorf("catalog entry %q has empty ZipName", f.Name)
		}
		if f.License == "" {
			t.Errorf("catalog entry %q has empty License", f.Name)
		}
	}
}

func TestCatalog_NoDuplicateNames(t *testing.T) {
	seen := map[string]bool{}
	for _, f := range Catalog {
		if seen[f.Name] {
			t.Errorf("duplicate catalog entry name: %q", f.Name)
		}
		seen[f.Name] = true
	}
}
