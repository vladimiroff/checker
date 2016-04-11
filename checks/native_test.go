package checks

import "testing"

func TestNativeCheckFileExists(t *testing.T) {
	testFileExists(t, NativeCheck{})
}

func TestNativeCheckFileContains(t *testing.T) {
	testFileContains(t, NativeCheck{})
}

func TestNativeCheckProcessIsRunning(t *testing.T) {
	testProcessIsRunning(t, NativeCheck{})
}
