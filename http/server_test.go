package main 

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "net"
    "net/http"
    "testing"
    "time"
)

func TestSimpleHttpServer(t *testing.T) {
    srv := &http.Server{
        Addr: "127.0.0.1:8080",
        Handler: http.TimeoutHandler(
            handlers.DefaultHandler(), 2*time.Minute, ""),
        IdleTimeout: 5 * time.Minute,
        ReadHeaderTimeout: time.Minute,
    }

    l, err := net.Listen("tcp", srv.Addr)
    if err != nil { t.Fatal(err) }

    go func() {
        err := srv.Serve(l)
        if err != http.ErrServerClosed {
            t.Error(err)
        }
    }()
}


