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

	buffer := make([]byte, 8)

	currentLineContent := ""
	for {
		n, err := fd.Read(buffer)
		if err != nil {
			if currentLineContent != "" {
				fmt.Printf("read: %s", currentLineContent)
			}

			if errors.Is(err, io.EOF) {
				return
			}

			log.Fatalf("failed to read from %s: %s\n", inputFilePath, err)
			panic(err)
		}

		chunks := strings.Split(string(buffer[:n]), "\n")

		seen := false
		for _, chunk := range chunks {
			if seen {
				fmt.Printf("read: %s\n", currentLineContent)
				currentLineContent = chunk
			} else {
				currentLineContent += chunk
				seen = true
			}
		}

	}
}
