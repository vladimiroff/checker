// Package checks defines a Checker and several implementations of it for
// performing system checks regardless of how are they being used.
package checks

// Checker defines all allowed checks.
type Checker interface {
	FileExists(path string) (bool, error)
	FileContains(path, substr string) (bool, error)
	ProcessIsRunning(name string) (bool, error)
}
