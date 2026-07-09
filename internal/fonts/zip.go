package fonts

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// maxExtractedFiles caps how many files a single font archive may
	// contain. Nerd Font family zips typically have well under this
	// many variant files; this guards against zip bombs with huge
	// file counts.
	maxExtractedFiles = 500

	// maxSingleFileBytes caps the size of any one extracted file.
	// Font files (ttf/otf) are rarely more than a few MB each.
	maxSingleFileBytes = 20 << 20 // 20 MiB

	// maxTotalExtractedBytes caps the sum of all extracted file sizes,
	// independent of the compressed zip size, guarding against
	// decompression-bomb style zips that are small on disk but expand
	// to an enormous size.
	maxTotalExtractedBytes = 300 << 20 // 300 MiB
)

// fontFileExtensions are the file types worth extracting from a font
// archive. Everything else (README, LICENSE, changelog) is skipped.
var fontFileExtensions = map[string]bool{
	".ttf": true,
	".otf": true,
}

// extractFontFiles extracts font files (.ttf/.otf) from a zip archive
// into destDir, applying a variant filter and returning the list of
// extracted file paths.
//
// Every extracted path is validated to stay within destDir — this
// guards against zip-slip, where a malicious archive entry named like
// "../../../.bashrc" could otherwise write outside the intended
// directory. Total extracted size and file count are also capped to
// guard against decompression-bomb archives.
func extractFontFiles(zipPath string, destDir string, variantFilter string, allVariants bool) ([]string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("could not open font archive: %w", err)
	}
	defer r.Close()

	if len(r.File) > maxExtractedFiles {
		return nil, fmt.Errorf("font archive contains %d files, exceeding the limit of %d", len(r.File), maxExtractedFiles)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create destination directory: %w", err)
	}

	var extracted []string
	var totalBytes int64

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(f.Name))
		if !fontFileExtensions[ext] {
			continue // skip README, LICENSE, etc.
		}

		// flatten directory structure — Nerd Font zips are usually
		// flat already, but be defensive about nested paths
		baseName := filepath.Base(f.Name)

		if !allVariants && variantFilter != "" {
			if !strings.Contains(strings.ToLower(baseName), strings.ToLower(variantFilter)) {
				continue
			}
		}

		destPath, err := sanitizeExtractPath(destDir, baseName)
		if err != nil {
			return nil, fmt.Errorf("archive entry %q rejected: %w", f.Name, err)
		}

		if f.UncompressedSize64 > maxSingleFileBytes {
			return nil, fmt.Errorf("archive entry %q (%d bytes) exceeds the per-file limit of %d bytes", f.Name, f.UncompressedSize64, int64(maxSingleFileBytes))
		}

		totalBytes += int64(f.UncompressedSize64)
		if totalBytes > maxTotalExtractedBytes {
			return nil, fmt.Errorf("font archive's total extracted size exceeds the limit of %d bytes", int64(maxTotalExtractedBytes))
		}

		if err := extractOneFile(f, destPath); err != nil {
			return nil, fmt.Errorf("could not extract %q: %w", f.Name, err)
		}

		extracted = append(extracted, destPath)
	}

	if len(extracted) == 0 {
		if variantFilter != "" && !allVariants {
			return nil, fmt.Errorf("no font files matched variant filter %q in this archive", variantFilter)
		}
		return nil, fmt.Errorf("no .ttf or .otf files found in font archive")
	}

	return extracted, nil
}

// sanitizeExtractPath resolves name against destDir and verifies the
// result stays within destDir, rejecting path traversal attempts
// (e.g. "../../etc/passwd") and absolute paths.
func sanitizeExtractPath(destDir, name string) (string, error) {
	cleaned := filepath.Clean(name)

	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid entry path (path traversal attempt)")
	}

	// filepath.IsAbs alone isn't enough here: on Windows it only treats
	// drive-letter (C:\...) or UNC paths as absolute, so a POSIX-style
	// rooted entry like "/etc/passwd" — which filepath.Clean turns into
	// "\etc\passwd" on Windows — would slip past IsAbs there. Zip
	// archives use forward slashes regardless of the platform that
	// built them, so we explicitly reject anything rooted at the
	// separator in addition to IsAbs's drive/UNC check.
	if filepath.IsAbs(cleaned) || strings.HasPrefix(cleaned, string(filepath.Separator)) {
		return "", fmt.Errorf("invalid entry path (absolute path)")
	}

	full := filepath.Join(destDir, cleaned)
	destClean := filepath.Clean(destDir)

	if full != destClean && !strings.HasPrefix(full, destClean+string(filepath.Separator)) {
		return "", fmt.Errorf("resolved path escapes destination directory")
	}

	return full, nil
}

// extractOneFile copies a single zip entry to destPath.
//
// Writes go to a temporary file in the same directory first, then are
// moved into place via os.Rename rather than truncating destPath
// in-place. This matters specifically for reinstalls on Windows: font
// files are a common real-time-scan target for antivirus software
// (font parsers have a history of CVEs), and a fresh .ttf can be
// briefly memory-mapped by the scanner immediately after creation —
// truncating it in place during that window fails with
// ERROR_USER_MAPPED_FILE. Renaming a new file over the old one avoids
// touching the old file's mapped content at all, and the retry loop
// absorbs the scanner's typically sub-second hold.
func extractOneFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	dir := filepath.Dir(destPath)
	tmp, err := os.CreateTemp(dir, ".cmdx-font-*.tmp")
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	// io.LimitReader as defense-in-depth beyond the UncompressedSize64
	// check in extractFontFiles, in case of a mismatched zip header.
	limited := io.LimitReader(rc, maxSingleFileBytes+1)
	written, err := io.Copy(tmp, limited)
	closeErr := tmp.Close()
	if err != nil {
		os.Remove(tmpPath)
		return err
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		return closeErr
	}
	if written > maxSingleFileBytes {
		os.Remove(tmpPath)
		return fmt.Errorf("file exceeded size limit during extraction")
	}

	// Retry the finalize step with exponential backoff. Windows
	// antivirus real-time scanning can hold a transient lock on a
	// freshly-written font file (font parsers have a history of CVEs,
	// so they're a common eager-scan target), which can surface as
	// either "user-mapped section open" or plain "Access is denied"
	// depending on how the scanner opened the file. An explicit
	// os.Remove before each rename attempt is included as a fallback:
	// some lock types block rename but permit delete, or vice versa,
	// so trying both gives the best chance of succeeding before the
	// scan naturally releases the file (typically well under a second
	// for a small font file).
	delay := 100 * time.Millisecond
	var lastErr error
	for attempt := 0; attempt < 10; attempt++ {
		if err := os.Rename(tmpPath, destPath); err == nil {
			return nil
		} else {
			lastErr = err
		}

		// fallback: try removing the existing destination first, then
		// rename again immediately — some Windows lock states permit
		// delete but not replace-via-rename.
		os.Remove(destPath)
		if err := os.Rename(tmpPath, destPath); err == nil {
			return nil
		} else {
			lastErr = err
		}

		time.Sleep(delay)
		if delay < 2*time.Second {
			delay *= 2
		}
	}

	os.Remove(tmpPath)
	return fmt.Errorf("could not finalize file after retries (likely locked by antivirus scanning — try again in a moment): %w", lastErr)
}
