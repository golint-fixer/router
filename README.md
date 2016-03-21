# pat [![Build Status](https://travis-ci.org/vinci-proxy/pat.png)](https://travis-ci.org/vinci-proxy/pat) [![GoDoc](https://godoc.org/github.com/vinci-proxy/pat?status.svg)](https://godoc.org/github.com/vinci-proxy/pat) [![Coverage Status](https://coveralls.io/repos/github/vinci-proxy/pat/badge.svg?branch=master)](https://coveralls.io/github/vinci-proxy/pat?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinci-proxy/pat)](https://goreportcard.com/report/github.com/vinci-proxy/pat)

Simple, idiomatic and fast pattern muxer designed for dynamic routing.

Originally based in [bmizerany/pat](https://github.com/bmizerany/pat).

## Installation

```bash
go get -u gopkg.in/vinci-proxy/pat.v0
```

## API

See [godoc reference](https://godoc.org/github.com/vinci-proxy/pat) for detailed API documentation.

## Examples

#### Router 

```go
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
```

#### Vhost like muxer router 

```go
package main

import (
  "fmt"
  "gopkg.in/vinci-proxy/mux.v0"
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

```

## License

[MIT](https://opensource.org/licenses/MIT).
