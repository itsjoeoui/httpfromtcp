// Package server implements a simple HTTP server.
package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/itsjoeoui/httpfromtcp/internal/request"
	"github.com/itsjoeoui/httpfromtcp/internal/response"
)

type Server struct {
	listener       net.Listener
	isServerClosed atomic.Bool
	handler        Handler
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(handler Handler, port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		handler:  handler,
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

	writer := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Failed to parse request: %v", err)

		err := writer.WriteStatusLine(response.StatusCodeBadRequest)
		if err != nil {
			log.Printf("Failed to write status line: %v", err)
		}

		body := []byte(err.Error())

		headers := response.GetDefaultHeaders(len(body))
		err = writer.WriteHeaders(headers)
		if err != nil {
			log.Printf("Failed to write headers: %v", err)
		}

		_, err = writer.WriteBody(body)
		if err != nil {
			log.Printf("Failed to write body: %v", err)
		}
		return
	}

	s.handler(writer, req)
}

func (s *Server) Close() error {
	s.isServerClosed.Store(true)

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}
