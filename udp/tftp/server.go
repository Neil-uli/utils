package tftp

import (
    "bytes"
    "errors"
    "fmt"
    "log"
    "net"
    "time"
)

type Server struct {
    Payload []byte
    Retries uint8 // number of times to retry a failed transmission
    Timeout time.Duration // duration to wait for an acknowledgement
}

func (s Server) ListenAndServe(addr string) error {
    conn, err := net.ListenPacket("udp", addr)
    if err != nil { return err }

    defer func() { _ = conn.Close() }()

    log.Printf("Listening on %s ... \n", conn.LocalAddr())
    
    return s.Serve(conn)
}

func (s *Server) Serve(conn net.PacketConn) error {
    fmt.Println("s")
    return nil



}




