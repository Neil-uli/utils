## Creating robust network applications by using the io package
In addition to interfaces common in go code, such as `io.Reader`  and `io.Writer`,
the `io` package provides several useful functions and utilities that make the creation
of robust network applications easy.
To proxy data between connections, log network traffic, and ping hosts when firewalls attempt
to keep you from doing so.

- io.Copy
- io.MultiWriter
- io.TeeReader

### Monitoring a Network Connection
The io package includes useful tools that allow you to do more with network data than just send 
and receive it using connections objects. For example, you could use io.MultiWriter to write a 
single payload to mutliple network connections. You could use io.TeeReader to log data read from 
a network connection.
