// Package handlers defines kite handlers for performing system checks.
package handlers

import (
	"github.com/koding/kite"
	"github.com/koding/kite/dnode"

	"github.com/vladimiroff/checker/checks"
)

// CheckRequest is used to describe every single check to be performed.
type CheckRequest struct {
	// Name holds the process name expected to be running when `Type` is
	// `"process_is_running"`.
	Name string `json:"name",omitempty`
	// Path holds absolute path to file expected to be existing or containing a
	// substring. Used only when `Type` is `"file_exists`" or
	// `"file_contains"`.
	Path string `json:"path",omitempty`
	// Type holds the type of performed check which essentially defines which
	// Checker method to be executed. Possible values are: `"file_exists`",
	// `"file_contains"` and `"process_is_running"`.
	Type string `json:"type"`
	// Substr is substring expected to be found in given file. Used only when
	// `Type` is `"file_contains"`.
	Substr string `json:"substring",omitempty`
}

// CheckResult holds the result of given check.
type CheckResult struct {
	Result bool  `json:"result"`
	Error  error `json:"error"`
}

// LocalCheck performs given checks with a selected checker.
//
// It expects two arguments:
//
// - Name of the check, which can be gathered from the Chercker handler.
//
// - JSON object with list of objects, unmarshable to map[string]CheckRequest.
func LocalCheck(request *kite.Request) (interface{}, error) {
	args, err := request.Args.SliceOfLength(2)
	if err != nil {
		return logAndFail(request, err)
	}

	// No need to check here if the first arguments is a valid string, because
	// it will fail anyways once we try to find that checker.
	checkerName, _ := args[0].String()
	rawChecks := args[1]

	checker, ok := checkers[checkerName]
	if !ok {
		return logAndFail(request, NoSuchChecker)
	}

	rawCheckRequests, err := rawChecks.Map()
	if err != nil {
		return logAndFail(request, err)
	}

	response := make(map[string]CheckResult)
	resultChans := make(map[string]chan *CheckResult)

	for name, rawCheckRequest := range rawCheckRequests {

		resultChans[name] = make(chan *CheckResult)

		go func(checker checks.Checker, raw *dnode.Partial, ch chan *CheckResult) {
			var request = &CheckRequest{}

			if err := raw.Unmarshal(request); err != nil {
				ch <- &CheckResult{Result: false, Error: err}
				return
			}

			checkResult, checkError := check(checker, request)
			ch <- &CheckResult{Result: checkResult, Error: checkError}
		}(checker, rawCheckRequest, resultChans[name])
	}

	for name, ch := range resultChans {
		response[name] = *<-ch
	}

	request.Context.Set("response", response)
	return response, nil
}

// check performs the given check request against a checker.
func check(checker checks.Checker, checkRequest *CheckRequest) (bool, error) {
	switch checkRequest.Type {
	case "file_exists":
		return checker.FileExists(checkRequest.Path)
	case "file_contains":
		return checker.FileContains(checkRequest.Path, checkRequest.Substr)
	case "process_is_running":
		return checker.ProcessIsRunning(checkRequest.Name)
	}

	return false, NoSuchCheckType
}
