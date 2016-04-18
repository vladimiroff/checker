package main

import (
	_ "expvar"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/vladimiroff/checker/handlers"
)

var expvarHost = "localhost:8123"

func main() {
	k := handlers.Instance

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
