package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var server *http.Server
var sigs chan os.Signal
var body string

type HandlerFunc func(http.ResponseWriter, *http.Request)

func logFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func respError(e error, w http.ResponseWriter) {
	if e != nil {
		w.Write([]byte(e.Error() + "\n"))
		log.Println(e)
	}
}

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.URL.Path == "/data" {
			body, err := ioutil.ReadAll(r.Body)
			respError(err, w)

			err = ioutil.WriteFile("/var/data/content", body, 0644)
			respError(err, w)
		}
	} else {
		if r.URL.Path == "/" {
			_, addrs, err := net.LookupSRV("", "", "kubia")
			respError(err, w)

			text := ""
			for _, addr := range addrs {
				resp, err := http.Get("http://" + addr.Target + ":8080" + "/data")
				respError(err, w)

				body := []byte("")
				if resp != nil {
					body, err = ioutil.ReadAll(resp.Body)
					respError(err, w)
				} else {
					w.Write([]byte("resp == nil\n"))
				}

				text += fmt.Sprintf("%s: %s\n", addr.Target, string(body))
			}

			w.Write([]byte(text))
			return
		}

		if r.URL.Path == "/data" {
			content, err := ioutil.ReadFile("/var/data/content")

			if os.IsNotExist(err) {
				host, err := os.Hostname()
				respError(err, w)
				w.Write([]byte(host + ": No data posted yet"))
			}

			w.Write(content)
			return
		}

		if r.URL.Path == "/sd" {
			w.WriteHeader(500)
			w.Write([]byte("Server shutdown!\n"))
			log.Println("Server shutdown!")
			sigs <- syscall.SIGINT
			return
		}
	}
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
