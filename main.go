package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/koding/kite"
	"github.com/koding/logging"

	"github.com/vladimiroff/checker/handlers"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		logging.Fatal("os.Hostname() failed: %s", err.Error())
	}

	k := kite.New(fmt.Sprintf("%s.checker", hostname), "1.0.0")
	k.HandleFunc("check", handlers.LocalCheck).DisableAuthentication()
	k.HandleFunc("checkers", handlers.Checkers)

	go k.Run()
	<-k.ServerReadyNotify()

	if _, err := k.Register(&url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%d/kite", k.Port()),
	}); err != nil {
		k.Log.Fatal("Can't register %s: %s", k.Kite().Name, err.Error())
	}

	<-k.ServerCloseNotify()
}
