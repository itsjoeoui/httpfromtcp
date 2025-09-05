# httpfromtcp

A simple HTTP server implementing the HTTP/1.1 spec

## Setup

Make sure you have [Go](https://go.dev/doc/install) and
[just](https://github.com/casey/just) installed.

Then, run the following command in the project's root folder to set things up:

```bash
just setup
```

You can find some example handlers in `cmd/httpserver/handlers`,
and they are wired up in `cmd/httpserver/main.go` (all the routes are here!).

To run it, simply execute:

```bash
just run
```

Then try sending a request to `http://127.0.0.1:42069/`!

More routes are available, such as:

```go
func handler(w *response.Writer, r *request.Request) {
  switch {
  case strings.HasPrefix(r.RequestLine.RequestTarget, "/yourproblem"):
    handlers.Handler400(w, r)
  case strings.HasPrefix(r.RequestLine.RequestTarget, "/myproblem"):
    handlers.Handler500(w, r)
  case strings.HasPrefix(r.RequestLine.RequestTarget, "/video"):
    handlers.HandlerVideo(w, r)
  case strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin"):
    handlers.HandlerHTTPBin(w, r)
  default:
    handlers.Handler200(w, r)
  }
}
```

## References

- [RFC 9112 - HTTP/1.1](https://datatracker.ietf.org/doc/html/rfc9112)
- [RFC 9110 - HTTP Semantics](https://datatracker.ietf.org/doc/html/rfc9110)
- [From TCP to HTTP | Full Course by @ThePrimeagen](https://www.youtube.com/watch?v=FknTw9bJsXM)
