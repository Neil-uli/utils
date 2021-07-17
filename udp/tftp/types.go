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
