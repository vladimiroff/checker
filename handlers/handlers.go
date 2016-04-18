// Package handlers defines kite handlers for performing system checks.
package handlers

import (
	"fmt"
	"os"

	"github.com/koding/kite"
	"github.com/koding/logging"
)

// Instance is a default Kite instance setup with current hostname.
var Instance *kite.Kite

// init sets up the default Kite instance
func init() {
	hostname, err := os.Hostname()
	if err != nil {
		logging.Fatal("os.Hostname() failed: %s", err.Error())
	}
	Instance = NewKite(hostname)
}

// NewKite sets up and returns a Kite instance.
func NewKite(hostname string) *kite.Kite {
	k := kite.New(fmt.Sprintf("%s.checker", hostname), "1.0.0")
	k.PreHandleFunc(LogRequest)
	k.PostHandleFunc(LogResponse)
	k.HandleFunc("local_check", LocalCheck)
	k.HandleFunc("checkers", Checkers)

	return k
}
