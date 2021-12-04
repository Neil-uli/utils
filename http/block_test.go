package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func blockIndefinitely(w http.ResponseWriter, r *http.Request) {
	select {}
}

//func TestBlockIndefinitely(t *testing.T) {
//ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))
//_, _ = http.Get(ts.URL)
//t.Fatal("client did not indefinitely block")
//}

func TestBlockIndefinitelyWithTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatal(err)
		}
		return
	}
	_ = resp.Body.Close()
}

// Control the cancellation of the context using a timeout
/*
   ctx, cancel := context.WithCancel(context.Background())
   timer := time.AfterFunc(5*time.Second, cancel)
   // Make the HTTP req, read the response headers, etc.
   // ...
   // Add 5 more seconds before reading the response body
   timer.Reset(5*time.Second)
*/
