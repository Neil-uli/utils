# utils-net


## [TFTP Server implementation](udp/tftp) 

Trivial File Transfer Protocol 
```sh
cd udp/tftp/example
go build tftp.go
sudo ./tftp
```

In other terminal use a TFTP client
```sh
tftp 
connect 127.0.0.1
binary
get file.jpeg
```

## [Unix Domain Sockets](unix_domain_socket) 

Stream based files can manage networks like tcp or udp but with a simple file.
Sockets for local interprcess communication that allows bidirectional data exhange
between processes running on the same machinel.
