package handlers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/koding/cache"
	"github.com/koding/kite"
	"github.com/koding/kite/dnode"
)

func runLocalKite() chan struct{} {
	final := make(chan struct{})
	k := Instance
	k.Config.Port = 6666

	go k.Run()

	go func() {
		<-final
		k.Close()
	}()

	return final
}

func TestRemoteCheckInvalidCalls(t *testing.T) {
	orig := mockCheckers()
	defer unmockCheckers(orig)

	args := [][]interface{}{
		[]interface{}{[]string{"localhost:6666"}, "too", "much", "arguments"},
		[]interface{}{[]string{"localhost:6666"}, "native", "this is not a map"},
		[]interface{}{"localhost:6666", 42, "this is not a map"},
		[]interface{}{[]string{"localhost:6666"}, 42, map[string]CheckRequest{"0": {}}},
	}

	for _, arg := range args {

		rawArgs, _ := json.Marshal(arg)

		_, err := RemoteCheck(&kite.Request{
			Args:      &dnode.Partial{Raw: rawArgs},
			LocalKite: kite.New("test", "0.0.0"),
			Context:   cache.NewLRU(2),
		})

		if err == nil {
			t.Errorf(noErrorMsg, arg[0], arg[1:len(arg)])
		}
	}
}

func TestRemoteCheck(t *testing.T) {
	orig := mockCheckers()
	defer unmockCheckers(orig)

	final := runLocalKite()
	defer close(final)

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
	}

	for i, arg := range argsTable {
		var (
			remoteResult map[string]*dnode.Partial
			result       map[string]CheckResult
			ok           bool
			checkName    = fmt.Sprintf("check_%d", i)
			checks       = map[string]CheckRequest{checkName: arg.checks}
			expResult    = map[string]CheckResult{checkName: arg.result}
		)

		rawArgs, _ := json.Marshal([]interface{}{[]string{"localhost:6666"}, arg.checker, checks})

		rawResult, err := RemoteCheck(&kite.Request{
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

		remoteResult, ok = rawResult.(map[string]*dnode.Partial)
		if !ok {
			t.Errorf(unmarshalError, arg.checker, checks, rawResult)
			continue
		}

		for _, rResult := range remoteResult {
			err := rResult.Unmarshal(&result)
			if err != nil {
				t.Errorf(unmarshalError, arg.checker, checks, rResult)
			}

			if !reflect.DeepEqual(result, expResult) {
				t.Errorf(diffResult, arg.checker, checks, result, expResult)
			}
		}

	}
}
