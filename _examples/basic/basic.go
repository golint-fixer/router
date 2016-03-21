package main

import (
	"fmt"
	"gopkg.in/vinci-proxy/pat.v0"
	"gopkg.in/vinci-proxy/vinci.v0"
	"net/http"
)

func main() {
	fmt.Printf("Server listening on port: %d\n", 3100)
	vs := vinci.NewServer(vinci.ServerOptions{Host: "localhost", Port: 3100})

	router := pat.New()
	router.Get("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, foo"))
	}))
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))

	vs.Use(router)
	vs.Forward("http://httpbin.org")

	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
