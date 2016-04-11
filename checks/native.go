package checks

import (
	"bufio"
	"os"
	"strings"

	"github.com/mitchellh/go-ps"
)

// NativeCheck implements the Checker interface using only NativeCheck Go implementation
// (i.e. without shelling out) which is intended to be platform-independed.
type NativeCheck struct{}

// FileExists checks if a file exists using Go's os package.
func (NativeCheck) FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	return !os.IsNotExist(err), err
}

// FileContains checks if given sub-string is contained into a file. Current
// implementation uses bufio.Scanner for that purpose, reading files line by
// line.
func (NativeCheck) FileContains(path, substr string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), substr) {
			return true, nil
		}
	}

	return false, nil
}

// ProcessIsRunning checks if given process name is currently
// running.
//
// - Darwin uses the sysctl syscall to retrieve the process table,
// via cgo.
//
// - On Unix uses the procfs at /proc to inspect the process tree.
//
// - On Windows uses the Windows API, and methods such as
// CreateToolhelp32Snapshot to get a point-in-time snapshot of the
// process table.
func (NativeCheck) ProcessIsRunning(name string) (bool, error) {
	procs, err := ps.Processes()
	if err != nil {
		return false, err
	}

	for _, proc := range procs {
		if proc.Executable() == name {
			return true, nil
		}
	}

	return false, nil
}
