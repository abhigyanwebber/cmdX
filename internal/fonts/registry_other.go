//go:build !windows

package fonts

// registerFont is a no-op on macOS/Linux. Fonts placed in the
// platform's per-user font directory (~/Library/Fonts on macOS,
// ~/.local/share/fonts on Linux) are discovered automatically by the
// font system without any registry-equivalent step — macOS's font
// system and Linux's fontconfig both scan these directories directly.
func registerFont(displayName, fontPath string) error {
	return nil
}

// unregisterFont is a no-op on macOS/Linux for the same reason.
func unregisterFont(displayName string) error {
	return nil
}

// registrySupported reports whether font registry registration is
// available on this platform.
func registrySupported() bool {
	return false
}
