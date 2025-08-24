package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	listenAddress = "127.0.0.1:42069"
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

		currentLineContent := ""
		for {
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContent != "" {
					lines <- currentLineContent
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
					lines <- currentLineContent
					currentLineContent = chunk
				} else {
					currentLineContent += chunk
					seen = true
				}
			}
		}
	}()

	return lines
}
