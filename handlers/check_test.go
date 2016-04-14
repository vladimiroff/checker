package handlers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/koding/cache"
	"github.com/koding/kite"
	"github.com/koding/kite/dnode"

	"github.com/vladimiroff/checker/checks"
)

const (
	noErrorMsg     = "LocalCheck(\"%s\", \"%s\") is expected to return an error, got nil instead"
	errorMsg       = "LocalCheck(\"%s\", \"%s\") is not expected to return an error, got \"%s\" instead"
	unmarshalError = "LocalCheck(\"%s\", \"%s\") result unmarshal error: %#v"
	diffResult     = "LocalCheck(\"%s\", \"%s\") is expected to be %v, got %v instead"
)

func mockCheckers() map[string]checks.Checker {
	orig := make(map[string]checks.Checker)
	for checker := range checkers {
		orig[checker] = checkers[checker]
		checkers[checker] = checks.MockCheck{}
	}
	return orig
}

func unmockCheckers(orig map[string]checks.Checker) {
	checkers = orig
}

func TestLocalCheckInvalidCalls(t *testing.T) {
	orig := mockCheckers()
	defer unmockCheckers(orig)

	args := [][]interface{}{
		[]interface{}{"too", "much", "arguments"},
		[]interface{}{"native", "this is not a map"},
	}

	for _, arg := range args {

		rawArgs, _ := json.Marshal(arg)

		_, err := LocalCheck(&kite.Request{
			Args:      &dnode.Partial{Raw: rawArgs},
			LocalKite: kite.New("test", "0.0.0"),
			Context:   cache.NewLRU(2),
		})

		if err == nil {
			t.Errorf(noErrorMsg, arg[0], arg[1:len(arg)])
		}
	}
}

func TestLocalCheck(t *testing.T) {
	orig := mockCheckers()
	defer unmockCheckers(orig)

	argsTable := []struct {
		hasError bool
		checker  string
		checks   CheckRequest
		result   CheckResult
	}{
		{
			false,
			"native",
			CheckRequest{Path: "./native.go", Type: "file_exists"},
			CheckResult{true, nil},
		},
		{
			false,
			"unix",
			CheckRequest{Path: "./native.go", Type: "file_exists"},
			CheckResult{true, nil},
		},
		{
			false,
			"native",
			CheckRequest{Path: "./native.go", Type: "file_contains", Substr: "func"},
			CheckResult{true, nil},
		},
		{
			false,
			"native",
			CheckRequest{Name: "go", Type: "process_is_running"},
			CheckResult{true, nil},
		},
		{
			false,
			"native",
			CheckRequest{Name: "go", Type: "foo"},
			CheckResult{false, NoSuchCheckType},
		},
		{
			true,
			"imaginary",
			CheckRequest{Path: "./native.go", Type: "file_exists"},
			CheckResult{},
		},
	}

	for i, arg := range argsTable {
		var (
			result    map[string]CheckResult
			ok        bool
			checkName = fmt.Sprintf("check_%d", i)
			checks    = map[string]CheckRequest{checkName: arg.checks}
			expResult = map[string]CheckResult{checkName: arg.result}
		)

		rawArgs, _ := json.Marshal([]interface{}{arg.checker, checks})

		rawResult, err := LocalCheck(&kite.Request{
			Args:      &dnode.Partial{Raw: rawArgs},
			LocalKite: kite.New("test", "0.0.0"),
			Context:   cache.NewLRU(2),
		})

		if arg.hasError && err == nil {
			t.Errorf(noErrorMsg, arg.checker, checks)
			continue
		} else if !arg.hasError && err != nil {
			t.Errorf(errorMsg, arg.checker, checks, err.Error())
			continue
		}

		if rawResult == nil {
			continue
		}

		result, ok = rawResult.(map[string]CheckResult)
		if !ok {
			t.Errorf(unmarshalError, arg.checker, checks, rawResult)
			continue
		}

		if !reflect.DeepEqual(result, expResult) {
			t.Errorf(diffResult, arg.checker, checks, arg.result, expResult)
		}
	}
}
