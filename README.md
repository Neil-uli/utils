# utils-net


## TFTP Server implementation 
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
