package tui

import (
	"errors"
	"testing"
	"time"
)

// Regression coverage for BUG-005: tea.NewProgram(...).Run() blocks
// indefinitely rather than erroring when stdin/stdout aren't a real
// terminal — which is exactly the environment `go test` itself runs
// in, so these tests exercise the real non-interactive branch without
// needing to fake or inject anything.

func TestIsInteractive_FalseUnderGoTest(t *testing.T) {
	// go test's stdin/stdout are not a real terminal, so this should
	// always be false in CI/normal test runs. If this ever starts
	// failing, it likely means the test is being run attached to a
	// real terminal, which is an unusual way to run `go test`.
	if IsInteractive() {
		t.Skip("stdin/stdout appear to be a real terminal in this test run — skipping, since the behavior this test guards only applies non-interactively")
	}
}

func TestRunSpinner_NonInteractiveRunsTaskAndReturnsQuickly(t *testing.T) {
	called := false
	done := make(chan error, 1)

	go func() {
		done <- RunSpinner("testing", "#FFFFFF", func() error {
			called = true
			return nil
		})
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("RunSpinner did not return within 3s — this is exactly the BUG-005 hang this test guards against")
	}

	if !called {
		t.Error("expected the wrapped task to have been called")
	}
}

func TestRunSpinner_NonInteractivePropagatesTaskError(t *testing.T) {
	wantErr := errors.New("boom")
	done := make(chan error, 1)

	go func() {
		done <- RunSpinner("testing", "#FFFFFF", func() error {
			return wantErr
		})
	}()

	select {
	case err := <-done:
		if !errors.Is(err, wantErr) {
			t.Errorf("expected task's error to propagate, got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("RunSpinner did not return within 3s")
	}
}

func TestRunThemeList_NonInteractiveReturnsErrorQuickly(t *testing.T) {
	done := make(chan error, 1)

	go func() {
		_, err := RunThemeList([]ThemeItem{{Name: "test"}}, "#FFFFFF", "#FFFFFF", "#FFFFFF")
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Error("expected a non-nil error when not running interactively (so callers fall back to a plain listing), got nil")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("RunThemeList did not return within 3s — this is exactly the BUG-005 hang this test guards against")
	}
}
