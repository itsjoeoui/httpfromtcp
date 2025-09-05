package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/itsjoeoui/httpfromtcp/internal/headers"
	"github.com/itsjoeoui/httpfromtcp/internal/request"
	"github.com/itsjoeoui/httpfromtcp/internal/response"
	"github.com/itsjoeoui/httpfromtcp/internal/server"
)

const port = 42069

func handler500(w *response.Writer, _ *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeInternalServerError)
	if err != nil {
		log.Printf("Failed to write status line: %v", err)
	}

	f, err := os.ReadFile("./cmd/httpserver/templates/500.html")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
	}

	h := response.GetDefaultHeaders(len(f))
	h.Override(headers.ContentTypeHeader, "text/html")
	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("Failed to write headers: %v", err)
	}

	_, err = w.WriteBody(f)
	if err != nil {
		log.Printf("Failed to write body: %v", err)
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeBadRequest)
	if err != nil {
		log.Printf("Failed to write status line: %v", err)
	}

	f, err := os.ReadFile("./cmd/httpserver/templates/400.html")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
	}

	h := response.GetDefaultHeaders(len(f))
	h.Override(headers.ContentTypeHeader, "text/html")
	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("Failed to write headers: %v", err)
	}

	_, err = w.WriteBody(f)
	if err != nil {
		log.Printf("Failed to write body: %v", err)
	}
}

func handlerVideo(w *response.Writer, _ *request.Request) {
	err := w.WriteStatusLine(response.StatusCodeOK)
	if err != nil {
		log.Printf("Failed to write status line: %v", err)
	}

	f, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		log.Printf("Failed to read file: %v, you can download it with 'just setup'", err)
	}

	h := response.GetDefaultHeaders(len(f))
	h.Override(headers.ContentTypeHeader, "video/mp4")
	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("Failed to write headers: %v", err)
	}

	_, err = w.WriteBody(f)
	if err != nil {
		log.Printf("Failed to write body: %v", err)
	}
}

func handler(w *response.Writer, r *request.Request) {
	switch {
	case strings.HasPrefix(r.RequestLine.RequestTarget, "/yourproblem"):
		handler400(w, r)
	case strings.HasPrefix(r.RequestLine.RequestTarget, "/myproblem"):
		handler500(w, r)
	case strings.HasPrefix(r.RequestLine.RequestTarget, "/video"):
		handlerVideo(w, r)
	case strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin"):
		route := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin")
		resp, err := http.Get(fmt.Sprintf("https://httpbin.org%s", route))
		if err != nil {
			log.Printf("Failed to fetch from httpbin: %v", err)
			return
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Printf("Failed to close httpbin response body: %v", err)
			}
		}()

		err = w.WriteStatusLine(response.StatusCodeOK)
		if err != nil {
			log.Printf("Failed to write status line: %v", err)
		}

		h := response.GetDefaultHeaders(0)
		h.Remove(headers.ContentLengthHeader)
		h.Override(headers.TransferEncodingHeader, "chunked")
		h.Set(headers.TrailerHeader, headers.XContentLengthHeader)
		h.Set(headers.TrailerHeader, headers.XContentSHA256)

		err = w.WriteHeaders(h)
		if err != nil {
			log.Printf("Failed to write headers: %v", err)
		}

		fullBody := make([]byte, 0)

		buffer := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				_, writeErr := w.WriteChunkedBody(buffer[:n])
				if writeErr != nil {
					log.Printf("Failed to write chunked body: %v", writeErr)
				}

				fullBody = append(fullBody, buffer[:n]...)
			}
			if err != nil {
				if err != io.EOF {
					log.Printf("Failed to read from httpbin response body: %v", err)
				}
				break
			}
		}

		_, err = w.WriteChunkedBodyDone()
		if err != nil {
			log.Printf("Failed to write chunked body done: %v", err)
		}

		trailers := headers.NewHeaders()
		sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
		trailers.Override(headers.XContentSHA256, sha256)
		trailers.Override(headers.XContentLengthHeader, fmt.Sprintf("%d", len(fullBody)))
		err = w.WriteTrailers(trailers)
		if err != nil {
			log.Printf("Failed to write trailers: %v", err)
		}

	default:
		err := w.WriteStatusLine(response.StatusCodeOK)
		if err != nil {
			log.Printf("Failed to write status line: %v", err)
		}

		f, err := os.ReadFile("./cmd/httpserver/templates/200.html")
		if err != nil {
			log.Printf("Failed to read file: %v", err)
		}

		h := response.GetDefaultHeaders(len(f))
		h.Override(headers.ContentTypeHeader, "text/html")
		err = w.WriteHeaders(h)
		if err != nil {
			log.Printf("Failed to write headers: %v", err)
		}

		_, err = w.WriteBody(f)
		if err != nil {
			log.Printf("Failed to write body: %v", err)
		}
	}
}

func main() {
	server, err := server.Serve(handler, port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer func() {
		err := server.Close()
		if err != nil {
			log.Fatalf("Failed to close server: %v", err)
		}
	}()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
