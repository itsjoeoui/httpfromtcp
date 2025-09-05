# httpfromtcp

A simple HTTP server implementing the HTTP/1.1 spec

## Setup

Make sure you have [Go](https://go.dev/doc/install) installed.

Then, clone the repository and navigate to the project directory:

```bash
just setup // https://github.com/casey/just
```

You can find some example handlers in `cmd/httpserver/handlers`,
and they are wired up in `cmd/httpserver/main.go`.

To run it, simply execute:

```bash
go run ./cmd/httpserver
```

Then try sending a request to `http://127.0.0.1:42069/`!

## References

- [RFC 9112 - HTTP/1.1](https://datatracker.ietf.org/doc/html/rfc9112)
- [RFC 9110 - HTTP Semantics](https://datatracker.ietf.org/doc/html/rfc9110)
- [From TCP to HTTP | Full Course by @ThePrimeagen](https://www.youtube.com/watch?v=FknTw9bJsXM)
