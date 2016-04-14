package handlers

import (
	"github.com/koding/kite"
)

// LogRequest is generic pre-handler for logging each request.
func LogRequest(request *kite.Request) (interface{}, error) {
	request.LocalKite.Log.Info(
		"%s(%s) called by %s",
		request.Method,
		request.Args,
		request.Client,
	)

	return nil, nil
}

// LogResponse is generic post-handler for logging each response.
func LogResponse(request *kite.Request) (interface{}, error) {
	result, err := request.Context.Get("response")
	if err != nil {
		return nil, err
	}

	request.LocalKite.Log.Info(
		"%s(%s) returned \"%v\"",
		request.Method, request.Args, result)

	return result, nil
}

// logAndFail logs and returns given error.
func logAndFail(request *kite.Request, err error) (interface{}, error) {
	request.LocalKite.Log.Error(
		"%s(%s) failed with \"%s\"!",
		request.Method, request.Args, err)
	return nil, err
}
