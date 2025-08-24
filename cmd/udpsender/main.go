package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// `nc -u -l 42069` to receive messages
const sendAddress = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", sendAddress)
	if err != nil {
		log.Fatalf("failed to resolve UDP address %s: %s\n", sendAddress, err.Error())
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("failed to dial UDP address %s: %s\n", sendAddress, err.Error())
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed to close UDP connection: %s\n", err.Error())
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s ", ">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read from stdin: %s\n", err.Error())
			os.Exit(1)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatalf("failed to write to UDP address %s: %s\n", sendAddress, err.Error())
			os.Exit(1)
		}

		log.Printf("Message sent to %s: %s", sendAddress, line)
	}
}
