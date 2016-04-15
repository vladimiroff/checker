package handlers

import (
	"github.com/koding/kite"

	"github.com/vladimiroff/checker/checks"
)

var checkers = map[string]checks.Checker{
	"native": checks.NativeCheck{},
	"unix":   checks.UnixShellCheck{},
}

// Checkers returns all implementations of checks.Checker available.
func Checkers(request *kite.Request) (interface{}, error) {
	response := make([]string, 0, len(checkers))
	for checker := range checkers {
		response = append(response, checker)
	}
	request.Context.Set("response", response)
	return response, nil
}
