package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		linesChan := getLinesChannel(conn)
		for line := range linesChan {
			fmt.Printf("read: %s\n", line)
		}

	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)

		buffer := make([]byte, 8)

		var sb strings.Builder
		for {
			n, err := f.Read(buffer)
			if err != nil {
				if sb.String() == "" {
					lines <- sb.String()
				}

				if errors.Is(err, io.EOF) {
					break
				}

				log.Fatalf("failed to read: %s\n", err.Error())
				panic(err)
			}

			chunks := strings.Split(string(buffer[:n]), "\n")

			seen := false
			for _, chunk := range chunks {
				if seen {
					lines <- sb.String()
					sb.Reset()
					sb.WriteString(chunk)
				} else {
					sb.WriteString(chunk)
					seen = true
				}
			}
		}
	}()

	return lines
}
