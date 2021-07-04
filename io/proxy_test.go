package main

import (
    "io"
    "net"
    "sync"
    "testing"
)

func proxy(from io.Reader, to io.Writer) error {
    fromWriter, fromIsWriter := from.(io.Writer)
    toReader, toIsReader := to.(io.Reader)

    if toIsReader && fromIsWriter {
        // Send replies since "from" and "to" implement the necessary interfaces
        go func() { _, _ = io.Copy(fromWriter, toReader) }()
    }
    _, err := io.Copy(to, from)
    return err
}

func TestProxy(t *testing.T) {
    var wg sync.WaitGroup

    // server listens for a "ping" message and responds with a "pong" message. All
    // other messages are echoed back to the client
    server, err := net.Listen("tcp", "127.0.0.1:")
    if err != nil {
        t.Fatal(err)
    }
    wg.Add(1)

    go func() {
        defer wg.Done()

    }()
}
