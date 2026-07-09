//go:build windows

package fonts

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// registerFont adds a per-user font registry entry so Windows
// applications recognize the font by name without requiring a login.
// Per-user fonts placed in the special Fonts folder are picked up by
// most apps automatically, but the registry entry ensures GDI-based
// name lookups (used by some older or native Win32 apps) also work.
func registerFont(displayName, fontPath string) error {
	k, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows NT\CurrentVersion\Fonts`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("could not open font registry key: %w", err)
	}
	defer k.Close()

	valueName := displayName + " (TrueType)"
	if err := k.SetStringValue(valueName, fontPath); err != nil {
		return fmt.Errorf("could not write font registry value: %w", err)
	}

	return nil
}

// unregisterFont removes a font's registry entry.
func unregisterFont(displayName string) error {
	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows NT\CurrentVersion\Fonts`,
		registry.SET_VALUE,
	)
	if err != nil {
		// key not existing is not an error condition worth failing on
		return nil
	}
	defer k.Close()

	valueName := displayName + " (TrueType)"
	// DeleteValue returns an error if the value doesn't exist — that's
	// fine, it just means this font was never registered.
	_ = k.DeleteValue(valueName)

	return nil
}

// registrySupported reports whether font registry registration is
// available on this platform.
func registrySupported() bool {
	return true
}
