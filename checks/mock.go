package checks

import (
	"errors"
	"os"
)

// MockCheck is an implementation of Checker that can be used for tests.
type MockCheck struct{}

// FileExists mocks the checks if a file exists.
func (MockCheck) FileExists(path string) (bool, error) {
	switch path {
	case "/tmp/exists":
		fallthrough
	case "./native.go":
		return true, nil
	case "/tmp/unreachable":
		fallthrough
	case "./native.exe":
		return false, &os.PathError{
			Op:   "stat",
			Path: path,
			Err:  errors.New("no such file or directory"),
		}
	}
	return false, nil
}

// FileContains mocks the check if given sub-string is contained into a file.
func (u MockCheck) FileContains(path, substr string) (bool, error) {
	exists, err := u.FileExists(path)
	if err != nil || !exists {
		return exists, err
	}

	switch substr {
	case "exists":
		fallthrough
	case "func":
		fallthrough
	case "FileContains(path, ":
		return true, nil
	}
	return false, nil
}

// ProcessIsRunning mocks the checks if given process name is currently
// running.
func (MockCheck) ProcessIsRunning(name string) (bool, error) {
	switch name {
	case "go":
		return true, nil
	default:
		return false, nil
	}

}
