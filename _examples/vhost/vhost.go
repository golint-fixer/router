package main

import (
	"fmt"
	"gopkg.in/vinxi/mux.v0"
	"gopkg.in/vinxi/router.v0"
	"gopkg.in/vinxi/vinxi.v0"
	"net/http"
)

func main() {
	fmt.Printf("Server listening on port: %d\n", 3100)
	vs := vinxi.NewServer(vinxi.ServerOptions{Host: "localhost", Port: 3100})

	router := pat.New()
	router.Get("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, foo"))
	}))
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))

	// Create a host header multiplexer
	muxer := mux.Host("localhost:3100")
	muxer.Use(router)

	vs.Use(muxer)
	vs.Forward("http://httpbin.org")

	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
