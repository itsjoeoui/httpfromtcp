package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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

	for {
		n, err := fd.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			log.Fatalf("failed to read from %s: %s\n", inputFilePath, err)
			panic(err)
		}

		fmt.Printf("read: %s\n", string(buffer[:n]))
	}
}
