package tftp

import (
    "bytes"
    "encoding/binary"
    "errors"
    "strings"
)
// TFTP limits datagram packets to 516 bytes or fewer to avoid fragmentation.
const (
    DatagramSize = 516 // max supported size
    BlockSize = DatagramSize - 4 // 4-byte header, 2 of them is an operation code 
)

type OpCode uint64

const (
    OpRRQ OpCode = iota + 1 // Read request
    _ // no WRQ support
    OpData
    OpAck
    OpErr
)

type ErrCode uint16

const (
    ErrUnknown ErrCode = iota
    ErrNotFound
    ErrAccessViolation
    ErrDiskFull
    ErrIllegalOp
    ErrUnknownID
    ErrFileExists
    ErrNoUser
)


type ReadReq struct {
    Filename string
    Mode string
}

// Although not used by our server, a client would make use of this method
func (q ReadReq) MarshalBinary() ([]byte, error) {
    mode := "octet"
    if q.Mode != "" {
        mode = q.Mode
    }

    // operation code + filename + 0 byte + mode + 0 byte
    cap := 2 + 2 + len(q.Filename) + 1 + len(q.Mode) + 1

    b := new(bytes.Buffer)
    b.Grow(cap)

    err := binary.Write(b, binary.BigEndian, OpRRQ) // write operation code
    if err != nil {
        return nil, err
    }
}
