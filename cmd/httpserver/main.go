package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/itsjoeoui/httpfromtcp/cmd/httpserver/handlers"
	"github.com/itsjoeoui/httpfromtcp/internal/request"
	"github.com/itsjoeoui/httpfromtcp/internal/response"
	"github.com/itsjoeoui/httpfromtcp/internal/server"
)

const port = 42069

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
