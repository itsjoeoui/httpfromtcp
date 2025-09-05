package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsjoeoui/httpfromtcp/internal/headers"
	"github.com/itsjoeoui/httpfromtcp/internal/request"
	"github.com/itsjoeoui/httpfromtcp/internal/response"
	"github.com/itsjoeoui/httpfromtcp/internal/server"
)

const port = 42069

func handler(w *response.Writer, r *request.Request) {
	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		err := w.WriteStatusLine(400)
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
	case "/myproblem":
		err := w.WriteStatusLine(500)
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
	default:
		err := w.WriteStatusLine(200)
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
