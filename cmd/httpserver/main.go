package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsjoeoui/httpfromtcp/internal/request"
	"github.com/itsjoeoui/httpfromtcp/internal/server"
)

const port = 42069

func handler(w io.Writer, r *request.Request) *server.HandlerError {
	switch r.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: 400,
			Message:    "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: 500,
			Message:    "Woopsie, my bad\n",
		}
	default:
		_, err := fmt.Fprintf(w, "All good, frfr\n")
		if err != nil {
			log.Printf("Failed to write response: %v", err)
		}
		return nil
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
