package checks

// Checker defines all allowed checks.
type Checker interface {
	FileExists(path string) (bool, error)
	FileContains(path, substr string) (bool, error)
	ProcessIsRunning(name string) (bool, error)
}
