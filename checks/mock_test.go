package checks

import "testing"

func TestMockCheckFileExists(t *testing.T) {
	testFileExists(t, MockCheck{})
}

func TestMockCheckFileContains(t *testing.T) {
	testFileContains(t, MockCheck{})
}

func TestMockCheckProcessIsRunning(t *testing.T) {
	testProcessIsRunning(t, MockCheck{})
}
