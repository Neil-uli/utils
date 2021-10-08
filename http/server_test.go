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

  "github.com/ulicod3/utils/http/handlers"
)

func TestSimpleHttpServer(t *testing.T) {
  srv := &http.Server{
    Addr: "127.0.0.1:8443",
    Handler: http.TimeoutHandler(
      handlers.DefaultHandler(), 2*time.Minute, ""),
      IdleTimeout: 5 * time.Minute,
      ReadHeaderTimeout: time.Minute,
    }

    l, err := net.Listen("tcp", srv.Addr)
    if err != nil { t.Fatal(err) }

    go func() {
      // https://github.com/FiloSottile/mkcert/ to get a key pair for dev purposes only
      // for production you can use https://github.com/FiloSottile/mkcert/
      err := srv.ServeTLS(l, "cert.pem", "key.pem")
      if err != http.ErrServerClosed {
        t.Error(err)
      }
    }()

    testCases := []struct {
      method   string
      body     io.Reader
      code     int
      response string
    }{
      {http.MethodGet, nil, http.StatusOK, "Hello, friend!"},
      {http.MethodPost, bytes.NewBufferString("<world>"), http.StatusOK, "Hello, &lt;world&gt;!"},
      {http.MethodHead, nil, http.StatusMethodNotAllowed, ""},
    }

    client := new(http.Client)
    path := fmt.Sprintf("https://%s/", srv.Addr)

    for i, c := range testCases {
      r, err := http.NewRequest(c.method, path, c.body)
      if err != nil {
        t.Errorf("%d: %v", i, err)
        continue
      }

      resp, err := client.Do(r)
      if err != nil {
        t.Errorf("%d: %v", i, err)
        continue
      }

      if resp.StatusCode != c.code {
        t.Errorf("%d: unexpected status code: %q", i, resp.Status)
      }

      b, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        t.Errorf("%d: %v", i, err)
        continue
      }

      _ = resp.Body.Close()
      if c.response != string(b) {
        t.Errorf("%d: expected %q; actual %q", i, c.response, b)
      }
    } 
  }


