// Package server implements a simple HTTP server.
package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/itsjoeoui/httpfromtcp/internal/request"
	"github.com/itsjoeoui/httpfromtcp/internal/response"
)

type Server struct {
	listener       net.Listener
	isServerClosed atomic.Bool
	handler        HandlerFunc
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he HandlerError) Write(w io.Writer) {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		log.Printf("Failed to write status line: %v", err)
	}

	body := []byte(he.Message)

	headers := response.GetDefaultHeaders(len(body))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		log.Printf("Failed to write headers: %v", err)
	}

	_, err = w.Write(body)
	if err != nil {
		log.Printf("Failed to write body: %v", err)
	}
}

type HandlerFunc func(w io.Writer, r *request.Request) *HandlerError

func Serve(handler HandlerFunc, port int) (*Server, error) {
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Failed to parse request: %v", err)
		handlerErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		handlerErr.Write(conn)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(buffer, req)
	if handlerError != nil {
		handlerError.Write(conn)
		return
	}

	err = response.WriteStatusLine(conn, response.StatusCodeOK)
	if err != nil {
		log.Printf("Failed to write status line: %v", err)
		return
	}

	b := buffer.Bytes()
	headers := response.GetDefaultHeaders(len(b))
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Printf("Failed to write headers: %v", err)
		return
	}

	_, err = conn.Write(b)
	if err != nil {
		log.Printf("Failed to write body: %v", err)
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
