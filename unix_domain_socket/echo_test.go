package echo

import (
    "bytes"
    "context"
    "fmt"
    "io/ioutil"
    "net"
    "os"
    "path/filepath"
    "testing"
)

// The streaming Unix domain socket works like TCP without the overhead assoiciated with TCP's
// acknowledgments, checksums, flow control, and so on. The operating system is responsible for
// implementing the streaming inter-process communication over Unix domain socket in lieu of TCP.
func TestEchoServer(t *testing.T) {
    dir, err := ioutil.TempDir("", "echo_unix")
    if err != nil {
        t.Fatal(err)
    }
    defer func() {
        if rErr := os.RemoveAll(dir); rErr != nil {
            t.Error(rErr)
        }
    }()

    ctx, cancel := context.WithCancel(context.Background())
    socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
    rAddr, err := streamingEchoServer(ctx, "unix", socket)
    if err != nil {
        t.Fatal(err)
    }
    defer cancel()

    err = os.Chmod(socket, os.ModeSocket|0666)
    if err != nil {
        t.Fatal(err)
    }

    conn, err := net.Dial("unix", rAddr.String())
    if err != nil {
        t.Fatal(err)
    }

    msg := []byte("ping")
    for i := 0; i < 3; i++ {
        _, err = conn.Write(msg)
        if err != nil {
            t.Fatal(err)
        }
    }

    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        t.Fatal(err)
    }

    expected := bytes.Repeat(msg, 3)
    if !bytes.Equal(expected, buf[:n]) {
        t.Fatalf("expected reply %q; actual reply %q", expected, buf[:n])
    }
}

// This echo server will comunicate using datagram based network types, such as udp and unixgram.
// Whether you're communicating over UDP or a unixgram socket, the server you'll write looks
// essentially the same. The difference is, you will need to clean up the socket file with a unixgram listener.
func TestEchoServerUnixDatagram(t *testing.T) {
    dir, err := ioutil.TempDir("", "echo_unixgram")
    if err != nil {
        t.Fatal(err)
    }

    defer func() {
        if rErr := os.RemoveAll(dir); rErr != nil {
            t.Error(rErr)
        }
    }()

    ctx, cancel := context.WithCancel(context.Background())
    sSocket := filepath.Join(dir, fmt.Sprintf("s%d.sock", os.Getpid()))
    serverAddr, err := datagramEchoServer(ctx, "unixgram", sSocket)
    if err != nil {
        t.Fatal(err)
    }
    defer cancel()

    err = os.Chmod(sSocket, os.ModeSocket|0622)
    if err != nil {
        t.Fatal(err)
    }

    cSocket := filepath.Join(dir, fmt.Sprintf("c%d.sock", os.Getpid()))
    client, err := net.ListenPacket("unixgram", cSocket)
    if err != nil {
        t.Fatal(err)
    }
    defer func() { _ = client.Close() }()

    err = os.Chmod(cSocket, os.ModeSocket|0622)
    if err != nil {
        t.Fatal(err)
    }

    msg := []byte("ping")
    for i := 0; i < 3; i++ {
        _, err = client.WriteTo(msg, serverAddr)
        if err != nil {
            t.Fatal(err)
        }
    }

    buf := make([]byte, 1024)
    for i := 0; i < 3; i++ {
        n, addr, err := client.ReadFrom(buf)
        if err != nil {
            t.Fatal(err)
        }

        if addr.String() != serverAddr.String() {
            t.Fatalf("received reply from %q instead of %q", addr, serverAddr)
        }

        if !bytes.Equal(msg, buf[:n]) {
            t.Fatalf("received reply %q; actual reply %q", msg, buf[:n])
        }
    }
}

// The sequence packet socket type is a hybrid that combines the session-oriented connections and reliability
// of TCP with clearly delineated datagrams of UDP.
// Discard unrequested data in each datagram. If you read 32 byte of a 50-byte datagram, the operating
// system discard the 18 unrequested bytes
func TestEchoServerUnixPackets(t *testing.T) {
    dir, err := ioutil.TempDir("", "echo_unixpacket")
    if err != nil {
        t.Fatal(err)
    }
    defer func() {
        if eErr := os.RemoveAll(dir); eErr != nil {
            t.Error(eErr)
        }
    }()

    ctx, cancel := context.WithCancel(context.Background())
    socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
    rAddr, err := streamingEchoServer(ctx, "unixpacket", socket)
    if err != nil {
        t.Fatal(err)
    }
    defer cancel()

    err = os.Chmod(socket, os.ModeSocket|0666)
    if err != nil {
        t.Fatal(err)
    }

    conn, err := net.Dial("unixpacket", rAddr.String())
    if err != nil {
        t.Fatal(err)
    }
    defer func() { _ = conn.Close() }()

    msg := []byte("ping")
    for i := 0; i < 3; i++ {
        _, err = conn.Write(msg)
        if err != nil {
            t.Fatal(err)
        }
    }

    // Read data from server
    buf := make([]byte, 1024)
    for i := 0; i < 3; i++ {
        n, err := conn.Read(buf)
        if err != nil {
            t.Fatal(err)
        }
        if !bytes.Equal(msg, buf[:n]) {
            t.Errorf("expected reply %q; actual reply %q", msg, buf[:n])
        }
    }
    // Discarding unread bytes
    //buf := make([]byte, 2)
    //for i := 0; i < 3; i++ {
    //n, err := conn.Read(buf)
    //if err != nil {
    //t.Fatal(err)
    //}
    //if !bytes.Equal(msg[:2], buf[:n]) {
    //t.Errorf("expected reply %q; actual reply %q", msg[:2], buf[:n])
    //}
    //}
}
