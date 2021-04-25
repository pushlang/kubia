package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var server *http.Server
var sigs chan os.Signal

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == `/sd` {
		w.WriteHeader(500)
		w.Write([]byte("Shutdown!\n"))
		sigs <- syscall.SIGINT
		return
	}
	w.Write([]byte("path: " + r.URL.Path + "\n"))
}

func main() {
	server = &http.Server{Addr: ":8080", Handler: new(HandlerFunc)}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\nawaiting signal")
	fmt.Println("\nsignal:", <-sigs)

	fmt.Println("exiting")
	shutdownServer()
	fmt.Println("main done")
}

func shutdownServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 10000000*time.Microsecond)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("\ndefer done")
	}
}
