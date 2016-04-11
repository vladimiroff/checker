package checks

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var (
	nativeCheck = NativeCheck{}
	notErrMsg   = "%s(\"%s\") is not expected to return an error. Got \"%s\", instead."
	errMsg      = "%s(\"%s\") is expected to return an error. Got nil, instead."
	notOkMsg    = "%s(\"%s\") is expected to return %t. Got %t instead."
)

// newTempFile creates new temporary file and returns its name.
// Make sure to delete this file, once done with it.
func newTempFile() (string, error) {
	tmpFile, err := ioutil.TempFile("", "checks_test")
	if err != nil {
		return "", err
	}

	fileInfo, err := tmpFile.Stat()
	if err != nil {
		return "", err
	}

	return fileInfo.Name(), nil
}

func TestFileExists(t *testing.T) {

	if ok, err := nativeCheck.FileExists("./native.go"); !ok {
		t.Errorf(
			"native.go is expected to be existing, but it's not. err: %s", err)
	}

	missingTempFile, err := newTempFile()
	if err != nil {
		t.Fatal("Can't create temp file: %s", err.Error())
	}
	os.Remove(missingTempFile)

	if ok, err := nativeCheck.FileExists(missingTempFile); ok {
		t.Errorf(
			"Temp file %s is expected to be missing, but it's not. err: %s",
			missingTempFile, err)
	}
}

func TestFileContains(t *testing.T) {
	var testArgs = []struct {
		path   string
		substr string
		result bool
		hasErr bool
	}{
		{"./native.go", "func", true, false},
		{"./native.go", "FileContains(path, ", true, false},
		{"./native.go", "Текст на български", false, false},
		{"./native.exe", "func", false, true},
	}

	for _, args := range testArgs {
		argsConcat := fmt.Sprintf("%s\", \"%s", args.path, args.substr)
		result, err := nativeCheck.FileContains(args.path, args.substr)
		if args.hasErr {
			if err == nil {
				t.Errorf(errMsg, "NativeCheck", argsConcat)
			}
		} else {
			if err != nil {
				t.Errorf(notErrMsg, "NativeCheck", argsConcat, err)
			}
		}
		if result != args.result {
			t.Errorf(notOkMsg, "NativeCheck", argsConcat, args.result, result)
		}
	}

}

func TestProcess(t *testing.T) {
	var testArgs = []struct {
		name   string
		result bool
	}{
		{"go", true},
		{"име на процес на български", false},
	}
	for _, args := range testArgs {
		result, err := nativeCheck.ProcessIsRunning(args.name)
		if err != nil {
			t.Errorf(notErrMsg, "ProcessIsRunning", args.name, err)
		}
		if result != args.result {
			t.Errorf(notOkMsg, "ProcessIsRunning", args.name, args.result, result)
		}
	}
}
