package checks

import (
	"os/exec"
	"syscall"
)

// UnixShellCheck implements the Checker interface using explicitly shell
// commands intended to work on Unix machines only.
//
// The nature of requied checks are easily doable with simple commands while
// check result being deducable simply by checking the exit status.
type UnixShellCheck struct{}

// exec executes a command and returns true if the command's exit status is 0.
func (UnixShellCheck) exec(command string, args ...string) (bool, error) {
	cmd := exec.Command(command, args...)
	if err := cmd.Start(); err != nil {
		return false, err
	}

	err := cmd.Wait()
	result := err == nil
	if exiterr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0. If it's 1 it should
		// be reported as a non-error in order to keep the same behaviour
		// between implementations.
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			if status.ExitStatus() == 1 {
				err = nil
			}
		}
	}
	return result, err
}

// FileExists checks if a file exists using the `stat` command.
func (u UnixShellCheck) FileExists(path string) (bool, error) {
	return u.exec("stat", path)
}

// FileContains checks if given sub-string is contained into a file, using
// `grep`  and returns true if the status code of it is 0 (zero).
func (u UnixShellCheck) FileContains(path, substr string) (bool, error) {
	return u.exec("grep", substr, path)
}

// ProcessIsRunning checks if given process name is currently
// running using `ps` .
func (u UnixShellCheck) ProcessIsRunning(name string) (bool, error) {
	return u.exec("ps", "-C", name)

}
