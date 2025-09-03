package main

import (
	"fmt"
	"log"
	"net"

	"github.com/itsjoeoui/httpfromtcp/internal/request"
)

// `nc localhost 42069` to establish a connection
const (
	listenAddress = ":42069"
)

func main() {
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("failed to listen on %s: %s\n", listenAddress, err.Error())
		panic(err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("failed to close listener: %s\n", err.Error())
			panic(err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("failed to accept connection: %s\n", err.Error())
		}
		log.Printf("accepted connection from %s\n", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed to close connection: %s\n", err.Error())
			panic(err)
		}
		log.Printf("closed connection from %s\n", conn.RemoteAddr().String())
	}()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Request line:\n")
	fmt.Printf("- Method: %s\n", req.RequestLine.Method)
	fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", req.RequestLine.HTTPVersion)
	fmt.Printf("Headers:\n")
	for k, v := range req.Headers {
		fmt.Printf("- %s: %s\n", k, v)
	}
}
