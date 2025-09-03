// Package server implements a simple HTTP server.
package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener       net.Listener
	isServerClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
	}

	go server.listen()

	return server, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isServerClosed.Load() {
				return
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}()

	response := "HTTP/1.1 200 OK\r\n" + // Status line
		"Content-Type: text/plain\r\n" + // Example header
		"Content-Length: 13\r\n" + // Content length header
		"\r\n" + // Blank line to separate headers from the body
		"Hello World!\n" // Body

	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}
}

func (s *Server) Close() error {
	s.isServerClosed.Store(true)

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}
