package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	fd, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("failed to open %s: %s\n", inputFilePath, err)
		panic(err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			log.Fatalf("failed to close %s: %s\n", inputFilePath, err)
			panic(err)
		}
	}()

	linesChan := getLinesChannel(fd)
	for line := range linesChan {
		fmt.Printf("read: %s\n", line)
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

				log.Fatalf("failed to read from %s: %s\n", inputFilePath, err)
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
