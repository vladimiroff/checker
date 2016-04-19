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

	for i, arg := range checksTable {
		var (
			results   map[string]RemoteCheckResult
			ok        bool
			checkName = fmt.Sprintf("check_%d", i)
			checks    = map[string]CheckRequest{checkName: arg.checks}
			expResult = map[string]CheckResult{checkName: arg.result}
		)

		rawArgs, _ := json.Marshal([]interface{}{
			[]string{"localhost:6666"}, arg.checker, checks})

		rawResult, err := RemoteCheck(&kite.Request{
			Args:      &dnode.Partial{Raw: rawArgs},
			LocalKite: kite.New("test", "0.0.0"),
			Context:   cache.NewLRU(2),
		})
		if err != nil {
			t.Errorf(errorMsg, arg.checker, checks, err)
		}

		results, ok = rawResult.(map[string]RemoteCheckResult)
		if !ok {
			t.Errorf(unmarshalError, arg.checker, checks, rawResult)
		}

		for _, result := range results {
			if arg.hasError && result.Error == nil {
				t.Errorf(noErrorMsg, arg.checker, checks)
				continue
			} else if !arg.hasError && result.Error != nil {
				t.Errorf(errorMsg, arg.checker, checks, result.Error)
				continue
			}

			if result.Error != nil {
				// We confirmed the error is correct and there's no need to
				// unmarshal and check nils.
				continue
			}

			if !reflect.DeepEqual(result.Results, expResult) {
				t.Errorf(diffResult, arg.checker, checks, result.Results, expResult)
			}
		}

	}
}
