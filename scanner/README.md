# Scanner in rust

## Summary
* Multithreading should be prefered when the program is CPU bound, async-await
when the program is I/O bound.
* Don't block the event loop
* Streams are async iterators
* Streams replace worker pools
* Always limit the number of concurrent tasks or the size of channels not to exhaust
resources
* If you are nesting async blocks, you are probably doing something wrong.


Software used to map attack surfaces is called scanners. Port, bulnerability,
sudomain, SQL injection scanner... They automate the long and fastidious task
that reconnaissance can be and prevent human errors (like forgetting a subdomain
or a server).

Keep in mind that scanners can be very noisy and thus may reveal your intentions,
be blocked by anti-spam systems or report incomlete data.

We will start with a simple scanner whoses purpose is to find subdomains of a target
and then will scan the most common ports for each subdomain. Then. as we go along,
we will add more features to find more interesting stuff, the automated way.

## Sharing data
You may want to share data between your tasks. As each task can be executed in a
different thread (processor), sharing data between `async` tasks are subject to the
same rules as sharing data between threads.

```sh
cargo run -- host
```

### Subdomains
We can check if the domains resolve by tunning the subdomains into a Stream. Thanks
to the combinators, easy to read.

```rust
let subdomains: Vec<Subdomain> = stream::iter(subdomains.into_iter())
    .map(|domain| Subdomain {
        domain,
        open_ports: Vec::new(),
    })
    .filter_map(|subdomain| {
        let dns_resolver = dns_resolver.clone();
        async move {
            if resolves(&dns_resolver, &subdomain).await {
                Some(subdomain)
        } else {
            None
        }
    }
    })
    .collect()
    .await;
    Ok(subdomains)
}

pub async fn resolves(dns_resolver: &DnsResolver, domain: &Subdomain) -> bool {
    dns_resolver.lookup_ip(domain.domain.as_str()).await.is_ok()
}
```

### Timeouts
When scanning a single port, we need a timeout `tokio::time::timeout` returns a
`Future<Output=Result>` we need to check that both the Result of TcpStream::connect
and timeout are Ok to be sure that the port is open.
```rust
// matches is a shortcut for Ok =>, Err =>
let is_open = matches!(
    tokio::time::timeout(timeout, TcpStream::connect(&socket_address)).await,
    Ok(Ok(_)),
);
```
