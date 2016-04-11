package checks

import "testing"

func TestUnixShellCheckFileExists(t *testing.T) {
	testFileExists(t, UnixShellCheck{})
}

func TestUnixShellCheckFileContains(t *testing.T) {
	testFileContains(t, UnixShellCheck{})
}

func TestUnixShellCheckProcessIsRunning(t *testing.T) {
	testProcessIsRunning(t, UnixShellCheck{})
}
