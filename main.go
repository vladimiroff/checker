package main

import (
	_ "expvar"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/koding/kite"
	"github.com/koding/logging"

	"github.com/vladimiroff/checker/handlers"
)

var expvarHost = "localhost:8123"

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		logging.Fatal("os.Hostname() failed: %s", err.Error())
	}

	k := kite.New(fmt.Sprintf("%s.checker", hostname), "1.0.0")
	k.PreHandleFunc(handlers.LogRequest)
	k.PostHandleFunc(handlers.LogResponse)
	k.HandleFunc("local_check", handlers.LocalCheck)
	k.HandleFunc("checkers", handlers.Checkers)

	go k.Run()
	<-k.ServerReadyNotify()

	sock, err := net.Listen("tcp", expvarHost)
	if err != nil {
		k.Log.Fatal("Registered expvar on %s failed with \"%s\"", expvarHost, err)
	}
	go func() {
		k.Log.Info("Registered expvar on %s", expvarHost)
		http.Serve(sock, nil)
	}()

	if _, err := k.Register(&url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%d/kite", k.Port()),
	}); err != nil {
		k.Log.Fatal("Can't register %s: %s", k.Kite().Name, err.Error())
	}

	<-k.ServerCloseNotify()
}
