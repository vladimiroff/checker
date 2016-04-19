package handlers

import (
	"fmt"
	"sync"

	"github.com/koding/kite"
	"github.com/koding/kite/dnode"
)

// RemoteCheckResult defines a check result returned from another kite.
type RemoteCheckResult struct {
	Address string
	Results map[string]CheckResult
	Error   error
}

// RemoteCheck takes a slice of machine addresses, checker and a check requests
// to be fan-outed to those machines.
func RemoteCheck(request *kite.Request) (interface{}, error) {
	var (
		wg         sync.WaitGroup
		resultChan = make(chan RemoteCheckResult)
		response   = make(map[string]RemoteCheckResult)
	)

	args, err := request.Args.SliceOfLength(3)
	if err != nil {
		return logAndFail(request, err)
	}

	machines, err := args[0].Slice()
	if err != nil {
		return logAndFail(request, err)
	}

	checkerName, err := args[1].String()
	if err != nil {
		return logAndFail(request, err)
	}

	checkRequests, err := args[2].Map()
	if err != nil {
		return logAndFail(request, err)
	}

	// dialCheck is dialing a remote machine with a check request and sends the
	// result through given result channel.
	dialCheck := func(ch chan<- RemoteCheckResult, addr string, r map[string]*dnode.Partial) {
		rawResult, err := dial(addr, checkerName, r)
		results := make(map[string]CheckResult)

		if err == nil {
			err = rawResult.Unmarshal(&results)
		}

		ch <- RemoteCheckResult{
			Address: addr,
			Results: results,
			Error:   err,
		}
		wg.Done()
	}

	// Fan out requests
	for _, machine := range machines {
		machineStr, err := machine.String()
		if err != nil {
			continue
		}
		wg.Add(1)
		go dialCheck(resultChan, machineStr, checkRequests)
	}

	// Wait for response from all machines and then close the result channel in
	// order to unblock the range over it below.
	go func(ch chan RemoteCheckResult) {
		wg.Wait()
		close(ch)
	}(resultChan)

	for result := range resultChan {
		response[result.Address] = result
	}

	request.Context.Set("response", response)
	return response, nil
}

func dial(addr, checker string, checks map[string]*dnode.Partial) (*dnode.Partial, error) {
	investigateWorker := Instance.NewClient(fmt.Sprintf("http://%s/kite", addr))
	investigateWorker.Dial()
	return investigateWorker.Tell("local_check", checker, checks)
}
