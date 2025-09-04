// Package server implements a simple HTTP server.
package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/itsjoeoui/httpfromtcp/internal/response"
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

	err := response.WriteStatusLine(conn, response.StatusCodeOK)
	if err != nil {
		log.Printf("Failed to write status line: %v", err)
		return
	}

	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Printf("Failed to write headers: %v", err)
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
