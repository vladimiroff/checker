package handlers

// Error is a generic error type used inside this package. It's essentially a
// string allowing its instances to be constants.
//
// Idea of Dave Cheney: http://dave.cheney.net/2016/04/07/constant-errors
type Error string

func (e Error) Error() string { return string(e) }

// NoSuchChecker is used when the requested implementation of
// checks.Checker is not available. The user is expected to call the Checkers()
// endpoint.
const NoSuchChecker = Error("no such checker")

// NoSuchCheckType is used when the Type value of a check request is not among
// supported.
const NoSuchCheckType = Error("no such check type")

// InvalidRemoteResponse is used when a result from remote check can't be
// unmarshalled to a boolean value.
const InvalidRemoteResponse = Error("error check gave invalid response")
