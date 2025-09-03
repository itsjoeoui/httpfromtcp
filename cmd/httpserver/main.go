package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsjoeoui/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port)
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
