package fonts

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// buildTestZip creates an in-memory zip archive with the given
// filename -> content entries, returning the raw zip bytes.
func buildTestZip(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for name, content := range entries {
		f, err := w.Create(name)
		if err != nil {
			t.Fatalf("could not create zip entry %q: %v", name, err)
		}
		if _, err := f.Write([]byte(content)); err != nil {
			t.Fatalf("could not write zip entry %q: %v", name, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("could not close zip writer: %v", err)
	}
	return buf.Bytes()
}

func writeTestZipFile(t *testing.T, entries map[string]string) string {
	t.Helper()
	data := buildTestZip(t, entries)
	f, err := os.CreateTemp(t.TempDir(), "test-*.zip")
	if err != nil {
		t.Fatalf("could not create temp zip file: %v", err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatalf("could not write temp zip file: %v", err)
	}
	f.Close()
	return f.Name()
}

// ── sanitizeExtractPath ───────────────────────────────────────────────────────

func TestSanitizeExtractPath_ValidRelativePath(t *testing.T) {
	dest := "/home/user/fonts"
	path, err := sanitizeExtractPath(dest, "MyFont-Regular.ttf")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	expected := filepath.Join(dest, "MyFont-Regular.ttf")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}

func TestSanitizeExtractPath_RejectsParentTraversal(t *testing.T) {
	dest := "/home/user/fonts"
	cases := []string{
		"../../../etc/passwd",
		"../../.bashrc",
		"..",
		"../sibling.ttf",
	}
	for _, c := range cases {
		if _, err := sanitizeExtractPath(dest, c); err == nil {
			t.Errorf("expected error for traversal path %q, got nil", c)
		}
	}
}

func TestSanitizeExtractPath_RejectsAbsolutePath(t *testing.T) {
	dest := "/home/user/fonts"
	if _, err := sanitizeExtractPath(dest, "/etc/passwd"); err == nil {
		t.Fatal("expected error for absolute path, got nil")
	}
}

func TestSanitizeExtractPath_AllowsNestedButSafePath(t *testing.T) {
	dest := "/home/user/fonts"
	// a subdirectory within dest is fine as long as it doesn't escape
	path, err := sanitizeExtractPath(dest, "subdir/font.ttf")
	if err != nil {
		t.Fatalf("expected no error for safe nested path, got: %v", err)
	}
	expected := filepath.Join(dest, "subdir", "font.ttf")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}

// ── extractFontFiles ──────────────────────────────────────────────────────────

func TestExtractFontFiles_ExtractsOnlyFontExtensions(t *testing.T) {
	zipPath := writeTestZipFile(t, map[string]string{
		"MyFont-Regular.ttf": "fake ttf data",
		"MyFont-Bold.otf":    "fake otf data",
		"README.md":          "not a font",
		"LICENSE":            "license text",
	})

	destDir := t.TempDir()
	extracted, err := extractFontFiles(zipPath, destDir, "", true)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(extracted) != 2 {
		t.Errorf("expected 2 font files extracted, got %d: %v", len(extracted), extracted)
	}
}

func TestExtractFontFiles_VariantFilter(t *testing.T) {
	zipPath := writeTestZipFile(t, map[string]string{
		"MyFont-Regular.ttf": "regular",
		"MyFont-Bold.ttf":    "bold",
		"MyFont-Italic.ttf":  "italic",
	})

	destDir := t.TempDir()
	extracted, err := extractFontFiles(zipPath, destDir, "Bold", false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(extracted) != 1 {
		t.Fatalf("expected 1 file matching 'Bold', got %d: %v", len(extracted), extracted)
	}
}

func TestExtractFontFiles_VariantFilterCaseInsensitive(t *testing.T) {
	zipPath := writeTestZipFile(t, map[string]string{
		"MyFont-BOLD.ttf": "bold upper",
	})
	destDir := t.TempDir()
	extracted, err := extractFontFiles(zipPath, destDir, "bold", false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(extracted) != 1 {
		t.Fatalf("expected case-insensitive match, got %d files", len(extracted))
	}
}

func TestExtractFontFiles_NoMatchingVariantErrors(t *testing.T) {
	zipPath := writeTestZipFile(t, map[string]string{
		"MyFont-Regular.ttf": "regular",
	})
	destDir := t.TempDir()
	_, err := extractFontFiles(zipPath, destDir, "Condensed", false)
	if err == nil {
		t.Fatal("expected error when no files match variant filter, got nil")
	}
}

func TestExtractFontFiles_NoFontFilesInArchive(t *testing.T) {
	zipPath := writeTestZipFile(t, map[string]string{
		"README.md": "just docs",
	})
	destDir := t.TempDir()
	_, err := extractFontFiles(zipPath, destDir, "", true)
	if err == nil {
		t.Fatal("expected error when archive has no font files, got nil")
	}
}

func TestExtractFontFiles_WritesReadableFiles(t *testing.T) {
	content := "this is fake but valid enough ttf content for a test"
	zipPath := writeTestZipFile(t, map[string]string{
		"Test-Regular.ttf": content,
	})
	destDir := t.TempDir()
	extracted, err := extractFontFiles(zipPath, destDir, "", true)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	data, err := os.ReadFile(extracted[0])
	if err != nil {
		t.Fatalf("could not read extracted file: %v", err)
	}
	if string(data) != content {
		t.Errorf("expected extracted content %q, got %q", content, string(data))
	}
}

func TestExtractFontFiles_InvalidZipFile(t *testing.T) {
	badZip := filepath.Join(t.TempDir(), "notazip.zip")
	os.WriteFile(badZip, []byte("this is not a zip file"), 0644)

	destDir := t.TempDir()
	_, err := extractFontFiles(badZip, destDir, "", true)
	if err == nil {
		t.Fatal("expected error for invalid zip file, got nil")
	}
}
