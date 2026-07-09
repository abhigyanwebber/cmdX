package fonts

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// withTempHome redirects os.UserHomeDir()'s source (HOME/USERPROFILE)
// to a temp directory for the duration of the test, so font
// installation and state tracking don't touch the real machine.
func withTempHome(t *testing.T) string {
	t.Helper()
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("LOCALAPPDATA", filepath.Join(tmpHome, "AppData", "Local"))
	return tmpHome
}

func fakeFontZipServer(t *testing.T, entries map[string]string) *httptest.Server {
	t.Helper()
	data := buildTestZip(t, entries)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Write(data)
	}))
}

func TestFontsDir_CreatesDirectory(t *testing.T) {
	withTempHome(t)
	dir, err := FontsDir()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected fonts directory to exist, got: %v", err)
	}
}

func TestIsInstalled_FalseInitially(t *testing.T) {
	withTempHome(t)
	if IsInstalled("firacode") {
		t.Error("expected IsInstalled to be false with no prior installs")
	}
}

func TestInstallFont_CustomURL_Success(t *testing.T) {
	withTempHome(t)
	srv := fakeFontZipServer(t, map[string]string{
		"TestFont-Regular.ttf": "fake regular ttf",
		"TestFont-Bold.ttf":    "fake bold ttf",
	})
	defer srv.Close()

	installed, err := InstallFont("my-custom-font", InstallOptions{
		URL:         srv.URL,
		AllVariants: true,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(installed) != 1 {
		t.Fatalf("expected 1 installed record, got %d", len(installed))
	}
	if len(installed[0].Files) != 2 {
		t.Errorf("expected 2 files installed, got %d", len(installed[0].Files))
	}
}

func TestInstallFont_TracksInstallState(t *testing.T) {
	withTempHome(t)
	srv := fakeFontZipServer(t, map[string]string{"Test-Regular.ttf": "data"})
	defer srv.Close()

	_, err := InstallFont("tracked-font", InstallOptions{URL: srv.URL, AllVariants: true})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !IsInstalled("tracked-font") {
		t.Error("expected font to be tracked as installed after InstallFont")
	}
}

func TestInstallFont_AlreadyInstalledWithoutForce(t *testing.T) {
	withTempHome(t)
	srv := fakeFontZipServer(t, map[string]string{"Test-Regular.ttf": "data"})
	defer srv.Close()

	InstallFont("dup-font", InstallOptions{URL: srv.URL, AllVariants: true})

	_, err := InstallFont("dup-font", InstallOptions{URL: srv.URL, AllVariants: true})
	if err == nil {
		t.Fatal("expected error reinstalling without --force, got nil")
	}
}

func TestInstallFont_ForceReinstalls(t *testing.T) {
	withTempHome(t)
	srv := fakeFontZipServer(t, map[string]string{"Test-Regular.ttf": "data"})
	defer srv.Close()

	InstallFont("force-font", InstallOptions{URL: srv.URL, AllVariants: true})

	_, err := InstallFont("force-font", InstallOptions{URL: srv.URL, AllVariants: true, Force: true})
	if err != nil {
		t.Fatalf("expected no error with --force, got: %v", err)
	}
}

func TestInstallFont_DryRunWritesNoFiles(t *testing.T) {
	withTempHome(t)
	srv := fakeFontZipServer(t, map[string]string{"Test-Regular.ttf": "data"})
	defer srv.Close()

	_, err := InstallFont("dryrun-font", InstallOptions{URL: srv.URL, AllVariants: true, DryRun: true})
	if err != nil {
		t.Fatalf("expected no error on dry run, got: %v", err)
	}
	if IsInstalled("dryrun-font") {
		t.Error("expected dry-run to not track font as installed")
	}
}

func TestInstallFont_UnknownCatalogName(t *testing.T) {
	withTempHome(t)
	_, err := InstallFont("not-a-real-catalog-font", InstallOptions{})
	if err == nil {
		t.Fatal("expected error for unknown catalog font with no URL, got nil")
	}
}

func TestInstallFont_DownloadSizeLimitEnforced(t *testing.T) {
	withTempHome(t)
	// build an entry that reports a large size via Content-Length
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := buildTestZip(t, map[string]string{"Big-Regular.ttf": "small actual content"})
		w.Header().Set("Content-Length", "999999999") // lie about size
		w.Write(data)
	}))
	defer srv.Close()

	_, err := InstallFont("oversized-font", InstallOptions{
		URL:              srv.URL,
		AllVariants:      true,
		MaxDownloadBytes: 1024, // 1KB cap, smaller than claimed content-length
	})
	if err == nil {
		t.Fatal("expected error when Content-Length exceeds MaxDownloadBytes, got nil")
	}
}

func TestRemoveFont_RemovesFilesAndState(t *testing.T) {
	withTempHome(t)
	srv := fakeFontZipServer(t, map[string]string{"Remove-Regular.ttf": "data"})
	defer srv.Close()

	installed, err := InstallFont("removable-font", InstallOptions{URL: srv.URL, AllVariants: true})
	if err != nil {
		t.Fatalf("setup: expected no error, got: %v", err)
	}
	filePath := installed[0].Files[0]

	if err := RemoveFont("removable-font"); err != nil {
		t.Fatalf("expected no error removing font, got: %v", err)
	}

	if IsInstalled("removable-font") {
		t.Error("expected font to no longer be tracked as installed")
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("expected font file to be deleted from disk")
	}
}

func TestRemoveFont_NotTracked(t *testing.T) {
	withTempHome(t)
	if err := RemoveFont("never-installed-font"); err == nil {
		t.Fatal("expected error removing untracked font, got nil")
	}
}

func TestListInstalledFonts_EmptyInitially(t *testing.T) {
	withTempHome(t)
	list, err := ListInstalledFonts()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 installed fonts initially, got %d", len(list))
	}
}

func TestListInstalledFonts_ReturnsInstalled(t *testing.T) {
	withTempHome(t)
	srv := fakeFontZipServer(t, map[string]string{"List-Regular.ttf": "data"})
	defer srv.Close()

	InstallFont("listed-font", InstallOptions{URL: srv.URL, AllVariants: true})

	list, err := ListInstalledFonts()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 installed font, got %d", len(list))
	}
	if list[0].Name != "listed-font" {
		t.Errorf("expected name 'listed-font', got %q", list[0].Name)
	}
}
